package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/frixuu/bearpush"
	"github.com/gin-gonic/gin"
)

func handleArtifactUpload(productId string, app *bearpush.Context, c *gin.Context) {

	// Is there a registered product with that id?
	p, ok := app.Products[productId]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   4,
			"message": "Resource does not exist.",
		})
		return
	}

	// Is there a file attached?
	file, err := c.FormFile("artifact")
	if err != nil {
		app.Logger.Warn(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   5,
			"message": fmt.Sprintf("Error while uploading: %s", err),
		})
		return
	}

	// Create a temporary directory for storing the uploaded file
	tempDir, err := os.MkdirTemp("", "bearpush-")
	if err != nil {
		app.Logger.Errorf("Could not create a temporary directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   6,
			"message": "Could not create a temporary directory for artifact.",
		})
		return
	}

	// Remember to clean up that directory after we're done with it
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			app.Logger.Warnf("Cannot remove directory %s: %v", tempDir, err)
		}
	}()

	// Save the file to the temp directory
	artifactPath := path.Join(tempDir, "artifact")
	err = c.SaveUploadedFile(file, artifactPath)
	if err != nil {
		app.Logger.Errorf("Cannot save artifact: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   7,
			"message": "Could not save the uploaded artifact.",
		})
		return
	}

	// Run the script associated with this product
	if p.Script != "" {
		cmd := exec.Command(p.Script)
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("ARTIFACT_PATH=%s", artifactPath),
		)

		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			app.Logger.Warnf("Cannot grab stdout pipe: %v", err)
		}

		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			app.Logger.Warnf("Cannot grab stderr pipe: %v", err)
		}

		if err := cmd.Start(); err != nil {
			app.Logger.Warnf("Cannot start the process: %s", err)
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   8,
				"message": "Pipeline associated with this resource could not be started.",
			})
			return
		}

		_, err = io.ReadAll(stdoutPipe)
		if err != nil {
			app.Logger.Errorf("Cannot read stdout: %v", err)
		}

		_, err = io.ReadAll(stderrPipe)
		if err != nil {
			app.Logger.Errorf("Cannot read stderr: %v", err)
		}

		if err := cmd.Wait(); err != nil {
			app.Logger.Warnf("Command failed: %v", err)
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   9,
				"message": "Pipeline associated with resource errored.",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Artifact for product %s processed successfully.", productId),
	})
}
