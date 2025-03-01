package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/snirkop89/mx-store/pkg/models"
)

var (
	currentCartOrderID uuid.UUID
	cartItems          []models.OrderItem
)

func (h *Handler) ShoppingHomepage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		OrderItems []models.OrderItem
	}{
		OrderItems: cartItems,
	}

	tmpl.ExecuteTemplate(w, "homepage", data)
}

func (h *Handler) ShoppingItemsView(w http.ResponseWriter, r *http.Request) {
	// Fake latency
	time.Sleep(2 * time.Second)

	products, err := h.Repo.Product.GetProducts("product_image != ''")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "shoppingItems", products)
}
