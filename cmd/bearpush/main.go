package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/frixuu/bearpush"
	"github.com/frixuu/bearpush/config/templates"
	"github.com/frixuu/bearpush/internal/util"
	"github.com/frixuu/bearpush/server"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {

	logger := CreateLogger()
	zap.ReplaceGlobals(logger.Desugar())
	defer logger.Sync()

	app := &cli.App{
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config-dir",
				Aliases: []string{"c"},
				Usage:   "Path to a directory containing configuration of this app.",
				Value:   bearpush.DefaultConfigDir,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "product",
				Aliases: []string{"p"},
				Usage:   "Options for product manipulation.",
				Subcommands: []*cli.Command{
					{
						Name:  "new",
						Usage: "Creates a new product template.",
						Action: func(c *cli.Context) error {

							productName := c.Args().First()
							if productName == "" {
								fmt.Println("Product name not specified!")
							}

							config, err := bearpush.LoadConfig(c.String("config-dir"))
							if err != nil {
								fmt.Println("Cannot load config")
								return err
							}

							dir := filepath.Join(config.Path, "products")
							err = os.MkdirAll(dir, 0740)
							if err != nil && !os.IsExist(err) {
								return err
							}

							productPath := filepath.Join(dir, productName+".yml")
							_, err = os.Stat(productPath)
							if err == nil || !os.IsNotExist(err) {
								fmt.Printf("A product named %s already exists.\n", productName)
								return os.ErrExist
							}

							file, err := os.Create(productPath)
							defer file.Close()
							if err != nil {
								fmt.Printf("Cannot open file %s for writing\n", productPath)
								return err
							}

							_, err = file.WriteString(templates.GenerateProductFile(productName))
							if err != nil {
								fmt.Printf("An error occurred while writing to file %s\n", productPath)
								return err
							}

							fmt.Printf("Configuration for new product %s scaffolded successfully.\n", productName)
							fmt.Printf("It has been saved in %s.\n", productPath)
							return nil
						},
					},
				},
			},
		},
		Action: func(c *cli.Context) error {

			config, err := bearpush.LoadConfig(c.String("config-dir"))
			if err != nil {
				logger.Error("Cannot load config")
				return err
			}

			logger.Infof("Config directory: %s", config.Path)
			app, err := bearpush.ContextFromConfig(config)
			if err != nil {
				logger.Errorf("Cannot create app context: %s\n", err)
				return err
			}

			app.Logger = logger
			logger.Infof("Loaded info about %d products", len(app.Products))
			for name, p := range app.Products {
				logger.Infof("  - '%s', with '%v' strategy", name, p.TokenSettings.Strategy)
			}

			gin.SetMode(gin.ReleaseMode)
			gin.DefaultWriter = io.Discard
			gin.DefaultErrorWriter = io.Discard

			router := gin.Default()
			router.MaxMultipartMemory = 8 << 20 // 8 MiB

			// Log all requests with Zap
			router.Use(ginzap.Ginzap(logger.Desugar(), time.RFC3339, false))
			// Log all panics with stacktraces
			router.Use(ginzap.RecoveryWithZap(logger.Desugar(), true))

			v1 := router.Group("/v1")
			{
				v1.POST(
					"/upload/:product",
					server.ValidateToken(app),
					func(c *gin.Context) {
						handleArtifactUpload(c.Param("product"), app, c)
					})
			}

			port := server.DeterminePort()
			logger.Info("The server will bind to ", port)

			srv := &http.Server{
				Addr:    port,
				Handler: router,
			}

			// Listen in a goroutine
			go server.Start(srv, logger.Desugar())

			util.WaitForInterrupt()
			logger.Info("Shutting down the server")
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				logger.Fatal("Server forced to shutdown: ", err)
			}

			return nil
		},
	}

	app.Setup()
	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}
