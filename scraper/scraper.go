package scraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"mvdan.cc/xurls"
)

const (
	strokeOrderURLFmt string = "http://www.strokeorder.info/mandarin.php?q=%s"
	translationURLFmt string = "https://bkrs.info/slovo.php?ch=%s"
)

func ensureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}

func downloadFile(distpath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	absPath, _ := filepath.Abs(filepath.Join("./resources", distpath))
	ensureDir(absPath)
	// Create the file
	out, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func extractURLPath(rawURL string) string {
	url, _ := url.Parse(rawURL)

	return url.Path
}

func scrapeTranslation(word string) {
	// Request the HTML page.
	res, err := http.Get(fmt.Sprintf(translationURLFmt, word))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("img.pointer").Each(func(i int, s *goquery.Selection) {
		onclickString, ok := s.Attr("onclick")

		if ok {
			rxRelaxed := xurls.Relaxed()
			audioLink := rxRelaxed.FindString(onclickString)
			path := extractURLPath(audioLink)
			if err := downloadFile(path, audioLink); err != nil {
				panic(err)
			}
			fmt.Printf("Found MP3 link (%s)\n", audioLink)
			return
		}

	})
}

func scrapeStrokeOrder(word string) {
	// Request the HTML page.
	res, err := http.Get(fmt.Sprintf(strokeOrderURLFmt, word))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		srcLink, ok := s.Attr("src")

		if ok && strings.Contains(srcLink, "http://bishun.strokeorder.info/characters/") {
			path := extractURLPath(srcLink)
			if err := downloadFile(path, srcLink); err != nil {
				panic(err)
			}
			fmt.Printf("Found image (%s)\n", srcLink)
			return
		}
	})
}

// Scrape : Runs process to scrape all necessary info about ieroglyph or word
func Scrape(word string) string {
	fmt.Println("Scrape start")

	scrapeTranslation(word)
	scrapeStrokeOrder(word)

	return word
}
