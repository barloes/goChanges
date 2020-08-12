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

package webserver

import (
	
	records "github.com/junhuiyara/goTest/records"

	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"strings"
	"errors"
	"net/url"
	"regexp"
)

type Item struct {
	Email string
	Url string
}

//checkEmail checks if the email inputted is a valid email.
func CheckEmail(email string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(email) > 254 || !rxEmail.MatchString(email) {
		return false
	}
	return true && len(email) > 0
}

//checkWebsite checks if the website inputted is a valid URL.
func CheckWebsite(website string) (bool, error) {
	u, err := url.Parse(website)
	if err != nil {
		err = errors.New("website inputted is not a valid URL")
		return false, err
	} else if u.Scheme == "" || u.Host == "" {
		err = errors.New("website inputted must be an absolute URL")
		return false, err
	} else if u.Scheme != "http" && u.Scheme != "https" {
		err = errors.New("website inputted must begin with http or https")
		return false, err
	}
	return true, nil
}


// TODO
func StartListen() {
	http.HandleFunc("/", handlerHome)
	http.HandleFunc("/save/", handlerSave)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// TODO
func handlerHome(w http.ResponseWriter, r *http.Request) {

	setContentType(w, r)

	URI := r.RequestURI
	if URI == "/" {
		URI = "./fe/home.html"
	} else {
		URI = "." + URI
	}

	page, _ := loadPage(URI)

	fmt.Fprintf(w, "%s", page)
}

// TODO
func handlerSave(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	website := r.FormValue("website")

	// CHECK INPUT
	emailOk := CheckEmail(email)
	websiteOk,err := CheckWebsite(website)

	if(err != nil){
		fmt.Println(err)
		return
	}

	// IF EMAIL OK, INSERT EMAIL
	if emailOk &&websiteOk{
		records.DBUpdate(email,website,"jundb")
	}
	// only if there are some data available print on page

	// DEBUG PRINTING
	fmt.Println("Email:", email)
	fmt.Println("Website:", website)
	fmt.Fprintf(w, "%s %s", email, website)
}

// TODO
func loadPage(filename string) (string, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func setContentType(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	contentType := "text/html"

	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "application/javascript"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png"
	}

	w.Header().Set("Content-Type", contentType)

}
