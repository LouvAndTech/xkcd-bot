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

And to allow a search within all the meme i'm using [`weaviate`](https://weaviate.io) a vectorial database which is filled with all the description and title.

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
##### *(You can use the [json file in the repo](https://raw.githubusercontent.com/LouvAndTech/xkcd-bot/main/src/data/xkcd.json) to fill your own weaviate instance to avoid having to download the ~2500 first)*
- #### Download the source code and run it yourself :
    *You need to know what you're doing in go and a bit of docker*
- #### Use the docker image to run it :
    *You can use the [docker-compose file in the repo](https://github.com/LouvAndTech/xkcd-bot/blob/main/utils/docker-compose.yml) to deploy it.*
    - [Docker hub link](https://hub.docker.com/repository/docker/louvandtech/xkcd-bot)

---
## Limits / Bugs :
*I've resolved all the limitation and bugs i've found, but if you find one feel free to open an issue*


---

## Analytics :
The bot is collecting some data to allow me to know how many people are using it and how many server are using it.

*To do that i only store the following data for each guild:*
- *the guild id*
- *the guild name*
- *the date when the bot was added to the server*
- *the number of each commands executed*
- *the date of the last command executed*
- *the date when the bot was removed from the server*