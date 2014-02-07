package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	wg := sync.WaitGroup{}

	urls := []string{
		"https://www.google.com",
		"https://www.bing.com",
		"https://www.yahoo.com",
		"https://duckduckgo.com/",
		"http://www.youtube.com/watch?v=hGlyFc79BUE",
	}

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			printResponseTime(u)
		}(url)
	}

	wg.Wait()

	log.Printf("All done in %s!", time.Now().Sub(start))
}

func printResponseTime(url string) {
	t0 := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("%s failed with error: %s\n", url, err)
	} else {
		_, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("%s failed with error: %s\n", url, err)
		} else {
			log.Printf("%s succeded in %s\n", url, time.Now().Sub(t0))
		}
	}
}
