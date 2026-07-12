package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"
)

//go:embed index.md.tmpl
var tmplSrc string
var tmpl = template.Must(template.New("release").Parse(tmplSrc))

const (
	baseURL = "https://api.discogs.com"
	// each releass makes 2 api calls, discogs api allows for 60 per minute, so techincally batch is 30 but to be on the safe side, we will use 28
	batchSize = 28
	userAgent = "DiscogsFetcher/1.0"
)

var client = &http.Client{Timeout: 10 * time.Second}

func fetchJSON[T any](url string) (T, error) {
	var zero T

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return zero, fmt.Errorf("creating request %s: %w", url, err)
	}
	r.Header.Set("User-Agent", userAgent)

	res, err := client.Do(r)
	if err != nil {
		return zero, fmt.Errorf("making request %s: %w", url, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("unexpected status %s for %s", res.Status, url)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return zero, fmt.Errorf("reading response %s: %w", url, err)
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return zero, fmt.Errorf("unmarshaling response %s: %w", url, err)
	}

	return result, nil
}

func fetchFile(url, path string) error {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request %s: %w", url, err)
	}
	r.Header.Set("User-Agent", userAgent)

	res, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("making request %s: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %s for %s", res.Status, url)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", path, err)
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return fmt.Errorf("writing to file %s: %w", path, err)
	}

	return nil
}

type ReleaseInfo struct {
	ReleaseID     int
	ReleaseArtist string
	ReleaseTitle  string
	ReleaseYear   int
	ReleaseAdded  string
	ReleaseURL    string
}

func getRelease(username, token, id string) error {
	fmt.Println("📥 Fetching release:", id)

	// if directory with the release already exists, early return
	matches, err := filepath.Glob("*-" + id)
	if err != nil {
		return err
	}
	if len(matches) > 0 {
		fmt.Println("Already exists, skipping:", id)
		return nil
	}

	// fetch general release, and user-specific release in parallel
	var wg sync.WaitGroup

	var release Release
	var releaseErr error
	var releaseUser ReleaseUser
	var releaseUserErr error

	wg.Go(func() {
		release, releaseErr = fetchJSON[Release](fmt.Sprintf("%s/releases/%s?token=%s", baseURL, id, token))
	})

	wg.Go(func() {
		releaseUser, releaseUserErr = fetchJSON[ReleaseUser](fmt.Sprintf("%s/users/%s/collection/releases/%s?token=%s", baseURL, username, id, token))
	})

	wg.Wait()

	if releaseErr != nil {
		return fmt.Errorf("fetching release: %w", releaseErr)
	}

	if releaseUserErr != nil {
		return fmt.Errorf("fetching user release: %w", releaseUserErr)
	}

	// check if the release exists in the user's collection
	if len(releaseUser.Releases) == 0 {
		return fmt.Errorf("release not found: %s", id)
	}

	releaseUserFirst := releaseUser.Releases[0]

	// make a folder for the release
	folderName := fmt.Sprintf("%s-%s", releaseUserFirst.DateAdded[:10], id)
	err = os.Mkdir(folderName, 0o755)
	if err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	// write a markdown file with the release information
	filePath := filepath.Join(folderName, "index.md")

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	var releaseArtists []string
	for _, artist := range releaseUserFirst.BasicInformation.Artists {
		releaseArtists = append(releaseArtists, artist.Name)
	}
	releaseAdded, err := time.Parse(time.RFC3339, releaseUserFirst.DateAdded)
	if err != nil {
		return fmt.Errorf("parsing date: %w", err)
	}
	releaseAddedFormatted := releaseAdded.Format("2006.01.02")

	data := ReleaseInfo{
		ReleaseID:     releaseUserFirst.ID,
		ReleaseArtist: strings.Join(releaseArtists, ", "),
		ReleaseTitle:  releaseUserFirst.BasicInformation.Title,
		ReleaseYear:   releaseUserFirst.BasicInformation.Year,
		ReleaseAdded:  releaseAddedFormatted,
		ReleaseURL:    release.URI,
	}

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	// get cover image url, prefer primary (break when found), then secondary
	var coverImageURL string
	for _, img := range release.Images {
		if img.Type == "primary" {
			coverImageURL = img.URI
			break
		}
		if img.Type == "secondary" || coverImageURL == "" {
			coverImageURL = img.URI
		}
	}

	if coverImageURL == "" {
		fmt.Println("no cover image found for release: ", id)
		return nil
	}

	// download image
	err = fetchFile(coverImageURL, filepath.Join(folderName, "cover.jpg"))
	if err != nil {
		fmt.Println("warning: could not download cover for", id, "-", err)
	}

	fmt.Println("✅ Release downloaded", id)
	fmt.Println("- - -")
	return nil
}

func main() {
	username := flag.String("username", "", "Discogs username to query")
	token := flag.String("token", "", "Discogs API token")
	IDs := flag.String("ids", "", "Comma-separated list of IDs to query")

	flag.Parse()

	if *username == "" {
		fmt.Println("No username provided")
		os.Exit(1)
	}

	if *token == "" {
		fmt.Println("No token provided")
		os.Exit(1)
	}

	if *IDs == "" {
		fmt.Println("No IDs provided")
		os.Exit(1)
	}

	count := 0
	for id := range strings.SplitSeq(*IDs, ",") {
		if count > 0 && count%batchSize == 0 {
			fmt.Println("Sleeping for 1 minute to avoid rate limit...")
			time.Sleep(time.Minute)
		}
		if err := getRelease(*username, *token, id); err != nil {
			fmt.Println("Error getting release: ", err)
		}
		count++
	}
}
