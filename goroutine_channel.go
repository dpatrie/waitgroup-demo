package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	start := time.Now()

	urls := []string{
		"https://www.google.com",
		"https://www.bing.com",
		"https://www.yahoo.com",
		"https://duckduckgo.com/",
		"http://www.youtube.com/watch?v=hGlyFc79BUE",
	}

	done := make(chan bool)

	for _, url := range urls {
		go func(u string) {
			printResponseTime(u)
			done <- true
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		<-done
	}

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
