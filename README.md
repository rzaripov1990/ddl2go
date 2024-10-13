# pg_to_go

This file was generated by ChatGPT from the source code :)

## Struct Generator from PostgreSQL

This project generates Go structs based on PostgreSQL database tables. It retrieves table and column information from information_schema, including column comments, and creates a Go struct for each table. Additionally, foreign key references are added as comments.

## How It Works

The program connects to a PostgreSQL database using the connection string from the PG_CONN_STR environment variable.
It fetches a list of tables from the schema and retrieves column information for each table, including:
* Column name.
* Data type.
* Foreign key references.
* Column comments (if any).

Based on this data, Go structs are generated for each table with corresponding JSON and DB tags.
If there is a foreign key, a comment is added indicating which table and column it references.

## Requirements

* Go 1.18+
* PostgreSQL 9.6+
* Installed Go libraries:
  * [github.com/jmoiron/sqlx](github.com/jmoiron/sqlx) for SQL handling.
  * [github.com/lib/pq](github.com/lib/pq) for the PostgreSQL driver.
  * [github.com/joho/godotenv/autoload](github.com/joho/godotenv/autoload) to auto-load environment variables from a `.env` file.
  * [github.com/google/uuid](github.com/google/uuid) for handling UUID types (if used in your database).

## Installation

Clone this repository.
Install the dependencies:
```bash
go mod tidy
```

Create a .env file in the project's root directory and add the following environment variables:
```bash
PACKAGE=<your-package-name>
PG_CONN_STR=<your-postgres-connection-string>
```
Example connection string:
```env
PG_CONN_STR="user=postgres password=example dbname=mydb sslmode=disable"
```

## Usage

Run the program with the following command:

```bash
go run main.go
```

The program will generate Go files in the directory specified by the PACKAGE environment variable. Each file will contain a Go struct corresponding to a table in the database.

## Example Output
Here is an example of a struct generated for tables:

```sql
CREATE TABLE customers (
    id INT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    phone VARCHAR(15)
);

CREATE TABLE orders (
    order_id INT PRIMARY KEY,
    customer_id INT,
    order_date DATE NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    CONSTRAINT fk_customer
        FOREIGN KEY (customer_id)
        REFERENCES customers(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
```

```go
package <your-package-name>

import (
    "time"
)

type Customers struct { 
	Id int `json:"id" db:"id"`  
	Name string `json:"name" db:"name"`  
	Email string `json:"email" db:"email"`  
	Phone *string `json:"phone" db:"phone"` 
}

type CustomersArr []Customers

type Orders struct { 
	OrderId int `json:"orderId" db:"order_id"`  
	CustomerId *int `json:"customerId" db:"customer_id"` // (ref to Customers.Id)  
	OrderDate time.Time `json:"orderDate" db:"order_date"`  
	Amount float64 `json:"amount" db:"amount"` 
}

type OrdersArr []Orders
```


### Supported Data Types

The program maps PostgreSQL data types to corresponding Go types. 

Examples of type mappings:

`integer`, `smallint`, `serial` → `int`

`bigint` → `int64`

`boolean` → `bool`

`uuid` → `uuid.UUID`

`timestamp`, `date`, `time` → `time.Time`

`text`, `varchar` → `string`

`json`, `bytea`, `xml` → `[]byte`

If a column can be `NULL`, the type will be wrapped in a pointer (e.g., `*int`).