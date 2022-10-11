package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

// struct untuk table auto migrasi
type Orders struct {
	Order_id      int       `from:"order_id" json:"order_id"`
	Customer_name string    `from:"customer_name" json:"customer_name"`
	Ordered_at    time.Time `from:"ordered_at" json:"ordered_at"`
}

type Items struct {
	Item_id     int    `from:"item_id" json:"item_id"`
	Item_code   string `from:"item_code" json:"item_code"`
	Description string `from:"description" json:"description"`
	Quantity    int    `from:"quantity" json:"quantity"`
	Order_id    int    `from:"order_id" json:"order_id"`
}

type Result struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	//koneksi ke databse
	db, err = gorm.Open("mysql", "root:@/orders_by?charset=utf8&parseTime=True")
	//cek koneksi
	if err != nil {
		log.Panicln("Koneksi Gagal", err)
	} else {
		log.Panicln("Koeneksi Berhasil")
	}
	//auto migrasi database
	db.AutoMigrate(&Orders{}, &Items{})
}
