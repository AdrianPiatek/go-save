package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type ChapterDownloadInfo struct {
	chapterDirPath string
	imageURLs      []string
}

func NewChapterDownloadInfo(urls []string, chapter string, title string) ChapterDownloadInfo {
	return ChapterDownloadInfo{
		imageURLs:      urls,
		chapterDirPath: fmt.Sprintf("%s/%s", title, chapter),
	}
}

func (i ChapterDownloadInfo) Download() {
	for _, url := range i.imageURLs {
		DownloadFromUrl(url, i.chapterDirPath)
	}
}

func main() {
	title := "wind-breaker"
	downloaderCount := 3
	c := colly.NewCollector()
	chaptersURLs := make(chan string, 100)
	downloadInfo := make(chan ChapterDownloadInfo, 100)
	downloadStatus := make(chan bool, 100)

	wg := &sync.WaitGroup{}
	wg.Add(2 + downloaderCount)

	go func() {
		defer wg.Done()
		for url := range chaptersURLs {
			err := c.Visit(url)
			if err != nil {
				fmt.Println(err)
			}
		}
		close(downloadInfo)
	}()

	for range downloaderCount {
		go chanExecutor(downloadInfo, downloadStatus, wg)
	}

	c.OnHTML("article.comic", func(element *colly.HTMLElement) {
		iURLs := element.ChildAttrs("div.separator img", "src")
		if len(iURLs) == 0 {
			iURLs = element.ChildAttrs("div.wp-block-image img", "src")
		}
		chapName := element.ChildText("header.entry-header")
		downloadInfo <- NewChapterDownloadInfo(iURLs, chapName, title)
	})

	c.OnHTML("li#ceo_latest_comics_widget-3", func(element *colly.HTMLElement) {
		urls := element.ChildAttrs("a", "href")
		go printStatus(downloadStatus, len(urls))
		for _, url := range urls {
			chaptersURLs <- url
		}
		close(chaptersURLs)
	})

	err := c.Visit("https://windbreakerwebtoon.com/")
	if err != nil {
		fmt.Println(err)
	}

	wg.Wait()
}

func chanExecutor(info chan ChapterDownloadInfo, status chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range info {
		i.Download()
		status <- true
	}
}

func DownloadFromUrl(url string, chapDir string) {
	split := strings.Split(url, "/")
	image := split[len(split)-1]
	filename := fmt.Sprintf("%s-%s", chapDir, image)
	if err := DownloadFile(filename, url); err != nil {
		fmt.Println(err)
	}
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func printStatus(status chan bool, i int) {
	toDownload := i
	downloaded := 0
	for range status {
		downloaded += 1
		fmt.Printf("%d/%d\n", downloaded, toDownload)
	}
}
