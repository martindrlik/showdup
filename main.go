// Command showdup shows files which are probably identical.
//
package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime/pprof"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	first      = flag.Int64("first", 512, "read only first 512 bytes")
	immediate  = flag.Bool("immediate", false, "print immediately; no order, but \"hash:\" prefix")
)

func main() {
	flag.Usage = usage
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	ch := make(chan File)
	go func() {
		if flag.NArg() == 0 {
			readFile(".", ch)
		}
		for i := 0; i < flag.NArg(); i++ {
			readFile(flag.Arg(i), ch)
		}
		close(ch)
	}()
	byHash := make(map[string][]string)
	for file := range ch {
		if file.err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", file.err)
			continue
		}
		files := append(byHash[file.hash], file.name)
		byHash[file.hash] = files
		if !*immediate {
			continue
		}
		switch len(files) {
		case 1:
		case 2:
			for _, name := range files {
				fmt.Fprintf(os.Stdout, "%v:%v\n", file.hash, name)
			}
		default:
			fmt.Fprintf(os.Stdout, "%v:%v\n", file.hash, file.name)
		}
	}
	if *immediate {
		return
	}
	for _, files := range byHash {
		for _, name := range files {
			fmt.Fprintf(os.Stdout, "%v\n", name)
		}
		fmt.Fprintf(os.Stdout, "\n")
	}
}

type File struct {
	err  error
	hash string
	name string
}

var bySize = make(map[int64][]string)

func readFile(name string, ch chan<- File) {
	f, err := os.Open(name)
	if err != nil {
		ch <- File{err: err}
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		ch <- File{err: err}
		return
	}
	if fi.IsDir() {
		readDir(name, f, ch)
		return
	}
	size := fi.Size()
	files, ok := bySize[size]
	if !ok {
		bySize[size] = []string{name}
		return
	}
	ch0 := make(chan File)
	if len(files) > 1 {
		go hash(name, f, ch0)
		ch <- <-ch0
		return
	}
	bySize[size] = append(files, name)
	go hash(name, f, ch0)
	go hash(files[0], nil, ch0)
	ch <- <-ch0
	ch <- <-ch0
}

func hash(name string, r io.Reader, ch chan<- File) {
	if r == nil {
		f, err := os.Open(name)
		if err != nil {
			ch <- File{err: err}
			return
		}
		defer f.Close()
		r = f
	}
	if *first < 0 || *first > 8*512 {
		*first = 512
	}
	b := make([]byte, *first)
	n, err := r.Read(b)
	if err != nil && err != io.EOF {
		ch <- File{err: err}
		return
	}
	hash := fmt.Sprintf("%x", md5.Sum(b[:n]))
	ch <- File{name: name, hash: hash}
}

func readDir(dir string, f *os.File, ch chan<- File) {
repeat:
	names, err := f.Readdirnames(64)
	if err != nil && err != io.EOF {
		ch <- File{err: err}
		return
	}
	for _, name := range names {
		readFile(path.Join(dir, name), ch)
	}
	if err == io.EOF {
		return
	}
	goto repeat
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: showdup [options] [file ...]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}
