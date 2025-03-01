package main

import (
	"fmt"
	"database/sql"
	"os"
	"github.com/dm1254/RSS/internal/config"
	"github.com/dm1254/RSS/internal/database"
	_ "github.com/lib/pq"
)
func main(){	
	new_config, err := config.Read()
	if err != nil{
		fmt.Printf("%s\n",err)	
	}
	db,err:= sql.Open("postgres",new_config.DBURL)
	if err != nil{
		fmt.Printf("%s",err)
		os.Exit(1)

	}

	dbQueries := database.New(db)
	s := &State{
		db: dbQueries,
		Config: &new_config,		
	}
	c := &commands{
		Handlers: make(map[string]func(*State,command)error),	
	}
	c.register("login",handlerLogin)
	c.register("register", registerUser)
	c.register("reset", resetDb)
	c.register("users", GetUsersInDb)
	c.register("agg", Agg)
	c.register("addfeed", middlewareLoggedIn(addfeed))
	c.register("feeds", ListFeeds)
	c.register("follow", middlewareLoggedIn(follow))
	c.register("following",middlewareLoggedIn(following))
	c.register("unfollow",middlewareLoggedIn(unfollow))
	c.register("browse", middlewareLoggedIn(browse))
	if len(os.Args) < 2{
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}


	commandName := os.Args[1]
	ArgumentName := os.Args[2:]
	cmd := command{
		Name: commandName,
		Args: ArgumentName,
	}

	if err := c.run(s, cmd); err != nil{
		fmt.Println(err)
	}

}


