go get -u github.com/gobuffalo/packr/packr
mkdir -p releases
GOOS=darwin GOARCH=amd64 packr build; mv ./GoSH ./releases/darwin-GoSH
GOOS=linux GOARCH=amd64 packr build; mv ./GoSH ./releases/linux-GoSH
GOOS=windows GOARCH=386 packr build; mv ./GoSH.exe ./releases/GoSH.exe
packr clean
