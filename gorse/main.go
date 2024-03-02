package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/zhenghaoz/gorse/client"
)

// Logger middleware
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			time.Since(start),
			r.UserAgent(),
		)
	})
}

func main() {
	r := mux.NewRouter()

	// Middlewares
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(cors.AllowAll().Handler)
	r.Use(Logger)

	// Create the client
	gorse := client.NewGorseClient("http://127.0.0.1:8088", "api_key")

	// Routes
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Welcome to port 6666")
	})

	// User

	// Insert user
	insertUser := func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			UserID  string   `json:"UserId"`
			Comment string   `json:"Comment"`
			Labels  []string `json:"Labels"`
		}
		json.NewDecoder(r.Body).Decode(&data)
		defer r.Body.Close()

		// Create a context
		ctx := context.Background()

		gorse.InsertUser(ctx, client.User{
			UserId:  data.UserID,
			Comment: data.Comment,
			Labels:  data.Labels,
		})

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	}

	r.HandleFunc("/user/insert", insertUser).Methods("POST")

	// Update user
	updateUser := func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			UserID  string   `json:"UserId"`
			Comment string   `json:"Comment"`
			Labels  []string `json:"Labels"`
		}
		json.NewDecoder(r.Body).Decode(&data)
		defer r.Body.Close()

		// Create a context
		ctx := context.Background()

		userPatch := client.UserPatch{
			Comment: &data.Comment,
			Labels:  data.Labels,
		}
		gorse.UpdateUser(ctx, data.UserID, userPatch)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	}

	r.HandleFunc("/user/update", updateUser).Methods("PATCH")

	// Item

	// insert item
	insertItem := func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			ItemID     string   `json:"ItemId"`
			Comment    string   `json:"Comment"`
			Categories []string `json:"Categories"`
			Labels     []string `json:"Labels"`
		}
		json.NewDecoder(r.Body).Decode(&data)
		defer r.Body.Close()

		// Create a context
		ctx := context.Background()

		gorse.InsertItem(ctx, client.Item{
			ItemId:     data.ItemID,
			Comment:    data.Comment,
			Categories: data.Categories,
			Labels:     data.Labels,
			Timestamp:  time.Now().Format(time.RFC3339),
		})

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	}

	r.HandleFunc("/item/insert", insertItem).Methods("POST")

	// Update item
	updateItem := func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			ItemID     string   `json:"ItemId"`
			Comment    string   `json:"Comment"`
			IsHidden   bool     `json:"IsHidden"`
			Categories []string `json:"Categories"`
			Labels     []string `json:"Labels"`
		}
		json.NewDecoder(r.Body).Decode(&data)
		defer r.Body.Close()

		// Create a context
		ctx := context.Background()

		timestamp := time.Now()
		itemPatch := client.ItemPatch{
			Comment:    &data.Comment,
			Categories: data.Categories,
			Labels:     data.Labels,
			Timestamp:  &timestamp,
		}
		gorse.UpdateItem(ctx, data.ItemID, itemPatch)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	}

	r.HandleFunc("/item/update", updateItem).Methods("PATCH")

	// Feedback

	// Insert feedbacks
	insertFeedbacks := func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			FeedbackType string `json:"FeedbackType"`
			UserID       string `json:"UserId"`
			ItemID       string `json:"ItemId"`
		}
		json.NewDecoder(r.Body).Decode(&data)
		defer r.Body.Close()

		// Create a context
		ctx := context.Background()

		gorse.InsertFeedback(ctx, []client.Feedback{
			{
				FeedbackType: data.FeedbackType,
				UserId:       data.UserID,
				ItemId:       data.ItemID,
				Timestamp:    time.Now().Format(time.RFC3339),
			},
		})

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	}

	r.HandleFunc("/feedback/insert", insertFeedbacks).Methods("POST")

	// Recommend

	// Get recommend
	getRecommend := func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			N      int    `json:"n"`
			UserID string `json:"userId"`
		}
		json.NewDecoder(r.Body).Decode(&data)
		defer r.Body.Close()

		// Create a context
		ctx := context.Background()

		rec, _ := gorse.GetRecommend(ctx, data.UserID, "", data.N)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rec)
	}

	r.HandleFunc("/recommend/get", getRecommend).Methods("GET")

	// Start server
	port := ":6666"
	fmt.Println("Listening on port", port)
	http.ListenAndServe(port, r)
}
