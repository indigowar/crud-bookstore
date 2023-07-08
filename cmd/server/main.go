package main

import (
	"bookstore/internal/delivery/handlers"
	"bookstore/internal/domain/services"
	storage "bookstore/internal/storages/book"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.Default()

	storage := storage.NewInMemoryBookStorage()

	bookSvc := services.NewBookService(logger, storage)

	bookEndpoint := handlers.NewBookEndpoint(bookSvc)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/books", func(r chi.Router) {
		r.Get("/", bookEndpoint.GetAllBooks)
		r.Get("/{id}", bookEndpoint.GetSpecificBook)
		r.Post("/", bookEndpoint.CreateBook)
		r.Delete("/{id}", bookEndpoint.DeleteBook)
		r.Put("/{id}", bookEndpoint.UpdateBook)
	})

	server := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(err.Error())
		log.Fatal("Failed to stop the server. Forced to shut it down.")
	}

	log.Println("Server is shot down.")
}
