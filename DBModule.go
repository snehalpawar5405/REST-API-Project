package main

import (
	"database/sql"
	"errors"
	"fmt"
	//"testing/quick"

	//"log"
	"net/http"
	"strconv"

	//"testing/quick"

	"github.com/gorilla/mux"
)

type product struct {
	Id       int     "json:'id'"
	Name     string  "json:'name'"
	Quantity int     "json: 'quantity'"
	Price    float64 "json: 'price'"
}

func getProducts(db *sql.DB) ([]product, error) {

	quary := "select * from products"

	row, err := db.Query(quary)

	if err != nil {
		return nil, err
	}

	var products []product

	for row.Next() {
		var prod product
		err := row.Scan(&prod.Id, &prod.Name, &prod.Quantity, &prod.Price)

		if err != nil {
			return nil, err
		}

		products = append(products, prod)
	}
	return products, nil
}

func getProduct(db *sql.DB, r *http.Request) (*product, error) {
	myvar := mux.Vars(r)

	key, err := strconv.Atoi(myvar["id"])

	if err != nil {
		return nil, errors.New("invalid product id")
	}
	var p product

	quary := fmt.Sprintf("select * from products where id=%v", key)

	row := db.QueryRow(quary)
	err = row.Scan(&p.Id, &p.Name, &p.Quantity, &p.Price)

	if err != nil {

		switch err {
		case sql.ErrNoRows:
			return nil, errors.New("Product not found")
		default:
			return nil, err
		}

	}
	return &p, nil

}

func (p *product) createProduct(db *sql.DB) error {

	quary := fmt.Sprintf("insert into products(name,quantity,price) values('%v',%d,%f)", p.Name, p.Quantity, p.Price)
	res, err := db.Exec(quary)
	if err != nil {
		return err
	}
	no, err := res.LastInsertId()

	if err != nil {
		return err
	}

	p.Id = int(no)
	return nil
}

func (p *product) productUpdate(db *sql.DB) error {

	quary := fmt.Sprintf("update products set name='%v', quantity= %v , price= %v where id=%v", p.Name, p.Quantity, p.Price, p.Id)

	result, err := db.Exec(quary)
	if err != nil {
		return err
	}

	rowaffected, err := result.RowsAffected()

	if rowaffected == 0 {
		return errors.New("No such row exist")
	}

	return nil
}

func (p *product) productDelete(db *sql.DB) error {

	query := fmt.Sprintf("delete from products where id=%v", p.Id)

	result, err := db.Exec(query)

	if err != nil {
		return err
	}

	rowaffected, err := result.RowsAffected()

	if rowaffected == 0 {
		return errors.New("No such id exist")
	}

	return nil
}
