package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"example.com/orders-api/internal/fetch"
	"example.com/orders-api/internal/store"
	"example.com/orders-api/utils"
)

func main() {
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var o store.Order
			json.NewDecoder(r.Body).Decode(&o)
			o.ID = utils.NewToken()
			if err := store.Save(&o); err != nil {
				log.Printf("save: %v", err)
				http.Error(w, err.Error(), 500)
				return
			}
			utils.WriteJSON(w, o)
			return
		} else {
			id := r.URL.Query().Get("id")
			o, err := store.Get(id)
			if err != nil {
				if err.Error() == "Order not found." {
					http.Error(w, "not found", 404)
					return
				}
				http.Error(w, err.Error(), 500)
				return
			}
			utils.WriteJSON(w, o)
		}
	})
	http.HandleFunc("/prices", func(w http.ResponseWriter, r *http.Request) {
		urls := strings.Split(r.URL.Query().Get("src"), ",")
		bodies, err := fetch.FetchAll(urls)
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		utils.WriteJSON(w, len(bodies))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
