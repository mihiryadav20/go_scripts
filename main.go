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

// displayUsage prints usage instructions and exits
func displayUsage() {
	fmt.Println(`
Usage: go run updateApplicationLink.go <slug> <application_link>

Example:
go run updateApplicationLink.go fpktnk.json https://new-link.com

Parameters:
  - slug: Required. The unique identifier for the document (e.g., fpktnk.json)
  - application_link: Required. The new URL for the application link
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
		displayUsage()
	}

	slug := args[0]
	applicationLink := args[1]
	if slug == "" || applicationLink == "" {
		log.Println("Error: Slug and application link are required")
		displayUsage()
	}

	fmt.Printf("Processing slug: %s\n", slug)
	fmt.Printf("New application link: %s\n", applicationLink)

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
			"data.en.application_link.value": applicationLink,
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

	fmt.Println("\nVerifying update - Application Link after update:")
	updatedLink := "Not found"
	if data, ok := updatedDoc["data"].(bson.M); ok {
		if en, ok := data["en"].(bson.M); ok {
			if appLink, ok := en["application_link"].(bson.M); ok {
				if val, ok := appLink["value"].(string); ok {
					updatedLink = val
				}
			}
		}
	}
	fmt.Printf("en.application_link.value: %s\n", updatedLink)
}
