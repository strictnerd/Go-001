package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func main() {
	g := new(errgroup.Group)
	var urls = []string{
		"https://github.com/",
		"https://www.baidu.com/",
	}

	for _, url := range urls {
		url := url
		g.Go(func() error {
			resp, err := http.Get(url)
			if err == nil {
				fmt.Println(resp.Header)
				resp.Body.Close()
			}
			return err
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}
