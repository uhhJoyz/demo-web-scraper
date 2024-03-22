package main

import (
  "fmt"
  "log"
  "strings"
  "time"
  // "sync"
  "math/rand"
)
import "github.com/gocolly/colly/v2"

func randomAgent() string {
  UserAgents := []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (HTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.18363",
  "Mozilla/5.0 (Linux; Android 13; SM-S901B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
  "Mozilla/5.0 (Linux; Android 13; Pixel 7 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
  "Mozilla/5.0 (iPhone14,6; U; CPU iPhone OS 15_4 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) Version/10.0 Mobile/19E241 Safari/602.1",
  "Mozilla/5.0 (Windows Phone 10.0; Android 4.2.1; Microsoft; RM-1127_16056) AppleWebKit/537.36(KHTML, like Gecko) Chrome/42.0.2311.135 Mobile Safari/537.36 Edge/12.10536",
  "Mozilla/5.0 (Linux; Android 5.0.2; SAMSUNG SM-T550 Build/LRX22G) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/3.3 Chrome/38.0.2125.102 Safari/537.36",
  "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36",
  "Mozilla/5.0 (PlayStation; PlayStation 5/2.26) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0 Safari/605.1.15"}
  randomIndex := rand.Intn(len(UserAgents))

  return UserAgents[randomIndex]
}

func main() {
  // create a new collector object
  c := colly.NewCollector(
    // visit set links, crawl through all links if not set
    colly.AllowedDomains("www.toscrape.com", "*.toscrape.com", "quotes.toscrape.com"),
    // set max depth of search
    colly.MaxDepth(5),
    // set user agent to random (will be static if not overridden earlier)
    colly.UserAgent(randomAgent()),
  )

  c.AllowURLRevisit = false
  // can't get async working rn, but should be very simple, just relies on the
  // WaitGroup class from the sync package
  // c.Async = true

  // a limit rule that simulates random delays when accessing webpages
  // helps not get blocked
  c.Limit(&colly.LimitRule{
    DomainGlob: "*",
    Parallelism: 2,
    Delay: 5 * time.Second,
  })

  // base url to start crawling
  url := "https://quotes.toscrape.com/page/1/"

  // on sending request, do this
  c.OnRequest(func(r *colly.Request) {
    // changes the user agent in the request on each request sent
    r.Headers.Set("User-Agent", randomAgent())
    // fmt.Println(r.URL)
  })

  // collect elements of class quote from each HTML response
  c.OnHTML(".quote", func(e *colly.HTMLElement) {
    quote := e.ChildText("span.text")
    quote = strings.TrimSpace(quote)
    // fmt.Println("Quote: ", quote)
  })

  // scrape links, visit once found
  c.OnHTML("a", func(e *colly.HTMLElement) {
    nextPage := e.Request.AbsoluteURL(e.Attr("href"))
    // fmt.Println("\nVisiting new page", nextPage, "\n")
    c.Visit(nextPage)
  })

  // demo portion for logging
  // no logging on the above call because I don't feel like rewriting it lol
  err := c.Visit(url)
  if err != nil {
    log.Fatal(err)
  }
}
