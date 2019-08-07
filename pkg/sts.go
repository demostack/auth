package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// GetSessionToken .
func GetSessionToken(region string, seconds int64) (*sts.Credentials, error) {
	// Create a new session.
	sess := session.Must(session.NewSession())

	// Create a STS client.
	svc := sts.New(sess, sess.Config.WithRegion(region))
	out, err := svc.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(seconds),
	})
	if err != nil {
		return nil, err
	}

	return out.Credentials, nil
}
