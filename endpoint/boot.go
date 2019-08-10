package endpoint

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/demostack/auth/pkg"

	"github.com/aws/aws-sdk-go/aws/endpoints"
)

// Boot .
func Boot(username, password string) *Core {
	// Get the namespace from the environment variable if it exists.
	namespace := "auth"
	if len(os.Getenv("NAMESPACE")) > 0 {
		namespace = os.Getenv("NAMESPACE")
	}

	// Get the region if set, otherwise default to us-east-1.
	region := endpoints.UsEast1RegionID
	if len(os.Getenv("AWS_REGION")) > 0 {
		region = os.Getenv("AWS_REGION")
	}
	if len(os.Getenv("AWS_DEFAULT_REGION")) > 0 {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}

	db := pkg.NewDynamoDB(region, namespace+"-verify")

	err := db.CreateVerifyTable()
	if err != nil {
		if !strings.Contains(err.Error(), "Table already exists") {
			fmt.Println("Error creating the verify table:", err)
		}
	}

	fromEmail := os.Getenv("FROMEMAIL")
	if len(fromEmail) == 0 {
		log.Fatalln("Error: you must set environment variable: FROMEMAIL (format: noreply@example.com)")
	}

	toEmail := os.Getenv("TOEMAIL")
	if len(toEmail) == 0 {
		log.Fatalln("Error: you must set environment variable: TOEMAIL (format: 2225557777@vtext.com)")
	}

	return &Core{
		DB:         db,
		region:     region,
		username:   username,
		password:   password,
		fromEmail:  fromEmail,
		toEmail:    toEmail,
		isSAMLocal: len(os.Getenv("AWS_SAM_LOCAL")) > 0,
	}
}
