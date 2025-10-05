package main

import (
	"slices"
	"strings"
	"crypto/rand"
	"errors"
	"math/big"
	"fmt"
)

func setSet(s *state, cmd command) (string, error) {
	set := ""
	
	if tag := slices.Index(cmd.arguments, "-s"); tag != -1 {
		if len(cmd.arguments) >= tag+2 && cmd.arguments[tag+1][0:1] != "-" {
			set = strings.Title(cmd.arguments[tag+1])

			switch set {
			case "A":
				set = "Alpha"
			case "B":
				set = "Beta"
			case "AL", "Al":
				set = "Arthurian Legends"
				//add more at apropiate release
			}

			if !slices.Contains(s.config.Sets, strings.Title(set)) && set != "All" && set != "Random" {
				fmt.Println("No such set in DB! Generating the random pack...")
				set = "Random"
			}

			if set == "Random" {
				randomSetNumber, _ := rand.Int(rand.Reader, big.NewInt(int64(len(s.config.Sets))))
				set = s.config.Sets[int(randomSetNumber.Int64())]
			}

			return set, nil

		} else {
			return "", errors.New("No set name given after -s tag")
		}
	} else {
		randomSetNumber, _ := rand.Int(rand.Reader, big.NewInt(int64(len(s.config.Sets))))
		set = s.config.Sets[int(randomSetNumber.Int64())]
	}

	return set, nil
}