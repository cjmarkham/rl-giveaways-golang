package main

import (
  "log"
  "net/http"
  "github.com/thedevsaddam/renderer"
  "github.com/gorilla/mux"
  "database/sql"
  _ "github.com/lib/pq"
  "strings"
  "github.com/rl-giveaways/logger"
)

var rend *renderer.Render
var db *sql.DB

func init () {
  rend = renderer.New()

  connStr := "user=carl dbname=rl_giveaways sslmode=disable"

  var err error
  db, err = sql.Open("postgres", connStr)
  logger.Error(err)
}

func main () {
  router := mux.NewRouter()

  router.HandleFunc("/", Index)
  router.HandleFunc("/{category:[a-z]+}", Category)
  http.Handle("/", router)

  log.Fatal(http.ListenAndServe(":8080", nil))
}

func Index (w http.ResponseWriter, r *http.Request) {
  tpls := []string{"view/layout.html", "view/index.html"}

  err := rend.Template(w, http.StatusOK, tpls, nil)
  logger.Error(err)
}

func Category (w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  category := vars["category"]

  rows, err := db.Query("SELECT name, painted, certified, image_url FROM items WHERE category=$1", category)
  logger.Error(err)

  defer rows.Close()

  type row struct {
    Name string
    Painted string
    Certified string
    Image_url string
  }

  items := []row{}

  for rows.Next() {
    var r row
    err := rows.Scan(&r.Name, &r.Painted, &r.Certified, &r.Image_url)
    logger.Error(err)

    items = append(items, r)
  }
  logger.Error(rows.Err())

  tpls := []string{"view/layout.html", "view/items.html"}

  variables := struct {
    Title string
    Items []row
  }{ strings.Title(category), items }

  err = rend.Template(w, http.StatusOK, tpls, variables)
  logger.Error(err)
}
