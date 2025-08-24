package utils

type Message struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data,omitempty"`
}
