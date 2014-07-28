package groupmebot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type CreateBotResponse struct {
	Meta struct {
		Code int `json:"code"`
	} `json:"meta"`
	Response struct {
		Bot GroupMeBot `json:"bot"`
	} `json:"response"`
}

type Message struct {
	Id          string                   `json:"id"`
	SourceGuid  string                   `json:"source_guid"`
	CreatedAt   int                      `json:"created_at"`
	UserId      string                   `json:"user_id"`
	GroupId     string                   `json:"group_id"`
	Name        string                   `json:"name"`
	AvatarUrl   string                   `json:"avatar_url"`
	Text        string                   `json:"text"`
	System      bool                     `json:"system"`
	FavoritedBy []string                 `json:"favorited_by"`
	Attachments []map[string]interface{} `json:"attachments"`
}

type PostMessageBody struct {
	Id   string `json:"bot_id"`
	Text string `json:"text"`
}

// Create a bot
func createBot(config BotConfig) (response *CreateBotResponse, err error) {
	var createEndpoint = "https://api.groupme.com/v3/bots?token=" + config.Token
	payload := map[string]interface{}{
		"bot": config,
	}
	payloadJson, err := json.Marshal(payload)
	resp, err := http.Post(createEndpoint, "application/json", bytes.NewReader(payloadJson))
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
func destroy(botId string, token string) (int, error) {
	var destroyEndpoint = "https://api.groupme.com/v3/bots/destroy?token=" + token
	payload := map[string]interface{}{
		"bot_id": botId,
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	resp, err := http.Post(destroyEndpoint, "application/json", bytes.NewReader(payloadJson))
	defer resp.Body.Close()
	check(err)
	return resp.StatusCode, nil
}

// Send a message
func (bot GroupMeBot) SendMessage(text string) (int, error) {
	var messageEndpoint = "https://api.groupme.com/v3/bots/post"
	payload := map[string]interface{}{
		"bot_id": bot.Id,
		"text":   text,
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	resp, err := http.Post(messageEndpoint, "application/json", bytes.NewReader(payloadJson))
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return resp.StatusCode, nil
}
