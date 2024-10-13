package models

import (
    "time"
    "github.com/google/uuid"
)

type ExampleTypes struct { 
	Id int `json:"id" db:"id"`  
	SmallValue *int `json:"smallValue" db:"small_value"`  
	LargeValue *int64 `json:"largeValue" db:"large_value"`  
	DecimalValue *float64 `json:"decimalValue" db:"decimal_value"`  
	NumericValue *float64 `json:"numericValue" db:"numeric_value"`  
	FloatValue *float32 `json:"floatValue" db:"float_value"`  
	DoubleValue *float32 `json:"doubleValue" db:"double_value"`  
	CharValue *string `json:"charValue" db:"char_value"` //comment  
	VarcharValue *string `json:"varcharValue" db:"varchar_value"`  
	TextValue *string `json:"textValue" db:"text_value"`  
	DateValue *time.Time `json:"dateValue" db:"date_value"`  
	TimeValue *time.Time `json:"timeValue" db:"time_value"`  
	TimestampValue *time.Time `json:"timestampValue" db:"timestamp_value"`  
	TimestamptzValue *time.Time `json:"timestamptzValue" db:"timestamptz_value"`  
	IntervalValue *time.Duration `json:"intervalValue" db:"interval_value"`  
	BooleanValue *bool `json:"booleanValue" db:"boolean_value"`  
	JsonValue []byte `json:"jsonValue" db:"json_value"`  
	JsonbValue []byte `json:"jsonbValue" db:"jsonb_value"`  
	XmlValue []byte `json:"xmlValue" db:"xml_value"`  
	UuidValue *uuid.UUID `json:"uuidValue" db:"uuid_value"`  
	InetValue any `json:"inetValue" db:"inet_value"`  
	CidrValue any `json:"cidrValue" db:"cidr_value"`  
	PointValue any `json:"pointValue" db:"point_value"`  
	LineValue any `json:"lineValue" db:"line_value"`  
	LsegValue any `json:"lsegValue" db:"lseg_value"`  
	BoxValue any `json:"boxValue" db:"box_value"`  
	ArrayValue any `json:"arrayValue" db:"array_value"`  
	TextArrayValue any `json:"textArrayValue" db:"text_array_value"`  
	ByteaValue []byte `json:"byteaValue" db:"bytea_value"`  
	MoneyValue *float64 `json:"moneyValue" db:"money_value"`  
	EnumValue any `json:"enumValue" db:"enum_value"` 
}

type ExampleTypesArr []ExampleTypes