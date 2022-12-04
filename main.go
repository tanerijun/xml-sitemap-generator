package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/tanerijun/html-link-parser-go/parser"
	"github.com/tanerijun/xml-sitemap-generator/queue"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: sitemapgen <link to website>")
		os.Exit(1)
	}

	links := crawlSite(os.Args[1])
	buildXMLSitemap(links)
}

// crawlSite crawls the whole website
// and return all the unique links
func crawlSite(siteUrl string) []string {
	urlSet := make(map[string]bool)
	q := queue.New([]string{siteUrl})

	for !q.Empty() {
		page := q.Dequeue()
		pageLinks := crawlPage(page)

		for _, link := range pageLinks {
			if urlSet[link] {
				continue
			}
			urlSet[link] = true
			q.Enqueue(link)
		}
	}

	urls := make([]string, 0, len(urlSet))
	for k := range urlSet {
		urls = append(urls, k)
	}

	return urls
}

// crawlPage crawls a webpage and return all the links
// in absolute format that're of the same domain.
func crawlPage(siteUrl string) []string {
	log.Println("Crawling", siteUrl)
	resp, err := http.Get(siteUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	baseUrl := &url.URL{
		Scheme: resp.Request.URL.Scheme,
		Host:   resp.Request.URL.Host,
	}

	parsedLinks, err := parser.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	links := cleanLinks(parsedLinks, baseUrl.String())
	return links
}

// cleanLinks return a slice of Link based on the input slice
// by converting relative path to absolute and
// ignores links that don't start with the base URL.
func cleanLinks(links []parser.Link, base string) []string {
	var res []string

	for _, link := range links {
		switch {
		case strings.HasPrefix(link.Href, "/"):
			res = append(res, base+link.Href)
		case strings.HasPrefix(link.Href, base):
			res = append(res, link.Href)
		}
	}

	return res
}

// buildXMLSitemap builds a sitemap from slice of links
func buildXMLSitemap(links []string) {

	type loc struct {
		Loc string `xml:"loc"`
	}

	type urlset struct {
		XMLName xml.Name `xml:"urlset"`
		Xmlns   string   `xml:"xmlns,attr"`
		Urls    []loc    `xml:"url"`
	}

	locs := make([]loc, 0, len(links))
	for _, link := range links {
		locs = append(locs, loc{link})
	}

	us := urlset{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  locs,
	}

	f, err := os.Create("sitemap.xml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.Write([]byte(xml.Header))

	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")
	if err := enc.Encode(us); err != nil {
		panic(err)
	}
}
