package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type ConfData struct {
	BotName            string `yaml:"bot_name"`
	IconEmoji          string `yaml:"icon_emoji"`
	MessageText        string `yaml:"message_text"`
	WebHookUri         string `yaml:"web_hook_uri"`
	ItunesAppId        string `yaml:"itunes_app_id"`
	DBPath             string `yaml:"db_path"`
	TmpDir             string `yaml:"tmp_dir"`
	GCSBucketId        string `yaml:"gcs_bucket_id"`
	AndroidPackageName string `yaml:"android_package_name"`
	MaxAttachmentCount int    `yaml:"max_attachment_count"`
	PlatformName       string
}

var (
	configFile  = flag.String("c", "config/config.yml", "location of config file")
	platformKey = flag.String("p", "ios", "platform key: ios or android")
)

func LoadConfig() *ConfData {
	d := ConfData{}
	buf, err := ioutil.ReadFile(*configFile)

	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(buf, &d); err != nil {
		log.Fatal(err)
	}

	d.PlatformName = *platformKey
	return &d
}
