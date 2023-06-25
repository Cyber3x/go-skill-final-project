package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func JMBAGHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(ConfigData.JMBAG)
}

type OpResponse struct {
	A      int `json:"a"`
	B      int `json:"b"`
	Result int `json:"result"`
}

func SumHandler(w http.ResponseWriter, req *http.Request) {
	queryData := req.URL.Query()

	// are there two params
	if len(queryData) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// does every param only have one value
	if len(queryData["a"]) != 1 || len(queryData["b"]) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	a, err1 := strconv.Atoi(queryData["a"][0])
	b, err2 := strconv.Atoi(queryData["b"][0])

	// are query params numbers
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// calculate and return
	// Here we could have integer overflow
	// but I don't know if i should worry about that in this project
	sum := a + b

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&OpResponse{a, b, sum})
}

func MultiplyHandler(w http.ResponseWriter, req *http.Request) {
	queryData := req.URL.Query()

	// are there two params
	if len(queryData) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// does every param only have one value
	if len(queryData["a"]) != 1 || len(queryData["b"]) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	a, err1 := strconv.Atoi(queryData["a"][0])
	b, err2 := strconv.Atoi(queryData["b"][0])

	// are query params numbers
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// calculate and return
	// Here we could have integer overflow
	// but I don't know if i should worry about that in this project
	sum := a * b

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&OpResponse{a, b, sum})
}

func FetchHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type Body struct {
		URL string `json:"url"`
	}

	bodyData, err := io.ReadAll(req.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := &Body{}
	err = json.Unmarshal(bodyData, &body)

	if err != nil || body.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := http.Get(body.URL)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.Header)
}

func JMBAGFileHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		if _, err := os.Stat(StudentFilePath); errors.Is(err, os.ErrNotExist) {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Target file not found\n")
			return
		}

		bodyRaw, err := os.ReadFile(StudentFilePath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain")
		w.Write(bodyRaw)

	} else if req.Method == http.MethodPost {
		bodyRaw, err := io.ReadAll(req.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		os.WriteFile(StudentFilePath, bodyRaw, 0644)

		w.WriteHeader(http.StatusCreated)
		return

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}
