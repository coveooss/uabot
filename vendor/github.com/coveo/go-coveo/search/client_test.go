package search_test

import (
	"github.com/coveo/go-coveo/search"
	"testing"
)

func TestClient(t *testing.T) {
	config := search.Config{}
	_, err := search.NewClient(config)
	if err != nil {
		t.Fatalf("unexpected error.  expected %v, actual %v", nil, err)
	}
}
