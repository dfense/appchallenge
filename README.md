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
