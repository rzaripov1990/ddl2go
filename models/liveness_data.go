package models

import (
    "github.com/google/uuid"
    "time"
)

type LivenessData struct { 
	Id uuid.UUID `json:"id" db:"id"`  
	CreationDate *time.Time `json:"creationDate" db:"creation_date"`  
	Biin string `json:"biin" db:"biin"`  
	PhoneNo *string `json:"phoneNo" db:"phone_no"`  
	Reason *string `json:"reason" db:"reason"`  
	ChannelId *string `json:"channelId" db:"channel_id"`  
	RedirectUrl *string `json:"redirectUrl" db:"redirect_url"`  
	CallbackUrl *string `json:"callbackUrl" db:"callback_url"`  
	NewPhoneNo *string `json:"newPhoneNo" db:"new_phone_no"`  
	BusinessKey *string `json:"businessKey" db:"business_key"`  
	PassedUrl *string `json:"passedUrl" db:"passed_url"`  
	FailedUrl *string `json:"failedUrl" db:"failed_url"`  
	CompletedDate *time.Time `json:"completedDate" db:"completed_date"`  
	Arh *bool `json:"arh" db:"arh"` 
}

type LivenessDataArr []LivenessData