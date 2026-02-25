package models

import (
	"encoding/json"
	"slices"

	api "codeberg.org/jessienyan/booruview"
)

func (ud UserData) ParseJSON() (UserDataJSON, error) {
	var parsed UserDataJSON
	if ud.Data == "" {
		return parsed, nil
	}

	decompressed := api.DecompressData([]byte(ud.Data))
	err := json.Unmarshal(decompressed, &parsed)
	return parsed, err
}

func (ud *UserData) Set(udJSON UserDataJSON) error {
	udJSON.SearchHistory = udJSON.SearchHistory[:min(len(udJSON.SearchHistory), SearchHistoryLimit)]
	slices.SortFunc(udJSON.SearchHistory, func(a, b SearchHistoryEntry) int {
		// Sort in descending order (newest first)
		return b.Date.Compare(a.Date)
	})

	marshalled, err := json.Marshal(udJSON)
	if err != nil {
		return err
	}

	compressed, err := api.CompressData(marshalled)
	if err != nil {
		return err
	}

	ud.Data = string(compressed)
	return nil
}
