package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type HttpTransport struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	httpDo     func(c *http.Client, req *http.Request) (*http.Response, error)
}

func (h HttpTransport) Request(req Request) ([]interface{}, error) {
	var raw []interface{}
	//log.Printf("HttpTransport.Request - visit %s", req.RefURL)
	rel, err := url.Parse(req.RefURL)
	if err != nil {
		log.Printf("HttpTransport.Request - url.Parse failed : %v", err)
		return nil, err
	}
	if req.Params != nil {
		rel.RawQuery = req.Params.Encode()
	}
	if req.Data == nil {
		req.Data = map[string]interface{}{}
	}

	b, err := json.Marshal(req.Data)
	if err != nil {
		log.Printf("HttpTransport.Request - json.Marshal failed : %v", err)
		return nil, err
	}

	body := bytes.NewReader(b)

	u := h.BaseURL.ResolveReference(rel)
	httpReq, err := http.NewRequest(req.Method, u.String(), body)
	for k, v := range req.Headers {
		httpReq.Header.Add(k, v)
	}
	if err != nil {
		log.Printf("HttpTransport.Request - http.NewRequest failed : %v", err)
		return nil, err
	}

	resp, err := h.do(httpReq, &raw)
	if err != nil {
		log.Printf("HttpTransport.Request - http.do failed : %v", err)
		if resp != nil {
			return nil, fmt.Errorf("could not parse response: %s", resp.Response.Status)
		} else {
			return nil, fmt.Errorf("%v", err)
		}
	}

	return raw, nil
}

// Do executes API request created by NewRequest method or custom *http.Request.
func (h HttpTransport) do(req *http.Request, v interface{}) (*Response, error) {
	log.Printf("visit : %s", req.URL)
	resp, err := h.httpDo(h.HTTPClient, req)
	if err != nil {
		log.Printf("HttpTransport.do - h.httpDo failed : %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)
	err = checkResponse(response)
	if err != nil {
		if response.String() == `{ "error": "ERR_RATE_LIMIT" }` {
			return nil, fmt.Errorf("rate limit")
		}

		return nil, err
	}

	if v != nil {
		err = json.Unmarshal(response.Body, v)
		if err != nil {
			log.Printf("HttpTransport.do - json.Unmarshal failed : %v", err)
			return response, err
		}
	}

	return nil, nil
}
