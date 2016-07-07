package search_test

import (
	"testing"

	"github.com/coveo/go-coveo/search"
)

func TestClient(t *testing.T) {
	config := search.Config{}
	_, err := search.NewClient(config)
	if err != nil {
		t.Fatalf("unexpected error.  expected %v, actual %v", nil, err)
	}
}
