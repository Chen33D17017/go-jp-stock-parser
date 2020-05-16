# go-jp-stock-parser
This is an simple project which can parse the Japanese stock price data into mysql database
## Usecase
Install the program by running "go install" or "go build" 
The tools will give you the helper like beneath:

``` 
  -end string
        End time in format '20xx-xx-xx'
  -file string
        Config file for storage, default with config.json (default "config.json")
  -start string
        Start time in format '20xx-xx-xx'
  -stock int
        The number of stock 
  -u    Update to today
```

To connect to MySQL, please assure the driver has been installed, input the config file
by format:

``` json
{
    "type": "mysql",
    "id" : "mysqlID",
    "password" : "*********",
    "database" : "stock_price",
    "dataTable": "stock_data",
    "nameTable": "stock_name",
    "categoryTable": "stock_category"
}
```

For the first time, Please run the DDL.sql in your MySQL database

## Requirement 
    - Go
    - Go MySQL Driver
    - MySQL
