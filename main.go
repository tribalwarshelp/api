package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	servermaphttpdelivery "github.com/tribalwarshelp/api/servermap/delivery/http"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/mode"

	httpdelivery "github.com/tribalwarshelp/api/graphql/delivery/http"
	"github.com/tribalwarshelp/api/graphql/resolvers"

	"github.com/tribalwarshelp/api/graphql/dataloaders"

	ennoblementrepo "github.com/tribalwarshelp/api/ennoblement/repository"
	ennoblementucase "github.com/tribalwarshelp/api/ennoblement/usecase"
	langversionrepo "github.com/tribalwarshelp/api/langversion/repository"
	langversionucase "github.com/tribalwarshelp/api/langversion/usecase"
	liveennoblementrepo "github.com/tribalwarshelp/api/liveennoblement/repository"
	liveennoblementucase "github.com/tribalwarshelp/api/liveennoblement/usecase"
	"github.com/tribalwarshelp/api/middleware"
	playerrepo "github.com/tribalwarshelp/api/player/repository"
	playerucase "github.com/tribalwarshelp/api/player/usecase"
	playerhistoryrepo "github.com/tribalwarshelp/api/playerhistory/repository"
	playerhistoryucase "github.com/tribalwarshelp/api/playerhistory/usecase"
	serverrepo "github.com/tribalwarshelp/api/server/repository"
	serverucase "github.com/tribalwarshelp/api/server/usecase"
	servermapucase "github.com/tribalwarshelp/api/servermap/usecase"
	serverstatsrepo "github.com/tribalwarshelp/api/serverstats/repository"
	serverstatsucase "github.com/tribalwarshelp/api/serverstats/usecase"
	triberepo "github.com/tribalwarshelp/api/tribe/repository"
	tribeucase "github.com/tribalwarshelp/api/tribe/usecase"
	tribehistoryrepo "github.com/tribalwarshelp/api/tribehistory/repository"
	tribehistoryucase "github.com/tribalwarshelp/api/tribehistory/usecase"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal(errors.Wrap(err, "cannot connect to redis"))
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatal(err)
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
	ennoblementRepo := ennoblementrepo.NewPGRepository(db)
	tribehistoryRepo := tribehistoryrepo.NewPGRepository(db)
	playerhistoryRepo := playerhistoryrepo.NewPGRepository(db)
	serverstatsRepo := serverstatsrepo.NewPGRepository(db)
	liveennoblementRepo := liveennoblementrepo.NewPGRepository(db, redisClient)

	serverUcase := serverucase.New(serverRepo)

	router := gin.Default()
	rest := router.Group("")
	servermaphttpdelivery.Attach(servermaphttpdelivery.Config{
		RouterGroup:   rest,
		MapUsecase:    servermapucase.New(villageRepo),
		ServerUsecase: serverUcase,
	})
	graphql := router.Group("")
	graphql.Use(middleware.DataLoadersToContext(serverRepo, dataloaders.Config{
		PlayerRepo:      playerRepo,
		TribeRepo:       tribeRepo,
		VillageRepo:     villageRepo,
		LangVersionRepo: langversionRepo,
	}))
	httpdelivery.Attach(httpdelivery.Config{
		RouterGroup: graphql,
		Resolver: &resolvers.Resolver{
			LangVersionUcase:     langversionucase.New(langversionRepo),
			ServerUcase:          serverUcase,
			TribeUcase:           tribeucase.New(tribeRepo),
			PlayerUcase:          playerucase.New(playerRepo),
			VillageUcase:         villageucase.New(villageRepo),
			EnnoblementUcase:     ennoblementucase.New(ennoblementRepo),
			LiveEnnoblementUcase: liveennoblementucase.New(liveennoblementRepo),
			TribeHistoryUcase:    tribehistoryucase.New(tribehistoryRepo),
			PlayerHistoryUcase:   playerhistoryucase.New(playerhistoryRepo),
			ServerStatsUcase:     serverstatsucase.New(serverstatsRepo),
		},
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
