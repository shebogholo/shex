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


### Usage
1. ``` shex -h ``` for help
2. ```shex -u https://shebogholo.com  -n 10 -d 1``` for 10 concurrent requests to https://shebogholo.com for 1 second

