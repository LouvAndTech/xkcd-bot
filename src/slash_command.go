package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

/* === slash command initialisation === */
var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "xkcb_daily",
			Description: "Get the xkcd of the day",
		},
		{
			Name:        "xkcd",
			Description: "Get a specific xkcd",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "number",
					Description: "The id of the xkcd",
					Required:    true,
				},
			},
		},
		{
			Name:        "xkcd_random",
			Description: "Get a random xkcd",
		},
		{
			Name:        "xkcd_search",
			Description: "Search for an xkcd by title",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "The title of the xkcd to search for",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"xkcb_daily": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			daily, err := fetchLastXKCD()
			if err != nil {
				log.Println(err)
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{formatXKCD(daily).MessageEmbed},
				},
			})
		},
		"xkcd": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			id := i.ApplicationCommandData().Options[0].IntValue()
			xkcd, err := fetchtXKCD(id)
			if err != nil {
				log.Println(err)
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{formatXKCD(xkcd).MessageEmbed},
				},
			})
		},
		"xkcd_random": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			rnd, err := fetchRandomXKCD()
			if err != nil {
				log.Println(err)
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{formatXKCD(rnd).MessageEmbed},
				},
			})
		},
		"xkcd_search": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			search := i.ApplicationCommandData().Options[0].StringValue()
			id, err := SearchXKCD(GetClient(cfg), search)
			if err != nil {
				log.Println(err)
				return
			}
			var xkcd XKCD
			if id != 0 {
				xkcd, err = fetchtXKCD(int64(id))
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				xkcd = XKCD{
					Num: 0,
				}
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{formatXKCD(xkcd).MessageEmbed},
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

/* === slash command format === */

func formatXKCD(data XKCD) *Embed {
	em := NewEmbed()
	if data.Num == 0 {
		em.AddTitle("No result found")
		return em
	}
	em.AddTitle(data.Title)
	em.AddImage(data.Image)
	em.AddDescription(
		"*" + data.Day + " / " + data.Month + " / " + data.Year + "*" + "\n ID : " + fmt.Sprint(data.Num))
	return em
}
