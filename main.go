package main

import (
	"log"
	"os"

	"database/sql"

	"github.com/andrefetch/gator/internal/config"
	"github.com/andrefetch/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	s := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)

	if len(os.Args) < 2 {
		log.Fatal("not enough args")
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
