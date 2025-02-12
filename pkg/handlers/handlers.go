package handlers

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-faker/faker/v3"
	"github.com/snirkop89/shopping-app/pkg/models"
	"github.com/snirkop89/shopping-app/pkg/repository"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var tmpl *template.Template

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{Repo: repo}
}

func init() {
	templatesDir := "./templates"
	pattern := filepath.Join(templatesDir, "**", "*.html")
	tmpl = template.Must(template.ParseGlob(pattern))
}

func (h *Handler) SeedProducts(w http.ResponseWriter, r *http.Request) {
	// Seed the random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	numProducts := 20

	productNames := []string{"Laptop", "Smartphone", "Tablet", "Headphones", "Speaker", "Camera", "TV", "Watch", "Printer", "Monitor"}

	titler := cases.Title(language.AmericanEnglish)

	for range numProducts {
		// Generate the random but more realistic product type
		productType := productNames[rng.Intn(len(productNames))]
		productName := titler.String(faker.Word()) + " " + productType

		product := models.Product{
			ProductName:  productName,
			Price:        float64(rng.Intn(100000)) / 100, // Random price betwen 0.00 and 999.99
			Description:  faker.Sentence(),
			ProductImage: "placeholder.jpeg",
		}

		err := h.Repo.Product.CreateProduct(&product)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Error creating product %s: %v", product.ProductName, err),
				http.StatusInternalServerError,
			)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Successfully added %d dummy products", numProducts)
}
