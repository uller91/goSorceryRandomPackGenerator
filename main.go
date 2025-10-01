package main

import _ "github.com/lib/pq"

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/uller91/goSorceryDraftDB/internal/apiInter"
	"github.com/uller91/goSorceryDraftDB/internal/database"
	"os"
	"strings"
)

type config struct {
	BaseUrl  string
	Sets     []string
	Types    []string
	Rarities []string
	ALSirs   []string
	MiniSets []string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("CONNECTION_STRING")

	var cfg config
	var st state

	cfg.BaseUrl = apiInter.BaseUrl
	st.config = &cfg

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	st.database = database.New(db) //database.Queries struct

	var cmds commands
	handlers := make(map[string]func(*state, command) error)
	descriptions := make(map[string]string)
	cmds.handlers = handlers
	cmds.descriptions = descriptions
	st.commands = &cmds

	cmds.register("help", handlerHelp, descriptionHelp)
	cmds.register("update", handlerUpdate, descriptionUpdate)
	cmds.register("reset", handlerReset, descriptionReset)
	cmds.register("generate", handlerGenerate, descriptionGenerate)

	//single command test
	/*
		cmd := command{name: "help"}
		err := cmds.run(&st, cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	*/

	args := os.Args[:]
	if len(args) < 2 {
		//err := fmt.Errorf("Not enough arguments")
		//fmt.Println(err.Error())
		fmt.Println("Welcome to the Sorcery TCG Random Pack Generator. To see the list of available commands use the \"help\" command.")
		os.Exit(0)
	}

	commandName := strings.ToLower(args[1])
	commandArgs := []string{}
	if len(args) > 2 {
		commandArgs = args[2:]
		for i := range commandArgs {
			commandArgs[i] = strings.ToLower(commandArgs[i])
		}
	}
	cmd := command{name: commandName, arguments: commandArgs}

	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
