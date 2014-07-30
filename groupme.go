package groupmebot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type createBotResponse struct {
	Meta struct {
		Code int `json:"code"`
	} `json:"meta"`
	Response struct {
		Bot GroupMeBot `json:"bot"`
	} `json:"response"`
}

// Message represents a new message posted to a group.
type Message struct {
	ID          string                   `json:"id"`
	SourceGuid  string                   `json:"source_guid"`
	CreatedAt   int                      `json:"created_at"`
	UserID      string                   `json:"user_id"`
	GroupID     string                   `json:"group_id"`
	Name        string                   `json:"name"`
	AvatarURL   string                   `json:"avatar_url"`
	Text        string                   `json:"text"`
	System      bool                     `json:"system"`
	FavoritedBy []string                 `json:"favorited_by"`
	Attachments []map[string]interface{} `json:"attachments"`
}

// Create a bot
func createBot(config BotConfig) (response *createBotResponse, err error) {
	var createEndpoint = "https://api.groupme.com/v3/bots?token=" + config.Token
	payload := map[string]interface{}{
		"bot": config,
	}
	payloadJSON, err := json.Marshal(payload)
	resp, err := http.Post(createEndpoint, "application/json", bytes.NewReader(payloadJSON))
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(b, &response)

	return
}

// Get a list of registered bots
func allBots(token string) (bots []GroupMeBot, err error) {
	var botsEndpoint = "https://api.groupme.com/v3/bots?token=" + token
	resp, err := http.Get(botsEndpoint)
	defer resp.Body.Close()
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading body")
		return nil, err
	}

	type Response struct {
		Meta struct {
			Code int `json:"code"`
		} `json:"meta"`
		Response []GroupMeBot `json:"response"`
	}

	var response Response
	err = json.Unmarshal(b, &response)
	return response.Response, nil
}

// Remove a bot that you have created
func destroy(botID string, token string) (int, error) {
	var destroyEndpoint = "https://api.groupme.com/v3/bots/destroy?token=" + token
	payload := map[string]interface{}{
		"bot_id": botID,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	resp, err := http.Post(destroyEndpoint, "application/json", bytes.NewReader(payloadJSON))
	defer resp.Body.Close()
	check(err)
	return resp.StatusCode, nil
}

// SendMessage sends a new message to the group.
func (bot GroupMeBot) SendMessage(text string) (int, error) {
	var messageEndpoint = "https://api.groupme.com/v3/bots/post"
	payload := map[string]interface{}{
		"bot_id": bot.ID,
		"text":   text,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	resp, err := http.Post(messageEndpoint, "application/json", bytes.NewReader(payloadJSON))
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return resp.StatusCode, nil
}
