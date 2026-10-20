package commands

import (
	"context"
	"fmt"
)

func handlerFeeds(state *State, _ Command) error {
	feeds, err := state.DB.GetFeeds(
		context.Background(),
	)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("* %s (%s) - %s\n", feed.Name, feed.Url, feed.UserName)
	}

	return nil
}
