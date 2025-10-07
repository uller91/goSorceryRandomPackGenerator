package main

import (
	"slices"
	"errors"
	"strconv"
)

func setNumber(cmd command) (int, error) {
	numberOfPacks := 1

	if tag := slices.Index(cmd.arguments, "-n"); tag != -1 {
		if len(cmd.arguments) >= tag+2 && cmd.arguments[tag+1][0:1] != "-" {
			packs, err := strconv.ParseInt(cmd.arguments[tag+1], 10, 64)
			if err != nil {
				return 0, errors.New("Number of packs (-n tag) can only be an integer!")
			}
			numberOfPacks = int(packs)
		} else {
			return 0, errors.New("No number of packs (or a neganive number) was given after -n tag")
		}
	}

	return numberOfPacks, nil
}