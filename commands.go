package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/uller91/goSorceryDraftDB/internal/apiInter"
	"github.com/uller91/goSorceryDraftDB/internal/database"
	"slices"
	"time"
)

type state struct {
	config   *config
	database *database.Queries
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

func addToCollection(origin *[]string, collection *[]string, item string) {
	if !slices.Contains(*origin, item) && !slices.Contains(*collection, item) {
		*collection = append(*collection, item)
	}
}

// Help
const (
	descriptionHelp = "Shows the list of commands (\"help\") or their description (\"help command_name\")"
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

// Update
const (
	descriptionUpdate = "Updates an internal card DB sending the API requiest to \"api.sorcerytcg.com\""
)

// current db size = 649
func handlerUpdate(s *state, cmd command) error {
	if len(cmd.arguments) != 0 {
		return errors.New("0 arguments are expected")
	}

	err := s.updateConfig()
	if err != nil {
		return err
	}

	var newSets []string

	fmt.Println("Initializing card DB update...")
	fmt.Println("")
	fmt.Printf("Sending the API requiest to %v...\n", s.config.BaseUrl)
	fmt.Println("")

	cards, err := apiInter.RequestCard(s.config.BaseUrl)
	if err != nil {
		return err
	}

	dbSize := len(cards)
	fmt.Printf("Cards found: %v\n", dbSize)
	fmt.Println("")

	fmt.Println("Updating card DB...")

	cardsAdded := 0
	cardsUpdated := 0

	for _, card := range cards {
		//add card
		paramCard := database.CreateCardParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: card.Name, Rarity: card.Guardian.Rarity, Type: card.Guardian.Type}
		cardCreated, err := s.database.CreateCard(context.Background(), paramCard)
		if err != nil {
			if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" {
				cardsUpdated += 1
			} else {
				return err
			}
		} else {
			cardsAdded += 1
		}

		//add set+card
		for _, set := range card.Sets {
			paramSets := database.CreateSetAndCardParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: set.Name, CardID: cardCreated.ID}
			_, err := s.database.CreateSetAndCard(context.Background(), paramSets)

			if err != nil {
				pqError, ok := err.(*pq.Error)
				if ok && pqError.Code == "23503" {
					continue
				} else {
					return err
				}
			}
			addToCollection(&s.config.Sets, &newSets, set.Name)
		}
	}

	//add new sets
	if newSets != nil {
		fmt.Println("")
		fmt.Println("New sets added to the DB:")

		for _, set := range newSets {
			paramSet := database.CreateSetParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: set}
			_, err := s.database.CreateSet(context.Background(), paramSet)
			if err != nil {
				return err
			}
			fmt.Println(set)
		}
	}

	fmt.Println("")
	fmt.Printf("Cards added in the DB: %v\n", cardsAdded)
	fmt.Printf("Cards already in the DB: %v\n", cardsUpdated)
	fmt.Println("")

	fmt.Println("SorceryDB update finished successfully")

	return nil
}

// Reset
const (
	descriptionReset = "Deletes all the entries from the DB"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.arguments) != 0 {
		return errors.New("0 arguments are expected")
	}

	err := s.database.SetlistReset(context.Background())
	if err != nil {
		return err
	}

	err = s.database.SetsReset(context.Background())
	if err != nil {
		return err
	}

	err = s.database.CardsReset(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Card DB successfully reset")

	return nil
}

// Open
const (
	descriptionGenerate = "Generates a random card pack from Sorcery TCG"
)

func handlerGenerate(s *state, cmd command) error {
	err := s.updateConfig()
	if err != nil {
		return err
	}

	if s.config.Sets == nil {
		return errors.New("The DB is empty! Use \"update\" command to fill the DB with cards")
	}

	//setting the set, tag -s
	set, err := setSet(s, cmd)
	if err != nil {
		return err
	}

	//setting number of cards in the pack, tag -p
	cardsInPack, err := setPack(cmd)
	if err != nil {
		return err
	}

	//"foils" in the pack, tag -f
	cardsInPack, _ = setFoil(cardsInPack, cmd)

	if set == "All" {
		return generateOnePackAll(s, cardsInPack)
	} else {
		return generateOnePack(s, set, cardsInPack)
	}

	return nil
}
