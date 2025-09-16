package apiInter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func RequestData(url string) ([]byte, error) {

	//request for a data
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, data)
	}
	if err != nil {
		return nil, err
	}

	//fmt.Printf("%s", body)
	return data, nil
}

func RequestCard(url string) (CardsData, error) {
	data, err := RequestData(url)
	if err != nil {
		return nil, err
	}

	//unmarshal data
	CardsData := CardsData{}
	err = json.Unmarshal(data, &CardsData)
	if err != nil {
		return nil, err
	}

	return CardsData, nil
}
