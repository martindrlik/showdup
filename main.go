// Command shows file duplication.
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

var (
	sizes = make(map[int64][]string)
	sums  = make(map[[md5.Size]byte]map[string]struct{})
)

func main() {
	flag.Parse()
	args := make(map[string]string)
	for i := 0; i < flag.NArg(); i++ {
		matches, err := filepath.Glob(flag.Arg(i))
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range matches {
			abs, err := filepath.Abs(m)
			if err != nil {
				log.Fatal(err)
			}
			args[abs] = m
		}
	}
	if flag.NArg() == 0 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		args[dir] = dir
	}
	for _, arg := range args {
		collect(arg)
	}
	for _, paths := range sizes {
		if len(paths) < 2 {
			continue
		}
		for _, p := range paths {
			sum(p)
		}
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

func collect(path string) {
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
	sizes[info.Size()] = append(sizes[info.Size()], path)
	return nil
}

func sum(path string) {
	h := md5.New()
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	defer f.Close()
	if _, err := io.Copy(h, f); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
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
}
