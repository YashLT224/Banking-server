package app

import (
	"Banking/domain"
	"Banking/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func sanityCheck() {
	if os.Getenv("Server_ADDRESS") == "" || os.Getenv("SERVER_PORT") == "" || os.Getenv("DB") == "" || os.Getenv("DBID") == "" || os.Getenv("DBPSWD") == "" {
		log.Fatal("Environment variables not defined")
	}
}
func Start() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	sanityCheck()

	router := mux.NewRouter()

	//wiring
	dbClient := getDBClient()
	customerRepositoryDb := domain.NewCustomerRepositoryDb(dbClient)
	ch := CustomerHandlers{service.NewCustomerService(customerRepositoryDb)}

	accountRepositoryDb := domain.NewAccountRepositoryDb(dbClient)
	ah := AccountHandler{service.NewAccountService(accountRepositoryDb)}

	router.HandleFunc("/customers", ch.getAllCustomers).Methods(http.MethodGet).Name("GetAllCustomers")
	router.HandleFunc("/customers/{customer_id:[0-9]+}", ch.GetCustomer).Methods(http.MethodGet).Name("GetCustomer")

	router.HandleFunc("/customers/{customer_id:[0-9]+}/account", ah.NewAccount).Methods(http.MethodPost).Name("NewAccount")
	router.HandleFunc("/customers/{customer_id:[0-9]+}/account/{account_id:[0-9]+}", ah.MakeTransaction).Methods(http.MethodPost).
		Name("NewTransaction")

	am := AuthMiddleware{domain.NewAuthRepository()}

	router.Use(am.authorizationHandler())
	address := os.Getenv("Server_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))
}

func getDBClient() *sqlx.DB {
	DBID := os.Getenv("DBID")
	DBPSWD := os.Getenv("DBPSWD")
	DB := os.Getenv("DB")
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", DBID, DBPSWD, DB))
	if err != nil {

		panic(err.Error())
	}
	db.SetConnMaxLifetime(time.Minute * 20)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	//defer db.Close()
	return db
}
