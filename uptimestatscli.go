package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	dbaddr := flag.String("dbaddr", "localhost", "Database address")
	dbpass := flag.String("dbpass", "", "Database password")
	dbuser := flag.String("dbuser", "uptimestats", "Database user")
	dbname := flag.String("dbname", "uptimestats", "Database name")
	adddomain := flag.String("add", "", "Add a domain")
	createdb := flag.Bool("createdb", false, "Create the database schema")
	flag.Parse()

	if *createdb {
		log.Println("connecting to database")
		db, err := dbconnection(*dbaddr, *dbname, *dbuser, *dbpass)
		if err != nil {
			panic(err)
		}
		log.Println("connection to database successful")
		stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS domains (id SERIAL PRIMARY KEY NOT NULL, name varchar(255) NOT NULL, enabled bool NOT NULL)")
		if err != nil {
			panic(err)
		}
		_, err = stmt.Exec()
		if err != nil {
			panic(err)
		}

		stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS response (id SERIAL PRIMARY KEY NOT NULL, domainid SERIAL REFERENCES domains(id) NOT NULL, time TIMESTAMP NOT NULL, responsetime FLOAT(32), responsecode INT)")
		if err != nil {
			panic(err)
		}
		_, err = stmt.Exec()
		if err != nil {
			panic(err)
		}
	}

	if *adddomain != "" {
		log.Println("connecting to database")
		db, err := dbconnection(*dbaddr, *dbname, *dbuser, *dbpass)
		if err != nil {
			panic(err)
		}
		log.Println("connection to database successful")
		stmt, err := db.Prepare("INSERT INTO domains(name, enabled) VALUES($1, true)")
		if err != nil {
			panic(err)
		}
		_, err = stmt.Exec(*adddomain)
		if err != nil {
			panic(err)
		}
		log.Println("domain", *adddomain, "added")
		db.Close()
	}
}

func dbconnection(DatabaseAddress, DatabaseName, DatabaseUser, DatabasePassword string) (*sql.DB, error) {
	DatabaseConnectionString := fmt.Sprintf("dbname=%s user=%s host=%s password=%s", DatabaseName, DatabaseUser, DatabaseAddress, DatabasePassword)
	db, err := sql.Open("postgres", DatabaseConnectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
