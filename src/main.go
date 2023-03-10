package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jkassis/pokemon/api"
)

func Report(resp api.Res, apiErr api.Err, err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else if apiErr != nil {
		out, err := json.MarshalIndent(apiErr, "", "  ")
		if err != nil {
			fmt.Printf("error marshaling response in test code: %v\n", err)
		} else {
			fmt.Println(string(out))
		}
		os.Exit(1)
	} else {
		out, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			fmt.Printf("error marshaling response in test code: %v\n", err)
		} else {
			fmt.Println(string(out))
		}
		os.Exit(0)
	}
}

func main() {
	var limit, page int
	flag.IntVar(&limit, "limit", 10, "limit results returned")
	flag.IntVar(&page, "page", 1, "page of results to receive")

	var query string
	queryDefault := "(types:fire or types:grass) hp:[90 to *] rarity:Rare"
	flag.StringVar(&query, "query", queryDefault, "query in query syntax")

	var fieldsStr string
	fieldsDefault := "name,type,hp,rarity"
	flag.StringVar(&fieldsStr, "fields", fieldsDefault, "fields to retrieve")

	flag.Parse()

	fields := strings.Split(fieldsStr, ",")

	api := &api.API{
		BaseURL: "https://api.pokemontcg.io/v2/",
		Headers: map[string]string{
			"X-Api-Key": "",
		},
	}

	resp, apiErr, err := api.CardsSearch(
		query,
		page,
		limit,
		[]string{"id"},
		fields)
	Report(resp, apiErr, err)
}
