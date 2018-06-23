package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

// IPQuery requests to an API with an IP address, return with location
func IPQuery(api, ip string) (loc string) {
	url := api + ip
	res, err := http.Get(url)

	if err != nil {
		log.Errorf("Request Error: %+v\n", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Errorf("Response Error: %+v\n", err)
	}

	loc = strings.TrimSpace(string(body))

	return fmt.Sprintf("IP: %v\nLocation: %v\nURL: %v\n", ip, loc, url)
}
