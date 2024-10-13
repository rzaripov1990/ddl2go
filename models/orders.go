package models

import (
    "time"
)

type Orders struct { 
	OrderId int `json:"orderId" db:"order_id"`  
	CustomerId *int `json:"customerId" db:"customer_id"` // (ref to Customers.Id)  
	OrderDate time.Time `json:"orderDate" db:"order_date"`  
	Amount float64 `json:"amount" db:"amount"` 
}

type OrdersArr []Orders