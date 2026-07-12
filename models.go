package main

// Release is the shape returned by GET /releases/{id}.
type Release struct {
	ID                int          `json:"id"`
	Status            string       `json:"status"`
	Year              int          `json:"year"`
	ResourceURL       string       `json:"resource_url"`
	URI               string       `json:"uri"`
	Artists           []Artist     `json:"artists"`
	ArtistsSort       string       `json:"artists_sort"`
	Labels            []Entity     `json:"labels"`
	Series            []Entity     `json:"series"`
	Companies         []Entity     `json:"companies"`
	Formats           []Format     `json:"formats"`
	DataQuality       string       `json:"data_quality"`
	Community         Community    `json:"community"`
	FormatQuantity    int          `json:"format_quantity"`
	DateAdded         string       `json:"date_added"`
	DateChanged       string       `json:"date_changed"`
	NumForSale        int          `json:"num_for_sale"`
	LowestPrice       float64      `json:"lowest_price"`
	MasterID          int          `json:"master_id"`
	MasterURL         string       `json:"master_url"`
	Title             string       `json:"title"`
	Country           string       `json:"country"`
	Released          string       `json:"released"`
	ReleasedFormatted string       `json:"released_formatted"`
	Identifiers       []Identifier `json:"identifiers"`
	Videos            []Video      `json:"videos"`
	Genres            []string     `json:"genres"`
	Styles            []string     `json:"styles"`
	Tracklist         []Track      `json:"tracklist"`
	ExtraArtists      []Artist     `json:"extraartists"`
	Images            []Image      `json:"images"`
	Thumb             string       `json:"thumb"`
	EstimatedWeight   int          `json:"estimated_weight"`
	BlockedFromSale   bool         `json:"blocked_from_sale"`
	IsOffensive       bool         `json:"is_offensive"`
}

// Artist covers both "artists" and "extraartists" entries.
type Artist struct {
	Name         string `json:"name"`
	ANV          string `json:"anv"`
	Join         string `json:"join"`
	Role         string `json:"role"`
	Tracks       string `json:"tracks"`
	ID           int    `json:"id"`
	ResourceURL  string `json:"resource_url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

// Entity covers "labels", "companies" and "series" — they share a shape.
type Entity struct {
	Name           string `json:"name"`
	Catno          string `json:"catno"`
	EntityType     string `json:"entity_type"`
	EntityTypeName string `json:"entity_type_name"`
	ID             int    `json:"id"`
	ResourceURL    string `json:"resource_url"`
	ThumbnailURL   string `json:"thumbnail_url"`
}

type Format struct {
	Name         string   `json:"name"`
	Qty          string   `json:"qty"`
	Descriptions []string `json:"descriptions"`
}

type Community struct {
	Have         int       `json:"have"`
	Want         int       `json:"want"`
	Rating       Rating    `json:"rating"`
	Submitter    UserRef   `json:"submitter"`
	Contributors []UserRef `json:"contributors"`
	DataQuality  string    `json:"data_quality"`
	Status       string    `json:"status"`
}

type Rating struct {
	Count   int     `json:"count"`
	Average float64 `json:"average"`
}

// UserRef covers "submitter" and "contributors" entries.
type UserRef struct {
	Username    string `json:"username"`
	ResourceURL string `json:"resource_url"`
}

type Identifier struct {
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type Video struct {
	URI         string `json:"uri"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Embed       bool   `json:"embed"`
}

type Track struct {
	Position     string   `json:"position"`
	Type         string   `json:"type_"`
	Title        string   `json:"title"`
	ExtraArtists []Artist `json:"extraartists"`
	Duration     string   `json:"duration"`
}

type Image struct {
	Type        string `json:"type"`
	URI         string `json:"uri"`
	ResourceURL string `json:"resource_url"`
	URI150      string `json:"uri150"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

// ReleaseUser is the shape returned by
// GET /users/{username}/collection/releases/{id} — a paginated wrapper
// around the collection entries for that release.
type ReleaseUser struct {
	Pagination Pagination          `json:"pagination"`
	Releases   []CollectionRelease `json:"releases"`
}

type Pagination struct {
	Page    int               `json:"page"`
	Pages   int               `json:"pages"`
	PerPage int               `json:"per_page"`
	Items   int               `json:"items"`
	URLs    map[string]string `json:"urls"`
}

type CollectionRelease struct {
	ID               int              `json:"id"`
	InstanceID       int              `json:"instance_id"`
	DateAdded        string           `json:"date_added"`
	Rating           int              `json:"rating"`
	BasicInformation BasicInformation `json:"basic_information"`
	FolderID         int              `json:"folder_id"`
	Notes            []Note           `json:"notes"`
}

type BasicInformation struct {
	ID          int      `json:"id"`
	MasterID    int      `json:"master_id"`
	MasterURL   string   `json:"master_url"`
	ResourceURL string   `json:"resource_url"`
	Thumb       string   `json:"thumb"`
	CoverImage  string   `json:"cover_image"`
	Title       string   `json:"title"`
	Year        int      `json:"year"`
	Formats     []Format `json:"formats"`
	Labels      []Entity `json:"labels"`
	Artists     []Artist `json:"artists"`
	Genres      []string `json:"genres"`
	Styles      []string `json:"styles"`
}

type Note struct {
	FieldID int    `json:"field_id"`
	Value   string `json:"value"`
}
