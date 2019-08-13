package main

import (
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
	txdb.Register("txdb_postgres", "postgres", url)
	var db *gorm.DB
	var err error
	for i := 0; i < 3; i++ {
		db, err = gorm.Open("txdb_postgres", "tx_1")
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
