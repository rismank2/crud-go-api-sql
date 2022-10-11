package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

func main() {
	//koneksi ke databse
	db, err = gorm.Open("mysql", "root:@/orders_by?charset=utf8&parseTime=True")
	//cek koneksi
	if err != nil {
		log.Panicln("Koneksi Gagal", err)
	} else {
		log.Panicln("Koeneksi Berhasil")
	}
}
