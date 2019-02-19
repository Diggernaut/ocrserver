package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/otiai10/marmoset"

	"github.com/natefinch/lumberjack"
	"github.com/Diggernaut/ocrserver/controllers"
	"github.com/Diggernaut/ocrserver/filters"
)

var (
	cfg    *viper.Viper
	apikey string
	port   string
	ip     string
)

func init() {
	// SET UP LOGGER
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/var/log/ocrserver.log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,   // max files
		MaxAge:     7,   // days
	})
	//log.SetOutput(os.Stdout)

	// READING CONFIG
	cfg = viper.New()
	cfg.SetConfigName("config")
	cfg.AddConfigPath("./")
	err := cfg.ReadInConfig()
	if err != nil {
		log.Fatalf("Error: cannot read config. Reason: %v\n", err)
	}
	apikey = cfg.GetString("apikey")
	ip = cfg.GetString("ocr_bind_ip")
	port = cfg.GetString("ocr_bind_port")
}


func main() {

	//marmoset.LoadViews("./app/views")

	r := marmoset.NewRouter()
	// API
	r.GET("/status", controllers.Status)
	r.POST("/base64", controllers.Base64)
	f := filters.AuthFilter{Apikey: apikey}
	r.Apply(f)
	//r.POST("/file", controllers.FileUpload)
	// Sample Page
	//r.GET("/", controllers.Index)
	//r.Static("/assets", "./app/assets")

	log.Println("OCR web server started")
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(ip + ":" + port, r); err != nil {
		logger.Println(err)
	}
}
