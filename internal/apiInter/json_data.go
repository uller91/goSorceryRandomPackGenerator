package apiInter

import (
	"time"
)

type CardsData []struct {
	Name     string `json:"name"`
	Guardian struct {
		Rarity     string `json:"rarity"`
		Type       string `json:"type"`
		RulesText  string `json:"rulesText"`
		Cost       any    `json:"cost"`
		Attack     any    `json:"attack"`
		Defence    any    `json:"defence"`
		Life       any    `json:"life"`
		Thresholds struct {
			Air   int `json:"air"`
			Earth int `json:"earth"`
			Fire  int `json:"fire"`
			Water int `json:"water"`
		} `json:"thresholds"`
	} `json:"guardian"`
	Sets     []struct {
		Name       string    `json:"name"`
	} `json:"sets"`
}

type CardsDataFull []struct {
	Name     string `json:"name"`
	Guardian struct {
		Rarity     string `json:"rarity"`
		Type       string `json:"type"`
		RulesText  string `json:"rulesText"`
		Cost       any    `json:"cost"`
		Attack     any    `json:"attack"`
		Defence    any    `json:"defence"`
		Life       any    `json:"life"`
		Thresholds struct {
			Air   int `json:"air"`
			Earth int `json:"earth"`
			Fire  int `json:"fire"`
			Water int `json:"water"`
		} `json:"thresholds"`
	} `json:"guardian"`
	Elements string `json:"elements"`
	SubTypes string `json:"subTypes"`
	Sets     []struct {
		Name       string    `json:"name"`
		ReleasedAt time.Time `json:"releasedAt"`
		Metadata   struct {
			Rarity     string `json:"rarity"`
			Type       string `json:"type"`
			RulesText  string `json:"rulesText"`
			Cost       any    `json:"cost"`
			Attack     any    `json:"attack"`
			Defence    any    `json:"defence"`
			Life       any    `json:"life"`
			Thresholds struct {
				Air   int `json:"air"`
				Earth int `json:"earth"`
				Fire  int `json:"fire"`
				Water int `json:"water"`
			} `json:"thresholds"`
		} `json:"metadata"`
		Variants []struct {
			Slug       string `json:"slug"`
			Finish     string `json:"finish"`
			Product    string `json:"product"`
			Artist     string `json:"artist"`
			FlavorText string `json:"flavorText"`
			TypeText   string `json:"typeText"`
		} `json:"variants"`
	} `json:"sets"`
}