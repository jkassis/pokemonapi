package ezhttp

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const otTracerName = "github.com/jkassis/pokemoncli/ezhttp"

const DefaultUserAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36"

func Get(ctx context.Context, url string, ps url.Values, headers map[string]string) ([]byte, error) {
	ctx, span := otel.Tracer(otTracerName).Start(ctx, "Get")
	defer span.End()

	// encode params
	if len(ps) > 0 {
		url += "?" + ps.Encode()
	}

	// setup the req
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	// make the req
	_, doSpan := otel.Tracer(otTracerName).Start(ctx, "http.DefaultClient.Do(req)")
	resp, err := http.DefaultClient.Do(req)
	doSpan.End()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// read the resp
	defer resp.Body.Close()
	_, readBodySpan := otel.Tracer(otTracerName).Start(ctx, "io.ReadAll(resp.Body)")
	body, err := io.ReadAll(resp.Body)
	readBodySpan.End()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// return
	return body, nil
}
