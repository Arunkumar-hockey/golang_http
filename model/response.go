package model

type Response struct   {
    Status      int     `json:"statuscode"`
    Message     string  `json:"message"`
    Data        Employee  `json:"data"` 
}

