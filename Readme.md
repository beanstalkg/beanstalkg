## Beanstalkg [![CircleCI](https://circleci.com/gh/vimukthi-git/beanstalkg.svg?style=svg)](https://circleci.com/gh/vimukthi-git/beanstalkg) [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/beanstalkg-chat/Lobby)

Beanstalkg is a go implementation of [Beanstalkd](https://github.com/kr/beanstalkd) **a fast, general-purpose work queue**. 
Idea is to support the same set of features and protocol with the addition of
high availability and failover built in. You can read the plan.md if interested in contributing. 

Right now it supports all the basic commands to run producers and workers. i.e "use", "put", "watch", "ignore", "reserve",  "delete", "release", "bury", "reserve-with-timeout". 

I wish to complete rest of the commands soon but any help is always appreciated.

### Advantages compared to beanstalkd

- Extensible design. For example you can replace the backend storage with anything you like, just implement a simple interface and plugin.
- Implemented in golang. More readable code with support for concurrency using awesome `go routines`.
- Support for clustering(coming soon :)


### User guide

Beanstalkg is currently only released as a docker image for users. Latest release is v0.0.3. Assuming you already have a 
working docker engine installation, you can start a Beanstalkg instance with following steps,

- Run command `docker run -p 11300:11300 beanstalkg/beanstalkg:v0.0.3`. This will start the beanstalkg server in the foreground.
 The server starts listening on port 11300.
- Now you can connect to the server with any client library available to [beanstalkd](https://github.com/kr/beanstalkd/wiki/Client-Libraries). 
 eg: Using [official go client](https://github.com/Beanstalkg/beanstalk)
    ```
    // Produce jobs:
    c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
    id, err := c.Put([]byte("hello"), 1, 0, 120*time.Second)
    
    // Consume jobs:
    c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
    id, body, err := c.Reserve(5 * time.Second)
    
    ```
Some introductory slides can be found [here](https://www.slideshare.net/VimukthiWickramasing/beanstalkg-76573237)

### Developer guide

Please install golang then with `GOPATH` set correctly run,

- `go get github.com/Beanstalkg/beanstalkg`
- Add `$GOAPTH/bin` to the $PATH and run `beanstalkg`

## Licensing

beanstalkg is licensed under the MIT License. See [LICENSE](https://github.com/vimukthi-git/beanstalkg/blob/master/LICENSE) for the full license text.
