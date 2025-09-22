package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/uller91/goSorceryDraftDB/internal/apiInter"
	"github.com/uller91/goSorceryDraftDB/internal/database"
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

// Help
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

// Update
const (
	descriptionUpdate = "Updates an internal DB sending the API requiest to api.sorcerytcg.com"
)

// current db size - 649
func handlerUpdate(s *state, cmd command) error {
	if len(cmd.arguments) != 0 {
		return errors.New("0 arguments are expected")
	}

	fmt.Println("Initializing the SorceryDB update...")
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

	fmt.Println("Updating the DB...")
	fmt.Println("")

	cardsAdded := 0
	cardsUpdated := 0

	for _, card := range cards {
		//add card
		paramCard := database.CreateCardParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: card.Name, Rarity: card.Guardian.Rarity, Type: card.Guardian.Type}
		cardCreated, err := s.database.CreateCard(context.Background(), paramCard)
		if err != nil {
			if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" {
				//fmt.Printf("%v card already exist in DB\n", card.Name)
				cardsUpdated += 1
			} else {
				return err
			}
		} else {
			fmt.Printf("\"%v\" added in the Cards table\n", cardCreated.Name)
			//fmt.Println(cardCreated.Rarity)
			//fmt.Println(cardCreated.Type)
			cardsAdded += 1
		}

		//add set+card
		for _, set := range card.Sets {
			paramSet := database.CreateSetAndCardParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: set.Name, CardID: cardCreated.ID}
			setCardCreated, err := s.database.CreateSetAndCard(context.Background(), paramSet)

			if err != nil {
				pqError, ok := err.(*pq.Error)
				if ok && pqError.Code == "23503" {
					continue
				} else {
					return err
				}
			} else {
				fmt.Printf("Combination of Set: \"%v\" and Card: \"%v\" added in the Sets table\n", setCardCreated.Name, cardCreated.Name)
			}
		}
	}

	fmt.Println("")
	fmt.Printf("Cards added in the DB: %v\n", cardsAdded)
	fmt.Printf("Cards already in the DB: %v\n", cardsUpdated)
	fmt.Println("")

	//Apprentice Wizard
	/*
		fmt.Println("Card found:")
		fmt.Println(cards[0].Name)
		fmt.Println(cards[0].Guardian.Rarity)
		fmt.Println(cards[0].Guardian.Type)
		fmt.Println(cards[0].Sets[0].Name)
	*/

	fmt.Println("SorceryDB update successfully finished")

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

	err := s.database.SetsReset(context.Background())
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			return pqError
		} else {
			return err
		}
	}

	err = s.database.CardsReset(context.Background())
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			return pqError
		} else {
			return err
		}
	}

	return nil
}
