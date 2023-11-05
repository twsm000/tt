package entities

import "encoding/json"

type Translation struct {
	ID          uint64 `json:"-"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
	Count       int    `json:"count"`
}

func (t Translation) String() string {
	data, _ := json.MarshalIndent(&t, "", "  ")
	return string(data)
}

type Translator interface {
	Translate(word string) (string, error)
}
