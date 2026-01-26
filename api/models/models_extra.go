package models

import "encoding/json"

func (ud UserData) ParseJSON() (UserDataJSON, error) {
	var parsed UserDataJSON
	err := json.Unmarshal([]byte(ud.Data), &parsed)
	return parsed, err
}

func (ud *UserData) Set(udJSON UserDataJSON) error {
	marshalled, err := json.Marshal(udJSON)
	if err != nil {
		return err
	}
	ud.Data = string(marshalled)
	return nil
}
