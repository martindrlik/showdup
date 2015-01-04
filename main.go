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
	"runtime/pprof"
)

var (
	sizes = make(map[int64][]string)

	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
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
	for size, paths := range sizes {
		if len(paths) < 2 {
			continue
		}
		if size > 512 {
			for _, paths := range sumAll(paths, 512) {
				if len(paths) > 1 {
					print(sumAll(paths, 0))
				}
			}
		} else {
			print(sumAll(paths, 0))
		}
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

func sumAll(paths []string, n int64) map[[md5.Size]byte][]string {
	sum := make(map[[md5.Size]byte][]string)
	for _, p := range paths {
		if s, ok := sumN(p, n); ok {
			sum[s] = append(sum[s], p)
		}
	}
	return sum
}

func sumN(path string, n int64) (s [md5.Size]byte, ok bool) {
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

func print(m map[[md5.Size]byte][]string) {
	for _, paths := range m {
		if len(paths) > 1 {
			for _, p := range paths {
				fmt.Printf("%s\n", p)
			}
			fmt.Print("\n")
		}
	}
}
