package models

import (
    "github.com/google/uuid"
)

type LivenessFiles struct { 
	Id int `json:"id" db:"id"`  
	LivenessId *uuid.UUID `json:"livenessId" db:"liveness_id"` // (ref to LivenessData.Id)  
	PhotoUuid *uuid.UUID `json:"photoUuid" db:"photo_uuid"`  
	VideoUuid uuid.UUID `json:"videoUuid" db:"video_uuid"` 
}

type LivenessFilesArr []LivenessFiles