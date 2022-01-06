
echo "start 1"
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build main.go  && ./main
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go  && ./main
CGO_ENABLED=0 GOOS=darwin GOARCH=arm go build main.go  && ./main
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build main.go  && ./main

echo "start 2"
CGO_ENABLED=0 GOOS=android GOARCH=amd64 go build main.go  && ./main
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go  && ./main
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go  && ./main
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go  && ./main
