package main

import (
	"cryptoGridProjectDemo/activityLogApp"
	"cryptoGridProjectDemo/internal/pkg/log"
	"cryptoGridProjectDemo/internal/v1/router"
	"cryptoGridProjectDemo/leverageApp"
	"cryptoGridProjectDemo/mainApp"
	"cryptoGridProjectDemo/orderApp"
	"cryptoGridProjectDemo/quoteApp"
	"cryptoGridProjectDemo/redisLib"
	"cryptoGridProjectDemo/wsLib"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"net/http"
	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func Init(rdb *redis.Client) {
	redisService := redisLib.RedisService{Rdb: rdb}
	err := redisService.DelTaskHearthBeat("TaskQuote")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	const config string = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s"
	pgSources := fmt.Sprintf(config,
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DATABASE"),
		os.Getenv("PG_SSLMODE"),
		os.Getenv("PG_TIME_ZONE"),
	)

	db, dbErr := gorm.Open(postgres.Open(pgSources), &gorm.Config{})

	if dbErr != nil {
		panic("failed to connect database")
	}

	//迁移 schema
	db.AutoMigrate(&activityLogApp.ActivityLog{})
	db.AutoMigrate(&leverageApp.LeverageSymbol{})
	db.AutoMigrate(&orderApp.OrderSymbol{})

	hub := wsLib.NewHub()
	go hub.Run()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	gin.SetMode(gin.DebugMode)
	Init(rdb)

	route := router.Default()
	//route.Use(cors.Default())
	route = activityLogApp.GetRoute(route, db, hub)
	route = quoteApp.GetRoute(route, db, hub, rdb)
	route = leverageApp.GetRoute(route, db, hub)
	route = mainApp.GetRoute(route, db, hub, rdb)
	route = orderApp.GetRoute(route, db, hub)

	go heartBeat(hub, rdb, db)

	fmt.Println("running local")
	log.Fatal(http.ListenAndServe(":8000", route))

}
