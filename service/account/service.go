package account

import (
	"encoding/csv"
	cr "errors"
	"github.com/shivasaicharanruthala/webapp/log"
	"github.com/shivasaicharanruthala/webapp/types"
	"os"

	"github.com/shivasaicharanruthala/webapp/errors"
	"github.com/shivasaicharanruthala/webapp/model"
	"github.com/shivasaicharanruthala/webapp/service"
	"github.com/shivasaicharanruthala/webapp/store"
)

type dataStore struct {
	accountStore store.Account
}

func New(accdb store.Account) service.Account {
	return &dataStore{accountStore: accdb}
}

func (a *dataStore) Insert(ctx *types.Context, account *model.Account) (*model.Account, error) {
	return nil, nil
}

func (a *dataStore) BulkInsert(ctx *types.Context, filepath string) error {
	cols := []string{}
	// Open the CSV file for reading
	csvFile, err := os.Open(filepath)
	if err != nil {
		return errors.NewCustomError(err)
	}
	defer func(csvFile *os.File) {
		err = csvFile.Close()
		if err != nil {
			lm := log.Message{Level: "ERROR", ErrorMessage: "Error closing the users.csv file"}
			ctx.Logger.Log(&lm)

			return
		}
	}(csvFile)

	// Create a CSV reader
	csvReader := csv.NewReader(csvFile)

	// Read 1st line of the csv file which gives the column names
	colNames, err := csvReader.Read()
	if err != nil {
		return errors.NewCustomError(err)
	}

	cols = append(cols, "id")
	cols = append(cols, colNames...)
	cols = append(cols, "account_created", "account_updated")

	// Loop through the CSV records and insert them into the database
	var rowsToInsert [][]string
	batchSize := 10
	for {
		record, err := csvReader.Read()
		if err != nil {
			break // End of file
		}

		// validate email, If not a valid email skip the row
		if er := model.ValidateEmail(record[2]); er != nil {
			continue
		}

		rowsToInsert = append(rowsToInsert, record)
		// If we have accumulated enough rows, insert them in bulk
		if len(rowsToInsert) == batchSize {
			// Perform the bulk insert (replace with your database-specific bulk insert query)
			er := a.accountStore.BulkInsert(ctx, cols, rowsToInsert)
			if er != nil {
				return er
			}

			// Clear the slice for the next batch
			rowsToInsert = nil
		}
	}

	// Insert any remaining rows (if less than batchSize)
	if len(rowsToInsert) > 0 {
		err = a.accountStore.BulkInsert(ctx, cols, rowsToInsert)
		if err != nil {
			return err
		}
	}

	lm := log.Message{Level: "INFO", ErrorMessage: "Bulk insert from users.csv completed successfully"}
	ctx.Logger.Log(&lm)

	return nil
}

func (a *dataStore) IsAccountExists(ctx *types.Context, email, pass string) (*model.User, error) {
	var user model.User

	if err := model.ValidateEmail(email); err != nil {
		return nil, errors.NewCustomError(errors.InvalidParam{Param: []string{"username"}}, 401)
	}

	account, err := a.accountStore.IsAccountExists(ctx, email)
	if err != nil {
		if err.Error() == "Username does not exists" {
			return nil, errors.NewCustomError(cr.New("Username does not exists"), 401)
		}

		return nil, err
	}

	err = model.VerifyPassword(*account.Password, pass)
	if err != nil {
		return nil, errors.NewCustomError(cr.New("Password Mismatch"), 401)
	}

	user.ID = account.ID
	user.Email = email

	return &user, nil
}
