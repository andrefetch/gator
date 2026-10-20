package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/andrefetch/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(state *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage : addfeed <name> <url>")
	}

	currentUser, err := state.DB.GetUser(
		context.Background(),
		state.Config.CurrentUserName,
	)
	if err != nil {
		return err
	}

	now := time.Now()

	feed, err := state.DB.CreateFeed(
		context.Background(),
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

	fmt.Printf("Created feed: %+v\n", feed)
	return nil
}
