package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Initailize DBmanager

//DBManger for controling the database
type DBManager struct {
	*sql.DB
}

//NewDBManager : Initailize for struct DBManager
func NewDBManager(configData config) (DBManager, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", configData.ID, configData.PW, configData.DB)
	db, err := sql.Open(configData.StoreType, dsn)
	if err != nil {
		return DBManager{}, err
	}

	dbm := DBManager{db}

	// check for connection
	err = dbm.Ping()
	if err != nil {
		return DBManager{}, err
	}

	return dbm, nil
}

func (dbm DBManager) insertStock(obj stock) error {
	insert, err := dbm.Query("INSERT INTO `stock_data`(stock_id, price_at, open, high, low, close, vol) VALUES(?,?,?,?,?,?,?)",
		obj.id, obj.TradeDate(), obj.dataSet[0], obj.dataSet[1], obj.dataSet[2], obj.dataSet[5], obj.dataSet[4])

	if err != nil {
		/* panic(err) */
		return fmt.Errorf("Stock Data insert err: %s", err)
	}
	defer insert.Close()

	return nil
}

func (dbm DBManager) newStock(si stockInfo, cf config) error {
	var ci int
	// Chech whether type in database
	ids, err := dbm.getCategoryID(si.category, cf.CategoryTable)
	if err != nil {
		return fmt.Errorf("newStock: check category exist err: %s", err.Error())
	}

	// Add new category into database if not exist
	if len(ids) == 0 {
		ci, err = dbm.newCategory(si, cf.CategoryTable)
		if err != nil {
			return fmt.Errorf("newStock: insert new category: %s", err.Error())
		}
	} else {
		ci = ids[0]
	}

	query := fmt.Sprintf("INSERT INTO %s(id, category_id, name) VALUES(?, ?, ?)", cf.NameTable)
	insert, err := dbm.Query(query, si.id, ci, si.name)
	if err != nil {
		return fmt.Errorf("newStock: Stock Data insert err: %s", err)
	}
	defer insert.Close()
	return nil
}

func (dbm DBManager) newCategory(si stockInfo, categoryTable string) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s(`category_name`) VALUES(?)", categoryTable)
	insert, err := dbm.Query(query, si.category)

	if err != nil {
		return 0, fmt.Errorf("newCategory: category data insert err: %s", err)
	}
	defer insert.Close()

	ids, err := dbm.getCategoryID(si.category, categoryTable)
	return ids[0], nil
}

func (dbm DBManager) getCategoryID(categoryName string, categoryTable string) ([]int, error) {
	rows, err := dbm.Query("SELECT id FROM `stock_category` WHERE `category_name`=?", categoryName)
	if err != nil {
		return nil, fmt.Errorf("getCategoryID: get category err %s", err.Error())
	}

	defer rows.Close()
	var id int
	ids := make([]int, 0)
	for rows.Next() {
		err = rows.Scan(&id)
		ids = append(ids, id)
		if err != nil {
			return nil, fmt.Errorf("newStock: check category err : %s", err.Error())
		}
	}
	return ids, nil
}

func (dbm DBManager) checkStockExist(si stockInfo, nameTable string) (bool, error) {
	rows, err := dbm.Query("SELECT COUNT(*) FROM `stock_name` WHERE `id`=?", si.id)
	if err != nil {
		return false, fmt.Errorf("checkStockExist: select query err %s", err.Error())
	}

	defer rows.Close()
	var c int
	for rows.Next() {
		err = rows.Scan(&c)
		if err != nil {
			return false, fmt.Errorf("newStock: check category err : %s", err.Error())
		}
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}
