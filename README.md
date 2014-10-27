Goldorf
=======

Goldorf is a command line tool inspired by [fresh](https://github.com/pilu/fresh) and used for watching .go file changes, and restarting the app in case of an update/delete/add operation.

File change event listening infrastructure depends on stable version (.v1) of [fsnotify](https://github.com/go-fsnotify/fsnotify)

## Installation

  Get the package with:

  `go get github.com/canthefason/goldorf`

  Install the binary under go/bin folder:

  `go install github.com/canthefason/goldorf`

  If not added please append your go/bin folder to PATH environment variable.

## Usage

  `cd /path/to/myapp`

Start goldorf:

  `goldorf`

Goldorf works like your native binary package. You can pass all the arguments that you are currently using.

##### Current app usage
  `myapp -c config -p 7000 -h localhost`

##### With goldorf
  `goldorf -c config -p 7000 -h localhost`

When you run the command it starts watching folders recursively, starting from the current working directory. It only watches .go and .tmpl files and ignores hidden folders.

##### Rootpackage (God mode on)
  `goldorf -c config -rootpackage github.com/username/somerootpackagename`

For the cases where your app depends on a few packages that are still WIP, and you want to watch all the changes including those packages, you can pass a root package name to goldorf.

Micro management of the watched packages are not supported.

## Name inspiration

Package gets its name from one of the old spectators of The Muppet Show: Waldorf.

## Author

* [Can Yucel](http://canthefason.com)

## License

The MIT License (MIT) - see LICENSE.md for more details


