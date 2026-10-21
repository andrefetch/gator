package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/andrefetch/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(state *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage : follow <url>")
	}

	ctx := context.Background()

	currentUser, err := state.DB.GetUser(
		ctx,
		state.Config.CurrentUserName,
	)
	if err != nil {
		return err
	}

	feed, err := state.DB.GetFeedByURL(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	now := time.Now()

	follow, err := state.DB.CreateFeedFollow(
		ctx,
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			UserID:    currentUser.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s is now following %s\n", follow.UserName, follow.FeedName)
	return nil
}
