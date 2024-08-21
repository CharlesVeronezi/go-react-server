package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/CharlesVeronezi/go-react-server.git/internal/api"
	"github.com/CharlesVeronezi/go-react-server.git/internal/store/pgstore/pgstore"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	//if err := godotenv.Load(); err != nil {
	//	panic(err)
	//}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("GS_DATABASE_USER"),
		os.Getenv("GS_DATABASE_PASSWORD"),
		os.Getenv("GS_DATABASE_HOST"),
		os.Getenv("GS_DATABASE_PORT"),
		os.Getenv("GS_DATABASE_NAME"),
	))

	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	handler := api.NewHandler(pgstore.New(pool))

	go func() {
		if err := http.ListenAndServe(":3000", handler); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
