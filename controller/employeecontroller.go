package controller

import(
	"context"
	"GolangHTTP/db"
	"GolangHTTP/model"
	"GolangHTTP/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"encoding/json"
	"net/smtp"
	"strings"
	"fmt"
	"time"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var employeeCollection *mongo.Collection = database.OpenCollection(database.Client, "employee")

var jwtKey = []byte("secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func CreateEmployee() httprouter.Handle {
	return func  (w http.ResponseWriter,r *http.Request,  p httprouter.Params) {
		switch r.Method {
		case http.MethodPost:
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			var employee model.Employee
			json.NewDecoder(r.Body).Decode(&employee)
			emailcount, emailErr := employeeCollection.CountDocuments(ctx, bson.M{"email": employee.Email})
			defer cancel()
			if emailErr != nil {
				log.Println(emailErr)
				w.WriteHeader(http.StatusBadRequest)
			}
	
			phonecount, phoneErr := employeeCollection.CountDocuments(ctx, bson.M{"phone": employee.Phone})
			defer cancel()
			if phoneErr != nil {
				log.Println(phoneErr)
				w.WriteHeader(http.StatusBadRequest)
			}
	
			if emailcount > 0 {
				w.WriteHeader(http.StatusBadRequest)
				value := model.Response{ 
					Status: http.StatusBadRequest,
					Message: "this email already exists"}
					json.NewEncoder(w).Encode(value)
				return
			}

			password := HashPassword(*employee.Password)
			employee.Password = &password

			if phonecount > 0 {
				w.WriteHeader(http.StatusBadRequest)
				value := model.Response{ 
					Status: http.StatusBadRequest,
					Message: "this phone number already exists"}
					json.NewEncoder(w).Encode(value)
				return
			}
			employee.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			employee.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			employee.ID = primitive.NewObjectID()
			employee.Employee_ID = employee.ID.Hex()
	
			resultInsertionNumber, insertErr := employeeCollection.InsertOne(ctx, employee)
			if insertErr != nil {
				log.Println(insertErr)
				w.WriteHeader(http.StatusBadRequest)
			}
	
			w.WriteHeader(http.StatusOK)
			value := model.Response{ 
				Status: http.StatusOK,
				Message: "API Success",
			     Data: employee}
				json.NewEncoder(w).Encode(value)
			//json.NewEncoder(w).Encode(employee)
			fmt.Println(resultInsertionNumber)
	
		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
		}
	}
}

func GetEmployee() httprouter.Handle {
	return func (w http.ResponseWriter,r *http.Request, p httprouter.Params) {
		switch r.Method {
		case http.MethodGet:
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			var employee model.Employee
			employeeId := p.ByName("employee_id")
			err := employeeCollection.FindOne(ctx, bson.M{"employee_id": employeeId}).Decode(&employee)
		defer cancel()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(employee)

		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
		}
	}
}

func SearchEmployeeByName() httprouter.Handle {
	return func(w http.ResponseWriter,r *http.Request, p httprouter.Params) {
		switch r.Method {
		case http.MethodGet:
			var searchEmployee []model.Employee
			queryParam := r.URL.Query().Get("name")

			if queryParam == "" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			defer cancel()

		searchQuery, err :=	employeeCollection.Find(ctx, bson.M{"name": bson.M{"$regex": queryParam}})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = searchQuery.All(ctx, &searchEmployee)
		if err != nil {
			log.Println(err)
			return
		}

		defer searchQuery.Close(ctx)

		if err := searchQuery.Err(); err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(searchEmployee)

		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
		}
	}
}

func GetAllEmployee() httprouter.Handle {
	return func (w http.ResponseWriter,r *http.Request, p httprouter.Params) {
		switch r.Method {
		case http.MethodGet:
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	
		result, err := employeeCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		var allemployees []bson.M
		if err = result.All(ctx, &allemployees); err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(allemployees)
		
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
		}
	}
}

func UpdateEmployee() httprouter.Handle {
	return func (w http.ResponseWriter,r *http.Request, p httprouter.Params) {
		switch r.Method {
		case http.MethodPatch:
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			var employee model.Employee
			employeeId :=  p.ByName("employee_id")
	
			json.NewDecoder(r.Body).Decode(&employee)
			 var updateObj primitive.D
			
			 if employee.Name != nil {
				 updateObj = append(updateObj, bson.E{"name", employee.Name})
			 }
	
			 if employee.Email != nil {
				updateObj = append(updateObj, bson.E{"email", employee.Email})
			}
	
			if employee.Phone != nil {
				updateObj = append(updateObj, bson.E{"phone", employee.Phone})
			}
	
			employee.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{"updated_at", employee.Updated_At})
			filter := bson.M{"employee_id": employeeId}
	
			_, err := employeeCollection.UpdateOne(
				ctx,
				filter,
				bson.D{
					{"$set", updateObj},
				},
				)
	
			if err != nil {
				log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			}
	
			defer cancel()
			w.WriteHeader(http.StatusOK)
		    json.NewEncoder(w).Encode(employee)

		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
		}
	}
}

