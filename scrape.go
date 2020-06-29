package main

import (
        "io/ioutil"
        "log"
        "net/http"
        "fmt"
        "time"
)

//automated script that runs d seconds
func doEvery(d time.Duration, f func(string)) {
	for x := range time.Tick(d) {
                _ = x
                f("https://www.adidas.com.sg/men-running-shoes?sort=price-low-to-high&v_size_en_sg=9_uk")
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
        
        fmt.Println(len(string(body)))
}

func main() {
	doEvery(2 * time.Second, getSiteWordCount)
}

