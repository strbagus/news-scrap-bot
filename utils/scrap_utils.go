package utils

import (
	"crypto/tls"
	"fmt"
	"gobot/models"
	"net/http"
	"os"
	"time"
	"github.com/PuerkitoBio/goquery"
)

func GetData() []models.NewsType {
	var news []models.NewsType
	url := os.Getenv("TARGET_URL")

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return nil
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request: ", err)
		return nil
	}

	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find(".event-post").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3").Text()
		link := s.Find("a").AttrOr("href", "")
		news = append(news, models.NewsType{Title: title, Link: link})
	})
	return news
}

func CompareData(oldData []models.NewsType, newData []models.NewsType) []models.NewsType {
	listNew := []models.NewsType{}
	for _, new := range newData {
		isNewExist := false
		for _, old := range oldData {
			if old.Title == new.Title {
				isNewExist = true
			}
		}
		if !isNewExist {
			listNew = append(listNew, new)
		}
	}
    return listNew
}
