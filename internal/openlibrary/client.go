package openlibrary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchEdition(ctx context.Context, isbn string) (*Edition, error) {
	url := fmt.Sprintf("https://openlibrary.org/isbn/%s.json", isbn)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch edition: status %d", resp.StatusCode)
	}

	var olEdition Edition
	if err := json.NewDecoder(resp.Body).Decode(&olEdition); err != nil {
		return nil, err
	}

	return &olEdition, nil
}
