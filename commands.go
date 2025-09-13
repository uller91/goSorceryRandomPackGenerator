package main

import (
	"errors"
	"fmt"
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

// not finished - error handling and descriptions with arguments
func handlerHelp(s *state, cmd command) error {
	fmt.Println("List of available commands:")
	for command, _ := range s.commands.handlers {
		fmt.Println(command)
		fmt.Println(s.commands.descriptions[command])
	}

	fmt.Println("")
	return nil
}
