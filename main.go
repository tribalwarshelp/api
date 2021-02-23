package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/go-pg/pg/extra/pgdebug"

	"github.com/gin-contrib/cors"
	servermaphttpdelivery "github.com/tribalwarshelp/api/servermap/delivery/http"

	"github.com/tribalwarshelp/shared/mode"

	httpdelivery "github.com/tribalwarshelp/api/graphql/delivery/http"
	"github.com/tribalwarshelp/api/graphql/resolvers"

	"github.com/tribalwarshelp/api/graphql/dataloaders"

	dailyplayerstatsrepo "github.com/tribalwarshelp/api/dailyplayerstats/repository"
	dailyplayerstatsucase "github.com/tribalwarshelp/api/dailyplayerstats/usecase"
	dailytribestatsrepo "github.com/tribalwarshelp/api/dailytribestats/repository"
	dailytribestatsucase "github.com/tribalwarshelp/api/dailytribestats/usecase"
	ennoblementrepo "github.com/tribalwarshelp/api/ennoblement/repository"
	ennoblementucase "github.com/tribalwarshelp/api/ennoblement/usecase"
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
	tribechangerepo "github.com/tribalwarshelp/api/tribechange/repository"
	tribechangeucase "github.com/tribalwarshelp/api/tribechange/usecase"
	tribehistoryrepo "github.com/tribalwarshelp/api/tribehistory/repository"
	tribehistoryucase "github.com/tribalwarshelp/api/tribehistory/usecase"
	versionrepo "github.com/tribalwarshelp/api/version/repository"
	versionucase "github.com/tribalwarshelp/api/version/usecase"
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
	log.Printf(os.Getenv("GIN_MODE"))
}

func main() {
	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		Addr:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		PoolSize: mustParseEnvToInt("DB_POOL_SIZE"),
	})
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("Database disconnecting:", err)
		}
	}()
	if strings.ToUpper(os.Getenv("LOG_DB_QUERIES")) == "TRUE" {
		db.AddQueryHook(pgdebug.DebugHook{
			Verbose: true,
		})
	}

	versionRepo, err := versionrepo.NewPGRepository(db)
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
	tribeChangeRepo := tribechangerepo.NewPGRepository(db)
	dailyPlayerStatsRepo := dailyplayerstatsrepo.NewPGRepository(db)
	dailyTribeStatsRepo := dailytribestatsrepo.NewPGRepository(db)

	serverUcase := serverucase.New(serverRepo)

	router := gin.Default()
	if mode.Get() == mode.DevelopmentMode {
		router.Use(cors.New(cors.Config{
			AllowOriginFunc: func(string) bool {
				return true
			},
			AllowCredentials: true,
			ExposeHeaders:    []string{"X-Access-Token", "X-Refresh-Token"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			AllowWebSockets:  true,
		}))
	}
	rest := router.Group("")
	servermaphttpdelivery.Attach(servermaphttpdelivery.Config{
		RouterGroup:   rest,
		MapUsecase:    servermapucase.New(villageRepo),
		ServerUsecase: serverUcase,
	})
	graphql := router.Group("")
	graphql.Use(middleware.DataLoadersToContext(middleware.DataLoadersToContextConfig{
		ServerRepo: serverRepo,
	},
		dataloaders.Config{
			PlayerRepo:  playerRepo,
			TribeRepo:   tribeRepo,
			VillageRepo: villageRepo,
			VersionRepo: versionRepo,
		}))
	graphql.Use(middleware.LimitWhitelist(middleware.LimitWhitelistConfig{
		IPAddresses: strings.Split(os.Getenv("LIMIT_WHITELIST"), ","),
	}))
	httpdelivery.Attach(httpdelivery.Config{
		RouterGroup: graphql,
		Resolver: &resolvers.Resolver{
			VersionUcase:          versionucase.New(versionRepo),
			ServerUcase:           serverUcase,
			TribeUcase:            tribeucase.New(tribeRepo),
			PlayerUcase:           playerucase.New(playerRepo),
			VillageUcase:          villageucase.New(villageRepo),
			EnnoblementUcase:      ennoblementucase.New(ennoblementRepo),
			TribeHistoryUcase:     tribehistoryucase.New(tribehistoryRepo),
			PlayerHistoryUcase:    playerhistoryucase.New(playerhistoryRepo),
			ServerStatsUcase:      serverstatsucase.New(serverstatsRepo),
			TribeChangeUcase:      tribechangeucase.New(tribeChangeRepo),
			DailyPlayerStatsUcase: dailyplayerstatsucase.New(dailyPlayerStatsRepo),
			DailyTribeStatsUcase:  dailytribestatsucase.New(dailyTribeStatsRepo),
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

func mustParseEnvToInt(key string) int {
	str := os.Getenv(key)
	if str == "" {
		return 0
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}
