package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHtmlTable(tournament string) {
	encodedTournament := url.QueryEscape(tournament)
	url := fmt.Sprintf("https://lol.fandom.com/wiki/Special:RunQuery/PickBanHistory?PBH%%5Bpage%%5D=%s&PBH%%5Bteam%%5D=&PBH%%5Btextonly%%5D%%5Bis_checkbox%%5D=true&PBH%%5Btextonly%%5D%%5Bvalue%%5D=&_run=&pfRunQueryFormName=PickBanHistory&wpRunQuery=&pf_free_text=", encodedTournament)

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("cant even fetch, you suck: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("f leaguepedia %v", response.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatalf("...? %v", err)
	}

	csvTable, err := os.Create(tournament + ".csv")
	if err != nil {
		log.Fatalf("you couldnt even create a csv file wtf so bad: %v", err)
	}
	defer csvTable.Close()

	writer := csv.NewWriter(csvTable)
	defer writer.Flush()

	doc.Find("table").First().Each(func(_ int, table *goquery.Selection) {
		headers := []string{}
		table.Find("th").Each(func(_ int, th *goquery.Selection) {
			headers = append(headers, strings.TrimSpace(th.Text()))
		})
		writer.Write(headers)

		table.Find("tbody tr").Each(func(_ int, row *goquery.Selection) {
			fields := []string{}
			row.Find("td").Each(func(_ int, cell *goquery.Selection) {
				fields = append(fields, strings.TrimSpace(cell.Text()))
			})
			writer.Write(fields)
		})
	})

	fmt.Println("Gotta Catch 'Em All!")

}

func main() {
	tournaments := []string{
		"LPL Spring Season",
		"LEC Spring Season",
		"LCS Spring Season",
		"PCS Spring Season",
		"LCO Split 1",
		"VCS Spring Season",
		"LJL Spring Season",
		"LLA Opening Season",
	}

	for _, tournament := range tournaments {
		getHtmlTable(tournament)
	}

}
