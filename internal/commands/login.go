package commands

import (
	"context"
	"fmt"
)

func handlerLogin(state *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is required to login.")
	}

	_, err := state.DB.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	if err := state.Config.SetUser(cmd.Args[0]); err != nil {
		return err
	}

	fmt.Printf("%s logged in\n", cmd.Args[0])
	return nil
}
