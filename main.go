package main 

import (
	"context"
	"errors"
	"database/sql"
	"os"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/josequiceno2000/gator/internal/config"
	"github.com/josequiceno2000/gator/internal/database"
	_ "github.com/lib/pq"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.DB.GetUser(context.Background(), s.CfgPointer.CurrentUsername)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("user '%s' does not exist", s.CfgPointer.CurrentUsername)
			}
			return fmt.Errorf("failed to get user: %w", err)
		}

		return handler(s, cmd, user)
	}
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return errors.New("unfollow: url argument is required")
	}

	url := cmd.Arguments[0]

	err := s.DB.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url: url,
	})

	if err != nil {
		return fmt.Errorf("unfollow: failed to unfollow feed: %w", err)
	}

	fmt.Printf("Unfollowed feed with URL: %s\n", url)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	user, err := s.DB.GetUser(context.Background(), s.CfgPointer.CurrentUsername)
	if err != nil {
		return fmt.Errorf("following: failed to get user: %w", err)
	}

	feedFollows, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("following, failed to get feed follows: %w", err)
	}

	for _, ff := range feedFollows {
		fmt.Println(ff.FeedName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return errors.New("follow: url argument is required")
	}

	url := cmd.Arguments[0]

	user, err := s.DB.GetUser(context.Background(), s.CfgPointer.CurrentUsername)
	if err != nil {
		return fmt.Errorf("follow: failed to get user: %w", err)
	}

	feed, err := s.DB.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("follow: failed to get feed: %w", err)
	}

	feedFollow, err := s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("follow: failed to create feed follow: %w", err)
	}

	fmt.Printf("Followed feed: %s by user: %s\n", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func handlerFeeds(s * state, cmd command) error {
	feeds, err := s.DB.GetFeedsWithUserNames(context.Background())
	if err != nil {
		return fmt.Errorf("feeds: failed to get feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %s, URL: %s, User: %s\n", feed.Name, feed.Url, feed.UserName)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return errors.New("addfeed: name and url arguments are required")
	}

	name := cmd.Arguments[0]
	url := cmd.Arguments[1]

	user, err := s.DB.GetUser(context.Background(), s.CfgPointer.CurrentUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("addfeed: user '%s' does not exist", s.CfgPointer.CurrentUsername)
		}
		return fmt.Errorf("addfeed: failed to get user: %w", err)
	}

	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
        ID:        uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Name:      name,
        Url:       url,
        UserID:    user.ID,
    })

    if err != nil {
        return fmt.Errorf("addfeed: failed to create feed: %w", err)
    }

	_, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("addfeed: failed to create feed follow: %w", err)
	}

	fmt.Printf("Feed created and dollowed: %+v\n", feed)
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("agg: failed to fetch feed %w", err)
	}
	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("users: failed to get users: %w", err)
	}

	for _, user := range users {
		if user == s.CfgPointer.CurrentUsername {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("register command requires a username argument")
	}

	username := cmd.Arguments[0]
	userID := uuid.New()
	now := time.Now().UTC()

	user, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID: userID,
		CreatedAt: now,
		UpdatedAt: now,
		Name: username,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user with name '%s' already exists", username)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	err = s.CfgPointer.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set current user in config: %w", err)
	}

	fmt.Printf("User '%s' registered successfully.\n", username)
	log.Printf("Registered user: %+v\n", user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.DB.DeleteAllUsers((context.Background()))
	if err != nil {
		return fmt.Errorf("reset: failed to delete all users: %w", err)
	}
	fmt.Println("reset: all users deleted successfully")
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Open db connection
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Error opening db connectio: %v", err)
	}
	defer db.Close()

	// Create database queries instance
	dbQueries := database.New(db)

	appState := state{DB:dbQueries, CfgPointer: &cfg}
	
	cmdRegistry := commands{}
	cmdRegistry.register("login", handlerLogin)
	cmdRegistry.register("register", handlerRegister)
	cmdRegistry.register("reset", handlerReset)
	cmdRegistry.register("users", handlerUsers)
	cmdRegistry.register("agg", handlerAgg)
	cmdRegistry.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmdRegistry.register("feeds", handlerFeeds)
	cmdRegistry.register("follow", middlewareLoggedIn(handlerFollow))
	cmdRegistry.register("following", middlewareLoggedIn(handlerFollowing))
	cmdRegistry.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	if len(os.Args) < 2 {
		fmt.Println("Error: not enough arguments provided")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{Name: cmdName, Arguments: cmdArgs}

	err = cmdRegistry.run(&appState, cmd)
	if err != nil {
		log.Fatalf("Command error: %v", err)
	}
}