package commands

import (
	"context"
	"fmt"
)

const RSSLink = "https://www.wagslane.dev/index.xml"

func handlerAgg(_ *State, _ Command) error {
	feed, err := fetchFeed(
		context.Background(),
		RSSLink,
	)
	if err != nil {
		return nil
	}

	fmt.Printf("%v\n", feed)
	return nil
}
