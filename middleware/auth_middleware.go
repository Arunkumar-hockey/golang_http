package middleware

// import(
// 	"fmt"
// 	"github.com/julienschmidt/httprouter"
// 	"net/http"
// 	"GolangHTTP/helper"
// 	"encoding/json"
// 	"GolangHTTP/model"
// )

// func Athentication() httprouter.Handle{
// 	return func(w http.ResponseWriter,r *http.Request,  p httprouter.Params) {
// 		clientToken := r.Header.Get("token")
// 		if clientToken == "" {
// 			w.WriteHeader(http.StatusBadRequest)
// 				value := model.Response{ 
// 					Status: http.StatusBadRequest,
// 					Message: "No authorization header provides"}
// 					json.NewEncoder(w).Encode(value)
// 				return
// 		}

// 		claims, err := helper.ValidateToken(clientToken)
// 		if err != "" {
// 			w.WriteHeader(http.StatusBadRequest)
// 				value := model.Response{ 
// 					Status: http.StatusBadRequest,
// 					Message: "Error validating the token"}
// 					json.NewEncoder(w).Encode(value)
// 				return
// 		}
// 		// w.Header().Set("email", claims.Email)
// 		// w.Header().Set("name", claims.Name)
// 		// w.Header().Set("phone", claims.Phone)
// 		// w.Header().Set("uid", claims.Uid)
// 		fmt.Println(claims)
// 	} 
// }