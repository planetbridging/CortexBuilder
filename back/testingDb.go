package main

import (
	"database/sql"
	"fmt"
)

func createTestTableAndData(db *sql.DB) error {
	fmt.Println("Checking if test tbl needs to be created (its only 100 items in 7 cols)")

	tableName := "test_table"

	// Check if the table exists
	if !tableExists(db, tableName) {
		fmt.Println("Creating tbl")
		// Create the table
		createTableStmt := `CREATE TABLE test_table (
            id_auto_increment INT PRIMARY KEY AUTO_INCREMENT,
            input1 INT,
            input2 INT,
            input3 INT,
            output1 INT,
            output2 INT,
            output3 INT
        )`

		_, err := db.Exec(createTableStmt)
		if err != nil {
			return err
		}
		fmt.Println("Created table:", tableName)
		// Insert sample data (1 to 100)
		baseValue := 0
		insertStmt := `INSERT INTO test_table (input1, input2, input3, output1, output2, output3) VALUES (?, ?, ?, ?, ?, ?)`
		stmt, err := db.Prepare(insertStmt)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer stmt.Close()

		for i := 1; i <= 100; i++ {
			baseValue += 6 // Increment the base values for each row
			_, err = stmt.Exec(
				baseValue+1, baseValue+2, baseValue+3,
				baseValue+4, baseValue+5, baseValue+6,
			)

			if err != nil {
				return err
			}
		}
		fmt.Println("Inserted 100 rows into table:", tableName)
	} else {
		fmt.Println("Table", tableName, "already exists.")
	}
	return nil
}
