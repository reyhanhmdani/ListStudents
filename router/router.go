package router

import (
	"ginDatabaseMhs/middleware"
	mhsService "ginDatabaseMhs/service"
	"github.com/gin-gonic/gin"
)

type RouteBuilder struct {
	dataService *mhsService.Handler
}

func NewRouteBuilder(dataService *mhsService.Handler) *RouteBuilder {
	return &RouteBuilder{dataService: dataService}
}

func (rb *RouteBuilder) RouteInit() *gin.Engine {

	r := gin.New()
	r.Use(middleware.RecoveryMiddleware(), middleware.Logger())
	//r.Use(gin.Recovery(), middleware.Logger(), middleware.BasicAuth())

	auth := r.Group("/", middleware.Authmiddleware())
	{
		auth.GET("/admin/viewUsers", rb.dataService.ViewAllUsers)
		auth.DELETE("admin/users/:user_id", rb.dataService.DeleteUser)
		auth.GET("/manage-data", rb.dataService.HandlerGetAll)
		auth.GET("/access", rb.dataService.Access)
		auth.POST("/create-form", rb.dataService.HandlerCreate)
		auth.GET("/manage-data/daftarMahasiswa/:id", rb.dataService.HandlerGetByID)
		auth.PUT("/manage-data/daftarMahasiswa/:id", rb.dataService.HandlerUpdate)
		auth.DELETE("/manage-data/daftarMahasiswa/:id", rb.dataService.HandlerDelete)
		auth.POST("/uploadS3/:id", rb.dataService.UploadFileS3AtchHandler)
		auth.POST("/uploadLocal/:id", rb.dataService.UploadLocalAtchHandler)
		auth.GET("/list-Search", rb.dataService.SearchHandler)
	}

	r.POST("/uploadBuckets", rb.dataService.UploadFileS3BucketsHandler)
	r.POST("/register", rb.dataService.Register)
	r.POST("/login", rb.dataService.Login)
	return r
}
