package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/jinzhu/configor"
)

var (
	ConfigPATH string
)

var Config struct {
	BotToken string `default:"895444309:AADdfntqx8sOXV4qxusM34qVGDDgH1Ls6-C"`
	BotAPIUrl string `default:"http://123.95.96.103"`
	BotServerPort string `default:":80"`

}
func init() {
	flag.StringVar(&ConfigPATH, "c", "config.yml", "config file path")
        flag.Parse()
        configor.Load(&Config, ConfigPATH)
}

func main() {
	log.Println("TGPush Started")

	http.HandleFunc("/telegram/webhook", telegramWebhookHandler)
	http.HandleFunc("/send", sendMessageWebhookHandler)

	if err := setTelegramWebhookPath(telegramWebhookURLGen()); err != nil {
		log.Fatalln("set telegram webhook fail", err.Error())
		return
	}

	if err := http.ListenAndServe(Config.BotServerPort, nil); err != nil {
		log.Fatalln("server run fail", err.Error())
	} else {
		log.Println("server stoped.")
	}
}
