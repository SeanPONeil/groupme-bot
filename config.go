package groupmebot

import (
	"encoding/json"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type BotConfig struct {
	Token       string `json:"token"`
	Name        string `json:"name"`
	GroupId     string `json:"group_id"`
	GroupName   string `json:"group_name"`
	CallbackUrl string `json:"callback_url"`
	AvatarUrl   string `json:"avatar_url"`
}

func (config BotConfig) String() string {
	b, err := json.MarshalIndent(config, " ", "\t")
	check(err)
	return string(b)
}
