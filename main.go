package main

import _ "github.com/lib/pq"
import (
	"fmt"
	"log"
	"os"

	"database/sql"
	"github.com/ananyabhardwaj10/gator/internal/config"
	"github.com/ananyabhardwaj10/gator/internal/database"
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
	
	s := state{
		db:dbQueries,
		cfg:&cfg, 
	}

	cmds := commands{
		CommandNames: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerFollowing)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("not enough arguments. exiting the program")
		os.Exit(1)
	} 
	
	cmd := command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}