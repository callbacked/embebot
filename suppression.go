package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func suppressEmbed(s *discordgo.Session, messageID, channelID string) {
	_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      messageID,
		Channel: channelID,
		Flags:   discordgo.MessageFlagsSuppressEmbeds,
	})

	if err != nil {
		log.Printf("Failed to suppress embed: %v", err)
	} else {
		log.Println("Successfully suppressed embed")
	}

}
