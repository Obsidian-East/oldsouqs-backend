package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"oldsouqs-backend/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func validateProduct(product models.Product) error {
	if product.Sku == "" {
		return fmt.Errorf("SKU is required")
	}
	if product.Title == "" {
		return fmt.Errorf("Title is missing")
	}
	if product.Price == 0.0 {
		return fmt.Errorf("Price is missing")
	}
	return nil
}

func CreateProduct(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	var product models.Product

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate product
	if err := validateProduct(product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// MongoDB collections
	productCollection := db.Collection("products")
	collectionCollection := db.Collection("collections")

	// Check if SKU is unique
	var existingProduct models.Product
	err := productCollection.FindOne(context.TODO(), bson.M{"sku": product.Sku}).Decode(&existingProduct)
	if err == nil {
		http.Error(w, "SKU already exists", http.StatusConflict)
		return
	}

	// Set timestamps
	timestamp := time.Now()
	product.CreatedAt = timestamp
	product.UpdatedAt = timestamp

	// Insert product into database
	result, err := productCollection.InsertOne(context.TODO(), product)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	// Get the inserted product ID
	insertedID := result.InsertedID.(primitive.ObjectID)

	// Process multiple tags and update collections accordingly
	if len(product.Tag) > 0 {
		for _, tag := range product.Tag {
			var collection models.Collection
			err := collectionCollection.FindOne(context.TODO(), bson.M{"collectionName": tag}).Decode(&collection)

			if err != nil { // Collection does not exist, create a new one
				newCollection := models.Collection{
					ID:             primitive.NewObjectID(),
					CollectionName: tag,
					ProductIds:     []primitive.ObjectID{insertedID},
					ShowCollection: true,
				}
				_, err := collectionCollection.InsertOne(context.TODO(), newCollection)
				if err != nil {
					http.Error(w, "Failed to create collection", http.StatusInternalServerError)
					return
				}
			} else { // Collection exists, update it
				_, err := collectionCollection.UpdateOne(
					context.TODO(),
					bson.M{"_id": collection.ID},
					bson.M{"$addToSet": bson.M{"productIds": insertedID}}, // Prevent duplicate IDs
				)
				if err != nil {
					http.Error(w, "Failed to update collection", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	// Send response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func GetProducts(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	isArabic := strings.Contains(r.URL.Path, "/ar")
	isAdmin := r.URL.Query().Get("isAdmin") == "true"

	collection := db.Collection("products")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var products []models.Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		http.Error(w, "Failed to parse products", http.StatusInternalServerError)
		return
	}

	if isAdmin {
		// Return full data for admin
		json.NewEncoder(w).Encode(products)
		return
	}

	// Format for user-facing API
	var response []map[string]interface{}
	for _, product := range products {
		response = append(response, formatProductResponse(product, isArabic, false))
	}
	json.NewEncoder(w).Encode(response)
}

func GetProduct(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	collection := db.Collection("products")
	var product models.Product
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	isArabic := strings.Contains(r.URL.Path, "/ar")
	isAdmin := r.URL.Query().Get("isAdmin") == "true"

	if isAdmin {
		json.NewEncoder(w).Encode(product)
		return
	}

	json.NewEncoder(w).Encode(formatProductResponse(product, isArabic, false))
}

func UpdateProduct(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := db.Collection("products")

	// Ensure SKU remains unique
	var existingProduct models.Product
	err = collection.FindOne(context.TODO(), bson.M{"sku": product.Sku, "_id": bson.M{"$ne": objID}}).Decode(&existingProduct)
	if err == nil {
		http.Error(w, "SKU already exists", http.StatusConflict)
		return
	}

	product.UpdatedAt = time.Now()
	var existing models.Product
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&existing)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Then apply updates manually:
	update := bson.M{
		// Only update fields that are non-zero or intentionally changed
		"sku":           product.Sku,
		"title":         product.Title,
		"titleAr":       product.TitleAr,
		"description":   product.Description,
		"descriptionAr": product.DescriptionAr,
		"price":         product.Price,
		"image":         product.Image,
		"tag":           product.Tag,
		"stock":         product.Stock,
		"updatedAt":     time.Now(),
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	productCollection := db.Collection("products")
	collectionCollection := db.Collection("collections")

	var product models.Product
	err = productCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	res, err := productCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil || res.DeletedCount == 0 {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	for _, tag := range product.Tag {
		filter := bson.M{"collectionName": tag}
		update := bson.M{
			"$pull": bson.M{
				"productIds": objID,
			},
		}
		_, err := collectionCollection.UpdateMany(context.TODO(), filter, update)
		if err != nil {
			log.Printf("Warning: failed to update collection with tag '%s': %v", tag, err)
		}
	}

	w.WriteHeader(http.StatusOK)
}

// Helper function to format product response based on language
func formatProductResponse(product models.Product, isArabic bool, isAdmin bool) map[string]interface{} {
	if isArabic {
		return map[string]interface{}{
			"id":          product.ID,
			"sku":         product.Sku,
			"title":       product.TitleAr,
			"description": product.DescriptionAr,
			"price":       product.Price,
			"image":       product.Image,
			"createdAt":   product.CreatedAt.Format(time.RFC3339),
			"updatedAt":   product.UpdatedAt.Format(time.RFC3339),
			"stock":       product.Stock,
			"tag":         product.Tag,
		}
	}
	if isAdmin {
		return map[string]interface{}{
			"id":            product.ID,
			"sku":           product.Sku,
			"title":         product.Title,
			"titleAr":       product.TitleAr,
			"description":   product.Description,
			"descriptionAr": product.DescriptionAr,
			"price":         product.Price,
			"image":         product.Image,
			"createdAt":     product.CreatedAt.Format(time.RFC3339),
			"updatedAt":     product.UpdatedAt.Format(time.RFC3339),
			"stock":         product.Stock,
			"tag":           product.Tag,
		}
	}
	return map[string]interface{}{
		"id":          product.ID,
		"sku":         product.Sku,
		"title":       product.Title,
		"description": product.Description,
		"price":       product.Price,
		"image":       product.Image,
		"createdAt":   product.CreatedAt.Format(time.RFC3339),
		"updatedAt":   product.UpdatedAt.Format(time.RFC3339),
		"stock":       product.Stock,
		"tag":         product.Tag,
	}
}

func GetProductsByIDs(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	var request struct {
		ProductIds []string `json:"productIds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(request.ProductIds) == 0 {
		http.Error(w, "No product IDs provided", http.StatusBadRequest)
		return
	}

	var objIDs []primitive.ObjectID
	for _, id := range request.ProductIds {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid product ID format", http.StatusBadRequest)
			return
		}
		objIDs = append(objIDs, objID)
	}

	collection := db.Collection("products")
	cursor, err := collection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": objIDs}})
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var products []models.Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		http.Error(w, "Failed to parse products", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}
