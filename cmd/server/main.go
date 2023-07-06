package main

import (
	"bookstore/internal/delivery/http/handlers"
	"bookstore/internal/services"
	inmem "bookstore/internal/storages/in_memory/book"
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

	storage := inmem.NewInMemoryBookStorage()

	bookSvc := services.NewBookService(logger, storage)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/books", func(r chi.Router) {
		r.Get("/", handlers.MakeGetAllBooksHandler(bookSvc))
		r.Get("/{id}", handlers.MakeGetSpecificBookHandler(bookSvc))
		r.Post("/", handlers.MakeCreateBookHandler(bookSvc))
		r.Delete("/{id}", handlers.MakeDeleteBookHandler(bookSvc))
		r.Put("/", handlers.MakeUpdateBookHandler(bookSvc))
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
