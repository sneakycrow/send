package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type FileMetaData struct {
	gorm.Model
	Location   string `json:"location"`
	Expiration int    `json:"expiration"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s3Client := gets3Client()
	spacesName := os.Getenv("SPACE_NAME")
	destinationFolder := os.Getenv("DESTINATION_FOLDER")
	endpoint := os.Getenv("ENDPOINT")
	router := gin.Default()

	router.POST("/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		fileName := fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(header.Filename))
		if err != nil {
			panic(err)
		}

		fileBytes, err := ioutil.ReadAll(file)

		contentType := http.DetectContentType(fileBytes)

		object := s3.PutObjectInput{
			Bucket:      aws.String(spacesName),
			Key:         aws.String(destinationFolder + "/" + fileName),
			Body:        strings.NewReader(string(fileBytes)),
			ACL:         aws.String("public-read"),
			ContentType: aws.String(contentType),
		}
		_, err = s3Client.PutObject(&object)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Could not send file to bucket", "error": err})
		} else {
			location := fmt.Sprintf("https://%s.%s/%s/%s", spacesName, endpoint, destinationFolder, fileName)

			t := time.Now().Add(time.Hour * 24)
			fileMeta := FileMetaData{
				Expiration: int(t.UnixNano()),
				Location:   location,
			}

			c.JSON(http.StatusAccepted, gin.H{"data": fileMeta})
		}

	})

	router.Run(":8080")
}

func addFilesToS3(fileDir string, size int64, r io.Reader) error {

	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, size)
	file.Read(buffer)

	return err
}

func gets3Client() *s3.S3 {
	key := os.Getenv("ACCESS_KEY")
	secret := os.Getenv("SECRET_KEY")
	endpoint := os.Getenv("ENDPOINT")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String("us-east-1"),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	return s3Client
}
