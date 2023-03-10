package niantic

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"regexp"

	"github.com/jkassis/pokemoncli/ezhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const otTracerName = "github.com/jkassis/pokemon/niantic"

// API to access the pokemon card service
type API struct {
	BaseURL string
	Headers map[string]string
}

type Res map[string]any

type Err map[string]any

// CardsSearch
type CardsSearchReq struct {
	Query    string `json:"query"`
	Page     int    `json:"page"`
	PageSize int    `json:"limit"`
	OrderBy  string `json:"order"`
	Fields   string `json:"fields"`
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

func (c *CardsSearchReq) Observe(span trace.Span, prefix string) {
	span.SetAttributes(attribute.String(prefix+"Query", c.Query))
	span.SetAttributes(attribute.Int(prefix+"Page", c.Page))
	span.SetAttributes(attribute.Int(prefix+"PageSize", c.PageSize))
	span.SetAttributes(attribute.String(prefix+"OrderBy", c.OrderBy))
	span.SetAttributes(attribute.String(prefix+"Fields", c.Fields))
}

func (c *CardsSearchReq) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func (a *API) CardsSearch(ctx context.Context, req *CardsSearchReq) (Res, Err, error) {
	ctx, span := otel.Tracer(otTracerName).Start(ctx, "CardSearch")
	defer span.End()

	req.Observe(span, "Req")

	ps := url.Values{}
	ps.Add("q", req.Query)
	ps.Add("page", fmt.Sprintf("%d", req.Page))
	ps.Add("pageSize", fmt.Sprintf("%d", req.PageSize))
	ps.Add("orderBy", req.OrderBy)
	ps.Add("select", req.Fields)
	res, apiErr, err := a.httpGet(ctx, "cards", ps, a.Headers)
	if apiErr != nil {
		// with more experience we can decide what apiErrs to log
		// most will probably be usage errors and since this is a CLI
		// we'll let the human deal with these...
		return res, apiErr, err
	}
	return res, apiErr, err
}

// private helpers
var spaces = regexp.MustCompile(`\s+`)

func (a *API) httpGet(ctx context.Context, endpoint string, ps url.Values, headers map[string]string) (Res, Err, error) {
	ctx, span := otel.Tracer(otTracerName).Start(ctx, "httpGet")
	defer span.End()

	if body, err := ezhttp.Get(ctx, a.BaseURL+endpoint, ps, a.Headers); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, nil, err
	} else {
		_, jsonSpan := otel.Tracer(otTracerName).Start(ctx, "httpGet")
		res := make(map[string]any, 0)
		err = json.Unmarshal(body, &res)
		jsonSpan.End()
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, nil, err
		}
		if _, ok := res["error"]; ok {
			// we defer the decision to record these errors to the endpoint
			return nil, Err(res["error"].(map[string]any)), nil
		}
		return res, nil, nil
	}
}
