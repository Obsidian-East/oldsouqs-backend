package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"oldsouqs-backend/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DiscountController wraps Mongo collections for DI
type DiscountController struct {
	Discounts *mongo.Collection
	Products  *mongo.Collection
}

// NewDiscountController constructor
func NewDiscountController(db *mongo.Database) *DiscountController {
	return &DiscountController{
		Discounts: db.Collection("discounts"),
		Products:  db.Collection("products"),
	}
}

// applyDiscount logic: saves original price and updates to discounted price
func (dc *DiscountController) applyDiscount(ctx context.Context, discount models.Discount) error {
	newPriceFunc := func(originalPrice float64, percentage float64) float64 {
		return originalPrice - (originalPrice * percentage / 100)
	}

	if discount.TargetType == "product" {
		filter := bson.M{"_id": discount.TargetID}

		var product models.Product
		err := dc.Products.FindOne(ctx, filter).Decode(&product)
		if err != nil {
			return fmt.Errorf("product not found: %w", err)
		}

		if product.OriginalPrice == nil || (product.OriginalPrice != nil && *product.OriginalPrice == 0.00) {
			originalPrice := product.Price
			newPrice := newPriceFunc(originalPrice, discount.Percentage)
			update := bson.M{"$set": bson.M{"price": newPrice, "originalPrice": originalPrice}}
			_, err = dc.Products.UpdateOne(ctx, filter, update)
			if err != nil {
				return fmt.Errorf("failed to update product price: %w", err)
			}
		}
	} else if discount.TargetType == "collection" {
		filter := bson.M{"productIds": discount.TargetID}
		cursor, err := dc.Products.Find(ctx, filter)
		if err != nil {
			return fmt.Errorf("failed to find products in collection: %w", err)
		}
		defer cursor.Close(ctx)

		var products []models.Product
		if err := cursor.All(ctx, &products); err != nil {
			return fmt.Errorf("failed to decode products: %w", err)
		}

		for _, product := range products {
			if product.OriginalPrice == nil {
				originalPrice := product.Price
				newPrice := newPriceFunc(originalPrice, discount.Percentage)
				update := bson.M{"$set": bson.M{"price": newPrice, "originalPrice": originalPrice}}
				_, err = dc.Products.UpdateOne(ctx, bson.M{"_id": product.ID}, update)
				if err != nil {
					fmt.Printf("Error updating product %s: %v\n", product.ID.Hex(), err)
				}
			}
		}
	}
	return nil
}

// revertDiscount logic: restores the original price and removes the originalPrice field
func (dc *DiscountController) revertDiscount(ctx context.Context, discountID primitive.ObjectID) error {
	var discount models.Discount
	err := dc.Discounts.FindOne(ctx, bson.M{"_id": discountID}).Decode(&discount)
	if err != nil {
		return err
	}

	if discount.TargetType == "product" {
		var product models.Product
		err := dc.Products.FindOne(ctx, bson.M{"_id": discount.TargetID}).Decode(&product)
		if err != nil {
			return err
		}

		if product.OriginalPrice != nil {
			update := bson.M{
				"$set": bson.M{
					"price":     *product.OriginalPrice,
					"updatedAt": time.Now(),
				},
				"$unset": bson.M{
					"originalPrice": "",
				},
			}

			_, err := dc.Products.UpdateByID(ctx, discount.TargetID, update)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CreateDiscount handles POST /discounts
func (dc *DiscountController) CreateDiscount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var discount models.Discount
	if err := json.NewDecoder(r.Body).Decode(&discount); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if discount.TargetType == "" || discount.TargetID.IsZero() || discount.Percentage <= 0 {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	discount.CreatedAt = time.Now()
	discount.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := dc.Discounts.InsertOne(ctx, discount)
	if err != nil {
		http.Error(w, "Failed to create discount", http.StatusInternalServerError)
		return
	}

	discount.ID = res.InsertedID.(primitive.ObjectID)

	// Apply discount logic
	if err := dc.applyDiscount(ctx, discount); err != nil {
		fmt.Printf("Warning: Failed to apply discount: %v\n", err)
	}

	json.NewEncoder(w).Encode(discount)
}

// GetDiscounts handles GET /discounts
func (dc *DiscountController) GetDiscounts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := dc.Discounts.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch discounts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var discounts []models.Discount
	if err := cursor.All(ctx, &discounts); err != nil {
		http.Error(w, "Error reading discounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discounts)
}

// UpdateDiscount handles PUT /discounts/{id}
func (dc *DiscountController) UpdateDiscount(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	discountID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updated models.Discount
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dc.revertDiscount(ctx, discountID); err != nil {
		fmt.Printf("Error reverting old discount: %v\n", err)
		http.Error(w, "Failed to revert old discount", http.StatusInternalServerError)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"targetType": updated.TargetType,
			"targetId":   updated.TargetID,
			"percentage": updated.Percentage,
			"updatedAt":  time.Now(),
		},
	}
	_, err = dc.Discounts.UpdateOne(ctx, bson.M{"_id": discountID}, update)
	if err != nil {
		http.Error(w, "Failed to update discount", http.StatusInternalServerError)
		return
	}

	updated.ID = discountID
	if err := dc.applyDiscount(ctx, updated); err != nil {
		fmt.Printf("Error applying new discount: %v\n", err)
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteDiscount handles DELETE /discounts/{id}
func (dc *DiscountController) DeleteDiscount(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	discountID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dc.revertDiscount(ctx, discountID); err != nil {
		fmt.Printf("Error reverting discount: %v\n", err)
		http.Error(w, "Failed to revert discount price", http.StatusInternalServerError)
		return
	}

	_, err = dc.Discounts.DeleteOne(ctx, bson.M{"_id": discountID})
	if err != nil {
		http.Error(w, "Failed to delete discount", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
