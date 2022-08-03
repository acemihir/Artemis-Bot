package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"strings"

	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
	"gopkg.in/ini.v1"
)

type client struct {
	AuthToken    string `ini:"auth_token"`
	GuildID      string `ini:"guild_id"`
	ActivityType int    `ini:"activity_type"`
	ActivityText string `ini:"activity_text"`
	ActivityUrl  string `ini:"activity_url"`
}

type commands struct {
	RetractAll bool `ini:"retract_all"`
	SubmitAll  bool `ini:"submit_all"`
}

type data struct {
	CacheExpiry   int    `ini:"cache_expiry"`
	EncryptionKey string `ini:"encryption_key"`
	HMACKey       string `ini:"hmac_key"`
}

type misc struct {
	HelpSpacingBase int `ini:"help_spacing_base"`
}

type IniConfig struct {
	AppMode  string   `ini:"app_mode"`
	Client   client   `ini:"client"`
	Commands commands `ini:"commands"`
	Data     data     `ini:"data"`
	Misc     misc     `ini:"misc"`
}

const defaultConfig string = `# Either development or production
app_mode = production

[client]
# The token used to authtenticate the bot
auth_token = ""
# The guild where the commands get registered (Only active when running in development mode)
guild_id = ""
# https://discord.com/developers/docs/game-sdk/activities#data-models-activitytype-enum
activity_type = 3
# The text displayed for the activity
activity_text = "the wind guide my arrows"
# Optional activity Url
activity_url = ""

[commands]
# Will throw error when there are no commands registered
retract_all = false
# Make sure the commands are deleted first
submit_all = false

[data]
# The TTL for cache entries
cache_expiry = 480
# The 32-byte base64 key used to encrypt/decrypt semi-private data
encryption_key = <ENC_KEY>
# HMAC key used while generating sha256 hashes
hmac_key = <HMAC_KEY>

[misc]
# Amount of base spacings in the help command
help_spacing_base = 9
`

var Cfg *IniConfig

func LoadConfig() {
	_, ex := os.Stat("app.ini")

	if os.IsNotExist(ex) {
		f, ex := os.Create("app.ini")
		if ex != nil {
			utils.Cout("[ERROR] Failed to create the configuration file: %v", utils.Red, ex)
			os.Exit(1)
		}

		defer f.Close()

		keyBuffer := make([]byte, 32)
		_, ex = rand.Read(keyBuffer)
		if ex != nil {
			utils.Cout("[ERROR] Could not generate encryption key: %v", utils.Red, ex)
			os.Exit(1)
		}

		key := base64.URLEncoding.EncodeToString(keyBuffer)
		contents := strings.Replace(defaultConfig, "<ENC_KEY>", key, 1)
		contents = strings.Replace(contents, "<HMAC_KEY>", utils.RandomString("", 16), 1)

		_, ex = f.WriteString(contents)
		if ex != nil {
			utils.Cout("[ERROR] Failed to write the default configuration: %v", utils.Red, ex)
			os.Exit(1)
		}

		utils.Cout("Please configure the bot in the app.ini file and start the bot again.", utils.Gray)
		os.Exit(0)
	}

	cfg, ex := ini.Load("app.ini")
	if ex != nil {
		utils.Cout("[ERROR] Could not load the configuration file: %v", utils.Red, ex)
		os.Exit(1)
	}

	Cfg = &IniConfig{
		AppMode: "production",
		Client: client{
			ActivityType: 3,
			ActivityText: "the wind guide my arrows",
			ActivityUrl:  "",
		},
		Data: data{
			CacheExpiry: 480,
		},
		Misc: misc{
			HelpSpacingBase: 9,
		},
	}
	ex = cfg.MapTo(Cfg)
	if ex != nil {
		utils.Cout("[ERROR] Could not map to struct: %v", utils.Red, ex)
		os.Exit(1)
	}
}
