package main

/**
 * first of all, i want to point out that I only use analytics to check how much compute power I need to run the bot.
 * and to check if the bot is running or not. I don't use it to track users or anything like that.
 *
 * To do that i only store the following data for each guild:
 * - the guild id
 * - the guild name
 * - the date when the bot was added to the server
 * - the number of each commands executed
 * - the date of the last command executed
 * - the date when the bot was removed from the server (1/1/2000 if it wasn't)
 */

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

type GuildData struct {
	ID                    string    `db:"ID"`
	Name                  string    `db:"name"`
	AddedAt               time.Time `db:"added_at"`
	CommandsUsageAPI      int       `db:"commands_usage_api"`
	CommandsUsageWeaviate int       `db:"commands_usage_weaviate"`
	LastCommand           time.Time `db:"last_command"`
	RemovedAt             time.Time `db:"removed_at"`
}

func Analytics_NewGuildData(dgGuild *discordgo.Guild) *GuildData {
	return &GuildData{
		ID:                    dgGuild.ID,
		Name:                  dgGuild.Name,
		AddedAt:               time.Now().UTC(),
		CommandsUsageAPI:      0,
		CommandsUsageWeaviate: 0,
		LastCommand:           time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		RemovedAt:             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

var pathAnalyticsDB = "data/analytics.db"

func InitAnalytics() error {
	//Check if the database file exists
	if _, err := os.Stat(pathAnalyticsDB); os.IsNotExist(err) {
		//if it doesn't exist, create it
		log.Println("Creating data/analytics.db")
		_, err := os.Create(pathAnalyticsDB)
		if err != nil {
			return err
		}
	}

	//Open the database and create the table if it doesn't exist
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS analytics
		( 
			ID         				VARCHAR, 
			name        			VARCHAR, 
			added_at    			DATE,
			commands_usage_api 		INT,
			commands_usage_weaviate INT,
			last_command 			DATE,
			removed_at 				DATE
		); 
	`)
	if err != nil {
		return err
	}
	return nil
}

func openDB() (*sql.DB, error) {
	return sql.Open("sqlite3", pathAnalyticsDB)
}

// Add a Guild row to the database
func Analytics_AddGuild(guild *GuildData) error {
	log.Printf("Guild '%v' as added the bot", guild.Name)
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	//Check if the guild is already in the database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM analytics WHERE ID = ?", guild.ID).Scan(&count)
	if err != nil {
		return err
	}

	//If the guild already exits in the db, reset the removed_at field
	if count > 0 {
		log.Printf("Guild '%v' already exists in the database, resetting removed_at field", guild.Name)
		_, err = db.Exec("UPDATE analytics SET removed_at = ? WHERE ID = ?", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), guild.ID)
		if err != nil {
			return err
		}
		return nil
	}

	//If the guild doesn't exist in the db, add it
	_, err = db.Exec("INSERT INTO analytics VALUES (?, ?, ?, ?, ?, ?, ?)", guild.ID, guild.Name, guild.AddedAt, guild.CommandsUsageAPI, guild.CommandsUsageWeaviate, guild.LastCommand, guild.RemovedAt)
	if err != nil {
		return err
	}
	return nil
}

func Analytics_RemoveGuild(guild *GuildData) error {
	log.Printf("Guild '%v' as removed the bot", guild.Name)
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	//Check if the guild is in the database (it should be this is just a security)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM analytics WHERE ID = ?", guild.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		//This isn't supposed to append if it does, add the guild to the database
		err := Analytics_AddGuild(guild)
		if err != nil {
			return err
		}
	}
	//Assuming that the guild is in the database, update the removed_at field
	_, err = db.Exec("UPDATE analytics SET removed_at = ? WHERE ID = ?", time.Now().UTC(), guild.ID)
	if err != nil {
		return err
	}
	return nil
}

func Analytics_UpdateCommandAPI(guildID string) error {
	return updateGuildFielf(guildID, "commands_usage_api")
}
func Analytics_UpdateCommandWeaviate(guildID string) error {
	return updateGuildFielf(guildID, "commands_usage_weaviate")
}

func updateGuildFielf(guildID string, field string) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	//Add 1 to the count and update the last_command field
	_, err = db.Exec("UPDATE analytics SET "+field+" = "+field+" + 1, last_command = ? WHERE ID = ?", time.Now().UTC(), guildID)
	if err != nil {
		return err
	}
	return nil
}
