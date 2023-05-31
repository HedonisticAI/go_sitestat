package http_handler

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/HedonisticAI/go_sitestat/db_handler"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DBInstance struct {
	Red  *redis.Client
	Post *gorm.DB
}

func getlist(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines, nil
}
func GetTime(db *DBInstance) {
	var Max time.Duration
	var Min time.Duration
	ctx := context.TODO()
	var sites []string
	sites, err := getlist(os.Getenv("FILE_NAME"))
	if err != nil {
		panic("could not read file")
	}

	for {
		Min = 0
		Max = 0

		for _, site := range sites {
			now := time.Now()
			val, err := http.Get("https://" + site)
			if err == nil && val.StatusCode == http.StatusOK {
				time := time.Since(now)
				err = db.Red.Set(ctx, site, time.String(), 0).Err()
				if err != nil {
					log.Println(err.Error())
				}
				if time >= Max {
					db.Red.Set(ctx, "Max is "+site, time.String(), 0)
					Max = time
				}
				if time <= Min {
					db.Red.Set(ctx, "Min is "+site, time.String(), 0)
					Min = time
				}
			}
			db.Red.Set(ctx, site, val.Status, 0)

		}
		log.Println("Data collected, sleeping")
		time.Sleep(1 * time.Minute)

	}

}
func (collection *DBInstance) Site(w http.ResponseWriter, r *http.Request) {
	var rec db_handler.Record
	ctx := context.TODO()
	params := mux.Vars(r)
	res := collection.Red.Get(ctx, params["sitename"])

	io.WriteString(w, res.String())
	collection.Post.Where(db_handler.Record{Site: params["sitename"]}).Assign(db_handler.Record{Counter: rec.Counter + 1}).FirstOrCreate(&rec)
}
func (collection *DBInstance) Max(w http.ResponseWriter, r *http.Request) {
	var rec db_handler.Record
	ctx := context.TODO()
	res := collection.Red.Get(ctx, "Max")

	io.WriteString(w, res.String())
	collection.Post.Where(db_handler.Record{Site: "Max"}).Assign(db_handler.Record{Counter: rec.Counter + 1}).FirstOrCreate(&rec)
}
func (collection *DBInstance) Min(w http.ResponseWriter, r *http.Request) {
	var rec db_handler.Record
	ctx := context.TODO()
	res := collection.Red.Get(ctx, "Min")

	io.WriteString(w, res.String())
	collection.Post.Where(db_handler.Record{Site: "Min"}).Assign(db_handler.Record{Counter: rec.Counter + 1}).FirstOrCreate(&rec)
}
func (collection *DBInstance) AdmSite(w http.ResponseWriter, r *http.Request) {
	var rec db_handler.Record
	params := mux.Vars(r)
	collection.Post.Where("name = ?", params["sitename"]).First(&rec)

	io.WriteString(w, strconv.Itoa(int(rec.Counter)))
}
func (collection *DBInstance) AdmMin(w http.ResponseWriter, r *http.Request) {
	var rec db_handler.Record
	collection.Post.Where("name = ?", "Min").First(&rec)
	io.WriteString(w, strconv.Itoa(int(rec.Counter)))

}
func (collection *DBInstance) AdmMax(w http.ResponseWriter, r *http.Request) {
	var rec db_handler.Record
	collection.Post.Where("name = ?", "Max").First(&rec)
	io.WriteString(w, strconv.Itoa(int(rec.Counter)))

}
