package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct{
    ID            primitive.ObjectID  `bson:"id"`
    Name          *string             `json:"name"`
    Email         *string             `json:"email"`
    Phone         *string             `json:"phone"`
    Password      *string             `json:"password"`
    OTP           *string             `json:"otp"`
    OTP_Expires   time.Time           `json:"otp_expires"`
    Created_At    time.Time           `json:"created_at"`
	Updated_At    time.Time           `json:"updated_at"`
	Employee_ID   string              `json:"employee_id"`
}

