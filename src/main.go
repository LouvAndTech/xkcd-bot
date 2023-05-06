package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

// Bot parameters TEMPLATE
/*var (
	GuildID        = flag.String("guild", "<GUILD_ID>", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "<ACESS_TOKEN>", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)*/

/* === Global variables === */
var s *discordgo.Session
var cron *gocron.Scheduler

var cfg = weaviate.Config{
	// Dv:
	Host: "localhost:8080",
	//Host:   "weaviate:8080",
	Scheme: "http",
}

func saveMissing() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	err := SaveMissingXkcd()
	if err != nil {
		panic(err)
	}
	err = UpdateDB(GetClient(cfg))
	if err != nil {
		panic(err)
	}
}

// init is called before main
func init() { flag.Parse() }

func init() {
	log.Println("Initializing bot...")
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	log.Println("Initializing storage...")
	err = InitStrorage()
	if err != nil {
		log.Fatalf("Cannot initialize storage: %v", err)
	}
	log.Println("Initializing database...")
	err = InitDB(cfg)
	if err != nil {
		log.Fatalf("Cannot initialize database: %v", err)
	}
	cron = gocron.NewScheduler(time.UTC)
	cron.Every(1).Day().At("00:00").Do(saveMissing)
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	/* Debug */
	// Add the 1000 xkcd to the database
	/*lastXKCD, err := fetchLastXKCD()
	if err != nil {
		log.Fatalf("Cannot get last XKCD: %v", err)
	}
	err = AddNewXKCDEntry(GetClient(cfg), lastXKCD)
	if err != nil {
		log.Fatalf("Cannot add XKCD: %v", err)
	}*/
	//Get the XKCD close to the number 1000
	/*xkcd, err := SearchXKCD(GetClient(cfg), "Overlapping")
	if err != nil {
		log.Fatalf("Cannot get XKCD: %v", err)
	}
	log.Printf("XKCD: %v", xkcd)*/

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
