package openlibrary_test // or just 'package openlibrary'

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/obligate5818/bookrecs/internal/openlibrary"
)

func TestUnmarshalEdition(t *testing.T) {
	data, err := os.ReadFile("../testdata/edition.json")
	if err != nil {
		t.Fatal(err)
	}

	var ed openlibrary.Edition
	if err := json.Unmarshal(data, &ed); err != nil {
		t.Fatal(err)
	}

	// Now assert on fields, for example:
	if ed.Title != "Metal from Heaven" {
		t.Errorf("unexpected title: %s", ed.Title)
	}
}
