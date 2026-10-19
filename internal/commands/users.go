package commands

import (
	"context"
	"fmt"
)

func handlerUsers(state *State, cmd Command) error {
	users, err := state.DB.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == state.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}
