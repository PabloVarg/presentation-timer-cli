package api

import (
	"fmt"
	"net/http"
	"time"
)

type APIClient struct {
	HTTPClient http.Client
}

type KeyValueRetriever = func(key string) (string, bool)

const (
	API_URL_KEY = "API_URL"
)

type PaginatedResponse[T any] struct {
	Data     T        `json:"data"`
	PageInfo struct{} `json:"page_info"`
}

func NewDefaultClient() http.Client {
	return http.Client{
		Timeout: 10 * time.Second,
	}
}

func (c APIClient) Url(get KeyValueRetriever) string {
	url, ok := get(API_URL_KEY)
	if !ok {
		panic(fmt.Sprintf("no %s present", API_URL_KEY))
	}

	return url
}
