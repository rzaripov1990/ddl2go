package models

import (
    "github.com/google/uuid"
    "time"
)

type KiscData struct { 
	Id uuid.UUID `json:"id" db:"id"`  
	CreationDate *time.Time `json:"creationDate" db:"creation_date"`  
	Biin *string `json:"biin" db:"biin"`  
	LivenessId *uuid.UUID `json:"livenessId" db:"liveness_id"`  
	ChannelId *string `json:"channelId" db:"channel_id"`  
	VendorId *int `json:"vendorId" db:"vendor_id"`  
	UniqKey string `json:"uniqKey" db:"uniq_key"`  
	Similarity *float32 `json:"similarity" db:"similarity"`  
	ErrorMsg *string `json:"errorMsg" db:"error_msg"`  
	VendorName *string `json:"vendorName" db:"vendor_name"` 
}

type KiscDataArr []KiscData