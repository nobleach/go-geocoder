package main

import (
	"database/sql"
	"encoding/json"
	"github.com/codegangsta/martini"
	_ "github.com/lib/pq"
	"log"
	//"net/http"
)

const connectString string = "user=services password=services dbname=services sslmode=disable host=localhost"

type Response struct {
	Path  int    `json:"path"`
	Row   int    `json:"row"`
	State string `json:"state"`
}

func main() {
	r := new(Response)
	m := martini.Classic()
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		log.Fatal(err)
	}

	m.Get("/", func() string {
		return "Welcome to Spatial Services"
	})

	m.Get("/geocode/:lat/:lng", func(params martini.Params) string {

		var state_abbr string
		var path int
		var row int
		err := db.QueryRow("SELECT states.state_abbr, pathrows.path, pathrows.row FROM states, pathrows WHERE st_contains(states.geom, ST_SetSRID(ST_MakePoint("+params["lng"]+","+params["lat"]+"), 4326)) AND st_contains(pathrows.geom, ST_SetSRID(ST_MakePoint("+params["lng"]+","+params["lat"]+"), 4326))").Scan(&state_abbr, &path, &row)
		switch {
		case err == sql.ErrNoRows:
			log.Printf("No state matches ?", state_abbr)
		case err != nil:
			log.Fatal(err)
		default:
			r.State = state_abbr
			r.Path = path
			r.Row = row
		}
		b, err := json.Marshal(r)
		if err != nil {
			log.Fatal(err)
		}

		return (string(b))
	})
	m.Run()
	// http.ListenAndServe(":8080", m)
}
