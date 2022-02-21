package util

import (
	"time"
)

// ToEpoch returns time in epoch and an error if any.
func ToEpoch(s string) (int64, error) {
	// Jan 2 15:04:05 2006 MST
	layout := "2006-01-02 15:04"
	t, err := time.Parse(layout, s)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}
