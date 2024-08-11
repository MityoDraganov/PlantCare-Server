// dropAllTables.go
package main

import (
	"log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func InitDBDrop() *gorm.DB {
	dsn := "sqlserver://server:P@ssw0rd123!@localhost:1433?database=Plant_Care"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	return db
}

func DropAllTables(db *gorm.DB) error {
	rows, err := db.Raw("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE'").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return err
		}
		tableNames = append(tableNames, tableName)
	}

	for _, tableName := range tableNames {
		if err := db.Migrator().DropTable(tableName); err != nil {
			return err
		}
		log.Printf("Dropped table: %s", tableName)
	}

	return nil
}

func main() {
	db := InitDBDrop()

	if err := DropAllTables(db); err != nil {
		log.Fatal("failed to drop all tables:", err)
	}

	log.Println("All tables dropped successfully!")
}
