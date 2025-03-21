// Copyright 2025 Francisco Oliveto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package insight

const basePath = "https://api.deps.dev/v3"

// Client is a client for the deps.dev API.
type Client struct {
	BasePath string // API endpoint base URL
}

// NewClient returns a new deps.dev API client.
func NewClient() *Client {
	return &Client{BasePath: basePath}
}
