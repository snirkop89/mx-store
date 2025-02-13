package handlers

import (
	"fmt"
	"html/template"
	"math"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
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

func (h *Handler) ProductsPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "products", nil)
}

func (h *Handler) AllProductsView(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "allProducts", nil)
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	products, err := h.Repo.Product.ListProducts(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalProducts, err := h.Repo.Product.GetTotalProductsCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalProducts) / float64(limit)))
	prevPage := page - 1
	nextPage := page + 1
	pageButtonsRange := makeRange(1, totalPages)

	data := struct {
		Products         []models.Product
		CurrentPage      int
		TotalPages       int
		Limit            int
		PreviousPage     int
		NextPage         int
		PageButtonsRange []int
	}{
		Products:         products,
		CurrentPage:      page,
		TotalPages:       totalPages,
		PreviousPage:     prevPage,
		NextPage:         nextPage,
		PageButtonsRange: pageButtonsRange,
	}

	// Fake latency
	// time.Sleep(4 * time.Second)
	tmpl.ExecuteTemplate(w, "productRows", data)
}

func makeRange(min, max int) []int {
	rangeArray := make([]int, max-min+1)
	for i := range rangeArray {
		rangeArray[i] = min + i
	}
	return rangeArray
}
