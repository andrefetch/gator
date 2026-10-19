package main

import (
	"github.com/andrefetch/gator/internal/config"
	"github.com/andrefetch/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
