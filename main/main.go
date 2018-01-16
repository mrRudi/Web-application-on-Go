package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-martini/martini"
	//"github.com/martini-contrib/render"

	"github.com/mrRudi/Web-application-on-Go/models"
	"github.com/mrRudi/Web-application-on-Go/util"
)

var posts map[string]*models.Post
var bell int

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/tindex.html", "templates/footer.html", "templates/handler.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Println(posts)
	fmt.Println("bell = ",bell)
	t.ExecuteTemplate(w, "index", posts)
}

func savePostHandler(w http.ResponseWriter, r *http.Request) {
	post := models.NewPost(r.FormValue("id"), r.FormValue("title"), r.FormValue("content"))
	posts[r.FormValue("id")] = post
	http.Redirect(w, r, "/", http.StatusFound)
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/twrite.html", "templates/footer.html", "templates/handler.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	id := util.GenerateId()
	t.ExecuteTemplate(w, "write", id)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/edit.html", "templates/footer.html", "templates/handler.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Println(r.Form)
	id := r.FormValue("id")
	post, found := posts[id]
	if !found {
		http.NotFound(w, r)
	}
	t.ExecuteTemplate(w, "edit", post)
}
func deletePostHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
	fmt.Println(params)
	id := params["id"]
	_, found := posts[id]
	if !found {
		http.NotFound(w, r)
	}
	delete(posts, id)
	http.Redirect(w, r, "/", 302)
}

func main() {
	posts = make(map[string]*models.Post, 0)
	id := util.GenerateId()
	posts[id] = models.NewPost(id, "ds1", "sdasd1")
	id = util.GenerateId()
	posts[id] = models.NewPost(id, "ds2", "sdasd2")
	bell = 0

	m := martini.Classic()

	fmt.Println("begin")

	staticOptions := martini.StaticOptions{Prefix:"assets"}
	m.Use(martini.Static("assets",staticOptions))

	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit", editHandler)
	m.Post("/SavePost", savePostHandler)
	m.Delete("/DeletePost:id", deletePostHandler)

	m.Use(func(r *http.Request) {
		if r.URL.Path == "/write"{
			bell++
		}
	})

	m.Run()
}
