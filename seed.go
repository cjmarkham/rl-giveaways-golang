package main

import (
  "log"
  "database/sql"
  "github.com/lib/pq"
)

var db *sql.DB

type item struct {
  Name string
  Painted string
  Certified string
  Image_url string
}

func init () {
  connStr := "user=carl dbname=rl_giveaways sslmode=disable"

  var err error
  db, err = sql.Open("postgres", connStr)
  if err != nil {
    log.Fatal(err)
  }

  db.QueryRow("DELETE FROM items")
  log.Printf("Truncated database")
}

func main () {
  seedWheels()
  seedBoosts()
}

func seedWheels () {
  wheels := []item{}
  wheels = append(
    wheels,
    item{
      "Ara-51",
      "Saffron",
      "Playmaker",
      "https://rocket-league.com/content/media/items/avatar/220px/00afc3099b1481307261.png",
    },
  )
  wheels = append(
    wheels,
    item{
      "Balla-Cara",
      "Saffron",
      "",
      "https://rocket-league.com/content/media/items/avatar/220px/51b60a6a0b1518047710.png",
    },
  )
  wheels = append(
    wheels,
    item{
      "Centro",
      "",
      "Playmaker",
      "https://rocket-league.com/content/media/items/avatar/220px/5595908d8e1527618402.png",
    },
  )
  wheels = append(
    wheels,
    item{
      "Chrono",
      "",
      "",
      "https://rocket-league.com/content/media/items/avatar/220px/5a5c1c67b51506795253.png",
    },
  )

  populate("wheels", wheels)
}

func seedBoosts () {
  boosts := []item{}
  boosts = append(
    boosts,
    item{
      "Cirrus",
      "",
      "",
      "https://rocket-league.com/content/media/items/avatar/220px/1beceb38bb1527614635.png",
    },
  )
  boosts = append(
    boosts,
    item{
      "Comet",
      "",
      "",
      "https://rocket-league.com/content/media/items/avatar/220px/574d0158961522777189.png",
    },
  )

  populate("boosts", boosts)
}

func populate (category string, items []item) {
  txn, _ := db.Begin()
  stmt, err := txn.Prepare(pq.CopyIn("items", "name", "painted", "certified", "image_url", "category"))
  if err != nil {
    log.Fatal(err)
  }

  for _, item := range items {
    _, err = stmt.Exec(item.Name, item.Painted, item.Certified, item.Image_url, category)
    if err != nil {
      log.Fatal(err)
    }
  }
  log.Printf("Seeded %s", category)

  err = stmt.Close()
  if err != nil {
    log.Fatal(err)
  }

  err = txn.Commit()
  if err != nil {
    log.Fatal(err)
  }
}
