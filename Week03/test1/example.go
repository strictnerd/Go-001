package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type result struct {
	path string
	sum  [md5.Size]byte
}

func main() {
	m, err := MD5All(context.Background(), "I:\\BaiduNetdiskDownload")
	if err != nil {
		log.Fatal("error")
	}
	for k, sum := range m {
		fmt.Printf("%s:\t%x\n", k, sum)
	}
}

func MD5All(ctx context.Context, root string) (map[string][md5.Size]byte, error) {
	g, ctx := errgroup.WithContext(ctx)

	paths := make(chan string)

	g.Go(func() error {
		defer close(paths)
		return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	})

	c := make(chan result)
	const countDegitser = 20
	for i := 0; i < countDegitser; i++ {
		g.Go(func() error {
			for path := range paths {
				file, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				select {
				case c <- result{path, md5.Sum(file)}:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
	}

	go func() {
		g.Wait()
		close(c)
	}()

	m := make(map[string][md5.Size]byte)
	for r := range c {
		m[r.path] = r.sum
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return m, nil
}
