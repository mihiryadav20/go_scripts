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

// displayYouTubeUsage prints usage instructions for YouTube link updater and exits
// displayHindiYouTubeUsage prints usage instructions for the Hindi YouTube updater and exits
func displayHindiYouTubeUsage() {
	fmt.Println(`
Usage: go run update_hindi_youtube.go <slug> <youtube_link>

Example:
go run update_hindi_youtube.go fpktnk.json https://www.youtube.com/watch?v=example

Parameters:
  - slug: Required. The unique identifier for the document (e.g., fpktnk.json)
  - youtube_link: Required. The new YouTube video link
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
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || len(args) < 2 {
		displayHindiYouTubeUsage()
	}

	slug := args[0]
	youtubeLink := args[1]
	if slug == "" || youtubeLink == "" {
		log.Println("Error: Slug and YouTube link are required")
		displayHindiYouTubeUsage()
	}

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

	// Create update object
	update := bson.M{
		"$set": bson.M{
			"data.hi.media.video": youtubeLink,
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
		if hi, ok := data["hi"].(bson.M); ok {
			if media, ok := hi["media"].(bson.M); ok {
				if val, ok := media["video"].(string); ok {
					updatedLink = val
				}
			}
		}
	}
	fmt.Printf("hi.media.video: %s\n", updatedLink)
}
