package commands

import (
	"context"
	"fmt"
)

func handlerFollowing(state *State, _ Command) error {

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
