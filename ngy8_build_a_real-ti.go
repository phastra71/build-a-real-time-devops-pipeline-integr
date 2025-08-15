package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v32/github"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BuildConfig struct {
	GitHubToken string `json:"github_token"`
	GitHubRepo  string `json:"github_repo"`
	MongoURI   string `json:"mongo_uri"`
}

func main() {
	ctx := context.Background()
	buildConfig := &BuildConfig{
		GitHubToken: "your-github-token",
		GitHubRepo:  "your-github-repo",
		MongoURI:   "your-mongo-uri",
	}

	// Create a GitHub client
	tc := &http.Client{}
	gl := &github.BasicAuthTransport{
		Username: "your-github-username",
		Password: buildConfig.GitHubToken,
	}
	gl.Client = tc
	client := github.NewClient(tc)

	// Create a MongoDB client
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(buildConfig.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(ctx)

	// Get the latest commit from GitHub
	commit, _, err := client.Repositories.GetCommit(ctx, "your-github-username", buildConfig.GitHubRepo, "master")
	if err != nil {
		log.Fatal(err)
	}

	// Insert the commit into MongoDB
	collection := mongoClient.Database("your-mongo-db").Collection("commits")
	_, err = collection.InsertOne(ctx, bson.M{"commit": commit})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Commit inserted into MongoDB:")
	fmt.Println(commit)
}