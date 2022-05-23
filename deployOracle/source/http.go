package source

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Get(url string) map[string]string {
	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	addresses := make(map[string]string)
	jsonErr := json.Unmarshal(body, &addresses)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return addresses
}
