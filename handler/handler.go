package handler

import (
	"dkubanyi/urlShortener/storage"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type handler struct {
	prefix      string
	accessToken string
	storage     storage.Service
}

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type errorResponse struct {
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
}

func New(prefix string, accessToken string, storage storage.Service) http.Handler {
	mux := http.NewServeMux()
	h := handler{prefix, accessToken, storage}

	mux.HandleFunc("/", h.redirect)
	mux.HandleFunc("/encode", responseHandler(h.encode))
	return mux
}

func responseHandler(h func(io.Writer, *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, status, err := h(w, r)
		if err != nil {
			data = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		if err != nil {
			err = json.NewEncoder(w).Encode(errorResponse{Error: data, Success: err == nil})
		} else {
			err = json.NewEncoder(w).Encode(response{Data: data, Success: err == nil})
		}

		if err != nil {
			log.Printf("could not encode response to output: %v", err)
		}
	}
}

func (h handler) encode(w io.Writer, r *http.Request) (interface{}, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("method %s not allowed", r.Method)
	}

	if r.Header.Get("access_token") != h.accessToken {
		return nil, http.StatusForbidden, fmt.Errorf("access_token is not correct")
	}

	var input struct{ URL string }
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("unable to decode JSON request body: %v", err)
	}

	url := strings.TrimSpace(input.URL)
	if url == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("URL is empty")
	}

	if !strings.Contains(url, "http") {
		url = "http://" + url
	}

	c, err := h.storage.Save(url)

	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Could not store in database: %v", err)
	}

	fmt.Printf("%v", h.prefix+"/"+c)

	return h.prefix + "/" + c, http.StatusCreated, nil
}

func (h handler) redirect(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	code := r.URL.Path[len("/"):]

	getParams := r.URL.Query().Encode()

	url, err := h.storage.Load(code)

	if getParams != "" {
		url = url + "?" + getParams
	}

	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("URL Not Found"))
		return
	}

	http.Redirect(rw, r, url, http.StatusMovedPermanently)
}
