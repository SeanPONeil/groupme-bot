package groupme-bot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

// GroupMeBot is an abstraction over the GroupMe Bot API, and
// allows for easy sending of messages and receiving new messages.
type GroupMeBot struct {
	BotConfig
	ID        string `json:"bot_id"`
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

// Run a server at the callback URL provided in BotConfig.
func (bot GroupMeBot) Run() error {
	if bot.CallbackURL == "" {
		return errors.New("Empty callback URL ")
	}
	callbackURL, err := url.Parse(bot.CallbackURL)
	if err != nil {
		return errors.New("Invalid callback url provided " + err.Error())
	}

	_, port, err := net.SplitHostPort(callbackURL.Host)
	if err != nil {
		port = "80"
	}

	http.HandleFunc(callbackURL.Path, bot.messageHandler)
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

// NewBot destroys any currently registered bots with the same
// name as BotConfig.Name, and then registers and returns a new
// GroupMeBot.
func NewBot(config BotConfig) (*GroupMeBot, error) {
	// Get all registered bots, delete any that have the same name
	bots, err := allBots(config.Token)
	if err != nil {
		return nil, err
	}
	for _, bot := range bots {
		if bot.Name == config.Name {
			_, err := destroy(bot.ID, config.Token)
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
	bot.AvatarURL = config.AvatarURL
	return bot, nil
}
