package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Res map[string]any
type Err map[string]any

type API struct {
	BaseURL string
	Headers map[string]string
}

func (a *API) httpGet(endpoint string, ps url.Values, headers map[string]string) (Res, Err, error) {
	if body, err := httpGet(a.BaseURL+endpoint, ps, a.Headers); err != nil {
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

func (a *API) CardsSearch(query string, page int, pageSize int, orderBy []string, fields []string) (Res, Err, error) {
	ps := url.Values{}
	ps.Add("q", query)
	ps.Add("page", fmt.Sprintf("%d", page))
	ps.Add("pageSize", fmt.Sprintf("%d", pageSize))
	if orderBy != nil {
		ps.Add("orderBy", strings.Join(orderBy, ","))
	}
	if fields != nil {
		ps.Add("select", strings.Join(fields, ","))
	}

	return a.httpGet("cards", ps, a.Headers)
}
