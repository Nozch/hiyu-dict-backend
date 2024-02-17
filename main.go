package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type AdjectiveMetaphors struct {
	ID		primitive.ObjectID `bson:"_id" json:"id"`
	Adjective string `json:"adjective"`
	Metaphors 	[]string `json:"metaphors"`
}
func main() {
	router := gin.Default()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() 

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/hiyuDictionary"))
	
	if (err != nil) {
		log.Fatalf("failed to connect to mongo: %v", err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("failed to disconnect from mongo: %v", err)
		}
	}()

	collection := client.Database("hiyuDictionary").Collection("metaphors")
	
	router.Use(func(c *gin.Context) {	

		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 既存のadjective - metaphorsを全て取得
	router.GET("/metaphors", func(c *gin.Context) {
		requestCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()


		var allAdjectiveMetaphors []AdjectiveMetaphors
		cur, err := collection.Find(requestCtx, bson.M{})
		
		if err != nil {
			log.Printf("Error finding documents: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching metaphors"})
			return
		}
		defer cur.Close(requestCtx)

		for cur.Next(requestCtx) {
			
			var adjectiveMetaphors AdjectiveMetaphors
			
			err := cur.Decode(&adjectiveMetaphors)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error while decoding metaphor"})
				return
			}
			
			allAdjectiveMetaphors = append(allAdjectiveMetaphors, adjectiveMetaphors)
		}

		if err := cur.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error in cursor"})
			return
		}

		c.JSON(http.StatusOK, allAdjectiveMetaphors)

	})

	router.GET("/adjectives/:adjective", func(c *gin.Context) {
		requestCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		adjective := c.Param("adjective")

		var result AdjectiveMetaphors
		err := collection.FindOne(requestCtx, bson.M{"adjective": adjective}).Decode((&result))

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, gin.H{"adjective": adjective, "metaphors": []string{}})
			return
		} else if err != nil {
			log.Printf("Error finding document: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"adjective": result.Adjective, "metaphors": result.Metaphors})
	})

	
	router.Run(":3000")
}