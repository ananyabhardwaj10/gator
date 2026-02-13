package main
import (
	"fmt"
	"log"

	"github.com/ananyabhardwaj10/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg.DBURL)
	err = cfg.SetUser("ananya")
	if err != nil {
		log.Fatal(err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)
}