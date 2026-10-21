package commands

import (
	"context"
	"fmt"

	"github.com/andrefetch/gator/internal/database"
)

func handlerFollowing(state *State, cmd Command, user database.User) error {

	ctx := context.Background()

	currentUser, err := state.DB.GetUser(
		ctx,
		state.Config.CurrentUserName,
	)
	if err != nil {
		return err
	}

	follows, err := state.DB.GetFeedFollowsForUser(
		ctx,
		currentUser.ID,
	)
	if err != nil {
		return err
	}

	for _, follow := range follows {
		fmt.Printf("* %s\n", follow.FeedName)
	}

	return nil
}
