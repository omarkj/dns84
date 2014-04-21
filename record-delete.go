package main

import (
	"fmt"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"log"
)

var cmdDeleteRecord = &Command {
	Run: runDeleteRecord,
	Name: "record-delete",
}

type DeleteResp struct {
	Ref bool
	Ok bool
	Empty bool
	Type string
}

var recordId string

func init() {
	cmdDeleteRecord.Flag.StringVarP(&zone, "zone", "z", "", "zone")
	cmdDeleteRecord.Flag.StringVarP(&recordId, "record", "r", "", "Record Id")
}

func runDeleteRecord(cmd *Command, args []string) {
	apiEndpoint := fmt.Sprintf("%s/domains/delentry/", apiURL)
	resp, _ := client.PostForm(apiEndpoint, url.Values{
		"entry": {recordId},
	})
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var f DeleteResp
	json.Unmarshal(body, &f)
	if f.Ok == true {
		fmt.Println("Deleted")
	} else {
		log.Fatalf("Could not delete")
	}
}
