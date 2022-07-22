package controller

import(
	"net/http"
	"log"
	"fmt"
	"os"
    "mime/multipart"
	"bytes"
	"path/filepath"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	SECRET_ID = "AKIAYMMEVVXKJ3ZNYL6E"
	SECRET_KEY = "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"
	BUCKET_NAME= "tcxdv9cjpydfp"
)

func UploadFileToS3(s *session.Session,file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
    size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	tempFileName := "filess/" + primitive.NewObjectID().Hex() + filepath.Ext(fileHeader.Filename)

	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(tempFileName),
		Body:   bytes.NewReader(buffer),
	})
	if err != nil {
		return "", err
	}

	return tempFileName, nil
}

func UploadFile() httprouter.Handle {
	return func(w http.ResponseWriter,r *http.Request,  p httprouter.Params) {
		switch r.Method {
		case http.MethodPost:
			maxSize := int64(1024000)    
			err := r.ParseMultipartForm(maxSize)
			if err != nil {
				log.Println(err)
				fmt.Println("Image too large")
				return
			}
	
			file, fileHeader, err := r.FormFile("file")
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Println(err)
				fmt.Fprintf(w, "Could not get uploaded file")
				return
			}
			defer file.Close()
	
	  s, err := session.NewSession(&aws.Config{
	  Region: aws.String("us-west-2"),
	  Credentials: credentials.NewStaticCredentials(
		  SECRET_ID,
		  SECRET_KEY,
		  ""),  //token
	  })
	  if  err != nil {
		  w.WriteHeader(http.StatusNotFound)
		  fmt.Println("Invalid AWS Credentials")
	  }
	
	  fileName, err := UploadFileToS3(s, file, fileHeader)
	  if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("Could not upload file....", err)
	  }
	  fmt.Println("true....")
	  fmt.Fprintf(w, "Image uploaded successfully: %v", fileName)
      w.WriteHeader(http.StatusOK)
      json.NewEncoder(w).Encode(fileName)
	
		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
			
		}
	}
}

func DownloadFile() httprouter.Handle{
	return func (w http.ResponseWriter,r *http.Request,  p httprouter.Params) {
		switch r.Method {
		case http.MethodGet:
			path := "filess"
			item := "filess/6247f57e9f99ea03a53746e4.jpeg"
		
			file, err := os.Create(path)
			if err != nil {
				fmt.Println("Path Error.....")
				fmt.Println("Unable to open file %q, %v", item, err)
			}
			fmt.Println("true.....")
			defer file.Close()
		
			sess, _ := session.NewSession(&aws.Config{
				Region: aws.String("us-west-2"),
				Credentials: credentials.NewStaticCredentials(
					SECRET_ID,
					SECRET_KEY,
					""),  //token
			})
		
			downloader := s3manager.NewDownloader(sess)
		
			numBytes, err := downloader.Download(file,
				&s3.GetObjectInput{
					Bucket: aws.String(BUCKET_NAME),
					Key:    aws.String(item),
				})
			if err != nil {
				fmt.Println("ERROR::::",  err)
			}
		
			fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(file.Name())

		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
		}
	}
}

func DeleteFile() httprouter.Handle{
	return func (w http.ResponseWriter,r *http.Request,  p httprouter.Params) {
		switch r.Method {
		case http.MethodDelete:
			sess, _ := session.NewSession(&aws.Config{
				Region: aws.String("us-west-2"),
				Credentials: credentials.NewStaticCredentials(
					SECRET_ID,
					SECRET_KEY,
					""),  //token
			})
			item := "filess/6247f57e9f99ea03a53746e4.jpeg"

			svc := s3.New(sess)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key: aws.String(item),
	}

	result, err := svc.DeleteObject(input)
	if err != nil {
		fmt.Println("ERROR::::",   err)
	}

	log.Printf("Result: %+v\n", result)
	w.WriteHeader(http.StatusOK)

		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
		}
	}
}


