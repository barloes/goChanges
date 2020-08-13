package records

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

type Item struct {
	Email string
	Url string
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
	   if a == str {
		  return true
	   }
	}
	return false
 }

func ListContent(tableName string) map[string]string{


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
	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	m := make(map[string]string)
	for _, i := range result.Items {
		item := Item{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		//key value pair for golang
		m[item.Url] = item.Email
	}
	fmt.Println(m["abc"])
	return m
}

func DBUpdate(url string,email string,tableName string){

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

	svc := dynamodb.New(sess)
	
	content := url
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

	//if website dont exist in db create new one and update with the email provided
	if len(item.Url) == 0 {
		email += ","
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":r": {
					S: aws.String(email),
				},
			},
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"url": {
					S: aws.String(content),
				},
			},
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression: aws.String("set email = :r"),
		}
	
		svc.UpdateItem(input)
		
	}else{

		fmt.Println("aaa")
		emailList := strings.Split(item.Email, ",")
		
		//if email already exist then return
		if(contains(emailList,email)){
			fmt.Println("email already existed!")
			return
		}
		emailList = append(emailList,email)
		fmt.Println(emailList)
		var updatedEmail string = "" 
		for _,eachEmail := range emailList{
			if len(eachEmail) == 0{
				continue
			}else{
				updatedEmail += eachEmail
				updatedEmail += ","
			}
		}
		fmt.Println(updatedEmail)

		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":r": {
					S: aws.String(updatedEmail),
				},
			},
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"url": {
					S: aws.String(content),
				},
			},
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression: aws.String("set email = :r"),
		}
	
		svc.UpdateItem(input)

	}

}