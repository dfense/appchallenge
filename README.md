# appchallenge
Code Challenge Exercise for seedif

## Motivation 
Given simple criteria to create a new service, that mashes up (2) other external REST service api calls from different sources.

This REST service combines the function of 
* chuck norris joke
* random name selection substituted for "chuck norris" 

There is only (1) endpoint of "/"

# build instructions
## requirements
* golang 1.12+ 

## go get
to build and run the binary from command line, this is the minimal steps you need to do
it utilizes the newer _"go modules"_ for dependency transitivity

```
go get github.com/dfense/appchallenge/cmd/namejoked
```
this should fetch all the code from github, place it in your default `$HOME/go/src directory`, and compile it into `$HOME/bin/namejoked`

then just run it from there with following command (assuming osx or linux)
```
~/go/bin/namejoked
```
then can open any browser or api call here: http://localhost:8082

## clone and build
if one is to be working in the repository, or would like the do the following:
* go vet
* build a docker container
* cross compile
please following the clone and build instructions below

```
$ git clone https://github.com/dfense/appchallenge.git
$ cd appchallenge
$ go build -o ./bin/namejoked github.com/dfense/appchallenge/cmd/namejoked
## then to run
$ ./bin/namejoked --loglevel debug ## <- or any log level, or none will default to info
```

## more build options
build exe for container, build a container and start it

```
$ GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/docker-namejoked github.com/dfense/appchallenge/cmd/namejoked
$ docker build -t appchallenge .
$ docker run -d --name namejoked -p 8080:8082 appchallenge
$ docker stop namejoked
$ docker start namejoked
```

if started with container, open any browser or api call with 

`localhost:8080 ## <- port is different for container`

to check code compliance with "go vet" or run "go test"

`
$ go test github.com/dfense/appchallenge/...
$ go vet github.com/dfense/appchallenge/...
`

## cheat sheet
a list of all these commands above and more are located in the build.txt file in root of project
