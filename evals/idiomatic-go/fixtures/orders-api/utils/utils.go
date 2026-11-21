package utils

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
)

func NewToken() string {
	return strconv.FormatInt(rand.Int63(), 36)
}

func WriteJSON(w http.ResponseWriter, v interface{}) {
	json.NewEncoder(w).Encode(v)
}
