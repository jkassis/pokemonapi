package niantic

import (
	"context"
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
		out, err := json.MarshalIndent(resp, "", "  ")
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

	req := CardsSearchReq{
		Query:    "(types:fire or types:grass) hp:[90 to *] rarity:Rare",
		Page:     1,
		PageSize: 10,
		OrderBy:  "id",
		Fields:   "name,type,hp,rarity",
	}
	resp, apiErr, err := api.CardsSearch(context.Background(), &req)
	Report(t, resp, apiErr, err)
}
