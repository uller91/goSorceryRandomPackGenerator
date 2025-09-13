package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/uller91/goSorceryDraftDB/internal/apiInter"
	"os"
)

type config struct {
	BaseUrl string
}

func main() {
	godotenv.Load()
	testVar := os.Getenv("TEST_VAR")
	fmt.Println(testVar)

	var cfg config
	var st state

	cfg.BaseUrl = apiInter.BaseUrl
	st.config = &cfg

	//db, err = ...
	//st.database = ...

	/*
		cards := apiInter.RequestCard(apiUrl)
		dbSize := len(cards)
		fmt.Println(cards[0])
		fmt.Println(cards[dbSize-1])
	*/

	var cmds commands
	handlers := make(map[string]func(*state, command) error)
	descriptions := make(map[string]string)
	cmds.handlers = handlers
	cmds.descriptions = descriptions
	st.commands = &cmds

	cmds.register("help", handlerHelp, "helps") //add description as a constant

	cmd := command{name: "help"}
	err := cmds.run(&st, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	/*
		args := os.Args[:]
		if len(args) < 2 {
			err := fmt.Errorf("Not enough arguments!")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		commandName := args[1]
		commandArgs := []string{}
		if len(args) > 2 {
			commandArgs = args[2:]
		}
		cmd := command{name: commandName, arguments: commandArgs}

		err = cmds.run(&st, cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	*/
}
