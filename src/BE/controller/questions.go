package controller

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
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

// Add a question or update the answer if question already exists
func AddQuestion(c *gin.Context, questionAdded string, answerAdded string) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var question models.Question
	question.ID = primitive.NewObjectID()
	question.Question = &questionAdded
	question.Answer = &answerAdded

	validationErr := validate.Struct(question)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}

	result, insertErr := questionCollection.InsertOne(ctx, question)
	if insertErr != nil {
		msg := fmt.Sprintf("question was not added")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

// get answer by the question
func GetResponseKMP(c *gin.Context) {

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
	//TODO
	//parse input dgn pemisah titik untuk baca lebih dari satu perintah/pertanyaan

	//TODO
	//CALCULATOR
	//REGEX dan cara parsing inputnya, trus bikin method baca input untuk return hasilnya

	//TODO
	//TANGGAL
	//REGEX dan cara parsing inputnya, trus bikin method baca input untuk return harinya

	// Add Question : "tambah pertanyaan .... jawaban ...." or "tambah pertanyaan .... dengan jawaban ...."
	regexAdd := regexp.MustCompile(`^tambah pertanyaan\s+(.+?)(?:\s+(dengan\s+)?jawaban\s+(.+))`)
	matchAdd := regexAdd.MatchString(question)
	parseAdd := regexAdd.FindStringSubmatch(question)

	if matchAdd {
		questionAdded := parseAdd[1]
		answerAdded := ""
		fmt.Println(len(parseAdd))
		if len(parseAdd) == 3 {
			answerAdded = parseAdd[2]
		} else if len(parseAdd) == 4 {
			answerAdded = parseAdd[3]
		}
		AddQuestion(c, questionAdded, answerAdded)
		return
	}

	// Delete Question : "hapus pertanyaan ...." or "hapus ...."
	regexDelete := regexp.MustCompile(`^hapus (?:(pertanyaan )?(.+?)$)`)
	matchDelete := regexDelete.MatchString(question)
	parseDelete := regexDelete.FindStringSubmatch(question)

	if matchDelete {
		questionDeleted := ""
		if len(parseDelete) == 2 {
			questionDeleted = parseDelete[1]
		} else if len(parseDelete) == 3 {
			questionDeleted = parseDelete[2]
		}
		DeleteQuestion(c, questionDeleted)
		return
	}

	//sementara masih manggil make KMP:
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

// update answer for an question
func UpdateAnswer(c *gin.Context) {

	questionID := c.Params.ByName("question")
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

	questionID := c.Params.ByName("answer")
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

// delete an question given the question
func DeleteQuestion(c *gin.Context, question string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := questionCollection.DeleteOne(ctx, bson.M{"question": question})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, result.DeletedCount)
}
