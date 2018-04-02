package handler

import (
	"../db"
	"net/http"
	"encoding/json"
	"fmt"
)

type obj map[string]interface{}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Хаю хай всем, я kolesa_web_bot!"))
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	err := db.PingDb()
	dbStatus := err == nil

	writeJsonResponse(w, obj{
		"status":   "ok",
		"db_check": dbStatus,
	})
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	writeJsonError(w, 404, nil, http.StatusNotFound)
}

func handleError(w http.ResponseWriter, err error) {
	fmt.Println(err)

	writeJsonError(w, 500, nil, http.StatusInternalServerError)
}

func writeJsonError(w http.ResponseWriter, errorCode int, details obj, statusCode int) {
	response, err := json.Marshal(obj{
		"status": statusCode,
		"error":  "",
	})

	if err != nil {
		handleError(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(statusCode)
		w.Write(response)
	}
}

func writeJsonResponse(w http.ResponseWriter, body interface{}) {
	response, err := json.Marshal(body)

	if err != nil {
		handleError(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(200)
		w.Write(response)
	}
}