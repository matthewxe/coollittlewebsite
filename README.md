# coollittlewebsite

My own personal site where I can put anything I want in it.

Written with the power of go baby

## Building

Get these dependencies:

- `go`
- `gcc`
- tailwindcss_3 (for updating)
- [air](https://github.com/air-verse/air) (for testing)

or just use this if you use nix:

```
nix-shell
```

then build main and run

```
go build cmd/main.go
./main
```

### Testing

1. Run air to live reload the go app

```
air
```

2. Run the tailwindcli to update it each time

```
cd web/
./build/tailwindcss.sh --watch
```

## [Structure](docs/structure.org)

- structure of the code

## [Webpages](docs/webpages.org)

- docs of available webpages and directories in the site
