package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	gormigrate "gopkg.in/gormigrate.v1"
)

var (
	migration_1 = func(tx *gorm.DB) error {
		type Person struct {
			gorm.Model
			Name string
		}

		return tx.CreateTable(&Person{}).Error
	}
	migration_2 = func(tx *gorm.DB) error {
		type Person struct {
			Age int
		}
		return tx.AutoMigrate(&Person{}).Error
	}
	return_nil = func(tx *gorm.DB) error {
		return nil
	}
)

type Person struct {
	gorm.Model
	Name string
	Age  int
}

func MigrateAll(gdb *gorm.DB) error {
	m := gormigrate.New(gdb, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID:       "first",
			Migrate:  migration_1,
			Rollback: return_nil,
		},
		{
			ID:       "second",
			Migrate:  migration_2,
			Rollback: return_nil,
		},
	})
	return m.Migrate()
}

func main() {
	url := os.Getenv("DATABASE_URL")
	txdb.Register("txdb", "postgres", url)
	s, err := sql.Open("txdb", "tx_1")
	if err != nil {
		panic(fmt.Sprintf("cannot open connection: %s", err))
	}
	var db *gorm.DB
	for i := 0; i < 3; i++ {
		db, err = gorm.Open("postgres", s)
		if err == nil {
			break
		}
		fmt.Printf("connection failed, retrying in 10 seconds. Reason: %s\n", err)
		time.Sleep(10 * time.Second)
	}
	if err != nil {
		panic(fmt.Sprintf("connection failed: %s", err))
	}

	defer db.Close()
	err = MigrateAll(db)
	if err != nil {
		panic(fmt.Sprintf("migration failed: %s", err))
	}
}
