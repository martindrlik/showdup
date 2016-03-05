// showdup displays identical files in specified directories.
// usage: showdup [directory...]

package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

const FirstNBytes = 128

var (
	size = make(map[int64][]string)
)

func main() {
	if len(os.Args) <= 1 {
		readDir(".")
	}
	for i := 1; i < len(os.Args); i++ {
		readDir(os.Args[i])
	}
	for _, same := range size {
		if len(same) < 2 {
			continue
		}
		readFiles(same)
	}
}

// readDir groups files in dir by file size.
func readDir(dir string) {
	f, err := os.Open(dir)
	if check(err, dir) {
		return
	}
	defer f.Close()
	for {
		fis, err := f.Readdir(100)
		if err != io.EOF && check(err, dir) {
			break
		}
		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}
			s := fi.Size()
			name := path.Join(dir, fi.Name())
			size[s] = append(size[s], name)
		}
		if err == io.EOF {
			break
		}
	}
}

// readFiles prints identical files. It groups files by first n bytes given by
// FirstNBytes then if there is multiple with the same first n bytes it
// continues to test these files and eventually prints if files are identical.
func readFiles(files []string) {
	type First [FirstNBytes]byte
	first := make(map[First][]string)
	readFirst := func(name string) {
		f, err := os.Open(name)
		if check(err, name) {
			return
		}
		defer f.Close()
		var b First
		_, err = f.Read(b[:]) // b[:] is hack to better avoid in code that matters
		if check(err, name) {
			return
		}
		first[b] = append(first[b], name)
	}
	for _, name := range files {
		readFirst(name)
	}
	for _, same := range first {
		if len(same) < 2 {
			continue
		}
		sumFiles(same)
	}
}

// sumFiles groups files by md5 checksum and eventually prints if files are
// identical.
func sumFiles(files []string) {
	sum := make(map[[md5.Size]byte][]string)
	for _, name := range files {
		b, err := ioutil.ReadFile(name)
		if check(err, name) {
			continue
		}
		s := md5.Sum(b)
		sum[s] = append(sum[s], name)
	}
	for _, same := range sum {
		if len(same) < 2 {
			continue
		}
		for _, name := range same {
			fmt.Println(name)
		}
		fmt.Println()
	}
}

func check(err error, name string) bool {
	if err == nil {
		return false
	}
	fmt.Fprintf(os.Stderr, "%s: %s: %v\n", os.Args[0], name, err)
	return true
}
