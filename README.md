Watcher [![GoDoc](https://godoc.org/github.com/canthefason/go-watcher?status.svg)](https://godoc.org/github.com/canthefason/go-watcher) [![Build Status](https://travis-ci.org/canthefason/go-watcher.svg?branch=master)](https://travis-ci.org/canthefason/go-watcher)
=======

Watcher is a command line tool inspired by [fresh](https://github.com/pilu/fresh) and used for watching .go file changes, and restarting the app in case of an update/delete/add operation.

Most of the existing file watchers have a configuration burden, and even though Go has a really short build time, this configuration burden makes your binaries really hard to run right away. With Watcher, we aimed simplicity in configuration, and tried to make it as simple as possible.

## Installation

  Get the package with:

  `go get github.com/canthefason/go-watcher`

  Install the binary under go/bin folder:

  `go install github.com/canthefason/go-watcher/cmd/watcher`

  After this step, please make sure that your go/bin folder is appended to PATH environment variable.

## Usage

  `cd /path/to/myapp`

Start watcher:

  `watcher`

Watcher works like your native package binary. You can pass all your existing package arguments to the Watcher, which really lowers the learning curve of the package, and makes it practical.

##### Current app usage
  `myapp -c config -p 7000 -h localhost`

##### With watcher
  `watcher -c config -p 7000 -h localhost`

As you can see nothing changed between these two calls. When you run the command, Watcher starts watching folders recursively, starting from the current working directory. It only watches .go and .tmpl files and ignores hidden folders and _test.go files.

##### Package dependency

By default Watcher recursively watches all files/folders under working directory. If you prefer to use it like `go run`, you can call watcher with -run flag anywhere you want (we assume that your GOPATH is properly set).

  `watcher -c config -run github.com/username/somerootpackagename`

For the cases where your main function is in another directory other than the dependant package, you can do this by passing a different package name to -watch parameter.

  `watcher -c config -run github.com/username/somerootpackagename -watch github.com/username`


##### Vendor directory
Since Globs and some optional folder arrays will make it harder to configure, we are not planning to have support for a configurable watched folder structure. Only configuration we have here is, by default we have excluded vendor/ folder from watched directories. If your intention is making some changes in place, you can set -watch-vendor flag as "true", and start watching vendor directory.

## Watcher in Docker

If you want to run Watcher in a containerized local environment, you can achieve this by using [canthefason/go-watcher](https://hub.docker.com/r/canthefason/go-watcher/) image in Docker Hub. There is an example project under [/docker-example](https://github.com/canthefason/go-watcher/tree/dockerfile-gvm/docker-examples) directoy. Let's try to dockerize this example code first.

In our example, we are creating a server that listens to port 7000 and responds to all clients with "watcher is running" string. The most essential thing to run your code in Docker is, mounting your project volume to a container. In the containerized Watcher, our GOPATH is set to /go directory by default, so you need to mount your project to this GOPATH.

  `docker run -v /path/to/hello:/go/src/hello -p 7000:7000 canthefason/go-watcher watcher -run hello`

Containerized Watcher also supports different versions of Go by leveraging [gvm](https://github.com/moovweb/gvm). Currently it only supports major versions right now. If you don't set anything, by default Watcher will pick version 1.7. If you want to change the Go version, you can use GO_VERSION environment variable. Currently it only supports 1.4, 1.5, 1.6, 1.7 at the moment

   `docker run -v /path/to/hello:/go/src/hello -e GO_VERSION=1.6 -p 7000:7000 canthefason/go-watcher watcher -run hello`

To provide a more structured repo, we also integrated a docker-compose manifest file. That file already handles volume mounting operation that and exposes the port to the localhost. With docker-compose the only thing that you need to do from the root, invoking `docker-compose up

#### Known Issues
On Mac OS X, when you make a tls connection, you can get a message like: x509: `certificate signed by unknown authority`

You can resolve this problem by setting CGO_ENABLED=0
https://github.com/golang/go/issues/14514
https://codereview.appspot.com/22020045

## Author

* [Can Yucel](http://canthefason.com)

## License

The MIT License (MIT) - see LICENSE.md for more details


