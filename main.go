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

	"github.com/Frixuu/BearPush/v2/config"
	"github.com/Frixuu/BearPush/v2/config/templates"
	"github.com/Frixuu/BearPush/v2/util"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

func main() {

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
								fmt.Printf("Error occured while writing to file %s\n", productPath)
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
				log.Println("Cannot load config")
				return err
			}

			log.Printf("Config directory: %s", config.Path)
			appContext, err := ContextFromConfig(config)
			if err != nil {
				log.Printf("Cannot create app context: %s\n", err)
				return err
			}

			for name, p := range appContext.Products {
				log.Printf("Loaded product %s, token strategy %v", name, p.TokenSettings.Strategy)
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
							log.Println(err)
						}

						stderrPipe, err := cmd.StderrPipe()
						if err != nil {
							log.Println(err)
						}

						if err := cmd.Start(); err != nil {
							log.Println(err)
						}

						_, err = io.ReadAll(stdoutPipe)
						if err != nil {
							log.Println(err)
						}

						_, err = io.ReadAll(stderrPipe)
						if err != nil {
							log.Println(err)
						}

						if err := cmd.Wait(); err != nil {
							log.Println(err)
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
		log.Fatal(err)
	}
}
