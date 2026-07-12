# Discogs user release downloader

## Build

```
go build -o drd
```

or

```
make build
```

## Use

Single:

```
./drd --username=XXX --token=XXX --ids=12345678
```

Or multiple:

```
./drd --username=XXX --token=XXX --ids=12345678,23456789
```

### Flags

- `--username` - Discogs username (required)
- `--token` - Discogs API token (required)
- `--ids` - comma-separated release IDs (technically not required, but it makes no sense to skip it)
