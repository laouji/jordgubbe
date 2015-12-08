package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConfData struct {
	BotName     string `yaml:"bot_name"`
	IconEmoji   string `yaml:"icon_emoji"`
	MessageText string `yaml:"message_text"`
	WebHookUri  string `yaml:"web_hook_uri"`
	ItunesAppId string `yaml:"itunes_app_id"`
	DBPath      string `yaml:"db_path"`
}

func LoadConfig() *ConfData {
	buf, err := ioutil.ReadFile("config/config_local.yml")

	d := ConfData{}
	if err != nil {
		buf, err := ioutil.ReadFile("config/config.yml")
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(buf, &d); err != nil {
			panic(err)
		}
		return &d
	}

	if err := yaml.Unmarshal(buf, &d); err != nil {
		panic(err)
	}

	return &d
}
