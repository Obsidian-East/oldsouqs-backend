package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"oldsouqs-backend/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateCollection - Add a new collection
func CreateCollection(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	var collection models.Collection

	if err := json.NewDecoder(r.Body).Decode(&collection); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure unique collection name
	collectionCollection := db.Collection("collections")
	var existingCollection models.Collection
	err := collectionCollection.FindOne(context.TODO(), bson.M{"collectionName": collection.CollectionName}).Decode(&existingCollection)
	if err == nil {
		http.Error(w, "Collection name already exists", http.StatusConflict)
		return
	}

	// Assign a new ObjectID
	collection.ID = primitive.NewObjectID()

	_, err = collectionCollection.InsertOne(context.TODO(), collection)
	if err != nil {
		http.Error(w, "Failed to create collection", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(collection)
}

// GetCollections - Retrieve all collections where ShowCollection = true
func GetCollections(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	collectionCollection := db.Collection("collections")

	cursor, err := collectionCollection.Find(context.TODO(), bson.M{"showCollection": true})
	if err != nil {
		http.Error(w, "Failed to retrieve collections", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var collections []models.Collection
	if err := cursor.All(context.TODO(), &collections); err != nil {
		http.Error(w, "Failed to parse collections", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(collections)
}

// GetCollectionByID - Retrieve a collection by ID (only if ShowCollection = true)
func GetCollectionByID(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	vars := mux.Vars(r)
	id := vars["id"]

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	collectionCollection := db.Collection("collections")
	var collection models.Collection
	err = collectionCollection.FindOne(context.TODO(), bson.M{"_id": objID, "showCollection": true}).Decode(&collection)
	if err != nil {
		http.Error(w, "Collection not found or hidden", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(collection)
}

// UpdateCollection - Modify an existing collection
func UpdateCollection(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	vars := mux.Vars(r)
	id := vars["id"]

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedCollection models.Collection
	if err := json.NewDecoder(r.Body).Decode(&updatedCollection); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collectionCollection := db.Collection("collections")
	_, err = collectionCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": updatedCollection},
	)
	if err != nil {
		http.Error(w, "Failed to update collection", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedCollection)
}

// DeleteCollection - Remove a collection
func DeleteCollection(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	vars := mux.Vars(r)
	id := vars["id"]

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	collectionCollection := db.Collection("collections")
	_, err = collectionCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "Failed to delete collection", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetProductsByCollection(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	vars := mux.Vars(r)
	collectionID := vars["id"]

	objID, err := primitive.ObjectIDFromHex(collectionID)
	if err != nil {
		http.Error(w, "Invalid collection ID format", http.StatusBadRequest)
		return
	}

	collectionCollection := db.Collection("collections")
	var collection models.Collection
	err = collectionCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&collection)
	if err != nil {
		http.Error(w, "Collection not found", http.StatusNotFound)
		return
	}

	if len(collection.ProductIds) == 0 {
		json.NewEncoder(w).Encode([]models.Product{}) // Return empty array if no products
		return
	}

	productCollection := db.Collection("products")
	cursor, err := productCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": collection.ProductIds}})
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
