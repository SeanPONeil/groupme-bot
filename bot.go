package groupmebot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

// Fluent API for creating GroupMe bots.
type GroupMeBot struct {
	BotConfig
	Id        string `json:"bot_id"`
	OnMessage func(Message)
}

func (bot GroupMeBot) messageHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	var message Message
	err = json.Unmarshal(b, &message)
	if err != nil {
		log.Fatal("Unable to parse message: " + err.Error())
	}
	bot.OnMessage(message)
	return
}

func (bot GroupMeBot) Run() error {
	if bot.CallbackUrl == "" {
		return errors.New("Empty callback URL ")
	}
	callbackUrl, err := url.Parse(bot.CallbackUrl)
	if err != nil {
		return errors.New("Invalid callback url provided " + err.Error())
	}

	_, port, err := net.SplitHostPort(callbackUrl.Host)
	if err != nil {
		port = "80"
	}

	http.HandleFunc(callbackUrl.Path, bot.messageHandler)
	log.Println("Serving on port " + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		return errors.New("ListenAndServe: " + err.Error())
	}
	return nil
}

func (bot GroupMeBot) String() string {
	b, err := json.MarshalIndent(bot, " ", "\t")
	if err != nil {
		return ""
	}
	return string(b)
}

func NewBot(config BotConfig) (*GroupMeBot, error) {
	// Get all registered bots, delete any that have the same name
	bots, err := allBots(config.Token)
	if err != nil {
		return nil, err
	}
	for _, bot := range bots {
		if bot.Name == config.Name {
			_, err := destroy(bot.Id, config.Token)
			if err != nil {
				return nil, err
			}
		}
	}

	// Create a new bot using the provided configuration
	response, err := createBot(config)
	if err != nil {
		return nil, err
	}
	bot := &response.Response.Bot
	bot.Token = config.Token
	bot.AvatarUrl = config.AvatarUrl
	return bot, nil
}
