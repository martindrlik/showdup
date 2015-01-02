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
	sums1 = make(map[[md5.Size]byte][]string)
	sums2 = make(map[[md5.Size]byte][]string)
)

func main() {
	flag.Parse()
	paths := make(map[string]string)
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
			paths[abs] = m
		}
	}
	if flag.NArg() == 0 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		paths[dir] = dir
	}
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		} else if !info.IsDir() {
			if err := walkFn(p, info, err); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
		} else if err := filepath.Walk(p, walkFn); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
	for _, paths := range sizes {
		if len(paths) < 2 {
			continue
		}
		for _, p := range paths {
			if s, ok := sumN(p, 512); ok {
				sums1[s] = append(sums1[s], p)
			}
		}
	}
	for _, paths := range sums1 {
		if len(paths) < 2 {
			continue
		}
		for _, p := range paths {
			if s, ok := sumN(p, -1); ok {
				sums2[s] = append(sums2[s], p)
			}
		}
	}
	for _, paths := range sums2 {
		for _, p := range paths {
			fmt.Printf("%s\n", p)
		}
		fmt.Printf("\n")
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

func sumN(path string, n int64) ([16]byte, bool) {
	s := [md5.Size]byte{}
	h := md5.New()
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return s, false
	}
	defer f.Close()
	if n > 0 {
		_, err = io.CopyN(h, f, n)
	} else {
		_, err = io.Copy(h, f)
	}
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return s, false
	}
	for i, b := range h.Sum(nil) {
		s[i] = b
	}
	return s, true
}
