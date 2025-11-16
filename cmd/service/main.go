package main

import (
	"context"
	"fmt"
	"go_avito_tech/internal/gateways/http"
	"go_avito_tech/internal/logger"
	"go_avito_tech/internal/repository/db"
	"go_avito_tech/internal/repository/postgres"
	"os"
	//"os/signal"
	//"syscall"
	"time"
)

//	func main() {
//		logger.Init()
//		defer logger.Sync()
//		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//		defer cancel()
//		if err := db.InitDB(ctx); err != nil {
//			panic(fmt.Errorf("failed to init db: %w", err))
//		}
//		defer db.ClosePool()
//		pool := db.GetPool()
//		usersRepo := postgres.NewUserRepository(pool)
//		teamsRepo := postgres.NewTeamRepository(pool)
//		pullRequestsRepo := postgres.NewPullRequestRepository(pool)
//		statsRepo := postgres.NewPgStatsRepository(pool)
//		config := http.Config{
//			Host: getEnv("HOST", "0.0.0.0"),
//			Port: uint16(getEnvInt("PORT", 8080)),
//		}
//		useCases := http.UseCases{
//			Users:  usersRepo,
//			Teams:  teamsRepo,
//			PullRs: pullRequestsRepo,
//			Stats:  statsRepo,
//		}
//		server := http.NewServer(config, useCases)
//		go func() {
//			fmt.Printf("Starting server at %s:%d\n", config.Host, config.Port)
//			if err := server.Run(); err != nil {
//				panic(fmt.Errorf("server failed: %w", err))
//			}
//		}()
//		stop := make(chan os.Signal, 1)
//		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
//		<-stop
//		fmt.Println("Shutting down...")
//	}
func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		var v int
		_, err := fmt.Sscan(val, &v)
		if err == nil {
			return v
		}
	}
	return defaultValue
}

func main() {
	logger.Init()
	defer logger.Sync()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := db.InitDB(ctx)
		cancel()
		if err == nil {
			break
		}
		fmt.Println("DB not ready, retrying:", err)
		time.Sleep(1 * time.Second)
	}
	defer db.ClosePool()
	pool := db.GetPool()
	usersRepo := postgres.NewUserRepository(pool)
	teamsRepo := postgres.NewTeamRepository(pool)
	pullRequestsRepo := postgres.NewPullRequestRepository(pool)
	statsRepo := postgres.NewPgStatsRepository(pool)
	config := http.Config{
		Host: getEnv("HOST", "0.0.0.0"),
		Port: uint16(getEnvInt("SERVER_PORT", 8080)), // FIX
	}
	useCases := http.UseCases{
		Users:  usersRepo,
		Teams:  teamsRepo,
		PullRs: pullRequestsRepo,
		Stats:  statsRepo,
	}
	server := http.NewServer(config, useCases)
	fmt.Printf("Starting server at %s:%d\n", config.Host, config.Port)
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("server failed: %w", err))
	}
}
