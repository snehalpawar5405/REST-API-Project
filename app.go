package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialise() error {

	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", dbuser, dbpassword, dbName)

	var err error
	app.DB, err = sql.Open("mysql", connectionString)

	if err != nil {
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)

	return nil
}

func (app *App) run(addr string) {
	app.HandleRoutes()
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

func sendResponce(w http.ResponseWriter, statusCode int, payload interface{}) {

	response, _ := json.Marshal(payload) // Convert payload to JSON

	w.Header().Set("Content-type", "application/json") // This line sets the HTTP response header to specify that the content type of the response is JSON.
	w.WriteHeader(statusCode)                          // Set the HTTP status code for the response
	w.Write(response)                                  // Write the JSON response to the HTTP response writer
}

func sendError(w http.ResponseWriter, StatusCode int, err string) {
	err_msg := map[string]string{"error": err}
	sendResponce(w, StatusCode, err_msg)
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {

	products, err := getProducts(app.DB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponce(w, http.StatusOK, products)
}

func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {

	products, err := getProduct(app.DB, r)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponce(w, http.StatusOK, products)

}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {

	var p product
	err := json.NewDecoder(r.Body).Decode(&p) //This creates a new JSON decoder that reads from the r.Body of an HTTP request.
	// r.Body represents the request body, which is where the JSON data is typically sent in an HTTP POST or PUT request.
	// .Decode(&p): This part of the code attempts to decode the JSON data read from the request body and store it in the variable p
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = p.createProduct(app.DB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponce(w, http.StatusOK, p)

}

func (app *App) productUpdate(w http.ResponseWriter, r *http.Request) {

	myvar := mux.Vars(r)
	key, err := strconv.Atoi(myvar["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var p product
	p.Id = key

	err = json.NewDecoder(r.Body).Decode(&p) //This creates a new JSON decoder that reads from the r.Body of an HTTP request.
	// r.Body represents the request body, which is where the JSON data is typically sent in an HTTP POST or PUT request.
	// .Decode(&p): This part of the code attempts to decode the JSON data read from the request body and store it in the variable p
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = p.productUpdate(app.DB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponce(w, http.StatusOK, p)

}

func (app *App) productDelete(w http.ResponseWriter, r *http.Request) {
	myvar := mux.Vars(r)

	var p product
	key, err := strconv.Atoi(myvar["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid id")
		return
	}

	p.Id = key

	err = p.productDelete(app.DB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponce(w, http.StatusOK, map[string]string{"result": "DONE,Successfuly deleted "})

}

func (app *App) HandleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.productUpdate).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.productDelete).Methods("DELETE")

}
