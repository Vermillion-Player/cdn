package controllers

import (
	"context"
	"cdn/models"
	"cdn/services"
	"time"
	"fmt"
	"mime"
    "path"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func validateUser(c *gin.Context, token string) {
	userCollection := getUserCollection()
	username, err := services.GetUserFromToken(token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error when trying to convert token",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	storedUser := models.User{}
	err = userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&storedUser)
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Incorrect username or password",
		})
		return
	}
}

func ListVideoController(c *gin.Context) {
    token := c.GetHeader("Authorization")
    validateUser(c, token)
    username, _ := services.GetUserFromToken(token)

    userDir := fmt.Sprintf("./uploads/%s", username)
    files, err := os.ReadDir(userDir)
    if err != nil {
        c.JSON(500, gin.H{
            "error": "No se pudo leer la carpeta del usuario",
        })
        return
    }

    var fileNames []string
    for _, file := range files {
        if !file.IsDir() {
            fileNames = append(fileNames, file.Name())
        }
    }

    c.JSON(200, gin.H{
        "videos": fileNames,
    })
}

func NewVideoController(c *gin.Context) {
	token := c.GetHeader("Authorization")
	validateUser(c, token)

	file, err := c.FormFile("video_url")
	if err != nil {
		c.JSON(400, gin.H{
			"message": "File is required",
		})
		return
	}

	username, _ := services.GetUserFromToken(token)

	filePath := fmt.Sprintf("./uploads/%s/%s", username, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to save file",
		})
		return
	}

	c.JSON(201, gin.H{
		"video": filePath,
	})
}

func GetVideoController(c *gin.Context) {
	token := c.GetHeader("Authorization")
	validateUser(c, token)
	username, _ := services.GetUserFromToken(token)
	videoName := c.Param("id")

	filePath := fmt.Sprintf("./uploads/%s/%s", username, videoName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
        c.JSON(404, gin.H{"error": "Video file not found"})
        return
    }
	mimeType := mime.TypeByExtension(path.Ext(videoName))
    if mimeType == "" {
        mimeType = "application/octet-stream"
    }

    c.Header("Content-Type", mimeType)
	c.Header("Content-Disposition", "inline")
    c.File(filePath)
}


func EditVideoController(c *gin.Context) {
    token := c.GetHeader("Authorization")
    validateUser(c, token)
    username, _ := services.GetUserFromToken(token)

    oldName := c.Param("id")
    var req struct {
        NewName string `json:"new_name"`
    }
    if err := c.ShouldBindJSON(&req); err != nil || req.NewName == "" {
        c.JSON(400, gin.H{"error": "Should be send the new name in the input 'new_name'"})
        return
    }

    oldPath := fmt.Sprintf("./uploads/%s/%s", username, oldName)
    newPath := fmt.Sprintf("./uploads/%s/%s", username, req.NewName)

    
    if _, err := os.Stat(oldPath); os.IsNotExist(err) {
        c.JSON(404, gin.H{"error": "Release video does not exist"})
        return
    }

    if err := os.Rename(oldPath, newPath); err != nil {
        c.JSON(500, gin.H{"error": "Cannot rename video file"})
        return
    }

    c.JSON(200, gin.H{
        "message":    "Update video name successfully",
        "old_name":   oldName,
        "new_name":   req.NewName,
        "video_path": newPath,
    })
}

func DeleteVideoController(c *gin.Context) {
    token := c.GetHeader("Authorization")
    validateUser(c, token)
    username, _ := services.GetUserFromToken(token)

    videoName := c.Param("id")
    filePath := fmt.Sprintf("./uploads/%s/%s", username, videoName)

    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        c.JSON(404, gin.H{"error": "El video no existe"})
        return
    }

    if err := os.Remove(filePath); err != nil {
        c.JSON(500, gin.H{"error": "Error when trying to delete video"})
        return
    }

    c.JSON(200, gin.H{
        "message":    "Delete video successfully",
        "video_name": videoName,
    })
}

func SearchVideoController(c *gin.Context) {
    token := c.GetHeader("Authorization")
    validateUser(c, token)
    username, _ := services.GetUserFromToken(token)

    query := c.Param("q")

    userDir := fmt.Sprintf("./uploads/%s", username)
    files, err := os.ReadDir(userDir)
    if err != nil {
        c.JSON(500, gin.H{"error": "Cannot read user directory"})
        return
    }

    var matches []string
    for _, file := range files {
        if !file.IsDir() && strings.Contains(strings.ToLower(file.Name()), strings.ToLower(query)) {
            matches = append(matches, file.Name())
        }
    }

    c.JSON(200, gin.H{
        "results": matches,
    })
}