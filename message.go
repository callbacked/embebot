package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	var responses []string
	for _, item := range matchConfig {
		if !getServerSettingFromDB(m.GuildID, strings.ReplaceAll(item.BaseLink, ".com", "")) {
			continue
		}
		pattern := regexp.MustCompile(item.Pattern)
		matches := pattern.FindAllString(m.Content, -1)
		for _, match := range matches {
			vxURL := strings.Replace(match, item.BaseLink, item.VxLink, 1)
			log.Printf("Sending %s link: %s", item.VxLink, vxURL)
			responses = append(responses, fmt.Sprintf("[⠀](%s)", vxURL))
		}
	}

	if len(responses) > 0 {
		attemptSuppressEmbed(s, m.ID, m.ChannelID)
		go verifyEmbedSuppression(s, m.ID, m.ChannelID)
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
			return
		}
	}
}
