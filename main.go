package main

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	db, err := NewDB()
	if err != nil {
		logger.Error("failed to create db", "error", err)
		return
	}

	e := NewEchoServer(db)
	defer db.Close()

	// アプリケーションの終了時に処理途中のものが強制終了しないように graceful shutdown を行う
	// ref: https://echo.labstack.com/docs/cookbook/graceful-shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func NewDB() (*sqlx.DB, error) {
	// TODO: ユーザー名、パスワード、タイムアウトなどの設定値は環境変数から取得する
	dsn := "user:password@tcp(127.0.0.1:3306)/shiharai_com_db?charset=utf8mb4&interpolateParams=false&readTimeout=0&timeout=0&tls=false&writeTimeout=0"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Join(errors.New("failed to open db"), err)
	}
	// connection pool の設定
	// TODO: load testing を行い適切な値を設定する
	// TODO: 設定値を環境変数から取得する
	// ref: https://github.com/go-sql-driver/mysql?tab=readme-ov-file#important-settings
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}

func NewEchoServer(db *sqlx.DB) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	return e
}
