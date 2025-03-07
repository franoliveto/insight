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

type Package struct {
	PackageKey PackageKey
	Versions   []Version
}

type PackageKey struct {
	System string
	Name   string
}

type Version struct {
	VersionKey  VersionKey
	PublishedAt string
	IsDefault   bool
}

type VersionKey struct {
	System  string
	Name    string
	Version string
}

const basePath = "https://api.deps.dev/v3"

// Client is a client for the deps.dev API.
type Client struct {
	client   *http.Client
	BasePath string // API endpoint base URL
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

// GetPackage returns information about a package, including a list of its available versions.
func (c *Client) GetPackage(system, name string) (*Package, error) {
	url := c.BasePath + "/systems/" + url.PathEscape(system) + "/packages/" + url.PathEscape(name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("%s", resp.Status)
		}
		return nil, fmt.Errorf("%s", string(data))
	}

	p := new(Package)
	if err := json.NewDecoder(resp.Body).Decode(p); err != nil {
		return nil, err
	}
	return p, nil
}
