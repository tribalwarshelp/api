package main

import (
	"context"
	"github.com/Kichiyaki/appmode"
	"github.com/Kichiyaki/goutil/envutil"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-contrib/cors"

	servermaphttpdelivery "github.com/tribalwarshelp/api/servermap/delivery/http"

	httpdelivery "github.com/tribalwarshelp/api/graphql/delivery/http"
	"github.com/tribalwarshelp/api/graphql/resolvers"

	"github.com/tribalwarshelp/api/graphql/dataloader"

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

	"github.com/Kichiyaki/ginlogrus"
	"github.com/Kichiyaki/go-pg-logrus-query-logger/v10"
	"github.com/gin-gonic/gin"
)

func init() {
	os.Setenv("TZ", "UTC")

	if appmode.Equals(appmode.DevelopmentMode) {
		godotenv.Load(".env.local")
	}

	setupLogger()
}

func main() {
	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		Addr:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		PoolSize: envutil.GetenvInt("DB_POOL_SIZE"),
	})
	defer func() {
		if err := db.Close(); err != nil {
			logrus.Fatalln("Couldn't close the db connections:", err)
		}
	}()
	if strings.ToUpper(os.Getenv("LOG_DB_QUERIES")) == "TRUE" {
		db.AddQueryHook(querylogger.Logger{
			Log:            logrus.NewEntry(logrus.StandardLogger()),
			MaxQueryLength: 5000,
		})
	}

	versionRepo, err := versionrepo.NewPGRepository(db)
	if err != nil {
		logrus.Fatal(err)
	}
	serverRepo, err := serverrepo.NewPGRepository(db)
	if err != nil {
		logrus.Fatal(err)
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

	router := gin.New()
	router.Use(gin.Recovery())
	if !envutil.GetenvBool("DISABLE_ACCESS_LOG") {
		router.Use(ginlogrus.Logger(logrus.StandardLogger()))
	}
	if appmode.Equals(appmode.DevelopmentMode) {
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
	graphql.Use(
		middleware.DataLoadersToContext(
			middleware.DataLoadersToContextConfig{
				ServerRepo: serverRepo,
			},
			dataloader.Config{
				PlayerRepo:  playerRepo,
				TribeRepo:   tribeRepo,
				VillageRepo: villageRepo,
				VersionRepo: versionRepo,
			},
		),
	)
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
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()
	logrus.Info("Server is listening on the port 8080")

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalln("Couldn't shutdown the server", err)
	}
	logrus.Println("Server exiting")
}

func setupLogger() {
	if appmode.Equals(appmode.DevelopmentMode) {
		logrus.SetLevel(logrus.DebugLevel)
	}

	timestampFormat := "2006-01-02 15:04:05"
	if appmode.Equals(appmode.ProductionMode) {
		customFormatter := new(logrus.JSONFormatter)
		customFormatter.TimestampFormat = timestampFormat
		logrus.SetFormatter(customFormatter)
	} else {
		customFormatter := new(logrus.TextFormatter)
		customFormatter.TimestampFormat = timestampFormat
		customFormatter.FullTimestamp = true
		logrus.SetFormatter(customFormatter)
	}
}
