package utils

import ( 
    "math/rand"
    "time"
	"fmt"

)        

var randowCodes = [...]byte{
    '1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
}        

func GenerateOTP() string{
    var r *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	var otp = ""
    for i := 0; i < 1; i++ {
        var pswrd []byte = make([]byte, 6)

        for j := 0; j < 6; j++ {
            index := r.Int() % len(randowCodes)
            pswrd[j] = randowCodes[index]
        }
		otp = fmt.Sprintf("%s\n", string(pswrd))     
		                         
    }    
	return otp
} 