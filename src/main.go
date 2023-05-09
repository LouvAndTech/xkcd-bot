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

// Bot parameters
var (
	BotToken       *string
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

/* === Global variables === */
var s *discordgo.Session
var cron *gocron.Scheduler

var cfg = weaviate.Config{
	// DEV CONFIG:
	//Host: "localhost:8080",
	Host:   "weaviate:8080",
	Scheme: "http",
}

func saveMissing() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	log.Println(">>> Saving missing XKCD for the day...")
	err := SaveMissingXkcd()
	if err != nil {
		panic(err)
	}
	err = UpdateDB(GetClient(cfg))
	if err != nil {
		panic(err)
	}
	log.Println(">>> Done!")
}

// init is called before main
func init() { flag.Parse() }

func init() {
	log.Println("Starting...")
	//Get the bot token from the environment variables
	if os.Getenv("TOKEN") == "" {
		log.Fatalf("You need to pass the bot token as an argument")
	}
	log.Println("Bot token:", os.Getenv("TOKEN"))
	BotToken = flag.String("token", os.Getenv("TOKEN"), "Bot access token")

	// Initialize the bot
	log.Println("Initializing bot...")
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	// Initialize the storage for persistent data
	log.Println("Initializing storage...")
	err = InitStrorage()
	if err != nil {
		log.Fatalf("Cannot initialize storage: %v", err)
	}

	// Initialize the database
	log.Println("Initializing database...")
	err = InitDB(cfg)
	if err != nil {
		log.Fatalf("Cannot initialize database: %v", err)
	}

	// Initialize the scheduler, add the job and start it
	log.Println("Initializing scheduler...")
	cron = gocron.NewScheduler(time.UTC)
	_, err = cron.Every(1).Day().At("00:00").Do(saveMissing)
	if err != nil {
		log.Fatal("Cannot add the job to the scheduler: ", err)
	}
	cron.StartAsync()

	//Initialize the analytics
	log.Println("Initializing analytics...")
	err = InitAnalytics()
	if err != nil {
		log.Fatalf("Cannot initialize analytics: %v", err)
	}
}

func main() {
	// Add the handelers
	Add_dsg_Handeler(s)

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// Register the commands globally
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
