package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDB implementation.
type DynamoDB struct {
	sess        *session.Session
	region      string
	tableSecret string
	tableVerify string
}

// NewDynamoDB returns a new instance of DynamoDB.
func NewDynamoDB(region string, tableVerify string) *DynamoDB {
	return &DynamoDB{
		sess:        session.Must(session.NewSession()),
		region:      region,
		tableVerify: tableVerify,
	}
}

// CreateVerifyTable creates the verify table.
func (db *DynamoDB) CreateVerifyTable() (err error) {
	svc := dynamodb.New(db.sess, db.sess.Config.WithRegion(db.region))
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(db.tableVerify),
	}

	_, err = svc.CreateTable(input)

	return
}

// SetVerify will set the verify in the database.
func (db *DynamoDB) SetVerify(id string, isVerified bool) (err error) {
	svc := dynamodb.New(db.sess, db.sess.Config.WithRegion(db.region))

	m := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
		"verified": {
			BOOL: aws.Bool(isVerified),
		},
	}

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item:      m,
		TableName: aws.String(db.tableVerify),
	})
	return
}

// Verified returns true if verified, true if found.
func (db *DynamoDB) Verified(id string) (bool, bool) {
	svc := dynamodb.New(db.sess, db.sess.Config.WithRegion(db.region))
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(db.tableVerify),
	}

	result, err := svc.GetItem(input)
	if err != nil || result.Item == nil {
		return false, false
	}

	if v, ok := result.Item["verified"]; ok {
		return aws.BoolValue(v.BOOL), true
	}

	return false, false
}

// DeleteVerify deletes old IDs.
func (db *DynamoDB) DeleteVerify(id string) (err error) {
	svc := dynamodb.New(db.sess, db.sess.Config.WithRegion(db.region))
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.tableVerify),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	_, err = svc.DeleteItem(input)
	return
}
