# xkcd - bot

## What ?

This is a discord bot that you can use to : 
- Get the last xkcd meme released
- Get a spécific one
- Get a Random one

 *and the best part :*
- Get a xkcd meme corresponding to a search querry in text

## How ? :

Under the hood it's using a `Golang` program and using [`discordgo`](https://github.com/bwmarrin/discordgo) librairy to interact with discord.

To get the data i'm using [xkcd api](https://xkcd.com/json.html).

And to allow a search within all the meme i'm using [`weaviate`](https://weaviate.io) a vectorial databse wich is filled with all the description and title.

All of this hosted on my `Docker Instance`.

## Why ? :

No real answer, just think we all could use a lilte bit more of xkcd in our boring life.

### Limit :
- For now it's only running on a private server but i plain to push it on a public one
- All new meme are not added to the db after first launch of the container
(plain to use gocron to do that)