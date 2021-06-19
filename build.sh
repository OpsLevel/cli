VERSION=$(git describe --tags --match "v[0-9].*" --always)
CGO_ENABLED=0 go build -o ./opslevel -ldflags="-X 'github.com/opslevel/cli/cmd.version=${VERSION}'"
chmod +x ./opslevel
mv ./opslevel /usr/local/bin/opslevel