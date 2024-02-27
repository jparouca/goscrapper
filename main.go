package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHtmlTable(tournament string, writer *csv.Writer) {
	encodedTournament := url.QueryEscape(tournament)
	url := fmt.Sprintf("https://lol.fandom.com/wiki/Special:RunQuery/PickBanHistory?PBH%%5Bpage%%5D=%s&PBH%%5Bteam%%5D=&PBH%%5Btextonly%%5D%%5Bis_checkbox%%5D=true&PBH%%5Btextonly%%5D%%5Bvalue%%5D=&_run=&pfRunQueryFormName=PickBanHistory&wpRunQuery=&pf_free_text=", encodedTournament)

	fmt.Printf("url %v\n", url)

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
	fmt.Printf("Scraping %s\n", tournament)

	// Get the region (tournament name)
	region := tournament

	doc.Find(".wikitable").Each(func(_ int, table *goquery.Selection) {
		headers := []string{}

		// Skip the first th element
		table.Find("th").Each(func(i int, th *goquery.Selection) {
			if i > 0 {
				headers = append(headers, strings.TrimSpace(th.Text()))
			}
		})

		// Add region as the first header
		headers = append([]string{"Region"}, headers...)
		writer.Write(headers)
		fmt.Printf("Headers: %v\n", headers)

		table.Find("tbody tr").Each(func(_ int, row *goquery.Selection) {
			fields := []string{}
			fields = append(fields, region) // Add region as the first field

			row.Find("td").Each(func(_ int, cell *goquery.Selection) {
				fields = append(fields, strings.TrimSpace(cell.Text()))
			})
			writer.Write(fields)
		})
	})
}

func mergeCSVFiles() {
	mergedFile, err := os.Create("pickbans.csv")
	if err != nil {
		log.Fatalf("Failed to create merged file: %v", err)
	}
	defer mergedFile.Close()

	writer := csv.NewWriter(mergedFile)
	defer writer.Flush()

	header := []string{"Region", "Week", "Blue", "Red", "Score", "Winner", "Patch", "BB1", "RB1", "BB2", "RB2", "BB3", "RB3", "BP1", "RP1-2", "BP2-3", "RP3", "RB4", "BB4", "RB5", "BB5", "RP4", "BP4-5", "RP5", "BR1", "BR2", "BR3", "BR4", "BR5", "RR1", "RR2", "RR3", "RR4", "RR5", "SB"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Failed to write header row: %v", err)
	}

	files, err := filepath.Glob("Data/*.csv")
	if err != nil {
		log.Fatalf("Failed to list CSV files: %v", err)
	}

	for _, file := range files {
		if file == "pickbans.csv" {
			continue
		}

		csvFile, err := os.Open(file)
		if err != nil {
			log.Printf("Failed to open CSV file %s: %v", file, err)
			continue
		}
		defer csvFile.Close()

		reader := csv.NewReader(csvFile)

		// Read and print the headers
		headers, err := reader.Read()
		if err != nil {
			log.Printf("Failed to read header row from %s: %v", file, err)
			continue
		}
		fmt.Printf("CSV File: %s\n", file)
		fmt.Printf("Headers: %v\n", headers)
		fmt.Printf("Number of Fields: %d\n", len(headers))

		_, err = reader.Read()
		if err != nil && err != io.EOF {
			log.Printf("Failed to read header row from %s: %v", file, err)
			continue
		}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Failed to read record from %s: %v", file, err)
				break
			}
			fmt.Printf("Record: %v\n", record)
			fmt.Printf("Number of Fields: %d\n", len(record))
			// Append the tournament name to the record
			record = append([]string{file[:len(file)-len(".csv")]}, record...)
			if err := writer.Write(record); err != nil {
				log.Printf("Failed to write record to merged file: %v", err)
				break
			}
		}
	}

	fmt.Println("All CSV files merged into pickbans.csv")
}

func main() {
	tournaments := []string{
		"CBLOL 2024 Split 1",
		"LCK 2024 Spring",
		"LPL 2024 Spring",
		"LEC 2024 Winter",
		"LCS 2024 Spring",
		"PCS 2024 Spring",
		"LCO 2024 Split 1",
		"VCS 2024 Spring",
		"LJL 2024 Spring",
		"LLA 2024 Opening",
	}

	mergedFile, err := os.Create("pickbans2.csv")
	if err != nil {
		log.Fatalf("Failed to create merged file: %v", err)
	}
	defer mergedFile.Close()

	writer := csv.NewWriter(mergedFile)
	defer writer.Flush()

	header := []string{"Region", "Phase", "Blue", "Red", "Score", "Winner", "Patch", "BB1", "RB1", "BB2", "RB2", "BB3", "RB3", "BP1", "RP1-2", "BP2-3", "RP3", "RB4", "BB4", "RB5", "BB5", "RP4", "BP4-5", "RP5", "BR1", "BR2", "BR3", "BR4", "BR5", "RR1", "RR2", "RR3", "RR4", "RR5", "SB", "VOD"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Failed to write header row: %v", err)
	}

	for _, tournament := range tournaments {
		getHtmlTable(tournament, writer)
	}
	fmt.Println("All CSV files merged into pickbans2.csv")

	// mergeCSVFiles()
}
