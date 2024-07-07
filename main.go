package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"gopkg.in/ini.v1"
)

type MatchConfig struct {
	Pattern  string `json:"pattern"`
	BaseLink string `json:"base_link"`
	VxLink   string `json:"vx_link"`
}

var (
	matchConfig []MatchConfig
	config      *ini.File
	dg          *discordgo.Session
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
	}

	configFile, err := os.ReadFile("match.json")
	if err != nil {
		log.Fatal("Error reading match.json file:", err)
	}

	err = json.Unmarshal(configFile, &matchConfig)
	if err != nil {
		log.Fatal("Error parsing match.json file:", err)
	}

	config, err = ini.Load("config.ini")
	if err != nil {
		log.Fatal("Error loading config.ini file:", err)
	}

	if config.Section("Settings").Key("EndpointOverride").MustBool(false) {
		log.Println("Endpoint overrides found, ignoring defaults")
		for i, item := range matchConfig {
			serviceName := strings.ReplaceAll(item.BaseLink, ".com", "")
			vxLinkOverrideKey := fmt.Sprintf("%s_vx_link", serviceName)
			vxLinkOverride := config.Section("vx_links").Key(vxLinkOverrideKey).String()
			if vxLinkOverride != "" && strings.ToLower(vxLinkOverride) != "default" {
				matchConfig[i].VxLink = vxLinkOverride
			}
		}
	}

	initDatabase()

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("No Discord bot token provided")
	}

	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		} else if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
			data := i.ApplicationCommandData()
			if data.Name == "enable" || data.Name == "disable" {
				for _, option := range data.Options {
					if option.Name == "service" {
						handleServiceAutocomplete(s, i, option)
						return
					}
				}
			}
		}
	})

	dg.Identify.Intents = discordgo.IntentsGuildMessages
}

func main() {
	err := dg.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}

	log.Printf("%s has connected to Discord!", dg.State.User.Username)

	err = dg.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: string(discordgo.StatusDoNotDisturb),
		Activities: []*discordgo.Activity{
			{
				Name: "your messages",
				Type: discordgo.ActivityTypeWatching,
			},
		},
	})
	if err != nil {
		log.Println("Error updating status:", err)
	}

	log.Println("Registering slash commands...")
	commands, err := dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, "", []*discordgo.ApplicationCommand{
		{
			Name:        "enable",
			Description: "Enable embed services",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "service",
					Description:  "Service to enable",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "disable",
			Description: "Disable embed services",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "service",
					Description:  "Service to disable",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "settings",
			Description: "Shows the current embed settings for this server",
		},
	})
	if err != nil {
		log.Fatalf("Error registering slash commands: %v", err)
	}
	log.Printf("Registered %d slash commands", len(commands))

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	select {}
}
