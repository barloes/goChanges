/*
This file is under GNU AFFERO GENERAL PUBLIC LICENSE

Permissions of this strongest copyleft license are conditioned
on making available complete source code of licensed works and
modifications, which include larger works using a licensed work,
under the same license. Copyright and license notices must be preserved.
Contributors provide an express grant of patent rights.
When a modified version is used to provide a service over a network,
the complete source code of the modified version must be made available.

Edoardo Ottavianelli, https://edoardoottavianelli.it

*/

package main

import (
	
	records "github.com/junhuiyara/goTest/records"

	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"log"
)

//checkEmail checks if the email inputted is a valid email.
func CheckEmail(email string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(email) > 254 || !rxEmail.MatchString(email) {
		return false
	}
	return true && len(email) > 0
}

//checkWebsite checks if the website inputted is a valid URL.
func CheckWebsite(website string) (bool) {
	//eg. https://www.google.com/
	//filename := key[12:len(key)]
	fmt.Println(website[0:12])
	fmt.Println(website[len(website)-4:len(website)])

	if(website[0:12] == "https://www." && website[len(website)-4:len(website)] == ".com"){
		fmt.Println("true")
		return true
	}else{
		fmt.Println("false")
		return false
	}
}


// TODO
func StartListen() {
	http.HandleFunc("/", handlerHome)
	http.HandleFunc("/submit", handlerSubmit)
    log.Fatal(http.ListenAndServe(":80", nil))
}

func handlerHome(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.URL.Path)
    p := "." + r.URL.Path
    if p == "./" {
        p = "fe/home.html"
    }
    http.ServeFile(w, r, p)
}

// TODO
func handlerSubmit(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	website := r.FormValue("url")

	// CHECK INPUT
	emailOk := CheckEmail(email)
	websiteOk := CheckWebsite(website)

	// IF EMAIL OK, INSERT EMAIL
	if emailOk &&websiteOk{
		records.DBUpdate(website,email,"jundb")
	}
	// only if there are some data available print on page

	if r.Method == "POST" {
		http.Redirect(w, r, "./", http.StatusSeeOther)
	}
}

// TODO
func loadPage(filename string) (string, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func main(){
	StartListen()
}