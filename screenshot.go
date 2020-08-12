package main

import (

	"log"
	"net/http"
	"fmt"
	"time"
	"os/exec"
	"os"
	"strings"
	"github.com/vitali-fedulov/images"
	"github.com/joho/godotenv"
	email "github.com/junhuiyara/goTest/email"
	records "github.com/junhuiyara/goTest/records"
)

//automated script that runs d seconds
func doEvery(d time.Duration, f func(string,string) bool) {

	m := make(map[string]string)
	m = records.ListContent("jundb")

	for x := range time.Tick(d) {
		_ = x

		for key, value := range m {
			//if image is different,send email to all the email
			//eg https://www.google.com/" to "google.com"
			filename := key[12:len(key)-1]
			filename += "-1366x768"
			fmt.Println(filename)
			if !f(key,filename){
				emailList := strings.Split(value, ",")

				for _,emailRecipient := range emailList{
					email.SendMailTo(emailRecipient,key)
				}
			}

		}
}
}

func GetReq(url string){

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
	// q.Add("old_word_count", strconv.Itoa(old_word_count))
	// q.Add("new_word_count", strconv.Itoa(new_word_count))
	req.URL.RawQuery = q.Encode()

	//fmt.Println(req.URL.String())
	http.Get(req.URL.String())

}

func sameImage(url string,filename string) bool{

	//command to get the image of the
    cmd := exec.Command("pageres",url,"--user-agent=XYZ/3.0")
    stdout, err := cmd.Output()

    if err != nil {
        fmt.Println(err.Error())
    }

	fmt.Print(string(stdout))
	
	// Open photos.
	Original_Path := filename + ".png"
	New_Path := filename +" (1).png"
	imgA, err := images.Open(Original_Path) 
	if err != nil {
		fmt.Println("image a not found")
		return true
	}
	imgB, err := images.Open(New_Path)
	if err != nil {
		fmt.Println("image b not found")
		return true
	}
	
	// Calculate hashes and image sizes.
	hashA, imgSizeA := images.Hash(imgA)
	hashB, imgSizeB := images.Hash(imgB)
	
	// Image comparison.
	if images.Similar(hashA, hashB, imgSizeA, imgSizeB) {
		//similar image then delete the (1).png file (with try)
		path := New_Path
		err := os.Remove(path)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Images are similar.")
		return true
	} else {
		//different image then delete the .png file and rename (1).png file to .png
		//then send api call
		path := Original_Path
		err := os.Remove(path)
		if err != nil {
			fmt.Println(err)
		}

		err1 := os.Rename(New_Path,Original_Path ) 
		if err1 != nil { 
			log.Fatal(err1) 
		} 

		fmt.Println("Images are distinct.")
		
		return false
	}
}

func main() {

	doEvery(5* time.Second, sameImage)
}
