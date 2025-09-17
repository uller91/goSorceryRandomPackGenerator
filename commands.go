package main

import (
	"errors"
	"fmt"
	"github.com/uller91/goSorceryDraftDB/internal/apiInter"
)

type state struct {
	config *config
	//database *database.Queries
	commands *commands
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	handlers     map[string]func(*state, command) error
	descriptions map[string]string
}

func (c *commands) run(s *state, cmd command) error {
	hndl, exists := c.handlers[cmd.name]
	if exists {
		err := hndl(s, cmd)
		if err != nil {
			return err
		}
	} else {
		return errors.New("No command with this name is registered")
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error, d string) {
	c.handlers[name] = f
	c.descriptions[name] = d
}

const (
	descriptionHelp = "Shows the list of commands (help) or their description (help command)"
)

func handlerHelp(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("List of available commands:")
		fmt.Println("")
		for command, _ := range s.commands.handlers {
			fmt.Println(command)
		}
		fmt.Println("")
		fmt.Println("To show the command description use \"help command\"")
	} else if len(cmd.arguments) == 1 {
		if discription, ok := s.commands.descriptions[cmd.arguments[0]]; ok {
			fmt.Printf("Command \"%s\"\n", cmd.arguments[0])
			fmt.Println(discription)
		} else {
			return errors.New("No command with this name is registered")
		}
	} else {
		return errors.New("Too many arguments")
	}

	return nil
}

// update description
const (
	descriptionUpdate = "Updates an internal DB sending the API requiest to api.sorcerytcg.com"
)

// current db size - 649
func handlerUpdate(s *state, cmd command) error {
	cards, err := apiInter.RequestCard(s.config.BaseUrl)
	if err != nil {
		return err
	}

	dbSize := len(cards)

	//fmt.Println(dbSize)
	//fmt.Println(cards[0])
	fmt.Println(cards[0].Name)
	fmt.Println(cards[0].Guardian.Rarity)
	fmt.Println(cards[0].Guardian.Type)
	fmt.Println(cards[0].Sets[0].Name)
	//fmt.Println(cards[dbSize-1])
	fmt.Println(cards[dbSize-1].Name)
	fmt.Println(cards[dbSize-1].Guardian.Rarity)
	fmt.Println(cards[dbSize-1].Guardian.Type)
	fmt.Println(cards[dbSize-1].Sets[0].Name)

	fmt.Println(cards[dbSize-50].Name)
	fmt.Println(cards[dbSize-50].Guardian.Rarity)
	fmt.Println(cards[dbSize-50].Guardian.Type)
	fmt.Println(cards[dbSize-50].Sets[0].Name)

	return nil
}
