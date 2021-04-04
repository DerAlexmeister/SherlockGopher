package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Metadata struct {
	img_id            int
	condition         bool
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
	//var metadata Metadata
	//var metadatas []Metadata
	db, err := gorm.Open(postgres.Open(postgresuri), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	xc := []map[string]interface{}{}
	result := db.Table("metadata").Find(&xc)
	fmt.Println(result)
	fmt.Println(xc)
	/*for k, v := range metadatas {
		fmt.Println(k, " img_id:", v.img_id, " condition:", v.condition, " datetime_original:", v.datetime_original, " model:", v.model, " make:", v.make, " maker_note:", v.maker_note, " software:", v.software, " gps_latitude:", v.gps_latitude, " gps_longitude:", v.gps_longitude)
		fmt.Println("--------")
	}

	fmt.Println("-----------------")
	db.Table("metadata").Find(&metadatas)
	fmt.Println(metadatas)
	fmt.Println(len(metadatas))
	fmt.Println(metadatas[1].make)
	fmt.Println("-----------------")

	tr := map[string]interface{}{}
	db.Table("metadata").Take(&tr)
	fmt.Println(tr)
	fmt.Println("-----------------")*/

}
