package main

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	maxRetries = 3
	retryDelay = 1 * time.Second
)

func suppressEmbed(s *discordgo.Session, messageID, channelID string) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := attemptSuppressEmbed(s, messageID, channelID)
		if err == nil {
			verified, err := verifyEmbedSuppression(s, messageID, channelID)
			if err == nil && verified {
				log.Printf("Successfully suppressed embed (attempt %d)", attempt+1)
				return
			}
		}

		if attempt < maxRetries-1 {
			log.Printf("Embed suppression attempt %d failed, retrying...", attempt+1)
			time.Sleep(retryDelay)
		}
	}
	log.Printf("Failed to suppress embed after %d attempts", maxRetries)
}

func attemptSuppressEmbed(s *discordgo.Session, messageID, channelID string) error {
	_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      messageID,
		Channel: channelID,
		Flags:   discordgo.MessageFlagsSuppressEmbeds,
	})
	return err
}

func verifyEmbedSuppression(s *discordgo.Session, messageID, channelID string) (bool, error) {
	message, err := s.ChannelMessage(channelID, messageID)
	if err != nil {
		return false, err
	}
	// checks the length of the embeds slice in the message (i.e. checks if the message has any embeds)
	return len(message.Embeds) == 0, nil
}
