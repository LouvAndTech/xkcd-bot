package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func InitStrorage() error {
	//check if the file exists
	if _, err := os.Stat("xkcd.json"); os.IsNotExist(err) {
		//if it doesn't exist, create it
		_, err := os.Create("xkcd.json")
		if err != nil {
			return err
		}
		//then write an empty array
		err = ioutil.WriteFile("xkcd.json", []byte("[]"), 0644)
		if err != nil {
			return err
		}
	}
	/* -- Then retrieve all the XKCDs -- */
	err := SaveMissingXkcd()
	if err != nil {
		return err
	}
	return nil
}

func SaveMissingXkcd() error {
	//get the last XKCD
	lastPublished, err := fetchLastXKCD()
	if err != nil {
		return err
	}
	//Get the data from the JSON
	file, err := ioutil.ReadFile("xkcd.json")
	if err != nil {
		return err
	}
	data := []XKCD_Short{}
	//Unmarshal the JSON
	err = json.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	// Compare the last id in the JSON and the last id published
	lastSaved := 0
	if len(data) > 0 {
		lastSaved = data[len(data)-1].Num
	}
	if lastSaved < lastPublished.Num {
		// Iterate from the last saved id to the last published id
		for i := lastSaved + 1; i <= lastPublished.Num; i++ {
			if i == 404 {
				continue
			}
			log.Printf("\rWIP : %d/%d", i, lastPublished.Num)
			//get the XKCD
			xkcd, err := fetchtXKCD(int64(i))
			if err != nil {
				return err
			}
			//convert it to a short XKCD
			xkcd_short := newShortXkcd(xkcd)
			//save it
			err = saveXKCD(xkcd_short)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func saveXKCD(new XKCD_Short) error {
	//Get the data from the JSON
	file, err := ioutil.ReadFile("xkcd.json")
	if err != nil {
		return err
	}
	data := []XKCD_Short{}
	// Unmarshal the JSON
	err = json.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	//Add the new XKCD
	data = append(data, new)

	//Save the data
	file, err = json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("xkcd.json", file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetAllXKCD() ([]XKCD_Short, error) {
	//Get the data from the JSON
	file, err := ioutil.ReadFile("xkcd.json")
	if err != nil {
		return nil, err
	}
	data := []XKCD_Short{}
	// Unmarshal the JSON
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
