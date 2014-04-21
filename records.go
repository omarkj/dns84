package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"regexp"
	"github.com/PuerkitoBio/goquery"
	"log"
)

type RecordJson struct {
	Ok bool
	Zone string
}

type Record struct {
	Id string
	Type string
	Priority string
	TTL string
	Host string
	Data string
	Target string
}

var cmdListRecords = &Command {
	Run: runListRecords,
	Name: "records",
}

var zone string

func init() {
	cmdListRecords.Flag.StringVarP(&zone, "zone", "z", "", "zone")
}

func runListRecords(cmd *Command, args []string) {
	apiEndpoint := fmt.Sprintf("%s/domains/ajax/6/?zone=%s", apiURL, zone)
	resp, _ := client.Get(apiEndpoint)
	
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var f RecordJson
	json.Unmarshal(body, &f)
	if f.Ok != true {
		log.Fatalf("Domain does not exist or does not belong to you")
	}
	records, _ := getRecords(zone, []byte(f.Zone))
	for record := range records {
		if records[record].Type == "MX" ||
			records[record].Type == "SRV" {
			fmt.Printf("%s\t%s\t%s\tIN\t%s\t%s %s\n", records[record].Id, records[record].Host,
				records[record].TTL, records[record].Type,
				records[record].Priority, records[record].Data)
		} else {
			fmt.Printf("%s\t%s\t%s\tIN\t%s\t%s\n", records[record].Id, records[record].Host,
				records[record].TTL, records[record].Type, records[record].Data)
		}
	}
}

func getRecords(zone string, body []byte) ([]Record, error) {
	bodyBuffer := bytes.NewBuffer(body)
	doc, _ := goquery.NewDocumentFromReader(bodyBuffer)
	records := []Record{}
	doc.Find("table tr").Each(
		func(i int, s *goquery.Selection) {
			recordTypeInput := s.Find("td input:first-child")
			recordTypeId, hasAttr := recordTypeInput.Attr("id")
			if hasAttr && recordTypeId != "host_new" {
				recordId := regexp.MustCompile("type_*").Split(recordTypeId, 2)[1]
				recordType, _ := recordTypeInput.Attr("value")
				record := Record{
					Id: recordId,
					Type: recordType,
				}
				host, hasHost := s.Find(fmt.Sprintf("input#host_%s", recordId)).Attr("value")
				if hasHost {
					if host == "@" {
						record.Host = zone
					} else {
						record.Host = fmt.Sprintf("%s.%s", host, zone)
					}
				}
				ttl, hasTtl := s.Find(fmt.Sprintf("select#ttl_%s option[selected=selected]", recordId)).Attr("value")
				if hasTtl {
					record.TTL = ttl
				}
				data, hasData := s.Find(fmt.Sprintf("input#rdata_%s", recordId)).Attr("value")
				if hasData {
					record.Data = data
				}
				if recordType == "MX" ||
					recordType == "SRV" {
					priority, hasPriority := s.Find(fmt.Sprintf("select#priority_%s option[selected]", recordId)).Attr("value")
					if hasPriority {
						record.Priority = priority
					}
				}
				if recordType == "SRV" {
					target, hasTarget := s.Find(fmt.Sprintf("input#target_%s", recordId)).Attr("value")
					if hasTarget {
						record.Data = target
					}
				}
				if recordType != "NS" {
					records = appendRecord(record, records)
				}
			}
		})
	return records, nil
}

func appendRecord(record Record, records []Record) ([]Record) {
	n := len(records)
	if n == cap(records) {
		newRecords := make([]Record, len(records), len(records)*2+1)
		copy(newRecords, records)
		records = newRecords
	}
	records = records[0 : n +1]
	records[n] = record
	return records
}
