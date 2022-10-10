package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"log"

	"github.com/brianvoe/gofakeit"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// createTestDatabase ensures tables, required indexes, test data, etc
func createTestDatabase(dbName string, schemaName string) (*gorm.DB, func(), string, error) {
	dialStringTemplate := "host=localhost dbname=%s sslmode=disable"

	adminDialString := fmt.Sprintf(dialStringTemplate, "postgres")
	adminDB, err := gorm.Open(postgres.Open(adminDialString), &gorm.Config{})
	if err != nil {
		return nil, nil, adminDialString, fmt.Errorf("unable to open Database (%s): %s", adminDialString, err)
	}

	defer func() {
		sdb, _ := adminDB.DB()
		sdb.Close()
	}()

	adminDB = adminDB.Exec(fmt.Sprintf(`create database "%s";`, dbName))
	if adminDB.Error != nil {
		return nil, nil, adminDialString, fmt.Errorf("unable to create Database %s (%s): %s", dbName, adminDialString, adminDB.Error)
	}

	dbAddr := fmt.Sprintf(dialStringTemplate, dbName)
	db, err := gorm.Open(postgres.Open(dbAddr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: fmt.Sprintf("%s.", schemaName), SingularTable: false}}) //nolint:exhaustivestruct
	if err != nil {
		return nil, nil, dbAddr, fmt.Errorf("unable to open Database: %w", err)
	}

	log.Println("Connected to Test Database")
	// api.db.Exec(fmt.Sprintf(`set search_path to %s, public`, schema))

	return db, func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("\nerror getting sql database: %v", err)
			return
		}

		sqlDB.Close()
		log.Printf("\nDisconnected from Test Database: %s", dbAddr)
	}, dbAddr, nil

}

// dropTestDatabase removes the test database
func dropTestDatabase(dbName string) error {
	dialStringTemplate := "host=localhost dbname=%s sslmode=disable"
	adminDialString := fmt.Sprintf(dialStringTemplate, "postgres")

	adminDB, err := gorm.Open(postgres.Open(adminDialString), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("unable to open Database (%s): %s", adminDialString, err)
	}

	defer func() {
		sdb, _ := adminDB.DB()
		sdb.Close()
	}()

	adminDB = adminDB.Exec(fmt.Sprintf(`drop database "%s";`, dbName))
	if adminDB.Error != nil {
		return fmt.Errorf("unable to drop Database (%s): %s", adminDialString, adminDB.Error)
	}
	return nil
}

// InitTestDatabase sets up a test-specific database.
// The caller is responsible for closing it out, by calling the returned closer.
// This is best done using a defer or in an AfterSuite() function
func InitTestDatabase(schemaName string) (*gorm.DB, func(), string, error) {
	dbName := fmt.Sprintf("testdb-%s", gofakeit.UUID())

	db, dbCloser, dbAddr, err := createTestDatabase(dbName, schemaName)
	if err != nil {
		dropTestDatabase(dbName)
		return nil, nil, "", fmt.Errorf("error creating test db: %w", err)
	}

	closer := func() {
		dbCloser()
		dropTestDatabase(dbName)
	}

	// run seed data script(s)
	if err := runSQLFile(db, filepath.Join("scripts", "seed.sql")); err != nil {
		defer closer()
		return nil, nil, "", fmt.Errorf("error running 00-create-tables.sql: %w", err)
	}

	return db, closer, dbAddr, nil
}

func runSQLFile(db *gorm.DB, path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf(`error reading "%s": %v`, path, err)
	}

	var commentRE = regexp.MustCompile(`(?m)^\s*--.*$`)
	cleanFile := commentRE.ReplaceAllString(string(file), "")
	statements := strings.Split(cleanFile, ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		db = db.Exec(stmt)
		if db.Error != nil {
			return fmt.Errorf(`error executing "%s": %v`, stmt, db.Error)
		}
	}
	return nil
}
