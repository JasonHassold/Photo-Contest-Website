package main

import (
	"fmt"
	"io"
	"io/ioutil"
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

		// Create votes.txt file
		votes, err := os.Create("site/assets/submissions/" + strconv.Itoa(id) + "/votes.txt")
		if err != nil {
			panic(err)
			return
		}
		defer votes.Close()

		fmt.Fprintf(votes, "0")

		// Save photo file
		out, err := os.Create("site/assets/submissions/" + strconv.Itoa(id) + "/picture.png")
		if err != nil {
			panic(err)
			return
		}
		defer out.Close()

		io.Copy(out, file)

		// Return to Success page
		http.Redirect(w, r, "/success", http.StatusFound)
	}
}

func vote(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		files, _ := ioutil.ReadDir("site/assets/submissions")
		num := len(files)
		nums := make([]int, 0)
		for i := 1; i <= num; i++ {
			nums = append(nums, i)
		}

		tmpl := template.Must(template.ParseFiles("site/vote.html", "site/partials/header.html", "site/partials/footer.html"))
		tmpl.ExecuteTemplate(w, "vote", nums)
	} else if r.Method == "POST" {
		var id = r.URL.Query().Get("id")

		votes, err := os.Open("site/assets/submissions/" + id + "/votes.txt");
		if err != nil {
			panic(err)
			return
		}
		
		current_votes, _ := ioutil.ReadAll(votes)
		string_votes := string(current_votes)
		int_votes, _ := strconv.Atoi(string_votes)
		int_votes += 1
		new_votes := strconv.Itoa(int_votes)
		votes.Close()

		err = os.Remove("site/assets/submissions/" + id + "/votes.txt")
	    if err != nil {
	        log.Fatal(err)
	    }

	    votes, err = os.Create("site/assets/submissions/" + id + "/votes.txt")
		if err != nil {
			panic(err)
			return
		}
		defer votes.Close()

		fmt.Fprintf(votes, new_votes)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("site/contact.html", "site/partials/header.html", "site/partials/footer.html"))
		tmpl.ExecuteTemplate(w, "contact", "")
	}
}

func success(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("site/success.html", "site/partials/header.html", "site/partials/footer.html"))
		tmpl.ExecuteTemplate(w, "success", "")
	}
}


func main() {
	
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("site/assets/"))))
	http.HandleFunc("/enter", enter)
	http.HandleFunc("/vote", vote)
	http.HandleFunc("/contact", contact)
	http.HandleFunc("/success", success)
	http.HandleFunc("/", index)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
