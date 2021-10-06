package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gocolly/colly"
	youtube "github.com/kkdai/youtube/v2"
)

func main() {

	urlArr := findIframeUrl()
	for i := 0; i < len(urlArr); i++ {
		if strings.Contains(urlArr[i], "https://www.youtube.com") {
			downloadVideo(urlArr[i])
		}
	}

}

func findIframeUrl() []string {
	var urlArr []string
	collector := colly.NewCollector(
		colly.AllowedDomains("firstsiteguide.com", "www.firstsiteguide.com", "youtube.com", "www.youtube.com"),
	)

	collector.OnHTML("iframe", func(element *colly.HTMLElement) {
		vdo_link := element.Attr("data-src")
		urlArr = append(urlArr, vdo_link)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	url := "https://firstsiteguide.com/what-is-blog/"
	collector.Visit(url)
	return urlArr
}

func downloadVideo(url string) {

	fmt.Println("url  :", url)
	url = strings.Replace(url, "https://www.youtube.com/embed/", "", 1)
	url = strings.Replace(url, "?feature=oembed", "", 1)
	url = strings.TrimSpace(url)
	fmt.Println(">>>", url, "<<<")

	videoID := url
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	stream, _, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		panic(err)
	}

	file, err := os.Create(videoID + ".mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}

	fmt.Println("Download Complete for", videoID+".mp4")

}
