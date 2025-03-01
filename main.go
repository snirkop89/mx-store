package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/snirkop89/mx-store/pkg/handlers"
	"github.com/snirkop89/mx-store/pkg/repository"
)

var db *sql.DB

func initDB() {
	var err error
	// Initialize the db variable
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("missing DSN in ENV")
	}
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Check the database connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()

	// Setup MySQL
	initDB()
	defer db.Close()

	// Setup Static folder for static files and images
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	repo := repository.NewRepository(db)
	handler := handlers.NewHandler(repo)

	// User shopping Routes
	r.HandleFunc("/", handler.ShoppingHomepage).Methods("GET")
	r.HandleFunc("/shoppingitems", handler.ShoppingItemsView).Methods("GET")
	r.HandleFunc("/cartitems", handler.CartView).Methods("GET")
	r.HandleFunc("/addtocart/{product_id}", handler.AddToCart).Methods("POST")
	r.HandleFunc("/gotocart", handler.ShoppingCartView).Methods("GET")
	r.HandleFunc("/updateorderitem", handler.UpdateOrderItemQuantity).Methods("PUT")

	// Admin Routes
	r.HandleFunc("/seed-products", handler.SeedProducts).Methods("POST")
	r.HandleFunc("/manageproducts", handler.ProductsPage).Methods("GET")
	r.HandleFunc("/allproducts", handler.AllProductsView).Methods("GET")
	r.HandleFunc("/products", handler.ListProducts).Methods("GET")
	r.HandleFunc("/products/{id}", handler.GetProduct).Methods("GET")
	r.HandleFunc("/createproduct", handler.CreateProductView).Methods("GET")
	r.HandleFunc("/products", handler.CreateProduct).Methods("POST")
	r.HandleFunc("/editproduct/{id}", handler.EditProductView).Methods("GET")
	r.HandleFunc("/products/{id}", handler.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", handler.DeleteProduct).Methods("DELETE")

	slog.Info("Starting server", "addr", ":5000")
	http.ListenAndServe(":5000", r)
}
