package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/andrefetch/gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(state *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is required to register")
	}

	now := time.Now()
	user, err := state.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.Args[0],
	})
	if err != nil {
		return err
	}

	if err := state.Config.SetUser(cmd.Args[0]); err != nil {
		return err
	}

	fmt.Printf("user %s was created", cmd.Args[0])
	fmt.Printf("created user: %v\n", user)
	return nil
}
