package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/vavilen84/go-martini-test-app/documents"
	"html/template"
	"log"
	"net/http"
	"time"
)

var postsCollection *mongo.Collection

func GenerateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func indexHandler(rnd render.Render) {
	posts := make([]documents.PostDocument, 0)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := postsCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		post := documents.PostDocument{}
		err := cur.Decode(&post)
		if err != nil {
			log.Fatal(err)
		} else {
			posts = append(posts, post)
		}
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	post := documents.PostDocument{}
	rnd.HTML(200, "write", post)
}

func deleteHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]
	if id == "" {
		rnd.Redirect("/")
	}
	filter := bson.D{{"_id", id}}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	postsCollection.DeleteOne(ctx, filter)
	rnd.Redirect("/")
}

func editHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]
	post := documents.PostDocument{}
	filter := bson.M{"_id": id}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := postsCollection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		rnd.Redirect("/")
		return
	}
	rnd.HTML(200, "write", post)
}

func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")
	postDocument := documents.PostDocument{id, title, content}
	fmt.Printf("%+v", postDocument)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if id != "" {
		filter := bson.M{"_id": id}
		update := bson.D{
			{"$set", bson.D{
				{"title", title},
				{"content", content},
			}},
		}
		_, err := postsCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			panic(err)
		}
	} else {
		postDocument.Id = GenerateId()
		_, err := postsCollection.InsertOne(ctx, postDocument)
		if err != nil {
			panic(err)
		}
	}
	rnd.Redirect("/")
}

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {
	fmt.Println("Listening on port :3000")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, "mongodb://mongodb:27017")
	if err != nil {
		panic(err)
	}
	postsCollection = client.Database("testing").Collection("post")
	m := martini.Classic()
	unescacpeFuncMap := template.FuncMap{"unescape": unescape}
	m.Use(render.Renderer(render.Options{
		Directory:       "template",                 // Specify what path to load the templates from.
		Layout:          "layout",                   // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions:      []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Funcs:           []template.FuncMap{unescacpeFuncMap},
		Charset:         "UTF-8",     // Sets encoding for json and html content-types. Default is "UTF-8".
		HTMLContentType: "text/html", // Output XHTML content type instead of default "text/html"
	}))
	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)
	m.Post("/SavePost", savePostHandler)
	m.Get("/edit/:id", editHandler)
	m.Get("/delete/:id", deleteHandler)
	m.Run()
}
