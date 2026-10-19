package main

import (
	"context"
	"fmt"
	"time"

	"github.com/andrefetch/gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {

	if len(cmd.args) == 0 {
		return fmt.Errorf("username is required to register")
	}

	now := time.Now()

	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Name:      cmd.args[0],
		},
	)

	if err != nil {
		return err
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("user %s was created", cmd.args[0])
	fmt.Printf("created user: %v\n", user)

	return nil

}
