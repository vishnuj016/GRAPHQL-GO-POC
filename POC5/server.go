package main

import (
	"bookapp/graph"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

const defaultPort = "8080"

// Define the WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("ui"))))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to Upgrade connection to websocket")
			return
		}

		defer conn.Close()

		fmt.Println("WS")

		// Kafka Consumer
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{"localhost:9092"},
			Topic:    "crud-events",
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		})

		defer reader.Close()

		for {
			// Read Kafka message
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Println(err)
				return
			}

			msg := string(m.Value)
			fmt.Println(msg)

			if err := conn.WriteMessage(websocket.TextMessage, m.Value); err != nil {
				log.Println(err)
				return
			}
		}
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
