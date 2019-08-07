package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// SendMessage .
func SendMessage(region, phone, message string) error {
	// Create a new session.
	sess := session.Must(session.NewSession())

	// Create a Pinpoint client from just a session.
	svc := sns.New(sess, sess.Config.WithRegion(region))
	_, err := svc.Publish(&sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phone),
	})

	return err
}
