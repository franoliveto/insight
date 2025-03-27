// Copyright 2025 Francisco Oliveto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package insight

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const basePath = "https://api.deps.dev/v3"

// Client is a client for the deps.dev API.
type Client struct {
	BasePath string // API endpoint base URL
}

// NewClient returns a new deps.dev API client.
func NewClient() *Client {
	return &Client{BasePath: basePath}
}

func (c *Client) get(path string, v any) error {
	url, _ := url.JoinPath(c.BasePath, path)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", string(data))
	}
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}
