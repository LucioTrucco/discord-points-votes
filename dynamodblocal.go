package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateLocalClient() *dynamodb.DynamoDB {

	creds := credentials.NewStaticCredentials("123", "123", "")
	awsConfig := &aws.Config{
		Credentials: creds,
	}
	awsConfig.WithRegion("us-east-1")
	awsConfig.WithEndpoint("http://localhost:8000")

	s, err := session.NewSession(awsConfig)
	if err != nil {
		panic(err)
	}
	dynamodbconn := dynamodb.New(s)
	return dynamodbconn
}

func CreateTableIfNotExists(d *dynamodb.DynamoDB, tableName string) {
	if tableExists(d, tableName) {
		log.Printf("table=%v already exists\n", tableName)
		return
	}
	_, err := d.CreateTable(buildCreateTableInput(tableName))
	if err != nil {
		log.Fatal("CreateTable failed", err)
	}
	log.Printf("created table=%v\n", tableName)
}

func buildCreateTableInput(tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
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
		TableName:   aws.String(tableName),
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	}
}

func tableExists(d *dynamodb.DynamoDB, name string) bool {
	tables, err := d.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatal("ListTables failed", err)
	}
	for _, n := range tables.TableNames {
		if *n == name {
			return true
		}
	}
	return false
}

func writeInDynamoDB(d *dynamodb.DynamoDB, item map[string]*dynamodb.AttributeValue, tableName string) (*dynamodb.PutItemOutput, error) {
	itemInput := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	out, err := d.PutItem(itemInput)
	if err != nil {
		fmt.Println(err)
	}
	return out, nil
}
