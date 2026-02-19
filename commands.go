package main
import (
	"errors"
	"fmt"
	"time"
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/gator/internal/config"
	"github.com/ananyabhardwaj10/gator/internal/database"
)

type state struct {
	db *database.Queries
	cfg *config.Config  `json:"cfg"`
}

type command struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type commands struct {
	CommandNames map[string]func(*state, command) error `json:"command_names"`
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("please pass a username.")
	}
	
	_, err := s.db.GetUser(context.Background(), cmd.Args[0]) 
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user does not exist")
		}
		return err 
	}

	err = s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("user has been set.")
	return nil 
}

func (c *commands) run(s *state, cmd command) error {
	f, exists := c.CommandNames[cmd.Name]
	if !exists {
		return errors.New("unknown command.") 
	}

	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.CommandNames[name] = f
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("please pass a name.")
	}

	id := uuid.New()
	now := time.Now()
	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: id,
		CreatedAt: now,
		UpdatedAt: now,
		Name: name,
	})
	_ = user 
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" {
			return fmt.Errorf("user already exists")
		}
		return err 
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err 
	}
	return nil 
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUser(context.Background())
	if err != nil {
		return fmt.Errorf("Error encountered while resetting the database") 
	}

	return nil 
}

func handlerGetUsers(s *state, cmd command) error {
	usersList, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error encountered: %v", err)
	}

	for _, user := range usersList {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil 
}

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	feed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err 
	}

	fmt.Println(feed)
	return nil 
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("please pass a valid feed name or url")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	userName := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return err 
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(), 
		CreatedAt: time.Now(), 
		UpdatedAt: time.Now(), 
		Name: name, 
		Url: url, 
		UserID: user.ID,
	})
	if err != nil {
		return err 
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err 
	}

	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("Created: %v\n", feed.CreatedAt)
	fmt.Printf("Updated: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)
	fmt.Printf("User: %s\n", userName)

	return nil 
}

func handlerFeeds(s *state, cmd command) error {
	feedList, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err 
	}

	for _, feed := range feedList {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err 
		}

		fmt.Printf("Feed Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("UserName: %s\n", user.Name)
	}
	return nil 
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("please pass a valid url")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err 
	}

	feed, err := s.db.GetFeedFromURL(context.Background(), cmd.Args[0])
	if err != nil {
		return err 
	}

	record, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err 
	}

	fmt.Printf("Feed Name: %s\n", record.FeedName)
	fmt.Printf("Current User: %s\n", record.UserName)

	return nil 
}

func handlerFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err 
	}
	feed_follow_list, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err 
	}

	for _, feed := range feed_follow_list {
		fmt.Println(feed.FeedName)
	}
	return nil 
}