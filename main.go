package main
import (
	"fmt"
	"log"
	"os"

	"github.com/ananyabhardwaj10/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	
	s := state{
		cfg:&cfg, 
	}

	cmds := commands{
		CommandNames: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

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
		fmt.Printf("Error encountered: %v\n", err)
		os.Exit(1)
	}
}