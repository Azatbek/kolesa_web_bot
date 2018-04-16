package handler

import (
	"../db"
	"net/http"
	"encoding/json"
	"fmt"
	"bytes"
	"strings"
	"io"
	"github.com/tealeg/xlsx"
	"strconv"
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

func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	var Buf bytes.Buffer

	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])

	io.Copy(&Buf, file)

	contents := Buf.Bytes()

	ReadFile(contents)

	Buf.Reset()

	return
}

func ReadFile (bytes []byte) {
	xlFile, err := xlsx.OpenBinary(bytes)
	if err != nil {
		panic(err)
	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			question := db.Questions{}
			values := []int{}
			texts  := []string{}

			for i := 0; i < len(row.Cells); i++ {
				question = db.Questions{
					Text: row.Cells[0].Value,
					Complexity: getIntFromString(row.Cells[1].Value),
					Category: getIntFromString(row.Cells[2].Value),
				}

				if i >= 3 {
					if i % 2 == 0 {
						values = append(values, getIntFromString(row.Cells[i].Value))
					} else {
						texts = append(texts, row.Cells[i].Value)
					}
				}
			}

			for j := 0; j < 4; j++ {
				variant := db.Variants{
					Text: texts[j],
					Value: values[j],
				}

				question.Variants = append(question.Variants, variant)
			}

			db.AddQuestionWithVariants(question)
		}
	}
}

func getIntFromString(str string) int {
	float, err := strconv.ParseFloat(str, 32)

	if err == nil {
		return int(float)
	}

	return -1
}