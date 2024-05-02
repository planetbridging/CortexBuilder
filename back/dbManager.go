package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Username string
	Password string
	Hostname string
	DBName   string
}

type DataTemp1 struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type DataTemp2 struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type DataTemp2Db struct {
	Id     int    `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type DataTemp3 struct {
	Category string `json:"Category"`
}

type DataTemp3Rename struct {
	Category string `json:"name"`
}

type User struct {
	ID                 string
	Email              string
	VerifiedEmail      bool
	Name               string
	GivenName          string
	FamilyName         string
	Picture            string
	Locale             string
	stripe_customer_id string
}

func ConnectToDB(config DBConfig) (*sql.DB, error) {
	// Register a custom TLS config that skips server's certificate verification.
	mysql.RegisterTLSConfig("custom", &tls.Config{
		InsecureSkipVerify: true,
	})

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?parseTime=true&tls=custom", config.Username, config.Password, config.Hostname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectToDBSmall(config DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", config.Username, config.Password, config.Hostname, config.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ListDatabases(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}

	return databases, nil
}

func CreateDatabase(db *sql.DB, dbName string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	return err
}

func EnsureDatabaseExists(db *sql.DB, dbName string) error {
	// Check if the database exists
	rows, err := db.Query(fmt.Sprintf("SHOW DATABASES LIKE '%s'", dbName))
	if err != nil {
		return err
	}
	defer rows.Close()

	// If the database exists, rows.Next() will return true
	if rows.Next() {
		return nil
	}

	// If the database does not exist, create it
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	return err
}

func isTableNameSafe(tableName string) bool {
	for _, ch := range tableName {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return false
		}
	}
	return true
}

func insertDataToNodes(db *sql.DB, tableName string, data []DataTemp1) error {
	// Check if the table exists
	rows, err := db.Query(fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName))
	if err != nil {
		return err
	}
	defer rows.Close()

	// If the table does not exist, create it
	if !rows.Next() {
		_, err = db.Exec(fmt.Sprintf(`CREATE TABLE %s (
			id INT PRIMARY KEY,
			name VARCHAR(255),
			category VARCHAR(255)
		)`, tableName))
		if err != nil {
			return err
		}

		// Insert each data item into the table
		for _, item := range data {
			_, err = db.Exec(fmt.Sprintf(`INSERT INTO %s (id, name, category) VALUES (?, ?, ?)`, tableName), item.ID, item.Name, item.Category)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func insertDataNodeConnections(db *sql.DB, data []DataTemp2, tblName string) error {
	tableExists := tableExists(db, tblName)
	fmt.Println(tableExists)
	if !tableExists {
		createTableStmt := `CREATE TABLE IF NOT EXISTS ` + tblName + ` (
			id INT AUTO_INCREMENT,
			source VARCHAR(255),
			target VARCHAR(255),
			PRIMARY KEY (id)
		);`

		_, err := db.Exec(createTableStmt)
		if err != nil {
			return err
		}

		createIndexStmt := `CREATE UNIQUE INDEX idx_source_target ON ` + tblName + ` (source, target);`
		_, err = db.Exec(createIndexStmt)
		if err != nil {
			return err
		}

		stmt, err := db.Prepare("INSERT INTO " + tblName + "(source, target) VALUES(?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, item := range data {
			_, err := stmt.Exec(item.Source, item.Target)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func tableExists(db *sql.DB, tblName string) bool {
	query := "SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	rows, err := db.Query(query, tblName)
	if err != nil {
		return false
	}
	defer rows.Close()

	return rows.Next()
}

func createView(db *sql.DB, viewStatement string, viewStatementName string, db_name string) error {
	// Check if view exists
	var exists bool
	err := db.QueryRow("SELECT COUNT(*) FROM information_schema.views WHERE table_schema = ? AND table_name = ?", db_name, viewStatementName).Scan(&exists)
	if err != nil {
		return err
	}

	// If view does not exist, create it
	if !exists {
		_, err := db.Exec(viewStatement)
		if err != nil {
			return err
		}
		fmt.Println("View " + viewStatementName + " created successfully.")
	} else {
		fmt.Println("View " + viewStatementName + " already exists.")
	}

	return nil
}

func selectAllToJsonOld(db *sql.DB, tableName string, model interface{}) ([]byte, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(sql.RawBytes)
	}

	var results []interface{}
	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		elem := reflect.New(reflect.TypeOf(model).Elem()).Interface()
		for i, column := range columns {
			field := reflect.ValueOf(elem).Elem().FieldByName(column)
			if field.IsValid() && field.Kind() == reflect.String {
				field.SetString(string(*values[i].(*sql.RawBytes)))
			}
		}

		results = append(results, elem)
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func selectAllToJson(db *sql.DB, tableName string, model interface{}) ([]byte, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(sql.RawBytes)
	}

	var results []interface{}

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		elem := reflect.New(reflect.TypeOf(model).Elem()).Interface()
		for i, column := range columns {
			field := reflect.ValueOf(elem).Elem().FieldByName(column)
			if field.Kind() == reflect.String {
				field.SetString(string(*values[i].(*sql.RawBytes)))
			} else if field.Kind() == reflect.Int {
				intVal, err := strconv.Atoi(string(*values[i].(*sql.RawBytes)))
				if err != nil {
					return nil, err
				}
				field.SetInt(int64(intVal))
			}
		}

		results = append(results, elem)
	}

	fmt.Println(results)

	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// thisone is pretty good just to get the as an array with passing the stupid types
func selectAll(db *sql.DB, tableName string) ([]map[string]interface{}, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		// Prepare a slice to hold the row data.
		columnsData := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))
		for i := range columnsData {
			columnPointers[i] = &columnsData[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Then build a map of column name -> column value for the current row.
		rowData := make(map[string]interface{})
		for i, colName := range columns {
			var v interface{}
			val := columnPointers[i].(*interface{})
			b, ok := (*val).([]byte)
			if ok {
				v = string(b)
			} else {
				v = *val
			}
			rowData[colName] = v
		}

		results = append(results, rowData)
	}

	return results, nil
}

// use this for god sakes the others are annoying
func selectAll2J(db *sql.DB, tableName string) (string, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return "", err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		// Prepare a slice to hold the row data.
		columnsData := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))
		for i := range columnsData {
			columnPointers[i] = &columnsData[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return "", err
		}

		// Then build a map of column name -> column value for the current row.
		rowData := make(map[string]interface{})
		for i, colName := range columns {
			var v interface{}
			val := columnPointers[i].(*interface{})
			b, ok := (*val).([]byte)
			if ok {
				v = string(b)
			} else {
				v = *val
			}
			rowData[colName] = v
		}

		results = append(results, rowData)
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func checkUserTbl(db *sql.DB) {
	// Create the table if it doesn't exist
	db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(255) PRIMARY KEY,
		email VARCHAR(255),
		verified_email BOOLEAN,
		name VARCHAR(255),
		given_name VARCHAR(255),
		family_name VARCHAR(255),
		picture VARCHAR(255),
		locale VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		stripe_customer_id  VARCHAR(255)
	)`)

}

