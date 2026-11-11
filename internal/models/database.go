package models

import (
	"log"

	_ "github.com/ncruces/go-sqlite3/embed" // SQLite engine inside the program (no install needed)
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"se69-lab04-starter/internal/utils"
)

// open the SQLite file and reset the data every time the API starts
func ConnectDB(path string) (*gorm.DB, error) {
	db, err := gorm.Open(gormlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // do not print SQL
	})
	if err != nil {
		return nil, err
	}
	log.Printf("database: opened SQLite file %s", path)

	// drop old tables, then create them again from the structs
	tables := []any{&Session{}, &Loan{}, &Equipment{}, &User{}}
	if err := db.Migrator().DropTable(tables...); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&User{}, &Equipment{}, &Loan{}, &Session{}); err != nil {
		return nil, err
	}
	log.Println("database: tables reset (users, equipment, loans, sessions)")

	if err := seedData(db); err != nil {
		return nil, err
	}
	log.Println("database: sample data added (4 users, 3 equipment, 1 loan)")
	return db, nil
}

// --- sample data ---

func seedData(db *gorm.DB) error {
	hash, err := utils.HashPassword("pass1234")
	if err != nil {
		return err
	}
	users := []User{
		{Username: "student1", Role: "student", PasswordHash: hash},
		{Username: "student2", Role: "student", PasswordHash: hash},
		{Username: "staff1", Role: "staff", PasswordHash: hash},
		{Username: "teacher1", Role: "teacher", PasswordHash: hash},
	}
	equipment := []Equipment{
		{Name: "Oscilloscope", Category: "measurement", Status: "borrowed"},
		{Name: "Raspberry Pi 5", Category: "board", Status: "available"},
		{Name: "Arduino Uno", Category: "board", Status: "available"},
	}
	loan := Loan{EquipmentID: 1, Borrower: "student2"}

	if err := db.Create(&users).Error; err != nil {
		return err
	}
	if err := db.Create(&equipment).Error; err != nil {
		return err
	}
	return db.Create(&loan).Error
}
