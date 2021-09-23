package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Metadata struct {
	neo4j_node_id     int
	img_url           string
	datetime_original string
	model             string
	make              string
	maker_note        string
	software          string
	gps_latitude      string
	gps_longitude     string
}

var postgresuri string = "host=0.0.0.0 user=gopher password=gopher dbname=metadata port=5432"

func main() {
	db, err := gorm.Open(postgres.Open(postgresuri), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	xc := []map[string]interface{}{}
	result := db.Table("metadata").Find(&xc)
	fmt.Println(result)
	fmt.Println(xc)

}
