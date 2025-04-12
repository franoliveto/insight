package insight

import (
	"testing"
)

type Options struct {
	Hash    string `url:"hash,omitempty"`
	System  string `url:"system,omitempty"`
	Version string `url:"version,omitempty"`
}

func TestAddOptions(t *testing.T) {
	testCases := []struct {
		u    string
		opts Options
		want string
	}{
		{"", Options{}, ""},
		{"/a", Options{}, "/a"},
		{"", Options{System: "npm", Version: "18.2.0"}, "?system=npm&version=18.2.0"},
		{"/a/b/c", Options{"ulXBPXrC/UTfnMgHRFVxmjPzdbk=", "npm", "18.2.0"},
			"/a/b/c?hash=ulXBPXrC%2FUTfnMgHRFVxmjPzdbk%3D&system=npm&version=18.2.0"},
	}

	for _, c := range testCases {
		got, err := addOptions(c.u, &c.opts)
		if err != nil {
			t.Errorf("addOptions failed: %v", err)
		}

		if got != c.want {
			t.Errorf("addOptions returned %q; want %q", got, c.want)
		}
	}
}
