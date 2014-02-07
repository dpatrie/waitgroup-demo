##Golang : Concurrency niceties##

We all have "aahaaa" moment when we learn new technologies. It happened to us when we started experimenting with Go's concurrency features. Let's pretend we want to measure the response time of a list of url's. Most programmer coming from mainstream language like PHP or Python will write something along theses lines:

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
    
    	for _, url := range urls {
    		printResponseTime(url)
    	}
    	log.Printf("All done in %s!", time.Now().Sub(start))
    }
    
    func printResponseTime(url string) {
    	log.Printf("Querying %s\n", url)
    
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
    
This will produce the following output:

    2014/02/06 13:07:56 Querying https://www.google.com
    2014/02/06 13:07:57 https://www.google.com succeded in 990.54153ms
    2014/02/06 13:07:57 Querying https://www.bing.com
    2014/02/06 13:07:58 https://www.bing.com succeded in 294.571785ms
    2014/02/06 13:07:58 Querying https://www.yahoo.com
    2014/02/06 13:07:59 https://www.yahoo.com succeded in 1.236019974s
    2014/02/06 13:07:59 Querying https://duckduckgo.com/
    2014/02/06 13:07:59 https://duckduckgo.com/ succeded in 339.576809ms
    2014/02/06 13:07:59 Querying http://www.youtube.com/watch?v=hGlyFc79BUE
    2014/02/06 13:08:00 http://www.youtube.com/watch?v=hGlyFc79BUE succeded in 559.076959ms
    2014/02/06 13:08:00 All done in 3.426691384s!
    
Well that was still pretty fast, but what if you wanted to send the request in parallels? Go newcomer's with their new "goroutine" super power will just slap the "go" keyword in front of "printResponseTime" inside the loop, as such:

    //...

    for _, url := range urls {
    	go printResponseTime(url)
    }

    //...

But if we try this, we get the following output:

    2014/02/06 13:10:50 All done in 6.428us!
    
Not quite the same result! Why is this? Well the reason is that **Go does not wait for your goroutines to complete before exiting the main thread**. You actually have to instruct the program that it has to wait for the  processing to be complete. Once they realize that, most people will then think: "Hey I need theses channels they've been talking about in the doc!". So after digging a few examples they might come up with something like this:

    //...
        
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
    
    //...

And this will work just fine! It will produce the same output than the original synchronous version. But that last waiting loop is still pretty weird right? There has to be a better way! Well that's exactly the purpose of [WaitGroup in the sync package](http://golang.org/pkg/sync/#WaitGroup). Using this new knowledge, we can then change the previous example:

    //...
    
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
    
    //...
    
    
The line count is still pretty much the same but I find that last example to be a lot easier to read. What do you think?

    
 


