package models

type Customers struct { 
	Id int `json:"id" db:"id"`  
	Name string `json:"name" db:"name"`  
	Email string `json:"email" db:"email"`  
	Phone *string `json:"phone" db:"phone"` 
}

type CustomersArr []Customers