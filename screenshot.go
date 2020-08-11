package main

import (

	"log"
	"net/http"
	"fmt"
	"time"
	"os/exec"
	"os"
	"github.com/vitali-fedulov/images"

	"github.com/joho/godotenv"
)

//automated script that runs d seconds
func doEvery(d time.Duration, f func(string)) {
	var urlArray [] string
	urlArray = append(urlArray,"https://www.adidas.com.sg/")

	for x := range time.Tick(d) {
		_ = x

		for _, url := range urlArray {
			f(url)
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

func checkImage(url string) {

	//command to get the image of the
    cmd := exec.Command("pageres",url,"--user-agent=XYZ/3.0")
    stdout, err := cmd.Output()

    if err != nil {
        fmt.Println(err.Error())
        return
    }

	fmt.Print(string(stdout))
	
	// Open photos.
	Original_Path := "adidas.com.sg-1366x768.png"
	New_Path := "adidas.com.sg-1366x768 (1).png"
	imgA, err := images.Open(Original_Path)
	if err != nil {
		panic(err)
	}
	imgB, err := images.Open(New_Path)
	if err != nil {
		panic(err)
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
			return
		}
		fmt.Println("Images are similar.")
	} else {
		//different image then delete the .png file and rename (1).png file to .png
		//then send api call
		path := Original_Path
		err := os.Remove(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		err1 := os.Rename(New_Path,Original_Path ) 
		if err1 != nil { 
			log.Fatal(err1) 
		} 
		GetReq(url)
		fmt.Println("Images are distinct.")
	}
}

func main() {
	doEvery(300* time.Second, checkImage)
}
