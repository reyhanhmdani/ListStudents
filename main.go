package main

import (
	"context"
	"ginDatabaseMhs/database"
	"ginDatabaseMhs/router"
	"ginDatabaseMhs/service"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func setupLogOutput() {
	f, _ := os.Create("gin-log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {

	setupLogOutput()

	ctx := context.Background()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	// ENV
	loadEnv()

	// pr
	// INITAL DATABASE
	db, err := database.Databaseinit(ctx)
	if err != nil {
		return
	}

	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Error running schema migration %v", err)
	}

	// Initialize AWS S3 client
	s3Config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.NewFromConfig(s3Config)

	// initial repo
	todoRepo := database.NewMahasiswaRepository(db, s3Client)
	todoService := service.NewMahasiswaService(todoRepo)
	routeBuilder := router.NewRouteBuilder(todoService)
	routeInit := routeBuilder.RouteInit()
	err = routeInit.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Mengambil nilai variabel lingkungan
	dbHost := os.Getenv("DB_HOST")
	dbRootPassword := os.Getenv("DB_PASS")
	dbDatabase := os.Getenv("DB_NAME")

	// Contoh penggunaan nilai variabel lingkungan
	log.Printf("DB Host: %s", dbHost)
	log.Printf("DB Root Password: %s", dbRootPassword)
	log.Printf("DB Database: %s", dbDatabase)
}
