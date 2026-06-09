package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	spew, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var output RSSFeed
	err = xml.Unmarshal(spew, &output)
	if err != nil {
		return nil, err
	}

	output.Channel.Title = html.UnescapeString(output.Channel.Title)
	output.Channel.Description = html.UnescapeString(output.Channel.Description)
	for i1, item := range output.Channel.Item {
		output.Channel.Item[i1].Title = html.UnescapeString(item.Title)
		output.Channel.Item[i1].Description = html.UnescapeString(item.Description)
	}
	return &output, nil
}
