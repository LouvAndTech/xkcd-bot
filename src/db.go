package main

import (
	"context"
	"fmt"
	"log"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

type XKCD_Short struct {
	Num        int    `json:"num"`
	Title      string `json:"title"`
	Alt        string `json:"alt"`
	Transcript string `json:"transcript"`
}

func newShortXkcd(xkcd XKCD) XKCD_Short {
	return XKCD_Short{
		Num:        xkcd.Num,
		Title:      xkcd.Title,
		Alt:        xkcd.Alt,
		Transcript: xkcd.Transcript,
	}
}

func InitDB(cfg weaviate.Config) error {
	// Create a client$
	client := GetClient(cfg)
	// Get the schema
	schema := GetSchema(client)
	if len(schema) == 0 {
		new_schema := createdSchema()
		// Create the schema
		err := client.Schema().ClassCreator().WithClass(new_schema).Do(context.Background())
		if err != nil {
			return fmt.Errorf("cannot create schema: %v", err)
		}

		// Add the XKCD entries
		all, err := GetAllXKCD()
		if err != nil {
			return fmt.Errorf("cannot get all XKCD: %v", err)
		}
		for _, entry := range all {
			err = AddNewXKCDEntry(client, entry)
			if err != nil {
				log.Printf("cannot add XKCD entry: %v", err)
				continue
			}
		}

	}
	return nil
}

func GetClient(cfg weaviate.Config) *weaviate.Client {
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	return client
}

func GetSchema(client *weaviate.Client) []*models.Class {
	schema, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		panic(err)
	}
	return schema.Classes
}

func AddNewXKCDEntry(client *weaviate.Client, entry XKCD_Short) error {
	obj := &models.Object{
		Class: "XKCD",
		Properties: map[string]any{
			"num":        entry.Num,
			"title":      entry.Title,
			"alt":        entry.Alt,
			"transcript": entry.Transcript,
		},
	}
	_, err := client.Batch().ObjectsBatcher().WithObjects(obj).Do(context.Background())
	//_, err := client.Data().Creator().WithClassName("XKCD").WithProperties(obj).Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot create object: %v", err)
	}
	return nil
}

func SearchXKCD(client *weaviate.Client, searchText string) (int, error) {
	title := graphql.Field{Name: "title"}
	num := graphql.Field{Name: "num"}
	_additional := graphql.Field{
		Name: "_additional", Fields: []graphql.Field{
			{Name: "certainty"},
			{Name: "distance"},
		},
	}

	concepts := []string{searchText}
	distance := float32(0.99)
	nearText := client.GraphQL().NearTextArgBuilder().
		WithConcepts(concepts).
		WithDistance(distance)

	res, err := client.GraphQL().Get().
		WithClassName("XKCD").
		WithFields(title, num, _additional).
		WithNearText(nearText).
		WithLimit(1).
		Do(context.Background())
	if err != nil {
		return 0, fmt.Errorf("cannot search objects: %v", err)
	}
	get := res.Data["Get"].(map[string]interface{})["XKCD"].([]interface{})
	log.Printf("Looking for : %v | %v\n", searchText, get)
	if len(get) == 0 {
		return 0, nil
	} else {
		num := int(get[0].(map[string]interface{})["num"].(float64))
		return num, nil
	}
}

// Func to return the schema in order to cleanup the code
func createdSchema() *models.Class {
	return &models.Class{
		Class:       "XKCD",
		Description: "Infos on an xkcd comic",
		Properties: []*models.Property{
			{
				DataType:    []string{"int"},
				Description: "The number of the comic",
				Name:        "num",
			},
			{
				DataType:    []string{"string"},
				Description: "The title of the comic",
				Name:        "title",
			},
			{
				DataType:    []string{"string"},
				Description: "The alt text of the comic",
				Name:        "alt",
			},
			{
				DataType:    []string{"string"},
				Description: "The transcript of the comic",
				Name:        "transcript",
			},
		},
	}
}
