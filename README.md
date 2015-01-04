# Showdup

Showdup is command to show file duplication.


## Installation

First, [install Go](http://golang.org/doc/install).
(Make sure you [set GOPATH](http://golang.org/doc/code.html).)

Second, download (or update) and build Showdup:

	$ go get -u github.com/martindrlik/showdup


## Usage

Run the `showdup` binary to show file duplication of the current directory:

	$ $GOPATH/bin/showdup

Specify path by passing arguments:

	$ $GOPATH/bin/showdup ./*.jpg /second/path

The argument can be pattern. The pattern syntax is defined [here](http://golang.org/pkg/path/filepath/#Match).

## Command options

`-first=N`: use only first N bytes

To find out that two files are identical you must read all bytes. It's practical to read only first N bytes and show that files are identical in first N bytes.
