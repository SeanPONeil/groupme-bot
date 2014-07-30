package groupmebot

import (
	"encoding/json"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// BotConfig describes a GroupMeBot. If CallbackUrl is set,
// new messages posted to the group can be received by calling
// GroupMeBot.Run()
type BotConfig struct {
	Token       string `json:"token"`
	Name        string `json:"name"`
	GroupID     string `json:"group_id"`
	GroupName   string `json:"group_name"`
	CallbackURL string `json:"callback_url"`
	AvatarURL   string `json:"avatar_url"`
}

func (config BotConfig) String() string {
	b, err := json.MarshalIndent(config, " ", "\t")
	check(err)
	return string(b)
}
