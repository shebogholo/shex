## shex - Load Testing CLI 


### Go Configuration
1. Add the following to ``` ~/.bashrc ``` on Ubuntu Machine
```bash
    export GOPATH=$HOME/go
    export GOBIN=$GOPATH/bin
    export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin
```
2. Then on terminal run ``` source ~/.bashrc ```


### Build and Install Go package
1. ``` go build -o bin/shex ```
2. ``` go install ```

### Compile 
1. ``` GOOS=darwin GOARCH=amd64 go build -o bin/shex-amd64-darwin main.go ``` 
2. ``` GOOS=linux GOARCH=amd64 go build -o bin/shex-amd64-linux main.go ```
3. ``` GOOS=windows GOARCH=amd64 go build -o bin/shex-amd64-windows.exe main.go ```

### Installation Script for Linux and Mac OS
```bash
    sudo curl -sSfL https://raw.githubusercontent.com/shebogholo/shex/main/install.sh | sh
```

### Usage
1. ``` shex -h ``` for help
2. ``` shex -u https://shebogholo.com  -n 10 -d 1``` for 10 concurrent requests to https://shebogholo.com for 1 second