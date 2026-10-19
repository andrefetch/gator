package commands

import (
	"context"
	"fmt"
)

func handlerReset(state *State, _ Command) error {
	if err := state.DB.Reset(context.Background()); err != nil {
		return err
	}

	fmt.Println("database reset successfully")
	return nil
}
