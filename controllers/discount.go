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

var discountCollection *mongo.Collection
var productCollection *mongo.Collection
var collectionCollection *mongo.Collection

// InitDiscountController initializes the collections needed for the discount logic.
func InitDiscountController(db *mongo.Database) {
	discountCollection = db.Collection("discounts")
	productCollection = db.Collection("products")
	collectionCollection = db.Collection("collections")
}

// applyDiscount logic: saves original price and updates to discounted price
func applyDiscount(ctx context.Context, discount models.Discount) error {
	filter := bson.M{}
	update := bson.M{}

	// Calculate the new price
	// Formula: new_price = original_price * (100 - percentage) / 100
	newPriceFunc := func(originalPrice float64, percentage float64) float64 {
		return originalPrice - (originalPrice * percentage / 100)
	}

	if discount.TargetType == "product" {
		filter = bson.M{"_id": discount.TargetID}

		var product models.Product
		err := productCollection.FindOne(ctx, filter).Decode(&product)
		if err != nil {
			return fmt.Errorf("product not found: %w", err)
		}

		if product.OriginalPrice == nil { // Only apply discount if one isn't already active
			originalPrice := product.Price
			newPrice := newPriceFunc(originalPrice, discount.Percentage)
			update = bson.M{"$set": bson.M{"price": newPrice, "originalPrice": originalPrice}}
			_, err = productCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				return fmt.Errorf("failed to update product price: %w", err)
			}
		}

	} else if discount.TargetType == "collection" {
		filter = bson.M{"productIds": discount.TargetID} // Using productIds field for collection
		cursor, err := productCollection.Find(ctx, filter)
		if err != nil {
			return fmt.Errorf("failed to find products in collection: %w", err)
		}
		defer cursor.Close(ctx)

		var products []models.Product
		if err = cursor.All(ctx, &products); err != nil {
			return fmt.Errorf("failed to decode products: %w", err)
		}

		for _, product := range products {
			if product.OriginalPrice == nil { // Only apply discount if one isn't already active
				originalPrice := product.Price
				newPrice := newPriceFunc(originalPrice, discount.Percentage)
				update = bson.M{"$set": bson.M{"price": newPrice, "originalPrice": originalPrice}}
				_, err = productCollection.UpdateOne(ctx, bson.M{"_id": product.ID}, update)
				if err != nil {
					// Log the error and continue to the next product
					fmt.Printf("Error updating product %s: %v\n", product.ID.Hex(), err)
				}
			}
		}
	}
	return nil
}

// revertDiscount logic: restores the original price and removes the originalPrice field
func revertDiscount(ctx context.Context, discountID primitive.ObjectID) error {
	var discount models.Discount
	err := discountCollection.FindOne(ctx, bson.M{"_id": discountID}).Decode(&discount)
	if err != nil {
		return fmt.Errorf("discount not found: %w", err)
	}

	filter := bson.M{}
	update := bson.M{}

	if discount.TargetType == "product" {
		filter = bson.M{"_id": discount.TargetID, "originalPrice": bson.M{"$exists": true}}
		// Set price to originalPrice and remove the originalPrice field
		update = bson.M{
			"$set":   bson.M{"price": "$originalPrice"},
			"$unset": bson.M{"originalPrice": ""},
		}
		_, err = productCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return fmt.Errorf("failed to revert product price: %w", err)
		}

	} else if discount.TargetType == "collection" {
		filter = bson.M{"productIds": discount.TargetID, "originalPrice": bson.M{"$exists": true}}
		cursor, err := productCollection.Find(ctx, filter)
		if err != nil {
			return fmt.Errorf("failed to find products in collection: %w", err)
		}
		defer cursor.Close(ctx)

		var products []models.Product
		if err = cursor.All(ctx, &products); err != nil {
			return fmt.Errorf("failed to decode products: %w", err)
		}

		for _, product := range products {
			if product.OriginalPrice != nil {
				// Set price to originalPrice and remove the originalPrice field
				update = bson.M{"$set": bson.M{"price": *product.OriginalPrice}, "$unset": bson.M{"originalPrice": ""}}
				_, err = productCollection.UpdateOne(ctx, bson.M{"_id": product.ID}, update)
				if err != nil {
					// Log the error and continue to the next product
					fmt.Printf("Error reverting product %s: %v\n", product.ID.Hex(), err)
				}
			}
		}
	}
	return nil
}

// CreateDiscount handles creating a new discount and applying it to products
func CreateDiscount(w http.ResponseWriter, r *http.Request) {
	var discount models.Discount
	if err := json.NewDecoder(r.Body).Decode(&discount); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if discount.Percentage < 0 || discount.Percentage > 100 {
		http.Error(w, "Invalid percentage value", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Insert the discount into the database
	res, err := discountCollection.InsertOne(ctx, discount)
	if err != nil {
		http.Error(w, "Failed to create discount", http.StatusInternalServerError)
		return
	}
	discount.ID = res.InsertedID.(primitive.ObjectID)

	// 2. Apply the discount to the targeted products
	if err := applyDiscount(ctx, discount); err != nil {
		// Log the error but still return the created discount
		fmt.Printf("Error applying discount: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discount)
}

// GetDiscounts fetches all discounts from the database.
func GetDiscounts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := discountCollection.Find(ctx, bson.M{})
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

// UpdateDiscount handles updating an existing discount and re-applying the new discount percentage
func UpdateDiscount(w http.ResponseWriter, r *http.Request) {
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

	// 1. Find and revert the old discount
	if err := revertDiscount(ctx, discountID); err != nil {
		fmt.Printf("Error reverting old discount: %v\n", err)
		http.Error(w, "Failed to revert old discount", http.StatusInternalServerError)
		return
	}

	// 2. Update the discount document
	update := bson.M{
		"$set": bson.M{
			"targetType": updated.TargetType,
			"targetId":   updated.TargetID,
			"percentage": updated.Percentage,
		},
	}
	_, err = discountCollection.UpdateOne(ctx, bson.M{"_id": discountID}, update)
	if err != nil {
		http.Error(w, "Failed to update discount", http.StatusInternalServerError)
		return
	}

	// 3. Apply the new discount
	updated.ID = discountID
	if err := applyDiscount(ctx, updated); err != nil {
		fmt.Printf("Error applying new discount: %v\n", err)
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteDiscount handles deleting a discount and reverting the price changes on products
func DeleteDiscount(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	discountID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Revert the price changes on the targeted products
	if err := revertDiscount(ctx, discountID); err != nil {
		fmt.Printf("Error reverting discount: %v\n", err)
		http.Error(w, "Failed to revert discount price", http.StatusInternalServerError)
		return
	}

	// 2. Delete the discount document
	_, err = discountCollection.DeleteOne(ctx, bson.M{"_id": discountID})
	if err != nil {
		http.Error(w, "Failed to delete discount", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
