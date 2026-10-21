package commands

import (
	"context"

	"github.com/andrefetch/gator/internal/database"
)

func middlewareLoggedIn(
	next func(*State, Command, database.User) error,
) func(*State, Command) error {
	return func(state *State, cmd Command) error {
		user, err := state.DB.GetUser(
			context.Background(),
			state.Config.CurrentUserName,
		)
		if err != nil {
			return err
		}

		return next(state, cmd, user)
	}
}
