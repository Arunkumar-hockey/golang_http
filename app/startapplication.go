package app

import(
	"github.com/julienschmidt/httprouter"
	"GolangHTTP/controller"
	 "log"
	 "net/http"
)

func StartApplication() {
	router := httprouter.New()
	router.POST("/create", controller.CreateEmployee())
	router.GET("/employee/:employee_id", controller.GetEmployee())
	router.GET("/allemployee", controller.GetAllEmployee())
	router.GET("/searchemployee", controller.SearchEmployeeByName())
	router.PATCH("/updateemployee/:employee_id", controller.UpdateEmployee())
	router.DELETE("/deleteemployee/:employee_id", controller.DeleteEmployee())
	router.POST("/employee/login", controller.Login())
	router.POST("/checkemployeeexist", controller.CheckEmployeeExistAndSendMail())
	router.POST("/otpverify", controller.VerifyOTP())
	router.POST("/updatepassword", controller.UpdatePassword())
	router.GET("/home", controller.Home())
	router.GET("/signout", controller.SignOut())
    router.POST("/uploadfile", controller.UploadFile())
	router.GET("/downloadfile", controller.DownloadFile())
	router.DELETE("/deletefile", controller.DeleteFile())
	
	log.Fatal(http.ListenAndServe(":8080", router))
}