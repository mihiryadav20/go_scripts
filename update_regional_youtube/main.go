package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// displayRegionalYouTubeUsage prints usage instructions for the Regional YouTube updater and exits
func displayRegionalYouTubeUsage() {
	fmt.Println(`
Usage: go run main.go <language> <slug> <youtube_link>

Example:
cd update_regional_youtube && go run main.go ta fpktnk.json https://www.youtube.com/watch?v=example

Parameters:
  - language: Required. The language code (e.g., te, as, kok, gu, ml, mr, mni, lus, or, pa, ta, bn, ks, kn)
  - slug: Required. The unique identifier for the document (e.g., fpktnk.json)
  - youtube_link: Required. The new YouTube video link

Supported language codes:
  - te: Telugu
  - as: Assamese
  - kok: Konkani
  - gu: Gujarati
  - ml: Malayalam
  - mr: Marathi
  - mni: Manipuri
  - lus: Mizo
  - or: Odia
  - pa: Punjabi
  - ta: Tamil
  - bn: Bengali
  - ks: Kashmiri
  - kn: Kannada
`)
	os.Exit(1)
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// MongoDB configuration
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set in .env file")
	}
	dbName := "Hull_Schemes"
	collectionName := "All_agri"

	// Parse command-line arguments
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || len(args) < 3 {
		displayRegionalYouTubeUsage()
	}

	language := args[0]
	slug := args[1]
	youtubeLink := args[2]

	// Validate language code
	validLanguages := map[string]string{
		"te":  "Telugu",
		"as":  "Assamese",
		"kok": "Konkani",
		"gu":  "Gujarati",
		"ml":  "Malayalam",
		"mr":  "Marathi",
		"mni": "Manipuri",
		"lus": "Mizo",
		"or":  "Odia",
		"pa":  "Punjabi",
		"ta":  "Tamil",
		"bn":  "Bengali",
		"ks":  "Kashmiri",
		"kn":  "Kannada",
	}

	langName, valid := validLanguages[language]
	if !valid {
		log.Printf("Error: Invalid language code '%s'", language)
		displayRegionalYouTubeUsage()
	}

	if slug == "" || youtubeLink == "" {
		log.Println("Error: Slug and YouTube link are required")
		displayRegionalYouTubeUsage()
	}

	fmt.Printf("Processing language: %s (%s)\n", language, langName)
	fmt.Printf("Processing slug: %s\n", slug)
	fmt.Printf("New YouTube link: %s\n", youtubeLink)

	// Connect to MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
		fmt.Println("\nMongoDB connection closed")
	}()

	fmt.Println("Connected to MongoDB")

	collection := client.Database(dbName).Collection(collectionName)

	// Find documents with this slug
	cursor, err := collection.Find(ctx, bson.M{
		"$or": []bson.M{
			{"filter.slug": slug},
			{"slug": slug},
		},
	})
	if err != nil {
		log.Fatalf("Error querying documents: %v", err)
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	if err := cursor.All(ctx, &docs); err != nil {
		log.Fatalf("Error decoding documents: %v", err)
	}

	if len(docs) == 0 {
		fmt.Printf("No documents found with slug: %s\n", slug)
		return
	}

	fmt.Printf("Found %d document(s) with slug: %s\n", len(docs), slug)

	// Prioritize document with specific ID
	var targetDoc bson.M
	for _, doc := range docs {
		if doc["_id"] == "775a846c8c5442458ea4860111b28c57" {
			targetDoc = doc
			break
		}
	}
	if targetDoc == nil && len(docs) > 0 {
		targetDoc = docs[0]
	}

	if targetDoc == nil {
		fmt.Println("Could not find a valid document to update")
		return
	}

	targetID := targetDoc["_id"].(string)
	fmt.Printf("Selected document with ID: %s\n", targetID)

	// Create update object with dynamic language field
	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("data.%s.media.video", language): youtubeLink,
		},
	}

	fmt.Println("Applying update:", update)

	// Update the document
	result, err := collection.UpdateOne(ctx, bson.M{"_id": targetID}, update)
	if err != nil {
		log.Fatalf("Error updating document: %v", err)
	}

	if result.ModifiedCount > 0 {
		fmt.Printf("Successfully updated document with ID: %s\n", targetID)
	} else {
		fmt.Println("Document found but no changes were made")
	}

	// Verify the update
	var updatedDoc bson.M
	if err := collection.FindOne(ctx, bson.M{"_id": targetID}).Decode(&updatedDoc); err != nil {
		log.Fatalf("Error verifying update: %v", err)
	}

	fmt.Println("\nVerifying update - YouTube link after update:")
	updatedLink := "Not found"
	if data, ok := updatedDoc["data"].(bson.M); ok {
		if langData, ok := data[language].(bson.M); ok {
			if media, ok := langData["media"].(bson.M); ok {
				if val, ok := media["video"].(string); ok {
					updatedLink = val
				}
			}
		}
	}
	fmt.Printf("%s.media.video: %s\n", language, updatedLink)
}