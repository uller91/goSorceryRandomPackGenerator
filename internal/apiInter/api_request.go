package apiInter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func RequestData(url string) []byte {

	//request for a data
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	data, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, data)
	}
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	//fmt.Printf("%s", body)
	return data
}

func RequestCard(url string) CardsData {
	data := RequestData(url)

	//unmarshal data
	CardsData := CardsData{}
	err := json.Unmarshal(data, &CardsData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	return CardsData
}
