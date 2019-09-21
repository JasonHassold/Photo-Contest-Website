package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"html/template"
)

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("site/index.html", "site/partials/header.html", "site/partials/footer.html"))
		tmpl.ExecuteTemplate(w, "index", "")
	}
}

func enter(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("site/enter.html", "site/partials/header.html", "site/partials/footer.html"))
		tmpl.ExecuteTemplate(w, "enter", "")
	} else if r.Method == "POST" {
		fmt.Println("POST")

		// Creates new id and folder with that id
		id := 1
		made := false
		for made == false {
			if _, err := os.Stat("site/assets/submissions/" + strconv.Itoa(id) + "/"); os.IsNotExist(err) && id <= 100 {
				err := os.Mkdir("site/assets/submissions/"+strconv.Itoa(id)+"/", 0755)
				if err != nil {
					panic(err)
				}
				made = true
			} else if id <= 100 {
				id++
			} else {
				return
			}
		}

		// Parse Form into Variables
		r.ParseMultipartForm(32 << 20)

		first_name := r.Form["first_name"][0]
		last_name := r.Form["last_name"][0]
		email := r.Form["email"][0]
		phone := r.Form["phone"][0]
		file, _, err := r.FormFile("photo")
		if err != nil {
			panic(err)
			return
		}
		defer file.Close()

		// Write info to info.txt file
		info, err := os.Create("site/assets/submissions/" + strconv.Itoa(id) + "/info.txt")
		if err != nil {
			panic(err)
			return
		}
		defer info.Close()

		fmt.Fprintf(info, "%s\r\n", first_name)
		fmt.Fprintf(info, "%s\r\n", last_name)
		fmt.Fprintf(info, "%s\r\n", email)
		fmt.Fprintf(info, "%s\r\n", phone)

		// Save photo file
		out, err := os.Create("site/assets/submissions/" + strconv.Itoa(id) + "/picture.png")
		if err != nil {
			panic(err)
			return
		}
		defer out.Close()

		io.Copy(out, file)

		// Return to Success page
		http.Redirect(w, r, "/success.html", http.StatusFound)
	}
}

func vote(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("site/contact.html", "site/partials/header.html", "site/partials/footer.html"))
		tmpl.ExecuteTemplate(w, "contact", "")
	}
}


func main() {
	
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("site/assets/"))))
	http.HandleFunc("/enter", enter)
	http.HandleFunc("/vote", vote)
	http.HandleFunc("/contact", contact)
	http.HandleFunc("/", index)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