func DeleteEmployee() httprouter.Handle {
	return func (w http.ResponseWriter,r *http.Request,  p httprouter.Params) {
		switch r.Method {
		case http.MethodDelete:
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		employeeId := p.ByName("employee_id")
		filter := bson.M{"employee_id": employeeId}

		_, deleteErr := employeeCollection.DeleteOne(ctx,filter)
		if deleteErr != nil {
			log.Println(deleteErr)
			w.WriteHeader(http.StatusBadRequest)
		}

		defer cancel()
		value := model.Response{ 
			Status:200,
			Message: "API Success !"}

		w.WriteHeader(http.StatusOK)
		    json.NewEncoder(w).Encode(value)

		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
		}
	}
}

func Login() httprouter.Handle {
	return func (w http.ResponseWriter,r *http.Request,  p httprouter.Params) {
		switch r.Method {
		case http.MethodPost:
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			var employee model.Employee
			var foundEmployee model.Employee

			json.NewDecoder(r.Body).Decode(&employee)

			err := employeeCollection.FindOne(ctx, bson.M{"email": employee.Email}).Decode(&foundEmployee)
			defer cancel()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				value := model.Response{ 
					Status: http.StatusInternalServerError,
					Message: "employee not found"}
					json.NewEncoder(w).Encode(value)
				return
			}

			passwordIsValid, msg := VerifyPassword(*employee.Password, *foundEmployee.Password)
			defer cancel()
			if !passwordIsValid {
				w.WriteHeader(http.StatusInternalServerError)
				value := model.Response{ 
					Status: http.StatusInternalServerError,
					Message: msg} 
					json.NewEncoder(w).Encode(value)
				return
			}

			expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Email: *foundEmployee.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	   http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
			Path: "/",
			HttpOnly: true,
		})

			w.WriteHeader(http.StatusOK)
			value := model.Response{ 
				Status: http.StatusOK,
				Message: "API Success",
			     Data: foundEmployee}
				json.NewEncoder(w).Encode(value)
		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
		}
	}
}

func CheckEmployeeExistAndSendMail() httprouter.Handle{
	return func(w http.ResponseWriter,r *http.Request,  p httprouter.Params) {
		switch r.Method {
		case http.MethodPost:
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			var employee model.Employee
			emailID := r.URL.Query().Get("email")

			err := employeeCollection.FindOne(ctx, bson.M{"email": emailID}).Decode(&employee)
		defer cancel()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
		}

		otp := strings.TrimSpace(utils.GenerateOTP())

		from := "tfpsmtp@gmail.com"
		password := "tfpsmtp@123"
	  
		// Receiver email address.
		to := []string{
		  emailID,
		}

		// smtp server configuration.
		smtpHost := "smtp.gmail.com"
		smtpPort := "587"
	  
		// Message.
		message := []byte(otp)
		
		// Authentication.
		auth := smtp.PlainAuth("", from, password, smtpHost)
		
		// Sending email.
		emailsentErr := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
		if emailsentErr != nil {
			fmt.Println("Failed to send.......")
		  fmt.Println(emailsentErr)
		  w.WriteHeader(http.StatusInternalServerError)
		  return
		}
		var updateObj primitive.D

		if otp != "" {
			updateObj = append(updateObj, bson.E{"otp", otp})
		}

		expirationTime := time.Now().Add(time.Minute * 5)
	    updateObj = append(updateObj, bson.E{"otp_expires", expirationTime})
		

		filter := bson.M{"employee_id": employee.Employee_ID}
	
			_, updateErr := employeeCollection.UpdateOne(
				ctx,
				filter,
				bson.D{
					{"$set", updateObj},
				},
				)

				if updateErr != nil {
					fmt.Println("Failed to update otp in database")
				}

		fmt.Println("Email Sent Successfully!")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(employee.Employee_ID)

		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
		}
	}
}

