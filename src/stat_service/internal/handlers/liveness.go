package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func liveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("liveliness probe")

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func StartLivenessHandler(port int) error {
	http.HandleFunc("/", liveHandler)

	log.Println("Stat service started")

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		return err
	}
	return nil
}
