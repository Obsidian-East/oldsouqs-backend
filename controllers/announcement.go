package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"oldsouqs-backend/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AnnouncementController struct {
	Announcements *mongo.Collection
}

func NewAnnouncementController(db *mongo.Database) *AnnouncementController {
	return &AnnouncementController{
		Announcements: db.Collection("announcements"),
	}
}

// CreateAnnouncement handles POST /announcements
func (ac *AnnouncementController) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	var announcement models.Announcement
	if err := json.NewDecoder(r.Body).Decode(&announcement); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	announcement.ID = primitive.NewObjectID()
	announcement.CreatedAt = time.Now()
	announcement.UpdatedAt = time.Now()

	_, err := ac.Announcements.InsertOne(context.Background(), announcement)
	if err != nil {
		http.Error(w, "Failed to create announcement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Announcement created")
}

// GetAnnouncements handles GET /announcements
func (ac *AnnouncementController) GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	cursor, err := ac.Announcements.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch announcements", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var announcements []models.Announcement
	if err := cursor.All(context.Background(), &announcements); err != nil {
		http.Error(w, "Error reading announcements", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(announcements)
}

// UpdateAnnouncement handles PUT /announcements/{id}
func (ac *AnnouncementController) UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updateData models.Announcement
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"message":   updateData.Message,
			"updatedAt": time.Now(),
		},
	}

	_, err = ac.Announcements.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		http.Error(w, "Failed to update announcement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Announcement updated")
}

// DeleteAnnouncement handles DELETE /announcements/{id}
func (ac *AnnouncementController) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	res, err := ac.Announcements.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "Failed to delete announcement", http.StatusInternalServerError)
		return
	}
	if res.DeletedCount == 0 {
		http.Error(w, "Announcement not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Announcement deleted")
}
