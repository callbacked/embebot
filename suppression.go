package main

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func attemptSuppressEmbed(s *discordgo.Session, messageID, channelID string) error {
	_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      messageID,
		Channel: channelID,
		Flags:   discordgo.MessageFlagsSuppressEmbeds,
	})
	return err
}

func verifyEmbedSuppression(s *discordgo.Session, messageID, channelID string) {
	time.Sleep(5 * time.Second) // sometimes embed slices do not load until after the fact, so it needs to verify
	message, err := s.ChannelMessage(channelID, messageID)
	if err != nil {
		log.Printf("Error fetching message: %v", err)
		return
	}
	// checks the length of the embeds slice in the message (i.e. checks if the message has any embeds)
	suppressed := len(message.Embeds) == 0
	flagsSuppressed := message.Flags&discordgo.MessageFlagsSuppressEmbeds != 0

	log.Printf("Verification: Embed count: %d, Suppressed by count: %v, Suppressed by flags: %v, Flags: %v",
		len(message.Embeds), suppressed, flagsSuppressed, message.Flags)
}
