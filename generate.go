package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"
)

type Card struct {
	Name   string
	Rarity string
	Type   string
	Sets   string
}

// func getRandomCardsFromCollection[T any](s *state, collection []T, quantity int) []T
func getRandomCardsFromCollection(s *state, collection []Card, quantity int) []Card {
	if quantity >= len(collection) {
		fmt.Println("The collection is too small to get random cards. Returning full collection")
		return collection
	}

	if quantity < 0 {
		fmt.Println("The quantity of cards to return is less then 0. Returning 0 cards...")
		fmt.Println("")
		quantity = 0
	}

	randomCollection := []Card{}

	for i := 0; i < quantity; i++ {
		randomCardNumber, _ := rand.Int(rand.Reader, big.NewInt(int64(len(collection))))
		cardNumber := int(randomCardNumber.Int64())
		randomCard := collection[cardNumber]
		collection = slices.Delete(collection, cardNumber, cardNumber+1)

		//checking for cards that shouldn't be drawn
		if slices.Contains(s.config.Exceptions, randomCard.Name) {
			//fmt.Println(randomCard.Name)
			i -= 1
			continue
		}

		randomCollection = append(randomCollection, randomCard)
	}

	return randomCollection
}

func generateOnePack(s *state, setName string, cardsInPack map[string]int) error {
	set, err := s.database.GetCardsBySet(context.Background(), setName)
	if err != nil {
		return err
	}

	cardsOrdinary := []Card{}
	cardsExceptional := []Card{}
	cardsElite := []Card{}
	cardsUnique := []Card{}

	for _, setCard := range set {
		card, err := s.database.GetCard(context.Background(), setCard.CardID)
		if err != nil {
			return err
		}

		rarity := card.Rarity

		cardClean := Card{
			Name:   card.Name,
			Rarity: rarity,
			Type:   card.Type,
		}

		switch rarity {
		case "Ordinary":
			cardsOrdinary = append(cardsOrdinary, cardClean)
		case "Exceptional":
			cardsExceptional = append(cardsExceptional, cardClean)
		case "Elite":
			cardsElite = append(cardsElite, cardClean)
		case "Unique":
			if setName == "Arthurian Legends" && slices.Contains(s.config.ALSirs, card.Name) {
				cardsElite = append(cardsElite, cardClean)
				continue
			}
			cardsUnique = append(cardsUnique, cardClean)
		default:
			return errors.New("Unknown rarity was found!")
		}
	}

	if slices.Contains(s.config.MiniSets, setName) {
		fmt.Printf("Can't generate pack for %s mini-set! Use \"generate -s all\" to generate pack to include all the cards\n", setName)
	} else {
		pack := getRandomCardsFromCollection(s, cardsOrdinary, cardsInPack["Ordinary"])
		pack = append(pack, getRandomCardsFromCollection(s, cardsExceptional, cardsInPack["Exceptional"])...)
		pack = append(pack, getRandomCardsFromCollection(s, cardsElite, cardsInPack["Elite"])...)
		pack = append(pack, getRandomCardsFromCollection(s, cardsUnique, cardsInPack["Unique"])...)

		fmt.Printf("Random pack from %s set:\n", setName)
		fmt.Println("")

		for _, card := range pack {
			fmt.Printf("%-25v | %-10v | %-15v\n", card.Name, card.Type, card.Rarity)
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
		sets, err := s.database.GetSetsByCard(context.Background(), card.ID)
		if err != nil {
			return err
		}

		var setNames []string
		for _, set := range sets {
			setNames = append(setNames, set.Name)
		}

		rarity := card.Rarity
		setName := strings.Join(setNames, " / ")

		cardClean := Card{
			Name:   card.Name,
			Rarity: rarity,
			Type:   card.Type,
			Sets:   setName,
		}

		switch rarity {
		case "Ordinary":
			cardsOrdinary = append(cardsOrdinary, cardClean)
		case "Exceptional":
			cardsExceptional = append(cardsExceptional, cardClean)
		case "Elite":
			cardsElite = append(cardsElite, cardClean)
		case "Unique":
			if setName == "Arthurian Legends" && slices.Contains(s.config.ALSirs, card.Name) {
				cardsElite = append(cardsElite, cardClean)
				continue
			}
			cardsUnique = append(cardsUnique, cardClean)
		default:
			return errors.New("Unknown rarity was found!")
		}
	}

	pack := getRandomCardsFromCollection(s, cardsOrdinary, cardsInPack["Ordinary"])
	pack = append(pack, getRandomCardsFromCollection(s, cardsExceptional, cardsInPack["Exceptional"])...)
	pack = append(pack, getRandomCardsFromCollection(s, cardsElite, cardsInPack["Elite"])...)
	pack = append(pack, getRandomCardsFromCollection(s, cardsUnique, cardsInPack["Unique"])...)

	fmt.Println("Random pack from all sets:")
	fmt.Println("")

	for _, card := range pack {
		fmt.Printf("%-25v | %-10v | %-15v | %-10v\n", card.Name, card.Type, card.Rarity, card.Sets)
	}

	return nil
}
