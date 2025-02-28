package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"net/http"
	"tg-sender/config"
	"tg-sender/lib/logger/sl"
)

type TgClient struct {
	baseURL string
	log     *slog.Logger
}

func NewTgClient(cfg config.TgConfig, log *slog.Logger) *TgClient {
	return &TgClient{baseURL: cfg.BaseURL + "/" + cfg.Token + "/", log: log}
}

type Command struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

func (tg *TgClient) SendMessageAll(userID uint64, text string) bool {
	url := tg.baseURL + "sendMessage"
	data := map[string]interface{}{
		"chat_id":    userID,
		"text":       text,
		"parse_mode": "markdown",
	}
	return tg.postRequest(url, data)
}

func (tg *TgClient) SetMyCommands(commands []Command) bool {
	url := tg.baseURL + "setMyCommands"
	data := map[string]interface{}{
		"commands": commands,
	}
	return tg.postRequest(url, data)
}

func (tg *TgClient) SendMessage(userID int, title, filePath string) bool {
	url := tg.baseURL + "sendMessage"
	data := map[string]interface{}{
		"chat_id":    userID,
		"text":       title,
		"parse_mode": "markdown",
		"reply_markup": map[string]interface{}{
			"inline_keyboard": [][]map[string]string{{{"text": "Open link", "url": filePath}}},
		},
	}
	return tg.postRequest(url, data)
}

func (tg *TgClient) SendDocument(userID int, filePath string) bool {
	url := tg.baseURL + "sendDocument"
	data := map[string]interface{}{
		"chat_id":  userID,
		"document": filePath,
	}
	return tg.postRequest(url, data)
}

func (tg *TgClient) SendAudio(userID int, title, filePath, photoURL string, duration int) bool {
	url := tg.baseURL + "sendAudio"
	data := map[string]interface{}{
		"chat_id":   userID,
		"audio":     filePath,
		"title":     title,
		"thumbnail": photoURL,
		"duration":  duration,
	}
	return tg.postRequest(url, data)
}

func (tg *TgClient) SendPhoto(userID uint64, title, filePath, photoURL string) bool {
	url := tg.baseURL + "sendPhoto"
	data := map[string]interface{}{
		"chat_id":    userID,
		"photo":      photoURL,
		"caption":    title,
		"parse_mode": "markdown",
		"reply_markup": map[string]interface{}{
			"inline_keyboard": [][]map[string]string{{{"text": "Open link", "url": filePath}}},
		},
	}
	return tg.postRequest(url, data)
}

func (tg *TgClient) GetWebhook() bool {
	url := tg.baseURL + "getWebhookInfo"
	resp, err := http.Get(url)
	if err != nil {
		tg.log.Info("Error making request:", sl.Err(err))
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	tg.log.Info("Webhook Info:", slog.String("body", string(body)))
	return resp.StatusCode == http.StatusOK
}

func (tg *TgClient) postRequest(url string, data interface{}) bool {
	jsonData, err := json.Marshal(data)
	if err != nil {
		tg.log.Error("Error marshalling JSON:", sl.Err(err))
		return false
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		tg.log.Error("Error making request:", sl.Err(err))
		return false
	}
	defer resp.Body.Close()

	tg.log.Info("Request to %s finished with status %d", slog.String("url", url), slog.Int("code", resp.StatusCode))
	return resp.StatusCode == http.StatusOK
}
