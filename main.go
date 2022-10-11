package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

type Orders struct {
	ID            int       `form:"order_id" json:"order_id"`
	Customer_name string    `from:"customer_name" json:"customer_name"`
	Ordered_At    time.Time `from:"ordered_at" json:"ordered_at"`
}

type Items struct {
	ID          int    `from:"item_id" json:"item_id"`
	Item_code   string `from:"item_code" json:"item_code"`
	Description string `from:"description" json:"description"`
	Quantity    int    `from:"quantity" json:"quantity"`
	Order_id    int    `from:"order_id" json:"order_id"`
}

// Result is an array of product
type Result struct {
	Orders  interface{} `json:"orders"`
	Items   interface{} `json:"items"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

func main() {
	//koneksi ke databse
	db, err = gorm.Open("mysql", "root:@/orders_by?charset=utf8&parseTime=True")
	if err != nil {
		log.Println("Koneksi Gagal", err)
	} else {
		log.Println("Koneksi Berhasil")
	}
	//migrasi tabel
	db.AutoMigrate(&Orders{}, &Items{})
	handleRequests()
}

func handleRequests() {
	log.Println("Server Aktif di http://127.0.0.1:7332")
	log.Println("CTRL+ C untuk Non Aktif Server")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		res := Result{Code: 404, Message: "Endpoint Tidak Ditemukan"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	myRouter.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		res := Result{Code: 403, Message: "Endpoint Tidak Tersedia"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/orders", createOrder).Methods("POST")
	myRouter.HandleFunc("/orders", getOrders).Methods("GET")
	myRouter.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	myRouter.HandleFunc("/orders/{id}", updateOrder).Methods("PUT")
	myRouter.HandleFunc("/orders/{id}", deleteOrder).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":7332", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ini Halaman Homepage !")
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var item Items
	json.Unmarshal(payloads, &item)
	db.Create(&item)

	var order Orders
	json.Unmarshal(payloads, &order)
	db.Create(&order)

	res := Result{Orders: order, Items: item, Message: "Success create orders"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get orders")

	order := []Orders{}
	db.Find(&order)

	item := []Items{}
	db.Find(&item)

	res := Result{Orders: order, Items: item, Message: "Success get orders"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var order Orders
	db.First(&order, productID)
	var item Items
	db.First(&item, productID)

	res := Result{Orders: order, Items: item, Message: "Success get orders"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var productUpdates Orders
	json.Unmarshal(payloads, &productUpdates)

	var order Orders
	db.First(&order, productID)
	db.Model(&order).Updates(productUpdates)

	var item Orders
	db.First(&item, productID)
	db.Model(&item).Updates(productUpdates)

	res := Result{Orders: order, Items: item, Message: "Success update orders"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var order Orders

	db.First(&order, productID)
	db.Delete(&order)

	res := Result{Message: "Success delete orders"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
