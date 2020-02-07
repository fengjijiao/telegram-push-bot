package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	simplejson "github.com/bitly/go-simplejson"
)

func sendMessageToTelegram(chatID string, message string) error {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", Config.BotToken)
	formData := url.Values{
		"chat_id":    {chatID},
		"text":       {message},
		"parse_mode": {"Markdown"},
	}
	resp, err := http.PostForm(endpoint, formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	sJSON, err := simplejson.NewJson(bJSON)
	if err != nil {
		return err
	}
	if sJSON.Get("ok").MustBool() {
		return nil
	} else {
		return errors.New(sJSON.Get("description").MustString())
	}
}

func setTelegramWebhookPath(path string) error {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=%s", Config.BotToken, url.QueryEscape(path))
	resp, err := http.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	sJSON, err := simplejson.NewJson(bJSON)
	if err != nil {
		return err
	}
	if sJSON.Get("ok").MustBool() {
		return nil
	} else {
		return errors.New(sJSON.Get("description").MustString())
	}
}

func sendMessageURLGen(chatid string) string {
	sign := stringSign(chatid, Config.BotToken)
	return fmt.Sprintf("%s/send?chatid=%s&sign=%s", Config.BotAPIUrl, chatid, sign)
}

func telegramWebhookURLGen() string {
	return fmt.Sprintf("%s/telegram/webhook", Config.BotAPIUrl)
}

func parseTelegramWebhook(data []byte) (chatid, text string, err error) {
	sJSON, err := simplejson.NewJson(data)
	if err != nil {
		return "", "", err
	}
	chatid = strconv.FormatInt(sJSON.Get("message").Get("chat").Get("id").MustInt64(), 10)
	text = sJSON.Get("message").Get("text").MustString()
	return chatid, text, nil
}
