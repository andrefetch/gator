package main

import (
	"context"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {

	if len(cmd.args) == 0 {
		return fmt.Errorf("username is required to login.")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("%s logged in\n", cmd.args[0])

	return nil
}
