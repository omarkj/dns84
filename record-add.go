package main

import (
	"strings"
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
	"net/url"
)

var cmdAddRecord = &Command {
	Run: runAddRecord,
	Name: "record-add",
}

type AddResp struct {
	WillLookLike string
	Newrow string
}

var (
	recordType string
	ttl string
	priority string
	host string
	data string
)

func init() {
	cmdAddRecord.Flag.StringVarP(&zone, "zone", "z", "", "Zone")
	cmdAddRecord.Flag.StringVarP(&recordType, "type", "t", "", "Record Type")
	cmdAddRecord.Flag.StringVarP(&ttl, "ttl", "e", "86400", "TTL")
	cmdAddRecord.Flag.StringVarP(&priority, "priority", "p", "10", "Priority")
	cmdAddRecord.Flag.StringVarP(&data, "data", "d", "", "Data")
	cmdAddRecord.Flag.StringVarP(&host, "host", "h", "", "Host")
}

func runAddRecord(cmd *Command, args []string) {
	apiEndpoint := fmt.Sprintf("%s/domains/newentry/", apiURL)
	resp, _ := client.PostForm(apiEndpoint, url.Values{
		"entry": {"new"},
		"zone": {zone},
		"type": {strings.ToUpper(recordType)},
		"host": {host},
		"ttl": {ttl},
		"priority": {priority},
		"rdata": {data},
	})
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var f AddResp
	json.Unmarshal(body, &f)
	if f.WillLookLike != "" {
		fmt.Println("Created")
	} else {
		log.Fatalf("Could not create")
	}
}
