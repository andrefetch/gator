package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/andrefetch/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(state *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid time_between_reqs: %w", err)
	}
	if timeBetweenRequests <= 0 {
		return fmt.Errorf("time_between_reqs must be positive")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for {
		if ctx.Err() != nil {
			return nil
		}
		if err := scrapeFeeds(ctx, state); err != nil && ctx.Err() == nil {
			log.Printf("error collecting feed: %v", err)
		}

		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func scrapeFeeds(ctx context.Context, state *State) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	feed, err := state.DB.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("get next feed: %w", err)
	}

	if err := state.DB.MarkFeedFetched(ctx, feed.ID); err != nil {
		return fmt.Errorf("mark feed fetched: %w", err)
	}

	fmt.Printf("Fetching feed: %s (%s)\n", feed.Name, feed.Url)
	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("fetch feed %q: %w", feed.Url, err)
	}
	for _, item := range rssFeed.Channel.Item {
		postURL := strings.TrimSpace(item.Link)
		if postURL == "" {
			log.Printf("skipping post %q from feed %q: missing URL", item.Title, feed.Url)
			continue
		}

		publishedAt, err := parsePublishedAt(item.PubDate)
		if err != nil {
			log.Printf("post %q from feed %q: %v", postURL, feed.Url, err)
		}

		now := time.Now().UTC()
		err = state.DB.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         postURL,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Printf("error saving post %q from feed %q: %v", postURL, feed.Url, err)
		}
	}

	return nil
}
