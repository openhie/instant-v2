mkdir -p bin
GOOS=darwin GOARCH=amd64 go build && mv ./cli ./bin/goinstant-macos \
&& GOOS=linux GOARCH=amd64 go build && mv ./cli ./bin/goinstant-linux \
&& GOOS=windows GOARCH=amd64 go build && mv ./cli.exe ./bin/goinstant.exe\
&& go clean
