package insights

import "testing"

func TestNewClient(t *testing.T) {
	c := NewClient()
	if got, want := c.BaseURL.String(), basePath; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

// TODO: add test for Client.get method.
