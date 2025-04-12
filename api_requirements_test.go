package insights

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TODO: add more and better tests.
func TestGetRequirements(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/systems/npm/packages/react/versions/18.2.0:requirements", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"npm":{}}`)
	})

	want := &Requirements{NPM: NPM{}}

	got, err := client.GetRequirements(context.Background(), "npm", "react", "18.2.0")
	if err != nil {
		t.Errorf("GetRequirements failed: %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Errorf("GetRequirements returned %+v; want %+v", got, want)
	}
}
