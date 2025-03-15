// Copyright 2025 Francisco Oliveto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package insight

import (
	"net/http"
)

const basePath = "https://api.deps.dev/v3"

// Client is a client for the deps.dev API.
type Client struct {
	BasePath string // API endpoint base URL
	client   *http.Client
}

// NewClient returns a new deps.dev API client using the provided http.Client as transport.
func NewClient(endpoint string, c *http.Client) *Client {
	if c == nil {
		c = http.DefaultClient
	}
	if endpoint == "" {
		endpoint = basePath
	}
	return &Client{
		BasePath: endpoint,
		client:   c,
	}
}
