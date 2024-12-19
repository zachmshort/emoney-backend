package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database // Update this to store *mongo.Database

func loadEnvFile() error {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	envPaths := []string{
		filepath.Join(currentDir, "..", ".env"), // From config dir
		".env",                                  // From root dir
	}

	var lastErr error
	for _, path := range envPaths {
		absPath, _ := filepath.Abs(path)
		if _, err := os.Stat(absPath); err == nil {
			log.Printf("Found .env file at: %s", absPath)
			if err := godotenv.Load(absPath); err == nil {
				log.Printf("Successfully loaded environment from: %s", absPath)
				return nil
			} else {
				lastErr = fmt.Errorf("found .env at %s but failed to load: %v", absPath, err)
			}
		} else {
			log.Printf("No .env file at: %s", absPath)
		}
	}

	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("no .env file found in search paths")
}

func ValidateEnv() error {
	err := loadEnvFile()
	if err != nil {
		log.Printf("Warning: %v", err)
		log.Println("Proceeding with system environment variables")
	}

	if wd, err := os.Getwd(); err == nil {
		log.Printf("Working directory: %s", wd)
	}

	required := []string{"DATABASE_URL", "DATABASE_NAME"}
	missing := []string{}

	for _, env := range required {
		if value := os.Getenv(env); value == "" {
			missing = append(missing, env)
			log.Printf("❌ Missing required environment variable: %s", env)
		} else {
			log.Printf("✅ Found environment variable: %s", env)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}

	log.Println("All required environment variables are present")
	return nil
}

func ConnectDB() {
	if err := ValidateEnv(); err != nil {
		log.Fatal(err)
	}

	mongoURI := os.Getenv("DATABASE_URL")
	dbName := os.Getenv("DATABASE_NAME")

	if mongoURI == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Connection error details: %v", err)
		log.Printf("MongoDB URI structure valid: %v", strings.HasPrefix(mongoURI, "mongodb+srv://"))
		log.Fatal("Failed to create MongoDB client: ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}

	log.Printf("Successfully connected to MongoDB!")

	// Store the database instance instead of client
	DB = client.Database(dbName)
}
