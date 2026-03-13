package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"google.golang.org/genai"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleGenerate)
	mux.HandleFunc("/healthz", handleHealthz)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to initialize genai client: %v", err), http.StatusInternalServerError)
		return
	}

	prompt := strings.TrimSpace(r.URL.Query().Get("q"))
	if prompt == "" {
		prompt = "Explain how AI works in a few words"
	}

	stream := client.Models.GenerateContentStream(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for chunk, err := range stream {
		if err != nil {
			http.Error(w, fmt.Sprintf("generation failed: %v", err), http.StatusInternalServerError)
			return
		}
		if _, err := io.WriteString(w, chunk.Text()); err != nil {
			log.Printf("failed to write response chunk: %v", err)
			return
		}
	}
}
