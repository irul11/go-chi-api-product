package main

import (
	"learn-go-chi/controller"
	"learn-go-chi/database"
	"learn-go-chi/models"
	"net/http"
	"sync"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func init() {
	database.LoadDatabase()

	controller.Channel = make(chan models.Message)
}

func main() {
	defer close(controller.Channel)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Route("/product", func(r chi.Router) {
		r.With(simpleMiddleware).Get("/", controller.GetProduct)
		r.With(simpleMiddleware).Get("/{productId}", controller.GetProductById)
		r.With(simpleMiddleware).Post("/", controller.CreateProduct)
		r.With(simpleMiddleware).Put("/{productId}", controller.UpdateProduct)
		r.With(simpleMiddleware).Delete("/{productId}", controller.DeleteProduct)
	})

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go controller.Exec(&wg, i)
	}

	http.ListenAndServe(":3000", r)
	wg.Wait()
}

func simpleMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
