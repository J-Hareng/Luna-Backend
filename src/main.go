package main

import (
	"context"
	"fmt"
	"os"
	email "server/src/api/Email"
	"server/src/api/db"
	filemanagement "server/src/api/file_management"
	"server/src/helper"
	"server/src/httpd"
	"server/src/httpd/security"
)

func main() {

	ctx := context.Background()
	defer ctx.Done()

	fmt.Println("Servers First initialisation")
	fmt.Println(os.Getenv("MONGOURI"))
	DB, err := db.New() // * setup datenbank
	if err != nil {
		helper.CustomError(err.Error())
		fmt.Println(err)
	}
	S3, err := filemanagement.New()
	if err != nil {
		helper.CustomError(err.Error())
		fmt.Println(err)
	}
	GKM := security.SelectGroupTokenMap{}
	EKM := security.EmailTokenMap{} // * Generate Email TokenMap
	E := email.GeneratEmail()       // * Inisialise email

	s := httpd.Init(ctx, DB, E, &EKM, &GKM, S3) // * Initialisiere server
	s.Run()                                     // * Starte server
}
