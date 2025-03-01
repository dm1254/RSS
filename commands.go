package main

import (
	"errors"
	"fmt"
	"time"
	"os"
	"context"
	"github.com/dm1254/RSS/internal/config"
	"github.com/dm1254/RSS/internal/database"
	"github.com/google/uuid"
	"strconv"
)
type State struct{
	db *database.Queries
	Config *config.Config

}

type command struct{
	Name string 
	Args []string
}

type commands struct{
	Handlers map[string]func(*State,command)error
}

func (c *commands) register(name string, f func(*State,command)error){
	if c.Handlers == nil{
		c.Handlers = make(map[string]func(*State,command)error)	
	}
	c.Handlers[name] = f

}

func (c *commands) run(s *State, cmd command)error{
	handler, exists := c.Handlers[cmd.Name]; 
	if !exists{
		return errors.New("Command doesnt exists")
	}

	return handler(s,cmd)
}

func middlewareLoggedIn(handler func(s *State,cmd command, user database.User) error) func(*State,command) error{
	return func(s *State,cmd command) error {
		user,err := s.db.GetUser(context.Background(),s.Config.Username)
		if err != nil{
			return err
		}
		return handler(s,cmd,user)
	}
}

func handlerLogin(s *State, cmd command) error{
	if cmd.Args == nil{
		return errors.New("Login expects a username")
	}

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Println(os.Stderr,err)
		os.Exit(1)
	}
	s.Config.SetUser(cmd.Args[0])
	fmt.Printf("%s is logged in", cmd.Args[0])
	os.Exit(0)
	return nil
}
func registerUser(s *State, cmd command) error{
	if cmd.Args == nil{
		return errors.New("No user stated")
	}
	Params := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.Args[0],
	}
	_,err := s.db.GetUser(context.Background(),cmd.Args[0])
	if err == nil{
		fmt.Println(os.Stderr,err)
		os.Exit(1)
		
	}
	user,err := s.db.CreateUser(context.Background(),Params)
	if err != nil{
		fmt.Printf("Database err:= %v\n",err)
		return errors.New("Error creating user")
	}
	s.Config.SetUser(cmd.Args[0])
	fmt.Println("User:  %+v was created", user)
	os.Exit(0)
	return nil
}

func resetDb(s *State,cmd command) error{
	err := s.db.Reset(context.Background())
	if err != nil{
		fmt.Println(os.Stderr,err)
		os.Exit(1)		
	}
	os.Exit(0)
	return nil
}

func GetUsersInDb(s *State, cmd command) error{
	var users []string
	users, err := s.db.GetUsers(context.Background())
	if err != nil{
		fmt.Println(os.Stderr,err)
		os.Exit(1)	
	}
	for _,user := range users{
		if s.Config.Username == user{
			fmt.Printf("%s (current)\n",user)
		}else{
			fmt.Println(user)
		}
	}

	os.Exit(0)
	return nil
}

func Agg(s *State,cmd command) error{
	if len(cmd.Args) == 0{
		return errors.New("Duration required")
	}
	time_between_req := cmd.Args[0]
	timeBetweenReq,err := time.ParseDuration(time_between_req)
	if err != nil{
		return err 
	}
	fmt.Printf("Collecting feeds every %v\n",timeBetweenReq)
	
	ticker := time.NewTicker(timeBetweenReq)
	for ; ; <- ticker.C{
		scrapeFeeds(s)	
	}
	return nil

}

func addfeed(s *State,cmd command, user database.User) error{

	Params := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.Args[0],
		Url:cmd.Args[1],
		UserID: user.ID,
	}
	usersFeed, err := s.db.CreateFeed(context.Background(), Params)	
	if err != nil{
		fmt.Printf("Error:%v\n",err)
		return errors.New("Error inputing user feed")
	}
	FeedParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    usersFeed.ID,
	}


	_, err = s.db.CreateFeedFollow(context.Background(), FeedParams)
	if err != nil{
		fmt.Println(err)
		return err
	}
	fmt.Println(usersFeed)
	return nil
}

func ListFeeds(s *State, cmd command) error {
	getFeeds,err := s.db.GetFeeds(context.Background())
	if err != nil{
		fmt.Printf("ERROR: %s",err)
		return errors.New("Error getting feeds")
	}
	for _,feed := range getFeeds{
		fmt.Println(feed)	
	}
	return nil
}

func follow(s *State, cmd command, user database.User) error{
	getFeedId, err := s.db.GetFeed(context.Background(), cmd.Args[0])
	if err != nil{
		return err
	}

	Params := database.CreateFeedFollowParams{
	ID:        uuid.New(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	UserID:    user.ID,
	FeedID:    getFeedId.ID,
	}


	recordFollow, err := s.db.CreateFeedFollow(context.Background(), Params)
	if err != nil{
		fmt.Println(err)
		return err
	}
	fmt.Println(recordFollow)
	return nil	


}

func following(s *State, cmd command, user database.User) error{


	feedFollowsForUser, err := s.db.GetFeedFollowForUser(context.Background(), user.ID)
	if err != nil{
		fmt.Println(err)
		return err
	}

	for _,follow := range feedFollowsForUser{
		fmt.Println(follow.FeedName)
	
		
	}
	return nil
}

func unfollow(s *State, cmd command, user database.User) error {
	Params := database.UnfollowFeedParams{
		UserID: user.ID,
		Url: cmd.Args[0],
	}
	err := s.db.UnfollowFeed(context.Background(),Params)
	if err != nil{
		return err
	}
	return nil

}

func browse(s *State, cmd command, user database.User) error{
	if len(cmd.Args) == 0{
		cmd.Args = append(cmd.Args, "2") 
	}
	limit,err := strconv.Atoi(cmd.Args[0])
	if err != nil{
		return err	
	}
	Params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
		

	}

	retrieve_post,err := s.db.GetPostsForUser(context.Background(),Params)
	if err != nil{
		return err 
	}

	for i, post := range retrieve_post{
		fmt.Printf("Post #%d\n",i+1)
		if post.Title.Valid{
			fmt.Printf("Title: %s\n",post.Title.String)
		}else{
			fmt.Println("Title: [No Title]")
		}
		fmt.Printf("Url: %s\n",post.Url)

		if post.Description.Valid{
			fmt.Printf("Description: %s\n", post.Description.String)
		}else{
			fmt.Println("Description: N/A")
		}
	}

	return nil
	
}
