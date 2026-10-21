package commands

import (
	"context"
	"fmt"

	"github.com/andrefetch/gator/internal/database"
)

func handlerUnfollow(state *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage : unfollow <url>")
	}

	ctx := context.Background()

	feed, err := state.DB.GetFeedByURL(
		ctx,
		cmd.Args[0],
	)
	if err != nil {
		return err
	}

	err = state.DB.DeleteFeedFollow(
		ctx,
		database.DeleteFeedFollowParams{
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s unfollowed %s\n", user.Name, feed.Name)
	return nil
}
