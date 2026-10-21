package commands

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
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

// Missing or unrecognized dates remain NULL so the post can still be saved.
func parsePublishedAt(value string) (sql.NullTime, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return sql.NullTime{}, nil
	}

	// Go otherwise treats unknown zone abbreviations as UTC. Resolve common
	// RSS abbreviations explicitly so their publication times sort correctly.
	zone := time.UTC
	fields := strings.Fields(value)
	zoneName := fields[len(fields)-1]
	if offset, ok := map[string]int{
		"EST": -5, "EDT": -4,
		"CST": -6, "CDT": -5,
		"MST": -7, "MDT": -6,
		"PST": -8, "PDT": -7,
	}[zoneName]; ok {
		zone = time.FixedZone(zoneName, offset*60*60)
	}

	for _, layout := range []string{
		time.RFC1123Z,
		time.RFC1123,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
		"2 Jan 2006 15:04:05 -0700",
		"2 Jan 2006 15:04:05 MST",
		time.RFC3339Nano,
		time.RFC822Z,
		time.RFC822,
		time.RFC850,
		time.ANSIC,
	} {
		if parsed, err := time.ParseInLocation(layout, value, zone); err == nil {
			return sql.NullTime{Time: parsed.UTC(), Valid: true}, nil
		}
	}

	return sql.NullTime{}, fmt.Errorf("unrecognized publication date %q", value)
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		feedURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "gator")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("returned status code is unexpected: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := 0; i < len(feed.Channel.Item); i++ {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}
