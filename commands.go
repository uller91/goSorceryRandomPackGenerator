package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/uller91/goSorceryDraftDB/internal/apiInter"
	"github.com/uller91/goSorceryDraftDB/internal/database"
	"math/big"
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
	var newTypes []string
	var newRarities []string

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
				//fmt.Printf("%v card already exist in DB\n", card.Name)
				cardsUpdated += 1
			} else {
				return err
			}
		} else {
			//fmt.Printf("\"%v\" added in the DB\n", cardCreated.Name)
			cardsAdded += 1

			addToCollection(&s.config.Types, &newTypes, cardCreated.Type)
			addToCollection(&s.config.Rarities, &newRarities, cardCreated.Rarity)
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
			} /* else {
				fmt.Printf("Combination of Set: \"%v\" and Card: \"%v\" added in the Sets table\n", setCardCreated.Name, cardCreated.Name)
			} */

			addToCollection(&s.config.Sets, &newSets, set.Name)
		}
	}

	//fmt.Println(newTypes)
	//fmt.Println(newRarities)
	//fmt.Println(newSets)

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

	//add new types
	if newTypes != nil {
		for _, tpe := range newTypes {
			paramType := database.CreateTypeParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: tpe}
			_, err := s.database.CreateType(context.Background(), paramType)
			if err != nil {
				return err
			}
		}
	}

	//add new rarities
	if newRarities != nil {
		for _, rarity := range newRarities {
			paramType := database.CreateRarityParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: rarity}
			_, err := s.database.CreateRarity(context.Background(), paramType)
			if err != nil {
				return err
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

	err := s.database.RaritylistReset(context.Background())
	if err != nil {
		return err
	}

	err = s.database.TypelistReset(context.Background())
	if err != nil {
		return err
	}

	err = s.database.SetlistReset(context.Background())
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
	if tag := slices.Index(cmd.arguments, "-f"); tag != -1 {
		//25% of "foil"
		foilProbability, _ := rand.Int(rand.Reader, big.NewInt(int64(4)))
		//"foil" distribution: 7/17 ordinary, 6/17 exceptional, 3/17 elite, 1/17 unique
		if foilProbability.Int64() == 0 {
			whichFoilProbability, _ := rand.Int(rand.Reader, big.NewInt(int64(17)))
			if whichFoilProbability.Int64() == 0 {
				cardsInPack["Ordinary"] -= 1
				cardsInPack["Unique"] += 1
			} else if whichFoilProbability.Int64() < 4 {
				cardsInPack["Ordinary"] -= 1
				cardsInPack["Elite"] += 1
			} else if whichFoilProbability.Int64() < 10 {
				cardsInPack["Ordinary"] -= 1
				cardsInPack["Exceptional"] += 1
			}
		}
	}

	if set == "All" {
		return generateOnePackAll(s, cardsInPack)
	} else {
		return generateOnePack(s, set, cardsInPack)
	}

	return nil
}
