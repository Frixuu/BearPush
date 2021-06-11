package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/Frixuu/BearPush/v2/config"
	"github.com/gin-gonic/gin"
)

func main() {

	cfgDir := flag.String("config-dir", config.DefaultConfigDir,
		"Path to a directory containing this application's configuration files.")

	flag.Parse()

	config := config.Load(*cfgDir)
	log.Printf("Config directory: %s", config.Path)

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("/v1")
	{
		v1.POST("/upload/:product", ValidateToken(), func(c *gin.Context) {

			product := c.Param("product")

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

			err = c.SaveUploadedFile(file, path.Join(tempDir, "artifact"))
			if err != nil {
				log.Printf("Cannot save artifact: %s", err)
				c.String(http.StatusInternalServerError,
					"Could not save the uploaded artifact. Check logs for details.")
				return
			}

			err = os.RemoveAll(tempDir)
			if err != nil {
				log.Printf("Cannot remove temporary directory %s: %s", tempDir, err)
			}

			c.String(http.StatusOK,
				fmt.Sprintf("Artifact for product %s processed successfully.", product))
		})
	}

	router.Run()
}
