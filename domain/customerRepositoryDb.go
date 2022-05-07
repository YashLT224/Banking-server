package domain

import (
	"Banking/errs"
	"Banking/logger"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type CustomerRepositoryDb struct {
	db *sqlx.DB
}

func (d CustomerRepositoryDb) FindAll(status string) ([]Customer, *errs.AppError) {
	// var rows *sql.Rows
	var err error
	customers := make([]Customer, 0)
	if status == "" {
		findAllSql := "select customer_id,name,city,zipcode,date_of_birth,status from customers"
		err = d.db.Select(&customers, findAllSql)
		// rows, err = d.db.Query(findAllSql)
	} else {
		findAllSql := "select customer_id,name,city,zipcode,date_of_birth,status from customers where status=?"
		err = d.db.Select(&customers, findAllSql, status)
		// rows, err = d.db.Query(findAllSql, status)
	}

	if err != nil {
		logger.Error("Error while querying customer table")
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	//	var customers []Customer
	// err = sqlx.StructScan(rows, &customers)
	// if err != nil {
	// 	logger.Error("Error while querying customer table")
	// 	return nil, errs.NewUnexpectedError("unexpected database error")
	// }
	// for rows.Next() {
	// 	var c Customer
	// 	err := rows.Scan(&c.Id, &c.Name, &c.City, &c.ZipCode, &c.DateofBirth, &c.Status)
	// 	if err != nil {
	// 		logger.Error("Error while querying customer table")
	// 		return nil, errs.NewUnexpectedError("unexpected database error")
	// 	}

	// 	customers = append(customers, c)
	// }
	return customers, nil

}

func (d CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	fmt.Println(id)
	customerSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers where customer_id = ?"
	// row := d.db.QueryRow(customerSql, id)
	var c Customer
	err := d.db.Get(&c, customerSql, id)
	// err := row.Scan(&c.Id, &c.Name, &c.City, &c.ZipCode, &c.DateofBirth, &c.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not Found")
		} else {

			logger.Error("Error while scanning customer" + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}

	}
	return &c, nil
}

func NewCustomerRepositoryDb(dbClient *sqlx.DB) CustomerRepositoryDb {

	return CustomerRepositoryDb{dbClient}
}
