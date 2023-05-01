package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
)

type XKCD struct {
	//date
	Year  string `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
	//info
	Num        int    `json:"num"`
	Title      string `json:"title"`
	Safe_title string `json:"safe_title"`
	//data
	Image      string `json:"img"`
	Alt        string `json:"alt"`
	Transcript string `json:"transcript"`
	//other
	News string `json:"news"`
}

func fetchtXKCD(nb int64) (XKCD, error) {
	num := ""
	if nb == 0 {
		num = ""
	} else {
		num = fmt.Sprint(nb) + "/"
	}

	response, err := http.Get("https://xkcd.com/" + num + "info.0.json")

	if err != nil {
		log.Print(err.Error())
		return XKCD{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err.Error())
		return XKCD{}, err
	}
	output := XKCD{}
	err = json.Unmarshal([]byte(responseData), &output)
	if err != nil {
		log.Print(err.Error())
		return XKCD{}, err
	}
	output = verifyXKCD(output)
	return output, nil
}

func fetchLastXKCD() (XKCD, error) {
	return fetchtXKCD(0)
}

func fetchRandomXKCD() (XKCD, error) {
	//Fetch Last XKCD to get the last id
	last, err := fetchLastXKCD()
	if err != nil {
		log.Print(err.Error())
		return XKCD{}, err
	}
	//Now using the max id we can get a random id (between 1 and max-1 to exclude today's)
	randomId := rand.Intn(last.Num-2) + 1
	//Fetch the random XKCD
	random, err := fetchtXKCD(int64(randomId))
	if err != nil {
		log.Print(err.Error())
		return XKCD{}, err
	}
	//Return it
	return random, nil
}

func verifyXKCD(xkcd XKCD) XKCD {
	// This is to remove all the non-ascii char
	re := regexp.MustCompile("[[:^ascii:]]")
	xkcd.Title = re.ReplaceAllLiteralString(xkcd.Title, "")
	xkcd.Alt = re.ReplaceAllLiteralString(xkcd.Alt, "")
	xkcd.Transcript = re.ReplaceAllLiteralString(xkcd.Transcript, "")

	return xkcd
}
