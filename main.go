package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
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

	// Load INI config
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
}

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("No Discord bot token provided")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
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

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	select {}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	var responses []string

	for _, item := range matchConfig {
		pattern := regexp.MustCompile(item.Pattern)
		matches := pattern.FindAllString(m.Content, -1)

		for _, match := range matches {
			vxURL := strings.Replace(match, item.BaseLink, item.VxLink, 1)
			log.Printf("Sending %s link: %s", item.VxLink, vxURL)
			responses = append(responses, fmt.Sprintf("[⠀](%s)", vxURL))
		}
	}

	if len(responses) > 0 {
		go suppressEmbed(s, m.ID, m.ChannelID)

		response := strings.Join(responses, " ")
		_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: response,
			Reference: &discordgo.MessageReference{
				MessageID: m.ID,
				ChannelID: m.ChannelID,
				GuildID:   m.GuildID,
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		})
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}
