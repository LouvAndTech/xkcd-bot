package main

import (
	"github.com/bwmarrin/discordgo"
)

//Class Embed

// VALUES
type Embed struct {
	*discordgo.MessageEmbed
}

// CONSTRUCTOR
// creacte an object
func NewEmbed() *Embed {
	return &Embed{&discordgo.MessageEmbed{}}
}

// METHODS
// Add Title
func (e *Embed) AddTitle(title string) *Embed {
	e.Title = title
	return e
}

// Add description
func (e *Embed) AddDescription(des string) *Embed {
	e.Description = des
	return e
}

// add image
func (e *Embed) AddImage(url string) *Embed {
	e.Image = &discordgo.MessageEmbedImage{
		URL: url,
	}
	return e
}

// Add a Thumbnail
func (e *Embed) AddThumbnail(url string) *Embed {
	e.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: url,
	}
	return e
}
