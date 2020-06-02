package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/tribalwarshelp/shared/mode"

	httpdelivery "github.com/tribalwarshelp/api/graphql/delivery/http"
	"github.com/tribalwarshelp/api/graphql/resolvers"

	"github.com/tribalwarshelp/api/graphql/dataloaders"

	langversionrepo "github.com/tribalwarshelp/api/langversion/repository"
	langversionucase "github.com/tribalwarshelp/api/langversion/usecase"
	"github.com/tribalwarshelp/api/middleware"
	playerrepo "github.com/tribalwarshelp/api/player/repository"
	playerucase "github.com/tribalwarshelp/api/player/usecase"
	serverrepo "github.com/tribalwarshelp/api/server/repository"
	serverucase "github.com/tribalwarshelp/api/server/usecase"
	triberepo "github.com/tribalwarshelp/api/tribe/repository"
	tribeucase "github.com/tribalwarshelp/api/tribe/usecase"
	villagerepo "github.com/tribalwarshelp/api/village/repository"
	villageucase "github.com/tribalwarshelp/api/village/usecase"

	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func init() {
	os.Setenv("TZ", "UTC")

	if mode.Get() == mode.DevelopmentMode {
		godotenv.Load(".env.development")
	}
}

func main() {
	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		Addr:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
	})
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("Database disconnecting:", err)
		}
	}()

	langversionRepo, err := langversionrepo.NewPGRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	serverRepo, err := serverrepo.NewPGRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	tribeRepo := triberepo.NewPGRepository(db)
	playerRepo := playerrepo.NewPGRepository(db)
	villageRepo := villagerepo.NewPGRepository(db)

	router := gin.Default()
	v1 := router.Group("")
	v1.Use(middleware.DataLoadersToContext(serverRepo, dataloaders.Config{
		PlayerRepo: playerRepo,
		TribeRepo:  tribeRepo,
	}))
	httpdelivery.Attach(v1, &resolvers.Resolver{
		LangVersionUcase: langversionucase.New(langversionRepo),
		ServerUcase:      serverucase.New(serverRepo),
		TribeUcase:       tribeucase.New(tribeRepo),
		PlayerUcase:      playerucase.New(playerRepo),
		VillageUcase:     villageucase.New(villageRepo),
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
