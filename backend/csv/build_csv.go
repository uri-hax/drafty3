package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

// model of suggestions table rows we'll be looking at
type SuggestionRow struct {
	IDUniqueID       int
	IDSuggestionType int
	Suggestion       string
	Confidence       int
	Active           int
}

// model for building the csprofs csv
type CSProfRecord struct {
	IDUniqueID int
	FullName   string
	University string
	JoinYear   string
	SubField   string
	Bachelors  string
	Doctorate  string
}

// main function to read flags and error accordingly if issues and call run function for logic
func main() {
	// get flags and parse them
	dbPath := flag.String("db", "", "Path to SQLite database file")
	outPath := flag.String("out", "", "Path to output CSV file")
	csvType := flag.String("csv_type", "", "Type of CSV to generate")
	flag.Parse()

	// make sure required flags are provided
	if *dbPath == "" {
		log.Fatal("missing required --db flag")
	}
	if *outPath == "" {
		log.Fatal("missing required --out flag")
	}
	if *csvType == "" {
		log.Fatal("missing required --csv_type flag")
	}

	// call the run function for the logic
	if err := run(*dbPath, *outPath, *csvType); err != nil {
		log.Fatalf("build_csv failed: %v", err)
	}
}

// run function to open the db and call the appropriate csv builder based on flags
func run(dbPath, outPath, csvType string) error {
	// open the db and eventually close it
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	// ping it to ensure we're connected
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	// call the appropriate csv builder based on csvType flag
	switch csvType {
	case "csprofs":
		return buildCSProfsCSV(db, outPath)
	default:
		return fmt.Errorf("unsupported csv_type: %s", csvType)
	}
}

// buildCSProfsCSV builds a csv file for the csprofs dataset
func buildCSProfsCSV(db *sql.DB, outPath string) error {
	// set up the query to get active suggestions of the relevant types
	query := `
		SELECT
			idUniqueID,
			idSuggestionType,
			suggestion,
			confidence,
			active
		FROM Suggestions
		WHERE active = 1
		  AND idSuggestionType IN (1, 2, 3, 5, 7, 9)
	`

	// get the rows after the query
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("query Suggestions: %w", err)
	}
	defer rows.Close()

	// set up a map to track the best suggestion for each (idUniqueID, idSuggestionType) pair
	best := make(map[string]SuggestionRow)

	// go through the rows and populate the best map
	for rows.Next() {
		var r SuggestionRow
		if err := rows.Scan(
			&r.IDUniqueID,
			&r.IDSuggestionType,
			&r.Suggestion,
			&r.Confidence,
			&r.Active,
		); err != nil {
			return fmt.Errorf("scan row: %w", err)
		}

		// call function to make the string key for the map
		key := makeKey(r.IDUniqueID, r.IDSuggestionType)
		existing, found := best[key]
		if !found || r.Confidence > existing.Confidence {
			best[key] = r
		}
	}

	// check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate rows: %w", err)
	}

	// map to hold the final records keyed by idUniqueID
	recordMap := make(map[int]*CSProfRecord)

	// go through the best suggestions and populate the recordMap
	for _, row := range best {
		// get or create the record for this idUniqueID
		rec, found := recordMap[row.IDUniqueID]
		if !found {
			rec = &CSProfRecord{IDUniqueID: row.IDUniqueID}
			recordMap[row.IDUniqueID] = rec
		}

		// normalize the suggestion value across string and string[]
		value := normalizeSuggestion(row.Suggestion)

		// populate the appropriate field based on idSuggestionType
		switch row.IDSuggestionType {
		case 1:
			rec.FullName = value
		case 2:
			rec.University = value
		case 3:
			rec.Bachelors = value
		case 5:
			rec.Doctorate = value
		case 7:
			rec.JoinYear = value
		case 9:
			rec.SubField = value
		}
	}

	// create a sorted list of idUniqueIDs for consistent output order
	ids := make([]int, 0, len(recordMap))
	for id := range recordMap {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	// create the output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	// create the output csv file
	file, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output csv: %w", err)
	}
	defer file.Close()

	// set up the csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// set the header row
	header := []string{
		"idUniqueID",
		"FullName",
		"University",
		"JoinYear",
		"SubField",
		"Bachelors",
		"Doctorate",
	}

	// write the header row
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	// go through the sorted ids and write the corresponding records to the csv
	for _, id := range ids {
		rec := recordMap[id]
		row := []string{
			strconv.Itoa(rec.IDUniqueID),
			rec.FullName,
			rec.University,
			rec.JoinYear,
			rec.SubField,
			rec.Bachelors,
			rec.Doctorate,
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("write row for idUniqueID=%d: %w", id, err)
		}
	}

	// flush the writer and check for errors
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush csv writer: %w", err)
	}

	// log success and return
	log.Printf("Wrote csprofs CSV to %s with %d rows", outPath, len(ids))
	return nil
}

// makeKey creates a string key for the map based on idUniqueID and idSuggestionType
func makeKey(idUniqueID, idSuggestionType int) string {
	return fmt.Sprintf("%d:%d", idUniqueID, idSuggestionType)
}

// normalizeSuggestion takes a raw suggestion string and normalizes it by trimming whitespace and handling JSON array strings and quoted strings
func normalizeSuggestion(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}

	// if it's a JSON array string keep it that way but trim whitespace from each element in the array
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		var arr []string
		// try to unmarshal it as a JSON array of strings
		if err := json.Unmarshal([]byte(s), &arr); err == nil {
			// trim whitespace from each element in the array
			for i := range arr {
				arr[i] = strings.TrimSpace(arr[i])
			}

			// marshal it back to a JSON string to ensure consistent formatting 
			normalized, err := json.Marshal(arr)
			if err == nil {
				return string(normalized)
			}
		}
	}

	// if it's a quoted JSON string then unwrap it
	var single string
	// try to unmarshal it as a JSON string
	if err := json.Unmarshal([]byte(s), &single); err == nil {
		return strings.TrimSpace(single)
	}

	// otherwise just return the trimmed raw string
	return s
}