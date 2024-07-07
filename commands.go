package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"enable":   handleEnableCommand,
	"disable":  handleDisableCommand,
	"settings": handleSettingsCommand,
}

func handleEnableCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	handleEmbedCommand(s, i, true)
}

func handleDisableCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	handleEmbedCommand(s, i, false)
}

func handleEmbedCommand(s *discordgo.Session, i *discordgo.InteractionCreate, enable bool) {
	options := i.ApplicationCommandData().Options
	service := options[0].StringValue()

	validServices := map[string]bool{
		"all": true, "reddit": true, "tiktok": true, "instagram": true, "twitter": true, "x": true,
	}

	if !validServices[service] {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Invalid Service",
						Color:       0xff0000,
						Description: "The service you requested is not valid.",
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Valid Options",
								Value: "all, reddit, tiktok, instagram, twitter",
							},
						},
					},
				},
			},
		})
		return
	}

	action := "enabled"
	if !enable {
		action = "disabled"
	}

	if service == "all" {
		for srv := range validServices {
			if srv != "all" {
				setServerSettingToDB(i.GuildID, srv, enable)
			}
		}
	} else {
		setServerSettingToDB(i.GuildID, service, enable)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Description: fmt.Sprintf("Successfully %s %s embed service(s) for this server.", action, service),
					Color:       0x82ff8c,
				},
			},
		},
	})
}

func handleSettingsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildID := i.GuildID

	serviceStatuses := make(map[string]bool)
	for _, item := range matchConfig {
		if getServerSettingFromDB(guildID, strings.ReplaceAll(item.BaseLink, ".com", "")) {
			serviceStatuses[strings.ReplaceAll(item.BaseLink, ".com", "")] = true
		} else {
			serviceStatuses[strings.ReplaceAll(item.BaseLink, ".com", "")] = false
		}
	}

	services := []string{"reddit", "tiktok", "instagram", "twitter", "x"}
	var enabledServices, disabledServices []string

	for _, service := range services {
		if serviceStatuses[service] {
			enabledServices = append(enabledServices, service)
		} else {
			disabledServices = append(disabledServices, service)
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Embed Settings for this Server",
					Color: 0x82ff8c,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Enabled Services",
							Value: strings.Join(enabledServices, ", "),
						},
						{
							Name:  "Disabled Services",
							Value: strings.Join(disabledServices, ", "),
						},
					},
				},
			},
		},
	})
}

func handleServiceAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate, option *discordgo.ApplicationCommandInteractionDataOption) {
	choices := []*discordgo.ApplicationCommandOptionChoice{
		{Name: "All", Value: "all"},
		{Name: "Reddit", Value: "reddit"},
		{Name: "TikTok", Value: "tiktok"},
		{Name: "Instagram", Value: "instagram"},
		{Name: "Twitter", Value: "twitter"},
		{Name: "X", Value: "x"},
	}

	if option.StringValue() != "" {
		filtered := make([]*discordgo.ApplicationCommandOptionChoice, 0)
		for _, choice := range choices {
			if strings.HasPrefix(strings.ToLower((choice.Name)), strings.ToLower(option.StringValue())) {
				filtered = append(filtered, choice)
			}
		}
		choices = filtered
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}
