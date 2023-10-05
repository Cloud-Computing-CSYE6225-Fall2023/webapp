package account

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shivasaicharanruthala/webapp/model"
)

func generateBulkInsertQuery(table string, cols []string, numRows int) string {
	var values []string

	for i := 0; i < numRows; i++ {
		var placeholders []string
		for position, _ := range cols {
			placeholders = append(placeholders, fmt.Sprintf("$%d", (len(cols)*i)+position+1))
		}

		values = append(values, fmt.Sprintf("(%s)", strings.Join(placeholders, ",")))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ON CONFLICT (email) DO NOTHING;", table, strings.Join(cols, ","), strings.Join(values, ","))

	return query
}

func flattenBulkUserRecord(rows [][]string) []interface{} {
	var oneDSlice []interface{}

	for _, row := range rows {
		//Insert account id
		oneDSlice = append(oneDSlice, uuid.New().String())

		// Insert column values
		for i, element := range row {
			value := element

			if i == 3 { // password need to be hashed
				value, _ = model.HashPassword(value)
			}

			oneDSlice = append(oneDSlice, value)
		}

		// Insert account_created, account_updated
		currentTimestamp := time.Now().UTC()
		currentTimestampUtc := currentTimestamp.Format(time.RFC3339)

		oneDSlice = append(oneDSlice, currentTimestampUtc, currentTimestampUtc)
	}

	return oneDSlice
}
