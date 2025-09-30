package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/uller91/goSorceryDraftDB/internal/database"
	"math/big"
	"slices"
)

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

func getRandomCardsFromCollection(collection []database.Card, quantity int) []database.Card {
	if quantity >= len(collection) {
		fmt.Println("The collection is too small to get random cards. Returning full collection")
		return collection
	}

	randomCollection := []database.Card{}

	for i := 0; i < quantity; i++ {
		randomCardNumber, _ := rand.Int(rand.Reader, big.NewInt(int64(len(collection))))
		cardNumber := int(randomCardNumber.Int64())
		randomCard := collection[cardNumber]
		randomCollection = append(randomCollection, randomCard)
		collection = slices.Delete(collection, cardNumber, cardNumber+1)
	}

	return randomCollection
}
