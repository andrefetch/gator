package commands

import (
	"fmt"

	"github.com/andrefetch/gator/internal/config"
	"github.com/andrefetch/gator/internal/database"
)

type State struct {
	DB     *database.Queries
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type handler func(*State, Command) error

type Commands struct {
	handlers map[string]handler
}

func New() *Commands {
	cmds := &Commands{handlers: make(map[string]handler)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	return cmds
}

func (c *Commands) Run(state *State, cmd Command) error {
	handler, exists := c.handlers[cmd.Name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}

	return handler(state, cmd)
}

func (c *Commands) register(name string, handler handler) {
	c.handlers[name] = handler
}
