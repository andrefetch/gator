package main

import (
	"log"
	"os"

	"database/sql"

	"github.com/andrefetch/gator/internal/commands"
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

	state := &commands.State{
		Config: &cfg,
		DB:     dbQueries,
	}
	cmds := commands.New()

	if len(os.Args) < 2 {
		log.Fatal("not enough args")
	}

	cmd := commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.Run(state, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
