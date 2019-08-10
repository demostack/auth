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

	// Create a SNS client from a session.
	svc := sns.New(sess, sess.Config.WithRegion(region))
	svc.SetSMSAttributes(&sns.SetSMSAttributesInput{
		Attributes: map[string]*string{
			"MonthlySpendLimit": aws.String("1"),
			//"DefaultSenderID":   aws.String("Demostack"),
			"DefaultSMSType": aws.String("Transactional"),
		},
	})
	_, err := svc.Publish(&sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phone),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"AWS.SNS.SMS.SMSType": &sns.MessageAttributeValue{
				StringValue: aws.String("Transactional"),
				DataType:    aws.String("String"),
			},
			"AWS.SNS.SMS.MaxPrice": &sns.MessageAttributeValue{
				StringValue: aws.String("0.01"),
				DataType:    aws.String("Number"),
			},
			/*"AWS.SNS.SMS.SenderID": &sns.MessageAttributeValue{
				StringValue: aws.String("Demostack"),
				DataType:    aws.String("String"),
			},*/
		},
	})

	return err
}
