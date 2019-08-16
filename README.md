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
it utilizes the newer "go modules" for dependency transitivity

`
go get github.com/dfense/appchallenge/cmd/namejoked
`

## clone and build
if one is to be working in the repository, or would like the do the following:
* go vet
* build a docker container
* cross compile

please follwing the clone and build instructions below

`
$ git clone https://github.com/dfense/appchallenge.git
$ cd appchallenge
$ go build -o ./bin/namejoked github.com/dfense/appchallenge/cmd/namejoked
## then to run
$ ./bin/namejoked --loglevel debug ## <- or any log level, or none will default to info
`

## utils
build a container and start it
`
$ docker build -t appchallenge .
$ docker run -d --name namejoked -p 8080:8082 appchallenge
`

open a browser or hit api with 
`localhost:8080`

to check code compliance with "go vet" or run "go test"
`
$ go test github.com/dfense/appchallenge/...
$ go vet github.com/dfense/appchallenge/...
`
