package main

import (
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello world"))
		if err != nil {
			slog.Error("error happened")
		}
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("listening on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("shutting down server")
			os.Exit(0)
		} else {
			log.Fatalln(err)
		}
	}
}
