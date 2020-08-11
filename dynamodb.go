
package main

// snippet-start:[dynamodb.go.list_tables.imports]
import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/joho/godotenv"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"strings"
	"os"
	"fmt"
	"log"
)

// snippet-end:[dynamodb.go.list_tables.imports]

type Item struct {
	Email string
	Url string
}

// func updateRecord(){

// }

// func listRecord(){
// 	result, err := svc.GetItem(&dynamodb.GetItemInput{
// 		TableName: aws.String(tableName),
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"url": {
// 				S: aws.String(content),
// 			},
// 		},
// 	})
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
		
// 	item := Item{}
	
// 	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
// 	if err != nil {
// 		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
// 	}
// }

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
    AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_SECRET_KEY := os.Getenv("AWS_SECRET_KEY")
	REGION := os.Getenv("REGION")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(REGION),
		Credentials: credentials.NewStaticCredentials(AWS_SECRET_ACCESS_KEY, AWS_SECRET_KEY,""),
	})
	if err != nil {
		log.Fatal(err)
	}

    // Create DynamoDB client
    svc := dynamodb.New(sess)
    // snippet-end:[dynamodb.go.list_tables.session]

	tableName := "jundb"
	content := "notthere"
	
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"url": {
				S: aws.String(content),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
		
	item := Item{}
	
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	//split record into list using delimeter ","
	urlList := strings.Split(item.Email, ",")

	for _,s := range urlList {
		fmt.Println(s)
	}
	if len(item.Url) == 0 {
		fmt.Println("item does not exist")
	}

	proj := expression.NamesList(expression.Name("url"), expression.Name("email"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}
	
	// Make the DynamoDB Query API call
	result1, err := svc.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	m := make(map[string]string)
	for _, i := range result1.Items {
		item := Item{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// Which ones had a higher rating than minimum?


		fmt.Println(item)

		//key value pair for golang
		m[item.Url] = item.Email
		fmt.Println(m)
		
	}
	fmt.Println(m["abc"])


	// newList := "1,2,3,4"
	// input := &dynamodb.UpdateItemInput{
	// 	ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
	// 		":r": {
	// 			S: aws.String(newList),
	// 		},
	// 	},
	// 	TableName: aws.String(tableName),
	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		"url": {
	// 			S: aws.String(content),
	// 		},
	// 	},
	// 	ReturnValues:     aws.String("UPDATED_NEW"),
	// 	UpdateExpression: aws.String("set email = :r"),
	// }
	
	// svc.UpdateItem(input)

}

//get the list of records every 5 seconds
//
