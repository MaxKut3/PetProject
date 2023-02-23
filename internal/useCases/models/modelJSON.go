package models

import "encoding/json"

type JSON struct {
	ID      int `json:"id"`
	Balance int `json:"balance"`
}

func Parser(v []byte, m JSON) error {
	err := json.Unmarshal(v, &m)
	if err != nil {
		return err
	}
	return nil
}
