package main

import (
	"crypto/rand"
	"math/big"
	"slices"
)

func setFoil(cardsInPack map[string]int, cmd command) (map[string]int, error) {
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
	return cardsInPack, nil
}
