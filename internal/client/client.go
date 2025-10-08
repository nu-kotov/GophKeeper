package client

import (
	"net/http"
	"time"
)

var baseURL = "http://localhost:8181"

var httpClient = &http.Client{
	Timeout: time.Second * 5,
}