func VerifyOTP() httprouter.Handle{
	return func(w http.ResponseWriter,r *http.Request,  p httprouter.Params) {
		switch r.Method {
		case http.MethodPost:
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			var employee model.Employee
			//var foundEmployee model.Employee
			employeeID := r.URL.Query().Get("employee_id")
			OTP := r.URL.Query().Get("otp")
			json.NewDecoder(r.Body).Decode(&employee)

		
			err := employeeCollection.FindOne(ctx, 
				bson.M{"employee_id": employeeID,"otp": OTP}).Decode(&employee)
			defer cancel()
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				value := model.Response{ 
					Status: http.StatusNotFound,
					Message: "employee not found"}
					json.NewEncoder(w).Encode(value)
				return
			}
			expirationOTPTime := employee.OTP_Expires
			currentTime := time.Now()

			if expirationOTPTime.Before(currentTime) {
				w.WriteHeader(http.StatusNotFound)
				value := model.Response{ 
					Status: http.StatusNotFound,
					Message: "otp expired"}
					json.NewEncoder(w).Encode(value)
				return
			}

			w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(employee.Employee_ID)
			
		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
		}
	}
}

func UpdatePassword() httprouter.Handle{
	return func(w http.ResponseWriter,r *http.Request,  p httprouter.Params) {
		switch r.Method {
		case http.MethodPost:
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			var employee model.Employee
			employeeID := r.URL.Query().Get("employee_id")
			password := HashPassword(r.URL.Query().Get("password"))
			employee.Password = &password

			json.NewDecoder(r.Body).Decode(&employee)
			 var updateObj primitive.D
			
			 if password != "" {
				 updateObj = append(updateObj, bson.E{"password", password})
			 }

			 employee.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			 updateObj = append(updateObj, bson.E{"updated_at", employee.Updated_At})
			 filter := bson.M{"employee_id": employeeID}
	
			_, err := employeeCollection.UpdateOne(
				ctx,
				filter,
				bson.D{
					{"$set", updateObj},
				},
				)
	
			if err != nil {
				log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			}
	
			defer cancel()
			w.WriteHeader(http.StatusOK)
		    json.NewEncoder(w).Encode(employee)
			
		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
		}
	}
}

func Home() httprouter.Handle{
	return func(w http.ResponseWriter,r *http.Request,  p httprouter.Params) {
			switch r.Method {
			case http.MethodGet:
				cookie, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		
			tokenStr := cookie.Value
		
			claims := &Claims{}
		
			tkn, err := jwt.ParseWithClaims(tokenStr, claims,
				func(t *jwt.Token) (interface{}, error) {
					return jwtKey, nil
				})
		
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		
			if !tkn.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		
			w.WriteHeader(http.StatusOK)
			value := model.Response{ 
				Status: http.StatusOK,
				Message: "API Success"}
				json.NewEncoder(w).Encode(value)

			default:
				w.WriteHeader(http.StatusBadRequest)
				log.Println("ERROR: Invalid HTTP Method")
			}
	}
}

func SignOut() httprouter.Handle{
	return func(w http.ResponseWriter,r *http.Request,  p httprouter.Params) {
 	c := http.Cookie{
		Name:   "token",
		MaxAge: -1}
	http.SetCookie(w, &c)

	w.Write([]byte("Old cookie deleted. Logged out!"))
	}
}


func HashPassword(password string) string{
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

 func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	 err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	 check := true
	 msg := ""

	 if err != nil {
		 msg = "email or password is incorrect"
		 check = false
	 }
	 return check, msg
 }

