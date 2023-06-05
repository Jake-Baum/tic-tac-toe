# Useful Commands

## Go

### Build for linux environment (Lambda) 
`$env:GOOS = "linux"`

`$env:GOARCH = "amd64"`

`$env:CGO_ENABLED = "0"`

`go build -o bin/handlers/games_handler /handlers/games_handler`

`../../../../bin/build-lambda-zip.exe -o bin/handlers/games_handler.zip bin/handlers/games_handler`