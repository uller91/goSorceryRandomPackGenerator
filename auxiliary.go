package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/uller91/goSorceryDraftDB/internal/database"
	"math/big"
	"slices"
)

type Card struct {
	Name   string
	Rarity string
	Type   string
	Set    string
}

func addToCollection(origin *[]string, collection *[]string, item string) {
	if !slices.Contains(*origin, item) && !slices.Contains(*collection, item) {
		*collection = append(*collection, item)
	}
}

func (s *state) updateConfig() error {
	//sets
	oldSets, err := s.database.GetSets(context.Background())
	if err != nil {
		return err
	}

	for _, st := range oldSets {
		s.config.Sets = append(s.config.Sets, st.Name)
	}

	/*
		//types
		oldTypes, err := s.database.GetTypes(context.Background())
		if err != nil {
			return err
		}

		for _, tp := range oldTypes {
			s.config.Types = append(s.config.Types, tp.Name)
		}
		fmt.Println(s.config.Types)

		//rarities
		oldRarities, err := s.database.GetRarities(context.Background())
		if err != nil {
			return err
		}

		for _, rt := range oldRarities {
			s.config.Rarities = append(s.config.Rarities, rt.Name)
		}
		fmt.Println(s.config.Rarities)
	*/

	s.config.Types = []string{"Avatar", "Minion", "Magic", "Aura", "Artifact", "Site"}
	s.config.Rarities = []string{"Ordinary", "Exceptional", "Elite", "Unique"}
	s.config.ALSirs = []string{"Dame Britomart", "Sir Agravaine", "Sir Balin", "Sir Bedivere", "Sir Bors the Younger", "Sir Gaheris", "Sir Gawain", "Sir Ironside", "Sir Kay", "Sir Lamorak", "Sir Morien", "Sir Pelleas", "Sir Perceval", "Sir Priamus", "Sir Tom Thumb", "Sir Tristan"}
	s.config.MiniSets = []string{"Dragonlord"}

	return nil
}

//uses generics now... T can be Card{} or database.Card{} at the moment
func getRandomCardsFromCollection[T any](collection []T, quantity int) []T {
	if quantity >= len(collection) {
		fmt.Println("The collection is too small to get random cards. Returning full collection")
		return collection
	}

	randomCollection := []T{}

	for i := 0; i < quantity; i++ {
		randomCardNumber, _ := rand.Int(rand.Reader, big.NewInt(int64(len(collection))))
		cardNumber := int(randomCardNumber.Int64())
		randomCard := collection[cardNumber]
		randomCollection = append(randomCollection, randomCard)
		collection = slices.Delete(collection, cardNumber, cardNumber+1)
	}

	return randomCollection
}

func generateOnePack(s *state, setName string, cardsInPack map[string]int) error {
	set, err := s.database.GetCardsBySet(context.Background(), setName)
	if err != nil {
		return err
	}

	cardsOrdinary := []database.Card{}
	cardsExceptional := []database.Card{}
	cardsElite := []database.Card{}
	cardsUnique := []database.Card{}

	for _, setCard := range set {
		card, err := s.database.GetCard(context.Background(), setCard.CardID)
		if err != nil {
			return err
		}

		rarity := card.Rarity

		switch rarity {
		case "Ordinary":
			cardsOrdinary = append(cardsOrdinary, card)
		case "Exceptional":
			cardsExceptional = append(cardsExceptional, card)
		case "Elite":
			cardsElite = append(cardsElite, card)
		case "Unique":
			if setName == "Arthurian Legends" && slices.Contains(s.config.ALSirs, card.Name) {
				cardsElite = append(cardsElite, card)
				continue
			}
			cardsUnique = append(cardsUnique, card)
		default:
			return errors.New("Unknown rarity was found!")
		}
	}

	if slices.Contains(s.config.MiniSets, setName) {
		fmt.Printf("There is no pack for %s set!\n", setName)
	} else {
		pack := getRandomCardsFromCollection(cardsOrdinary, cardsInPack["Ordinary"])
		pack = append(pack, getRandomCardsFromCollection(cardsExceptional, cardsInPack["Exceptional"])...)
		pack = append(pack, getRandomCardsFromCollection(cardsElite, cardsInPack["Elite"])...)
		pack = append(pack, getRandomCardsFromCollection(cardsUnique, cardsInPack["Unique"])...)

		fmt.Printf("Random pack from %s set:\n", setName)
		fmt.Println("")

		for _, card := range pack {
			fmt.Printf("%v - %v - %v\n", card.Name, card.Type, card.Rarity)
		}
	}

	return nil
}

func generateOnePackAll(s *state, cardsInPack map[string]int) error {
	cards, err := s.database.GetAllCards(context.Background())
	if err != nil {
		return err
	}

	cardsOrdinary := []Card{}
	cardsExceptional := []Card{}
	cardsElite := []Card{}
	cardsUnique := []Card{}

	for _, card := range cards {
		set, err := s.database.GetSetByCard(context.Background(), card.ID)
		if err != nil {
			return err
		}

		rarity := card.Rarity
		setName := set.Name
		if setName == "Alpha" || setName == "Beta" {
			setName = "Alpha / Beta"
		}

		cardClean := Card{
			Name: card.Name,
			Rarity: rarity,
			Type: card.Type,
			Set: setName,
		}


		switch rarity {
		case "Ordinary":
			cardsOrdinary = append(cardsOrdinary, cardClean)
		case "Exceptional":
			cardsExceptional = append(cardsExceptional, cardClean)
		case "Elite":
			cardsElite = append(cardsElite, cardClean)
		case "Unique":
			cardsUnique = append(cardsUnique, cardClean)
		default:
			return errors.New("Unknown rarity was found!")
		}
	}

	pack := getRandomCardsFromCollection(cardsOrdinary, cardsInPack["Ordinary"])
	pack = append(pack, getRandomCardsFromCollection(cardsExceptional, cardsInPack["Exceptional"])...)
	pack = append(pack, getRandomCardsFromCollection(cardsElite, cardsInPack["Elite"])...)
	pack = append(pack, getRandomCardsFromCollection(cardsUnique, cardsInPack["Unique"])...)

	fmt.Println("Random pack from all sets:")
	fmt.Println("")

	for _, card := range pack {
		fmt.Printf("%v - %v - %v - %v\n", card.Name, card.Type, card.Rarity, card.Set)
	}

	return nil
}
