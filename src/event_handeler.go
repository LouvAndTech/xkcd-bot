package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func Add_dsg_Handeler(s *discordgo.Session) {
	// Add Handeler for the ready event
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	// Add Handeler for the Guild Create event
	s.AddHandler(func(s *discordgo.Session, r *discordgo.GuildCreate) {
		log.Printf("Joined guild: %v", r.Name)
		err := Analytics_AddGuild(Analytics_NewGuildData(r.Guild))
		if err != nil {
			log.Printf("Cannot add guild to analytics: %v", err)
		}
	})
	// Add Handeler for the Guild Delete event
	s.AddHandler(func(s *discordgo.Session, r *discordgo.GuildDelete) {
		log.Printf("Left guild: %v", r.Name)
		err := Analytics_RemoveGuild(Analytics_NewGuildData(r.Guild))
		if err != nil {
			log.Printf("Cannot remove guild from analytics: %v", err)
		}
	})
}
