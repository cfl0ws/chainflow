package stargazer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

// A MissesBlock containe start & end heights and conunt.
type MissesBlock struct {
	StartHeight string `json:"startHeight"`
	EndHeight   string `json:"endHeight"`
	Count       string `json:"count"`
}

// MissGroups contain Misses blocks
type MissGroups struct {
	MissGroups []MissesBlock `json:"missGroups"`
}

// GetMissedGroups function will return a metrics block on each request.
func GetMissedGroups(address string) ([]MissesBlock, error) {

	url := "https://sgapiv2.certus.one/v1/validator/" + address + "/groupedMisses"

	missed := MissGroups{}
	// Get stats by calling the api endpoint
	r, err := myClient.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// close connection
	defer r.Body.Close()

	// read the message body of the response
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Fatal(readErr)
		return nil, err
	}

	// unmarshal the json response
	json.Unmarshal(body, &missed)
	return missed.MissGroups, nil
}
