package openlibrary

import (
	"encoding/json"
	"strings"
	"time"
)

// DateTimeWrapper unmarshals {"type": "...", "value": "..."} JSON objects, extracting the datetime string.
type DateTimeWrapper struct {
	Value time.Time `json:"value"`
	Type  string    `json:"type"`
}

func (d *DateTimeWrapper) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Append "Z" if missing timezone
	timeStr := aux.Value
	if len(timeStr) > 0 && timeStr[len(timeStr)-1] != 'Z' && !strings.ContainsAny(timeStr[len(timeStr)-6:], "+-") {
		timeStr += "Z"
	}

	parsed, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return err
	}
	d.Type = aux.Type
	d.Value = parsed
	return nil
}
