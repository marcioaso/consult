package db

import (
	"fmt"
	"log"
	"strings"

	"github.com/marcioaso/consult/app/model"
)

func (d *DB) GetKLineForBacktest() []model.BybitResponse {
	defer d.Close()

	rows, err := d.Instance.Query("SELECT id, t, v, o, c, h, l, s, sn FROM bybit_responses order by t asc")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	// Slice to store results
	var responses []model.BybitResponse

	// Iterate over rows
	for rows.Next() {
		var response model.BybitResponse
		err := rows.Scan(&response.Id, &response.T, &response.V, &response.O, &response.C, &response.H, &response.L, &response.S, &response.SN)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		responses = append(responses, response)
	}

	// Check for errors from iteration
	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}
	return responses
}

func (d *DB) InsertKLines(data []model.BybitResponse) (bool, error) {
	defer d.Close()

	// Number of records to insert
	numRecords := len(data)

	// Prepare the SQL query for bulk insert
	insertQuery := `INSERT INTO bybit_responses (t, v, o, c, h, l, s, sn) VALUES `

	// Generate fake data and build the query
	values := make([]interface{}, 0, numRecords*8) // Reserve space for values
	queryParts := make([]string, 0, numRecords)    // Store query placeholders

	for _, record := range data {
		queryParts = append(queryParts, "(?, ?, ?, ?, ?, ?, ?, ?)")
		values = append(values, record.T, record.V, record.O, record.C, record.H, record.L, record.S, record.SN)
	}

	// Combine the query
	insertQuery += fmt.Sprintf("%s ON DUPLICATE KEY UPDATE id=id;", strings.Join(queryParts, ", "))

	// Execute the query
	_, err := d.Instance.Exec(insertQuery, values...)
	if err != nil {
		log.Fatalf("Error inserting bulk data: %v", err)
		return false, err
	}

	fmt.Printf("Successfully inserted %d records into the table.\n", numRecords)
	return true, nil
}
