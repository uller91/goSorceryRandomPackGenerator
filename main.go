package main

import (
	"fmt"
	"github.com/uller91/goSorceryDraftDB/internal/apiInter"
)

func main() {
	fmt.Printf("Hello World\n")

	apiUrl := apiInter.BaseUrl
	cards := apiInter.RequestCard(apiUrl)
	dbSize := len(cards)
	fmt.Println(cards[0])
	fmt.Println(cards[dbSize-1])
}
