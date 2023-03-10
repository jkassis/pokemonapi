package niantic

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"regexp"

	"github.com/jkassis/pokemoncli/ezhttp"
)

// API to access the pokemon card service
type API struct {
	BaseURL string
	Headers map[string]string
}

type Res map[string]any

type Err map[string]any

// CardsSearch
type CardsSearchReq struct {
	Query    string
	Page     int
	PageSize int
	OrderBy  string
	Fields   string
	flagSet  *flag.FlagSet
}

func (c *CardsSearchReq) Init() {
	c.flagSet = flag.NewFlagSet("cardsearch", flag.ContinueOnError)
	c.flagSet.IntVar(&c.PageSize, "limit", 10, "limit results returned")
	c.flagSet.IntVar(&c.Page, "page", 1, "page of results to receive")

	queryDefault := "(types:fire or types:grass) hp:[90 to *] rarity:Rare"
	c.flagSet.StringVar(&c.Query, "query", queryDefault, "query in query syntax")

	fieldsDefault := "name,type,hp,rarity"
	c.flagSet.StringVar(&c.Fields, "fields", fieldsDefault, "fields to retrieve")

	orderByDefault := "id"
	c.flagSet.StringVar(&c.OrderBy, "order", orderByDefault, "fields to order results by")
}

func (c *CardsSearchReq) Parse(cmd string) error {
	args := spaces.Split(cmd, -1)
	return c.flagSet.Parse(args)
}

func (c *CardsSearchReq) Usage() {
	c.flagSet.Usage()
}

func (a *API) CardsSearch(req *CardsSearchReq) (Res, Err, error) {
	ps := url.Values{}
	ps.Add("q", req.Query)
	ps.Add("page", fmt.Sprintf("%d", req.Page))
	ps.Add("pageSize", fmt.Sprintf("%d", req.PageSize))
	ps.Add("orderBy", req.OrderBy)
	ps.Add("select", req.Fields)
	return a.httpGet("cards", ps, a.Headers)
}

// private helpers
var spaces = regexp.MustCompile(`\s+`)

func (a *API) httpGet(endpoint string, ps url.Values, headers map[string]string) (Res, Err, error) {
	if body, err := ezhttp.Get(a.BaseURL+endpoint, ps, a.Headers); err != nil {
		return nil, nil, err
	} else {
		res := make(map[string]any, 0)
		err = json.Unmarshal(body, &res)
		if err != nil {
			return nil, nil, err
		}
		if _, ok := res["error"]; ok {
			return nil, Err(res["error"].(map[string]any)), nil
		}
		return res, nil, nil
	}
}
