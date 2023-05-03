package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	Algorithm "BE/String-Matching-Algorithm"
	"BE/server/models"
	"BE/server/routes"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

// Client Database instance
var Client *mongo.Client = routes.DBinstance()

var questionCollection *mongo.Collection = routes.OpenCollection(Client, "questions")

// add a question
func AddQuestion(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var question models.Question

	if err := c.BindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	validationErr := validate.Struct(question)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}
	question.ID = primitive.NewObjectID()

	result, insertErr := questionCollection.InsertOne(ctx, question)
	if insertErr != nil {
		msg := fmt.Sprintf("question was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}
	defer cancel()

	c.JSON(http.StatusOK, result)
}

// get all questions
func GetQuestions(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var questions []bson.M

	cursor, err := questionCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	if err = cursor.All(ctx, &questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	fmt.Println(questions)

	c.JSON(http.StatusOK, questions)
}

// get answer by the question
func GetAnswerByQuestion(c *gin.Context) {

	question := c.Params.ByName("question")

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var result []bson.M
	var questions []bson.M

	cursor, err := questionCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	if err = cursor.All(ctx, &questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	if questions != nil {
		for _, elmt := range questions {
			fmt.Println(elmt)
			if Algorithm.KmpSearch(question, string(elmt["question"].(string))) != -1 {
				result = append(result, elmt)
				fmt.Println("kmp exact match")
				break
			} else {
				fmt.Println(Algorithm.LongestCommonSubstring(question, string(elmt["question"].(string))))
				if Algorithm.LongestCommonSubstring(question, string(elmt["question"].(string))) >= 85.0 {
					result = append(result, elmt)
					fmt.Println("lcs match > 85%")
					break
				}
			}
		}

		if len(result) != 1 {
			flag := bson.M{"answer": "Pertanyaan tidak ditemukan, mungkin maksud anda: \n"}
			result = append(result, flag)
			fmt.Println("adding flag")
			max := Algorithm.LongestCommonSubstring(question, string(questions[0]["question"].(string)))

			for i := 1; i < len(questions); i++ {
				if max < Algorithm.LongestCommonSubstring(question, string(questions[i]["question"].(string))) {
					max = Algorithm.LongestCommonSubstring(question, string(questions[i]["question"].(string)))
				}
			}

			for i := 0; i < len(questions); i++ {
				if Algorithm.LongestCommonSubstring(question, string(questions[i]["question"].(string))) == max {
					if len(result) == 4 {
						break
					} else {
						result = append(result, questions[i])
						fmt.Println("adding recommendation")
					}
				}
			}
		}
	}
	c.JSON(http.StatusOK, result)
}

// get an question by its id
func GetQuestionById(c *gin.Context) {

	questionID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(questionID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var question bson.M

	if err := questionCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	fmt.Println(question)

	c.JSON(http.StatusOK, question)
}

// update answer for an question
func UpdateAnswer(c *gin.Context) {

	questionID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(questionID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	type Answer struct {
		Server *string `json:"server"`
	}

	var answer Answer

	if err := c.BindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	result, err := questionCollection.UpdateOne(ctx, bson.M{"_id": docID},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "server", Value: answer.Server}}},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, result.ModifiedCount)

}

// update the question
func UpdateQuestion(c *gin.Context) {

	questionID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(questionID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var question models.Question

	if err := c.BindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	validationErr := validate.Struct(question)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}

	result, err := questionCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docID},
		bson.M{
			"question": question.Question,
			"answer":   question.Answer,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, result.ModifiedCount)
}

// delete an question given the id
func DeleteQuestion(c *gin.Context) {

	orderID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(orderID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := questionCollection.DeleteOne(ctx, bson.M{"_id": docID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, result.DeletedCount)

}
