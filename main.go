package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./view/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("uploadfile")
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		path := fmt.Sprintf("./img/%s", file.Filename)
		err = saveImg(file, path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"FileName": file.Filename,
			"Size":     file.Size,
		})
	})

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func saveImg(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	f, err := os.OpenFile(fmt.Sprintf("./img/%s", file.Filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, src)

	return err	
}