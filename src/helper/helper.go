package helper

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"server/src/httpd/bodymodels"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func GetEnvVar(key string) string {
	// load .env file
	s := os.Getenv(key)
	fmt.Println("ENV" + s)
	if s != "" {
		return s
	}
	panic("!!!--------------!!!\nENV ERROR: " + key + " not found in .env file.\n!!!--------------!!!")
}

func CustomError(err string) {

	panic("!!!--------------!!!\nCUSTOM ERROR" + err + " .\n!!!--------------!!!")
}

type Status struct {
	STATUS string
}

func CustomErrorApi(c *gin.Context, err error) {
	if err != nil {
		fmt.Print("Custom error : ")
		fmt.Println(err)

		c.JSON(http.StatusBadRequest, bodymodels.NewLunaResponse_ERROR("Custom Error: "+err.Error()))
		c.Abort()
		return
	}
}

type InsertOption[T any] struct {
	Property T
}

func Prittylog(des string, data any, info string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.WithFields(logrus.Fields{
		des: data,
	}).Info(info)
}
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
func RemoveFileRaw(path string) error {
	filePath := GetEnvVar("FILEURL") + path
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
