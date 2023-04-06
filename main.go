package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jubila-teresa/assignment/models"
)

// Using in-memory map to store products for simplicity
var products = make(map[uint64]models.Product)
var productID uint64

// Using in-memory map to store orders for simplicity
var orders = make(map[uint64]models.Order)
var orderID uint64

func main() {
	router := gin.Default()

	router.GET("/orders", getOrders)
	router.GET("/orders/{id}", getOrderByID)
	router.GET("/products", getProducts)
	router.GET("/products/{id}", getProductByID)

	router.POST("/orders", createOrder)
	router.POST("/products", addProducts)

	router.PATCH("/orders/:id", updateOrderByID)

	router.Run("localhost:8080")
}

// Get all the products
func getProducts(c *gin.Context) {
	fmt.Println("Getting product details.")
	if len(products) == 0 {
		fmt.Println("No products added.")
		c.Status(http.StatusNoContent)
		return
	}
	resp := []models.ProductResponse{}
	for i, item := range products {
		resp = append(resp, models.ProductResponse{
			ID:          i,
			ProductName: item.ProductName,
			Price:       item.Price,
			Category:    item.Category,
			Quantity:    item.Quantity,
		})
	}
	c.IndentedJSON(http.StatusOK, resp)
}

// Get Product with matching Order
func getProductByID(c *gin.Context) {
	fmt.Println("Getting product details")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Id"})
		return
	}

	if resp, ok := products[id]; ok {
		c.IndentedJSON(http.StatusOK, resp)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
}

// Add a new product to the catalog
func addProducts(c *gin.Context) {
	fmt.Println("Adding product details.")
	var item models.Product
	if err := c.BindJSON(&item); err != nil {
		return
	}

	if !item.Category.IsValid() {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid category"})
		return
	}
	productID++
	products[productID] = item
	resp := models.ProductResponse{
		ID:          productID,
		ProductName: item.ProductName,
		Price:       item.Price,
		Category:    item.Category,
		Quantity:    item.Quantity,
	}
	c.IndentedJSON(http.StatusCreated, resp)
}

// Get all orders
func getOrders(c *gin.Context) {
	fmt.Println("Getting order details.")
	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	resp := []models.OrderResponse{}
	for i, item := range orders {
		resp = append(resp, models.OrderResponse{
			ID:           i,
			DispatchDate: item.DispatchDate,
			Status:       item.Status,
			Products:     item.Products,
			Total:        item.Total,
		})
	}
	c.IndentedJSON(http.StatusOK, resp)
}

// Get order with matching Order
func getOrderByID(c *gin.Context) {
	fmt.Println("Getting order details")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Id"})
		return
	}

	if resp, ok := orders[id]; ok {
		c.IndentedJSON(http.StatusOK, resp)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Order not found"})
}

// Update order with matching Order
func updateOrderByID(c *gin.Context) {
	fmt.Println("Updating order status")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Id"})
		return
	}
	var r models.UpdateOrder
	if err = c.BindJSON(&r); err != nil {
		fmt.Println("Invalid request.")
		return
	}

	if !r.Status.IsValid() {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid status"})
		return
	}

	dispatchDate, err := time.Parse("02-01-2006", r.DispatchDate)
	if r.Status == models.Dispatched && err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Dispatch date"})
		return
	}

	if resp, ok := orders[id]; ok {
		resp.Status = r.Status
		if r.Status == models.Dispatched {
			resp.DispatchDate = &dispatchDate
		}
		orders[id] = resp
		c.IndentedJSON(http.StatusOK, resp)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Order not found"})
}

// Create a new order
func createOrder(c *gin.Context) {
	fmt.Println("Creating new order.")
	var r models.OrderRequest
	if err := c.BindJSON(&r); err != nil {
		return
	}

	sort.Slice(r.Products, func(i, j int) bool {
		return r.Products[i].ProductID < r.Products[j].ProductID
	})

	for i, item := range r.Products {
		if i+1 < len(r.Products) && r.Products[i].ProductID == r.Products[i+1].ProductID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate entries found"})
			return
		}
		if item.Quantity > 10 || products[item.ProductID].Quantity < item.Quantity {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Quantity"})
			return
		}
	}

	prem := 0
	orderID++
	resp := models.Order{Status: models.Placed, Products: []models.Product{}}
	for _, i := range r.Products {
		temp := products[i.ProductID]
		temp.Quantity -= i.Quantity
		products[i.ProductID] = temp
		if products[i.ProductID].Category == models.Premium {
			prem++
		}
		prd := models.Product{ProductName: products[i.ProductID].ProductName,
			Price: products[i.ProductID].Price, Category: products[i.ProductID].Category, Quantity: i.Quantity}
		resp.Products = append(resp.Products, prd)
		resp.Total += products[i.ProductID].Price * float32(i.Quantity)
	}
	if prem > 2 {
		fmt.Println("Adding discount.")
		resp.Total = 0.9 * resp.Total
	}
	orders[orderID] = resp
	res := models.OrderResponse{ID: orderID, Status: resp.Status, Products: resp.Products, Total: resp.Total}
	c.IndentedJSON(http.StatusCreated, res)
}
