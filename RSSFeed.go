package main 
import (
	"io"
	"net/http"
	"encoding/xml"
	"html"
	"context"


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

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){
	client := &http.Client{}
	req,err := http.NewRequestWithContext(ctx,"GET", feedURL,nil)
	if err != nil{
		return &RSSFeed{},err
	}
	req.Header.Set("User-Agent", "gator")
	
	res, err := client.Do(req)
	if err != nil{
		return &RSSFeed{},err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil{
		return &RSSFeed{},err
	}
	rssfeed := RSSFeed{}
	if err = xml.Unmarshal(data, &rssfeed); err != nil{
		return &RSSFeed{},err
	}
	rssfeed.Channel.Title = html.UnescapeString(rssfeed.Channel.Title)
	rssfeed.Channel.Description = html.UnescapeString(rssfeed.Channel.Description)
	for i := range rssfeed.Channel.Item{
		rssfeed.Channel.Item[i].Title  = html.UnescapeString(rssfeed.Channel.Item[i].Title)
		rssfeed.Channel.Item[i].Description  = html.UnescapeString(rssfeed.Channel.Item[i].Description)
	
	}
	return &rssfeed,nil
	
	
}
