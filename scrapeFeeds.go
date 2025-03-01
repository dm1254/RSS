package main
import (

	"time"
	"database/sql"
	"fmt"
	"context"
	"github.com/google/uuid"
	"github.com/dm1254/RSS/internal/database"
	_ "github.com/lib/pq"
)

func scrapeFeeds(s *State) error{
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil{
		return err 
	}
	err = s.db.MarkFeedFetched(context.Background(),nextFeed.ID)
	if err != nil{
		return err
	}
	fetchFeed, err := FetchFeed(context.Background(), nextFeed.Url)
	if err != nil{
		return err
	}
	
	for _,feed := range fetchFeed.Channel.Item{
		pubDate, err := time.Parse(time.RFC1123Z, feed.PubDate)
		if err != nil{
			fmt.Printf("Error parsing date '%s': %v. Using current time.\n", feed.PubDate, err)
			pubDate = time.Now()
		}
		Params := database.CreatePostsParams{
			ID:uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: sql.NullString{String:feed.Title, Valid:feed.Title != ""},
			Url: feed.Link,
			Description:  sql.NullString{String:feed.Description, Valid:feed.Description != ""},
			PublishedAt: pubDate,
			FeedID: nextFeed.ID,
	}
		_,err = s.db.CreatePosts(context.Background(),Params)
		if err != nil{	
			return err
			}

	}

	return nil
}
