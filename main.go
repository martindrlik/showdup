// Command shows file duplicity.
package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var sums = make(map[[md5.Size]byte]map[string]struct{})

func main() {
	flag.Parse()
	for i := 0; i < flag.NArg(); i++ {
		matches, err := filepath.Glob(flag.Arg(i))
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range matches {
			report(m)
		}
	}
	if flag.NArg() == 0 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		report(dir)
	}
	for _, paths := range sums {
		if len(paths) > 1 {
			for p := range paths {
				fmt.Printf("%s\n", p)
			}
			fmt.Printf("\n")
		}
	}
}

func report(path string) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	if !info.IsDir() {
		if err := walkFn(path, info, err); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		return
	}
	if err := filepath.Walk(path, walkFn); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

func walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil
	}
	if info.IsDir() {
		return nil
	}
	h := md5.New()
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil
	}
	defer f.Close()
	if _, err := io.Copy(h, f); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil
	}
	s := [md5.Size]byte{}
	for i, b := range h.Sum(nil) {
		s[i] = b
	}
	paths, ok := sums[s]
	if !ok {
		paths = make(map[string]struct{})
		sums[s] = paths
	}
	paths[path] = struct{}{}
	return nil
}
