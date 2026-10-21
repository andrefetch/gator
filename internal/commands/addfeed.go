package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/andrefetch/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(state *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	ctx := context.Background()

	currentUser, err := state.DB.GetUser(
		ctx,
		state.Config.CurrentUserName,
	)
	if err != nil {
		return err
	}

	now := time.Now()

	feed, err := state.DB.CreateFeed(
		ctx,
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Name:      cmd.Args[0],
			Url:       cmd.Args[1],
			UserID:    currentUser.ID,
		},
	)
	if err != nil {
		return err
	}

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

	fmt.Printf(
		"feed %s added and followed by %s\n",
		follow.FeedName,
		follow.UserName,
	)

	return nil
}
