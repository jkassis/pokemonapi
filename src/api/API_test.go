package api

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Report(t *testing.T, resp Res, apiErr Err, err error) {
	if err != nil {
		t.Error(err)
	} else if apiErr != nil {
		t.Error(apiErr)
	} else {
		out, err := json.Marshal(resp)
		if err != nil {
			t.Error(fmt.Errorf("error marshaling response in test code: %v", err))
		} else {
			fmt.Println(string(out))
		}
	}
}

func TestCardGet(t *testing.T) {
	api := &API{
		BaseURL: "https://api.pokemontcg.io/v2/",
		Headers: map[string]string{
			"X-Api-Key": "",
		},
	}

	// resp, apiErr, err := api.CardsSearch("hp:[90 to *] rarity:Rare", 1, 10, nil, nil)
	resp, apiErr, err := api.CardsSearch(
		"(types:fire or types:grass) hp:[90 to *] rarity:Rare",
		1,
		10,
		[]string{"id"},
		[]string{"name", "type", "hp", "rarity"})
	Report(t, resp, apiErr, err)
}
