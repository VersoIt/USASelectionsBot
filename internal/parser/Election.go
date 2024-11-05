package parser

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strconv"
)

type Election struct {
	rootUrl                  string
	firstCandidateContainer  string
	secondCandidateContainer string
}

func NewElection(rootUrl, firstCandidateContainer, secondCandidateContainer string) *Election {
	return &Election{rootUrl, firstCandidateContainer, secondCandidateContainer}
}

func (e *Election) Parse() (int, int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", e.rootUrl, nil)
	if err != nil {
		return 0, 0, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return 0, 0, err
	}

	firstCandidateRes, err := strconv.Atoi(doc.Find(e.firstCandidateContainer).First().Text())
	if err != nil {
		return 0, 0, err
	}

	secondCandidateRes, err := strconv.Atoi(doc.Find(e.secondCandidateContainer).First().Text())
	if err != nil {
		return 0, 0, err
	}

	return firstCandidateRes, secondCandidateRes, nil
}
