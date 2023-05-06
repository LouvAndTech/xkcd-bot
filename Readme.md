# xkcd - bot
## Intro : 

### - What ?

This is a discord bot that you can use to : 
- Get the last xkcd meme released
- Get a spécific one
- Get a Random one

 *and the best part :*
- Get a xkcd meme corresponding to a search querry in text

### - How ? :

Under the hood it's using a `Golang` program and using [`discordgo`](https://github.com/bwmarrin/discordgo) librairy to interact with discord.

To get the data i'm using [xkcd api](https://xkcd.com/json.html).

And to allow a search within all the meme i'm using [`weaviate`](https://weaviate.io) a vectorial databse wich is filled with all the description and title.

All of this hosted on my `Docker Instance`.

### - Why ? :

No real answer, just think we all could use a lilte bit more of xkcd in our boring life.

---
## How to use it : 
### User side :
- #### Invite the bot on your server :
    [Invite link](https://discord.com/api/oauth2/authorize?client_id=1102198415439429693&permissions=0&scope=bot)
    *You need to be admin of the server* 
- #### Try the bot on it's support server :
    [Join link](https://discord.gg/jyPPTFXs)

### Developer side :

*i won't explain how to do it but feel free to use or modify my code if you want*
- #### Download the source code and run it yourself :
    *You need to know what you're doing in go and a bit of docker*
- #### **(WIP)**  ~~Download the docker image and run it :~~
    *It's not possible for now because the bot token need to be hardcoded* [ ~~Docker hub link~~](https://hub.docker.com/repository/docker/louvandtech/xkcd-bot)

---
## Limit :
- All new meme are not added to the db after first launch of the container
(plain to use gocron to do that)