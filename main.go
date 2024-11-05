package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	article "github.com/Ethanol48/medium-api-library/article"
	user "github.com/Ethanol48/medium-api-library/user"
	"github.com/Ethanol48/medium-api-library/utilities"
	"golang.org/x/net/html"
)

type ApiResponse struct {
    Message string `json:"message,omitempty"`
    Data    any    `json:"data,omitempty"` // Use `any` for flexible data types or replace with specific types
}


func main()  {
  mux := http.NewServeMux()

  go utilities.SpinUp("../medium-api-library/testing/")


  mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
    // url parameter
    usr := r.URL.Query().Get("usr")
    if (usr == "") {
      w.WriteHeader(422)
      fmt.Fprint(w, "we couldn't find a user in your request :(")
      return
    }

    response := ApiResponse{
      Message: usr,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)

    json.NewEncoder(w).Encode(response)
  })


  /* user funcs */
  mux.HandleFunc("GET /user/metadata", func(w http.ResponseWriter, r *http.Request) {
    // url parameter
    usr := r.URL.Query().Get("usr")
    if (usr == "") {
      w.WriteHeader(422)
      fmt.Fprint(w, "we couldn't find a user in your request :(")
      return
    }

    // create a user link
    userLink := fmt.Sprintf("https://medium.com/%s", usr)
    metadata := user.GetUserMetadata(userLink)

    fmt.Printf("metadata: %v\n", metadata)
  })


  /* article funcs */
  mux.HandleFunc("GET /article/html", func(w http.ResponseWriter, r *http.Request) {
    // url parameter
    link := r.URL.Query().Get("link")
    if (link == "") {
      w.WriteHeader(422)
      fmt.Fprint(w, "we couldn't find a link in your request :(")
      return
    }

    // TODO: validation for link

    // art := article.GetArticle(link)
    art := article.GetArticle("http://localhost:8080")
    // w.WriteHeader(200)
    //
    // fmt.Fprint(w, art.ToHTML())

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)


    resp := ApiResponse{
      Message: "",
      Data: html.UnescapeString(art.ToHTML()),
    }

    json.NewEncoder(w).Encode(resp)
  })

  mux.HandleFunc("GET /article/markdown", func(w http.ResponseWriter, r *http.Request) {
    // url parameter
    link := r.URL.Query().Get("link")
    if (link == "") {
      w.WriteHeader(422)
      fmt.Fprint(w, "we couldn't find a link in your request :(")
      return
    }

    // TODO: validation for link
    art := article.GetArticle("http://localhost:8080")

    resp := ApiResponse{
      Message: "",
      Data: html.UnescapeString(art.ToMarkdown()),
    }

    json.NewEncoder(w).Encode(resp)
  })

  mux.HandleFunc("GET /article/metadata", func(w http.ResponseWriter, r *http.Request) {

    // url parameter
    link := r.URL.Query().Get("link")
    if (link == "") {
      w.WriteHeader(422)
      fmt.Fprint(w, "we couldn't find a link in your request :(")
      return
    }

    // TODO: validation for link
    art := article.GetArticle("http://localhost:8080")

    type MetadataResponse struct {
      Title string     `json:"title"`;
      Tags []string    `json:"tags"`;
      ReadTime string  `json:"readtime"`;
      Published string `json:"published"`;
    }

    resp := MetadataResponse{
      Title: art.Title,
      Tags: art.Tags,
      ReadTime: art.ReadTime,
      Published: art.Published,
    }

    json.NewEncoder(w).Encode(resp)
  })


  fmt.Println("Serving @ http://localhost:1234 !!!")
  if err := http.ListenAndServe("localhost:1234", mux); err != nil {
    fmt.Println(err.Error())
  }
}
