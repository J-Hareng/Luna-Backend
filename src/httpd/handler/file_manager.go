package handler

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"server/src/api/db"
	"server/src/api/db/models"
	filemanagement "server/src/api/file_management"
	"server/src/helper"
	"server/src/httpd/bodymodels"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileManager struct{}

type File struct {
	Head models.File
	Body io.Reader
}
type FileBlob struct {
	Head models.File
	Body string
}

func UploadFile(db *db.DB, s3 *filemanagement.S3Client) gin.HandlerFunc {

	return func(c *gin.Context) {
		fileRaw, err := c.FormFile("file")
		helper.CustomErrorApi(c, err)
		name := c.PostForm("name")
		username := c.PostForm("username")
		userId, err := primitive.ObjectIDFromHex(c.PostForm("userId"))
		helper.CustomErrorApi(c, err)
		tags := c.PostForm("tags")

		collection, err := primitive.ObjectIDFromHex(c.PostForm("coll"))
		helper.CustomErrorApi(c, err)

		path := helper.RandomString(6) + fileRaw.Filename
		file := models.CreateFile(strings.Split(tags, "§"), name, path, models.UserLink{ID: userId, NAME: username})

		_, err = db.AddFileToCollection(file, collection)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add file to collection"})
		}
		src, err := fileRaw.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
			return
		}
		defer src.Close()

		// Create a temporary in-memory file
		tempFile, err := ioutil.TempFile("", "*")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file"})
			return
		}

		// Copy the uploaded file to the temporary file
		_, err = io.Copy(tempFile, src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy uploaded file"})
			return
		}
		// Use the temporary file
		err = s3.UploadFile(tempFile, path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to S3"})
			return
		}
		c.JSON(http.StatusOK, file)
	}

}
func DownloadFile(db *db.DB, s3 *filemanagement.S3Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var val bodymodels.DonwloadFileMod
		if err := c.ShouldBindJSON(&val); err != nil {
			fmt.Print(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//! NEED TO INSTALL USER AUTH HERE
		tempfile, err := s3.DownloadFile(val.File.Path)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
			return
		}
		defer tempfile.Close()

		fmt.Println(val)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+val.File.Name)
		c.Header("Content-Type", "application/octet-stream")
		c.DataFromReader(http.StatusOK, -1, "application/octet-stream", tempfile, map[string]string{})
	}
}
func GetMultibleFiles(db *db.DB, s3 *filemanagement.S3Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var val bodymodels.GetMultibleFilesMod
		if err := c.ShouldBindJSON(&val); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		files := make([]FileBlob, 0, len(val.Files))

		fileChan := make(chan FileBlob)
		errchan := make(chan error)
		defer close(fileChan)
		defer close(errchan)

		for _, file := range val.Files {
			go func(file models.File, fchan chan FileBlob, errchan chan error) {

				tempfile, err := s3.DownloadFile(file.Path)
				if err != nil {
					errchan <- err
					return
				}
				defer tempfile.Close()

				// Read the file into a byte slice
				bytes, err := ioutil.ReadAll(tempfile)
				if err != nil {
					errchan <- err
					return
				}

				// Encode the byte slice to Base64
				encoded := base64.StdEncoding.EncodeToString(bytes)

				// Add the encoded file to the files slice
				fchan <- FileBlob{Head: file, Body: encoded}
			}(file, fileChan, errchan)
		}
		for i := 0; i < len(val.Files); i++ {
			select {
			case file := <-fileChan:
				files = append(files, file)
			case err := <-errchan:
				c.JSON(http.StatusInternalServerError, bodymodels.NewLunaResponse_ERROR(err.Error()))
				return
			}
		}
		c.JSON(http.StatusOK, bodymodels.NewLunaResponse_OK(files, ""))

	}
}

func GetFile(db *db.DB, s3 *filemanagement.S3Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var val bodymodels.DonwloadFileMod
		if err := c.ShouldBindJSON(&val); err != nil {
			c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR_INVALID_PAYLOAD())
			return
		}
		fmt.Println(val.File.Path)

		tempfile, err := s3.DownloadFile(val.File.Path)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
			return
		}
		defer tempfile.Close()

		// Get the file extension
		fileExt := filepath.Ext(val.File.Path)
		fmt.Println(tempfile)
		// Read the file data into a byte slice
		fileData, err := ioutil.ReadAll(tempfile)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file data"})
			return
		}
		fmt.Println("File data size:", len(fileData))

		// Set headers for file download
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+val.File.Name)
		c.Header("Content-Type", getContentType(fileExt))

		// Serve the file
		c.Data(http.StatusOK, getContentType(fileExt), fileData)
	}
}
func RemoveFile(db *db.DB, S3 *filemanagement.S3Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var val bodymodels.RemoveFileFromCollectionMod
		if err := c.ShouldBindJSON(&val); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := S3.DeleteFile(val.File.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from S3"})
		}
		// Remove the file from the database
		_, err = db.RemoveFileFromCollection(val.File, val.CollId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove file from database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File removed successfully"})
	}
}
func getContentType(ext string) string {
	switch ext {
	case ".jpg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}
