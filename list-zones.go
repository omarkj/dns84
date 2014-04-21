package main

import (
	"fmt"
	"bytes"
	"sync"
	"encoding/json"
	"io/ioutil"
	"github.com/PuerkitoBio/goquery"
)

var cmdListZones = &Command {
	Run: runListZones,
	Name: "list-zones",
}

type ZoneStatus struct {
	Zone string
	CheckOk bool
	Message string
	Ok bool
}

func runListZones(cmd *Command, args []string) {
	apiEndpoint := fmt.Sprintf("%s/domains/", apiURL)
	resp, _ := client.Get(apiEndpoint)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	domains, _ := getZones(body)
	zoneInfoChannel := make(chan ZoneStatus, len(domains))
	var wg sync.WaitGroup
	wg.Add(len(domains))
	for domainIdx := range domains {
		go func(domain string, respChannel chan ZoneStatus, wg *sync.WaitGroup) {
			respChannel <- getZoneInfo(domain)
			wg.Done()
		}(domains[domainIdx], zoneInfoChannel, &wg)
	}
	wg.Wait()
	for i := 0; i < len(domains); i++ {
		msg := <- zoneInfoChannel
		if msg.CheckOk == true {
			fmt.Printf("%s\t%s\n", msg.Zone, "ok")
		} else {
			fmt.Printf("%s\t%s\n", msg.Zone, "not ok")
		}
	}
}

func getZones(body []byte) ([]string, error) {
	bodyBuffer := bytes.NewBuffer(body)
	doc, _ := goquery.NewDocumentFromReader(bodyBuffer)
	domains := []string{}
	doc.Find("table.table tr td:first-child").Each(
		func(i int, s *goquery.Selection) {
			domains = appendDomain(s.Text(), domains)
		})
	return domains, nil
}

func getZoneInfo(domain string) (ZoneStatus) {
	apiEndpoint := fmt.Sprintf("%s/domains/ZoneStatus/%s", apiURL, domain)
	resp, _ := client.Get(apiEndpoint)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var zoneStatus ZoneStatus
	json.Unmarshal(body, &zoneStatus)
	zoneStatus.Zone = domain
	return zoneStatus
}

func appendDomain(domain string, domains []string) ([]string) {
	n := len(domains)
	if n == cap(domains) {
		newDomains := make([]string, len(domains), len(domains)*2+1)
		copy(newDomains, domains)
		domains = newDomains
	}
	domains = domains[0 : n +1]
	domains[n] = domain
	return domains
}
