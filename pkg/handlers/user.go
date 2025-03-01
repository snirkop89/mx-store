package handlers

import (
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func (h *Handler) CartView(w http.ResponseWriter, r *http.Request) {
	data := struct {
		OrderItems []models.OrderItem
		Message    string
		AlertType  string
		TotalCost  float64
	}{
		OrderItems: cartItems,
		Message:    "",
		AlertType:  "",
		TotalCost:  getTotalCartCost(),
	}

	tmpl.ExecuteTemplate(w, "cartItems", data)
}

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := uuid.Parse(vars["product_id"])
	if err != nil {
		http.Error(w, "Invlalid product ID", http.StatusBadRequest)
		return
	}

	// Generate a new order id for the session if one does not exist
	if currentCartOrderID == uuid.Nil {
		currentCartOrderID = uuid.New()
	}

	var exists bool
	for _, item := range cartItems {
		if item.ProductID == productID {
			exists = true
			break
		}
	}

	product, err := h.Repo.Product.GetProductByID(productID)
	if err != nil {
		http.Error(w, "Failed to get product", http.StatusInternalServerError)
		return
	}

	var cartMessage string
	var alertType string
	if !exists {
		// Create a new order item
		newOrderItem := models.OrderItem{
			OrderID:   currentCartOrderID,
			ProductID: productID,
			Quantity:  1,
			Product:   *product,
		}

		// Add new order items to the array
		cartItems = append(cartItems, newOrderItem)
		cartMessage = product.ProductName + " successfully added"
		alertType = "success"
	} else {
		cartMessage = product.ProductName + " already exists in cart"
		alertType = "danger"
	}

	data := struct {
		OrderItems []models.OrderItem
		Message    string
		AlertType  string
		TotalCost  float64
	}{
		OrderItems: cartItems,
		Message:    cartMessage,
		AlertType:  alertType,
		TotalCost:  getTotalCartCost(),
	}

	tmpl.ExecuteTemplate(w, "cartItems", data)
}

func (h *Handler) ShoppingCartView(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "shoppingCart", cartItems)
}

func (h *Handler) UpdateOrderItemQuantity(w http.ResponseWriter, r *http.Request) {
	// Get product ID and action from URL parameters
	var cartMessage string
	var refreshCartList bool // Signals a refresh of cart items when an item is removed

	productID, err := uuid.Parse(r.URL.Query().Get("product_id"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	action := r.URL.Query().Get("action")

	// find the order item
	itemIndex := -1
	for i, item := range cartItems {
		if item.ProductID == productID {
			itemIndex = i
			break
		}
	}
	if itemIndex == -1 {
		http.Error(w, "Product not found in order", http.StatusNotFound)
		return
	}

	// Update quantity based on action
	switch action {
	case "add":
		cartItems[itemIndex].Quantity++
	case "subtract":
		cartItems[itemIndex].Quantity--
		if cartItems[itemIndex].Quantity == 0 {
			// Remove items if quantity is 0
			cartItems = slices.Delete(cartItems, itemIndex, itemIndex+1)
			refreshCartList = true
		}
	case "remove":
		// Remove items regardless of quantity
		cartItems = slices.Delete(cartItems, itemIndex, itemIndex+1)
		refreshCartList = true
	default:
		cartMessage = "Invalid Action"
	}

	data := struct {
		OrderItems       []models.OrderItem
		Message          string
		AlertType        string
		TotalCost        float64
		Action           string
		RefreshCartItems bool
	}{
		OrderItems:       cartItems,
		Message:          cartMessage,
		AlertType:        "info",
		TotalCost:        getTotalCartCost(),
		Action:           action,
		RefreshCartItems: refreshCartList,
	}

	tmpl.ExecuteTemplate(w, "updateShoppingCart", data)
}

func getTotalCartCost() float64 {
	totalCost := 0.0
	for _, item := range cartItems {
		totalCost += float64(item.Quantity) * item.Product.Price
	}
	return totalCost
}
