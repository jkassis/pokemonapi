package ezhttp

import (
	"io"
	"net/http"
	"net/url"
)

const DefaultUserAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36"

func Get(url string, ps url.Values, headers map[string]string) ([]byte, error) {
	// encode params
	if len(ps) > 0 {
		url += "?" + ps.Encode()
	}

	// setup the req
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	// make the req
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// read the resp
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// return
	return body, nil
}
