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
	// sometimes embed slices do not load until after the fact, so it needs to verify
	const totalDuration = 6 * time.Second
	const checkInterval = 5 * time.Second

	startTime := time.Now()
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	timeout := time.After(totalDuration)

	for {
		select {
		case <-ticker.C:
			message, err := s.ChannelMessage(channelID, messageID)
			if err != nil {
				log.Printf("Error fetching message: %v", err)
				continue
			}

			embedCount := len(message.Embeds)
			flagsSuppressed := message.Flags&discordgo.MessageFlagsSuppressEmbeds != 0

			log.Printf("Verification at %v: Embed count: %d, Suppressed by flags: %v, Flags: %v",
				time.Since(startTime), embedCount, flagsSuppressed, message.Flags)

			if embedCount > 0 || !flagsSuppressed {
				log.Printf("Embed detected or flags not set at %v. Attempting to suppress.", time.Since(startTime))
				if err := attemptSuppressEmbed(s, messageID, channelID); err != nil {
					log.Printf("Suppression attempt failed: %v", err)
				} else {
					log.Println("Suppression attempt sent successfully")
				}
			} else {
				log.Printf("No embeds detected and suppression flags set at %v.", time.Since(startTime))
			}

		case <-timeout:
			finalMessage, err := s.ChannelMessage(channelID, messageID)
			if err != nil {
				log.Printf("Error fetching message for final check: %v", err)
				return
			}

			if len(finalMessage.Embeds) == 0 && finalMessage.Flags&discordgo.MessageFlagsSuppressEmbeds != 0 {
				log.Println("Embeds successfully suppressed after final check")
				attemptSuppressEmbed(s, messageID, channelID)
			} else {
				log.Printf("Final check: Embed count: %d, Suppressed by flags: %v, Flags: %v",
					len(finalMessage.Embeds), finalMessage.Flags&discordgo.MessageFlagsSuppressEmbeds != 0, finalMessage.Flags)
				log.Println("Failed to fully suppress embeds after extended monitoring")
			}
			return
		}
	}
}
