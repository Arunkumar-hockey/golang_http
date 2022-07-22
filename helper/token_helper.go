package helper

// import (
// 	"context"
// 	"fmt"
// 	"GolangHTTP/db"
// 	"log"
// 	"os"
// 	"time"
// 	jwt "github.com/dgrijalva/jwt-go"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	//"go.mongodb.org/mongo-driver/mongo/options"
// )

// type SignedDetails struct {
// 	Email   string
// 	Name    string
// 	Phone   string
// 	Uid     string
// 	jwt.StandardClaims
// }

// var employeeCollection *mongo.Collection = database.OpenCollection(database.Client, "employee")

// var SECRET_KEY string = os.Getenv("SECRET_KEY")

// func GenerateAllTokens(email string, name string, phone string, uid string)  (signedToken string,signedRefreshToken string, err error) {
// 	claims := &SignedDetails{
// 		Email: email,
// 		Name: name,
// 		Phone: phone,
// 		Uid: uid,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(), 
// 		},
// 	}

// 	refreshClaims := &SignedDetails{
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
// 		},
// 	}

// 	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
// 	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

// 	if err != nil {
// 		log.Panic(err)
// 		return
// 	}

// 	return token, refreshToken, err
// }

// func UpdateAllTokens(signedToken string, signedRefreshToken string, employeeId string) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 	var updateObj primitive.D

// 	updateObj = append(updateObj, bson.E{"token", signedToken})
// 	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

// 	Updated_At, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
// 	updateObj = append(updateObj, bson.E{"updated_at", Updated_At})

// 	filter := bson.M{"employee_id": employeeId}
	
// 			_, err := employeeCollection.UpdateOne(
// 				ctx,
// 				filter,
// 				bson.D{
// 					{"$set", updateObj},
// 				},
// 				)
// 				defer cancel()

// 				if err != nil {
// 					log.Panic(err)
// 					return
// 				}
// 				return
// }

// func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
// 	token, err := jwt.ParseWithClaims(
// 		signedToken,
// 		&SignedDetails{},
// 		func(t *jwt.Token) (interface{}, error) {
// 			return []byte(SECRET_KEY), nil
// 		},
// 	)

// 	claims, ok := token.Claims.(*SignedDetails)
// 	if !ok {
// 		msg = fmt.Sprintf("token is expired")
// 		msg = err.Error()
// 		return
// 	}

// 	if claims.ExpiresAt <time.Now().Local().Unix() {
// 		msg = fmt.Sprintf("token is expired")
// 		msg = err.Error()
// 			return
// 	}
// 	return claims, msg
// }