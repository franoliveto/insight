// Copyright 2025 Francisco Oliveto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package insight

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const basePath = "https://api.deps.dev/v3/"

// Client is a client for the deps.dev API.
type Client struct {
	// Base URL for API requests.
	BaseURL *url.URL
}

// NewClient returns a new deps.dev API client.
func NewClient() *Client {
	u, _ := url.Parse(basePath)
	return &Client{BaseURL: u}
}

func (c *Client) get(ctx context.Context, path string, v any) error {
	// path must not have a leading slash.
	path = strings.TrimPrefix(path, "/")

	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Error messages are just text/plain.
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%d %v", resp.StatusCode, err)
		}
		return fmt.Errorf("%d %s", resp.StatusCode, string(data))
	}
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return err
	}
	return nil
}
