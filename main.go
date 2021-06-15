package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/Frixuu/BearPush/config"
	"github.com/Frixuu/BearPush/config/templates"
	"github.com/Frixuu/BearPush/util"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

func main() {

	logger := CreateLogger()
	defer logger.Sync()

	app := &cli.App{
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config-dir",
				Aliases: []string{"c"},
				Usage:   "Path to a directory containing configuration of this app.",
				Value:   config.DefaultConfigDir,
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

							config, err := config.Load(c.String("config-dir"))
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

			config, err := config.Load(c.String("config-dir"))
			if err != nil {
				logger.Error("Cannot load config")
				return err
			}

			logger.Infof("Config directory: %s", config.Path)
			appContext, err := ContextFromConfig(config)
			if err != nil {
				logger.Errorf("Cannot create app context: %s\n", err)
				return err
			}

			for name, p := range appContext.Products {
				logger.Infof("Loaded product %s, token strategy %v", name, p.TokenSettings.Strategy)
			}

			router := gin.Default()
			router.MaxMultipartMemory = 8 << 20 // 8 MiB

			router.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})

			v1 := router.Group("/v1")
			{
				v1.POST("/upload/:product", ValidateToken(appContext), func(c *gin.Context) {

					product := c.Param("product")
					p, ok := appContext.Products[product]
					if !ok {
						c.JSON(http.StatusBadRequest, gin.H{
							"error":   4,
							"message": "Resource does not exist.",
						})
						return
					}

					file, err := c.FormFile("artifact")
					if err != nil {
						log.Println(err)
						c.String(http.StatusBadRequest, fmt.Sprintf("Error while uploading: %s", err))
						return
					}

					tempDir, err := ioutil.TempDir("", "bearpush-")
					if err != nil {
						log.Println(err)
						c.String(http.StatusInternalServerError,
							"Could not create a temporary directory for artifact. Check logs for details.")
						return
					}
					defer util.TryRemoveDir(tempDir)

					artifactPath := path.Join(tempDir, "artifact")
					err = c.SaveUploadedFile(file, artifactPath)
					if err != nil {
						log.Printf("Cannot save artifact: %s", err)
						c.String(http.StatusInternalServerError,
							"Could not save the uploaded artifact. Check logs for details.")
						return
					}

					if p.Script != "" {
						cmd := exec.Command(p.Script)
						cmd.Env = append(os.Environ(),
							fmt.Sprintf("ARTIFACT_PATH=%s", artifactPath),
						)

						stdoutPipe, err := cmd.StdoutPipe()
						if err != nil {
							logger.Errorf("Cannot grab stdout pipe: %s\n", err)
						}

						stderrPipe, err := cmd.StderrPipe()
						if err != nil {
							logger.Errorf("Cannot grab stderr pipe: %s\n", err)
						}

						if err := cmd.Start(); err != nil {
							logger.Errorf("Cannot start: %s\n", err)
						}

						_, err = io.ReadAll(stdoutPipe)
						if err != nil {
							logger.Errorf("Cannot read stdout: %s\n", err)
						}

						_, err = io.ReadAll(stderrPipe)
						if err != nil {
							logger.Errorf("Cannot read stderr: %s\n", err)
						}

						if err := cmd.Wait(); err != nil {
							logger.Errorf("Command failed: %s\n", err)
							c.JSON(http.StatusUnprocessableEntity, gin.H{
								"error":   8,
								"message": "Pipeline associated with resource errored.",
							})
							return
						}
					}

					c.String(http.StatusOK,
						fmt.Sprintf("Artifact for product %s processed successfully.", product))
				})
			}

			return router.Run()
		},
	}

	app.Setup()
	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}