func insertUser(db *sql.DB, user User) error {
	// Create the table if it doesn't exist
	/*_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(255) PRIMARY KEY,
		email VARCHAR(255),
		verified_email BOOLEAN,
		name VARCHAR(255),
		given_name VARCHAR(255),
		family_name VARCHAR(255),
		picture VARCHAR(255),
		locale VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		stripe_customer_id  VARCHAR(255),
	)`)
	if err != nil {
		return err
	}*/

	// Insert the user
	_, err := db.Exec(`INSERT INTO users (id, email, verified_email, name, given_name, family_name, picture, locale,stripe_customer_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?,?)
		ON DUPLICATE KEY UPDATE last_login=CURRENT_TIMESTAMP`,
		user.ID, user.Email, user.VerifiedEmail, user.Name, user.GivenName, user.FamilyName, user.Picture, user.Locale, user.stripe_customer_id)
	return err
}

func getUserByEmail(db *sql.DB, email string) (*User, error) {
	// Query the database for the user with the given email
	row := db.QueryRow(`SELECT id, email, verified_email, name, given_name, family_name, picture, locale, stripe_customer_id FROM users WHERE email = ?`, email)

	// Create a User struct to hold the retrieved values
	var user User

	// Scan the result into the User struct
	err := row.Scan(&user.ID, &user.Email, &user.VerifiedEmail, &user.Name, &user.GivenName, &user.FamilyName, &user.Picture, &user.Locale, &user.stripe_customer_id)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows were returned, the user does not exist
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func executeSQLQuery(db *sql.DB, sqlQuery string) ([]map[string]interface{}, error) {
	// Execute the SQL query
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Create a slice of maps to store the query results
	results := []map[string]interface{}{}

	// Fetch the data rows
	for rows.Next() {
		// Create a map to store each row's data
		rowData := make(map[string]interface{})

		// Create a slice of interface{} to store scan results
		scanArgs := make([]interface{}, len(columns))
		for i := range columns {
			// Create a temporary variable to hold the scan result
			var temp interface{}
			scanArgs[i] = &temp

			// Scan the result into the temporary variable
			if err := rows.Scan(&temp); err != nil {
				// Log the error and continue or return an error
				fmt.Printf("Error scanning row: %v\n", err)
				continue // or return nil, err
			}

			// Assign the temporary variable to the map
			rowData[columns[i]] = temp
		}

		// Append the row data to the results
		results = append(results, rowData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func ExecuteAndPrintQueryResults(db *sql.DB, sqlQuery string) {
	// Execute the SQL query
	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	// Create a slice of interface{} to store scan results
	scanArgs := make([]interface{}, len(columns))
	for i := range columns {
		scanArgs[i] = new(interface{})
	}

	// Print column names
	fmt.Println("Columns:")
	for _, col := range columns {
		fmt.Printf("%s\t", col)
	}
	fmt.Println()

	// Fetch and print the data rows
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}

		for _, val := range scanArgs {
			// Convert bytes to string
			strVal, ok := (*val.(*interface{})).([]byte)
			if ok {
				fmt.Printf("%s\t", string(strVal))
			} else {
				fmt.Printf("%v\t", *val.(*interface{}))
			}
		}
		fmt.Println()
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func insertJsonDataWithoutDateTimeStamp(db *sql.DB, tableName string, jsonData string) error {
	// Parse the JSON string into a map
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return err
	}

	// Get the current year
	currentYear := time.Now().Year()

	// Generate the table name with the current year
	tableNameWithYear := fmt.Sprintf("%s_%d", tableName, currentYear)

	// Check if the table exists
	tableExists := tableExists(db, tableNameWithYear)

	if !tableExists {
		// Create the table with columns based on the JSON keys
		createTableStmt := fmt.Sprintf("CREATE TABLE %s (idCount INT AUTO_INCREMENT, ", tableNameWithYear)

		for key := range data {
			createTableStmt += fmt.Sprintf("%s TEXT, ", key)
		}

		createTableStmt += "PRIMARY KEY (idCount))"

		_, err := db.Exec(createTableStmt)
		if err != nil {
			return err
		}
	}

	// Prepare the insert statement
	insertStmt := fmt.Sprintf("INSERT INTO %s (", tableNameWithYear)

	var columns []string
	var placeholders []string
	var values []interface{}

	for key, value := range data {
		columns = append(columns, key)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	insertStmt += strings.Join(columns, ", ") + ") VALUES (" + strings.Join(placeholders, ", ") + ")"

	// Execute the insert statement
	_, err = db.Exec(insertStmt, values...)
	if err != nil {
		return err
	}

	return nil
}

func insertJsonData(db *sql.DB, tableName string, jsonData string) error {
	// Parse the JSON string into a map
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return err
	}

	// Get the current year
	currentYear := time.Now().Year()

	// Generate the table name with the current year
	tableNameWithYear := fmt.Sprintf("%s_%d", tableName, currentYear)

	// Check if the table exists
	tableExists := tableExists(db, tableNameWithYear)

	if !tableExists {
		// Create the table with columns based on the JSON keys
		createTableStmt := fmt.Sprintf("CREATE TABLE %s (idCount INT AUTO_INCREMENT, ", tableNameWithYear)

		for key := range data {
			createTableStmt += fmt.Sprintf("%s TEXT, ", key)
		}

		createTableStmt += "created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, "
		createTableStmt += "PRIMARY KEY (idCount))"

		_, err := db.Exec(createTableStmt)
		if err != nil {
			return err
		}
	}

	// Prepare the insert statement
	insertStmt := fmt.Sprintf("INSERT INTO %s (", tableNameWithYear)

	var columns []string
	var placeholders []string
	var values []interface{}

	for key, value := range data {
		columns = append(columns, key)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	insertStmt += strings.Join(columns, ", ") + ") VALUES (" + strings.Join(placeholders, ", ") + ")"

	// Execute the insert statement
	_, err = db.Exec(insertStmt, values...)
	if err != nil {
		return err
	}

	return nil
}

func getTableColumnNames(db *sql.DB, tableName string) ([]string, error) {
	// Prepare the SQL query to retrieve column names
	query := fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = '%s'", tableName)

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Fetch the column names
	columnNames := make([]string, 0)
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return nil, err
		}
		columnNames = append(columnNames, columnName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columnNames, nil
}
func getTableData(db *sql.DB, tableName string) ([][]string, error) {
	// Prepare the SQL query to retrieve table data
	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Fetch the table data
	tableData := make([][]string, 0)
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Create a slice of interface{} to hold the row values
	values := make([]interface{}, len(columns))
	rowData := make([]string, len(columns))
	for i := range values {
		values[i] = &rowData[i]
	}

	// Iterate over the rows and scan the values into rowData
	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		// Remove quotes from each element in rowData
		for i, value := range rowData {
			rowData[i] = removeQuotes(value)
		}
		// Make a copy of the rowData to store in tableData
		copiedData := make([]string, len(rowData))
		copy(copiedData, rowData)
		tableData = append(tableData, copiedData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tableData, nil
}

func removeQuotes(s string) string {
	return strings.ReplaceAll(s, "\"", "")
}
