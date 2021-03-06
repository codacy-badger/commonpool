package main

import (
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/config"
	_ "github.com/commonpool/backend/docs"
	"github.com/commonpool/backend/handler"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/router"
	"github.com/commonpool/backend/store"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
)

var (
	d  *gorm.DB
	rs resource.Store
	as auth.Store
	e  *echo.Echo
)

// @title resources api
// @version 1.0
// @description resources api
// @termsOfService http://swagger.io/terms
// @host 127.0.0.1:8585
// @basePath /api/v1
func main() {

	appConfig, err := config.GetAppConfig(os.LookupEnv, ioutil.ReadFile)
	if err != nil {
		panic(err)
	}

	r := router.NewRouter()

	r.GET("/swagger/*", echoSwagger.WrapHandler)

	db := getDb(appConfig)
	store.AutoMigrate(db)
	rs = store.NewResourceStore(db)
	as = store.NewAuthStore(db)

	v1 := r.Group("/api/v1")
	authorization := auth.NewAuth(v1, appConfig, "/api/v1", as)

	h := handler.NewHandler(rs, as, authorization)

	h.Register(v1)

	r.Logger.Info("Secure Cookies", appConfig.SecureCookies)

	r.Logger.Fatal(r.Start("0.0.0.0:8585"))
}

func getDb(appConfig *config.AppConfig) *gorm.DB {
	cs := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", appConfig.DbHost, appConfig.DbUsername, appConfig.DbPassword, appConfig.DbName, appConfig.DbPort)
	db, err := gorm.Open(postgres.Open(cs), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
