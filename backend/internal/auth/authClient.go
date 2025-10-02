package auth

import (
	"net/http"
	"time"
)

var Client = &http.Client{
	Timeout: 10 * time.Second,
}
