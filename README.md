# Showdup

Command showdup shows files which are probably identical.


## Installation

First, [install Go](http://golang.org/doc/install).
(Make sure you [set GOPATH](http://golang.org/doc/code.html).)

Second, download (or update) and build Showdup:

	$ go get -u github.com/martindrlik/showdup


## Usage

Run the `showdup` binary to show identical files in the current directory:

	$ $GOPATH/bin/showdup

Specify path by passing arguments:

	$ $GOPATH/bin/showdup ./*.jpg /second/path

The argument can be pattern. The pattern syntax is defined [here](http://golang.org/pkg/path/filepath/#Match).

## Command options

`first=512`: read only first 512 bytes
`immediate=false`: print immediately; no order, but "hash:" prefix
