## Beanstalkg [![CircleCI](https://circleci.com/gh/vimukthi-git/beanstalkg.svg?style=svg)](https://circleci.com/gh/vimukthi-git/beanstalkg)

Beanstalkg is a golang implementation of [beanstalkd](https://github.com/kr/beanstalkd). Idea is to support the same set of features and protocol with the addition of
high availability and failover built in. You can read the plan.md if interested in contributing. 

Right now it supports all the basic commands to run producers and workers. i.e "use", "put", "watch", "ignore", "reserve",  "delete", "release", "bury" 
except for "reserve-with-timeout". 

I wish to complete rest of the commands soon but any help is always appreciated.

### Advantages compared to beanstalkd

- Extensible design. For example you can replace the backend storage with anything you like, just implement a simple interface and plugin.
- Implemented in golang. More readable code with support for concurrency using awesome `go routines`.
- Support for clustering(coming soon :)

### Running Locally

Please install golang(binaries will be released in the future) and then clone this repo and run,

- `go install`
- `go run main.go`

## Licensing

beanstalkg is licensed under the Apache License, Version 2.0. See [LICENSE]((https://github.com/vimukthi-git/beanstalkg/blob/master/LICENSE)) for the full license text.
