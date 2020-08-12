package main

import (
        "io/ioutil"
        "log"
        "net/http"
        "fmt"
        "time"
        "context"
        "os"
        "strconv"

        "go.mongodb.org/mongo-driver/bson"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"

        "github.com/joho/godotenv"
)

type Website struct {
	Url string
	Word_count  int
}

//automated script that runs d seconds
func doEvery(d time.Duration, f func(string)) {
        var urlArray [] string
        urlArray = append(urlArray,"https://www.adidas.com.sg/")
        urlArray = append(urlArray,"https://www.example.com/")

	for x := range time.Tick(d) {
                _ = x

                for _, url := range urlArray {
                        f(url)
                }
	}
}

func getSiteWordCount(url string){
	client := &http.Client{}
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
                log.Fatalln(err)
        }

        req.Header.Set("User-Agent", "XYZ/3.0")

        resp, err := client.Do(req)
        if err != nil {
                log.Fatalln(err)
        }

        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Fatalln(err)
        }
        word_count := len(string(body))
        fmt.Println(word_count)

        connectDB(url,word_count)
}

func connectDB(url string, word_count int){

        err := godotenv.Load(".env")
        if err != nil {
          log.Fatalf("Error loading .env file")
        }

        user := os.Getenv("mongo_user")
        pass := os.Getenv("mongo_pass")
        db_name := os.Getenv("db_name")
        collection_name := os.Getenv("collection_name")

        client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://" + user + ":" + pass + "@cluster0-nlvco.mongodb.net/<dbname>?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
        }
        
        ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
        }
        fmt.Println("Connected to MongoDB!")
        
        collection := client.Database(db_name).Collection(collection_name)
        
        var result Website
        filter := bson.D{{"url", url}}
	collection.FindOne(context.TODO(), filter).Decode(&result)
        
        if result.Url == url && result.Word_count != word_count{
                fmt.Println("result is changed")

                GetReq(url,result.Word_count,word_count)

                result, err := collection.UpdateOne(
                        context.TODO(),
                        filter,
                        bson.D{
                                {"$set", bson.D{{"word_count", word_count}}},
                        },
                )
                if err != nil {
                        log.Fatal(err)
                }
                _ = result

        } else if (Website {}) == result {
                new_website := Website{url, word_count}

                // Insert a single document
                collection.InsertOne(context.TODO(), new_website)
        } else {
                fmt.Println("result is the same")
        }

        client.Disconnect(ctx)
	fmt.Println("Connection to MongoDB closed.")
  
}

func GetReq(url string,old_word_count int,new_word_count int){

        err := godotenv.Load(".env")
        if err != nil {
          log.Fatalf("Error loading .env file")
        }

        api_gateway := os.Getenv("api_gateway")

        req, err := http.NewRequest("GET", api_gateway, nil)
        if err != nil {
                log.Print(err)
                os.Exit(1)
        }

        q := req.URL.Query()
        q.Add("website_url", url)
        q.Add("old_word_count", strconv.Itoa(old_word_count))
        q.Add("new_word_count", strconv.Itoa(new_word_count))
        req.URL.RawQuery = q.Encode()

        //fmt.Println(req.URL.String())
        http.Get(req.URL.String())

}

func main() {
        doEvery(3* time.Second, getSiteWordCount)
}

