package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/bootdotdev/gator/internal/config"
	"github.com/bootdotdev/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("We've got a problem.")
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := &state{db: dbQueries, cfg: &cfg}
	c := commands{registeredCommands: make(map[string]func(*state, command) error)}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerListUsers)
	c.register("agg", handlerAgg)
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("feeds", handlerListFeeds)
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	c.register("following", middlewareLoggedIn(handlerListFeedFollows))

	inputs := os.Args
	nInputs := len(inputs)
	if nInputs < 2 {
		log.Fatal("Not enough args.")
	}
	cmd := command{Name: inputs[1], Args: inputs[2:]}
	err = c.run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
