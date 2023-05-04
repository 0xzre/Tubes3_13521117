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
	"go.mongodb.org/mongo-driver/mongo"
)

// Client Database instance
var validate = validator.New()
var Client *mongo.Client = routes.DBinstance()
var questionCollection *mongo.Collection = routes.OpenCollection(Client, "questions")

// get response by user input using KMP
func GetResponseKMP(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	question := c.Params.ByName("question")
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
	regexAdd := regexp.MustCompile(`^(?:\s+)?tambah(?:kan)?(?:\s+)?pertanyaan(?:(?:\s+)?(.+?)?(?:(?:\s+dengan)?\s+jawaban(?:nya)?))?(?:\s+)?(.+?)?(?:\s*|\b)$`)
	matchAdd := regexAdd.MatchString(question)
	parseAdd := regexAdd.FindStringSubmatch(question)

	if matchAdd {
		questionAdded := ""
		answerAdded := ""

		// If user not adding the answer
		if parseAdd[1] == "" {
			questionAdded = parseAdd[2]
		} else {
			questionAdded = parseAdd[1]
			answerAdded = parseAdd[2]
		}
		AddQuestion(c, questionAdded, answerAdded, questions)
		return
	}

	// Delete Question : "hapus pertanyaan ...." or "hapus ...."
	regexDelete := regexp.MustCompile(`^(?:\s+)?(?:meng)?hapus(?:lah)?(?:kan)?(?:\s+)?(?:(?:pertanyaan(?:\s+)?)?(.+?)(?:\s*|\b)$)`)
	matchDelete := regexDelete.MatchString(question)
	parseDelete := regexDelete.FindStringSubmatch(question)

	if matchDelete {
		// questionDeleted := ""
		fmt.Println(len(parseDelete))
		// DeleteQuestion(c, questionDeleted)
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

// Add a question or update the answer if question already exists
func AddQuestion(c *gin.Context, questionAdded string, answerAdded string, questions []bson.M) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var result []bson.M
	var question models.Question
	question.Question = &questionAdded
	question.Answer = &answerAdded

	validationErr := validate.Struct(question)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}

	for _, elmt := range questions {
		if Algorithm.KmpSearch(questionAdded, string(elmt["question"].(string))) != -1 ||
			Algorithm.LongestCommonSubstring(questionAdded,
				string(elmt["question"].(string))) >= 85.0 {

			if answerAdded == "" {
				if string(elmt["question"].(string)) == "" {
					flag := bson.M{"answer": "Pertanyaan " + questionAdded + " sudah ada namun belum tersimpan jawaban, silakan update jawaban"}
					result = append(result, flag)
				} else {
					flag := bson.M{"answer": "Pertanyaan " + questionAdded + " sudah ada dan telah tersimpan jawaban: " + string(elmt["question"].(string))}
					result = append(result, flag)
					c.JSON(http.StatusOK, result)
					return
				}
			} else {
				flag := bson.M{"answer": "Pertanyaan " + questionAdded + " sudah ada! Jawaban diupdate menjadi: " + answerAdded}
				result = append(result, flag)
			}
			DeleteQuestion(c, string(elmt["question"].(string)))
			temp := string(elmt["question"].(string))
			question.Question = &temp
			_, insertErr := questionCollection.InsertOne(ctx, question)
			if insertErr != nil {
				msg := "question was not added"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				fmt.Println(insertErr)
				return
			}
			c.JSON(http.StatusOK, result)
			return
		}
	}

	_, insertErr := questionCollection.InsertOne(ctx, question)
	if insertErr != nil {
		msg := "question was not added"
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}

	flag := bson.M{"answer": "Berhasil menambahkan pertanyaan " + questionAdded + " dengan jawaban " + answerAdded}
	result = append(result, flag)
	c.JSON(http.StatusOK, result)
}

// update answer for an question
// func UpdateAnswer(c *gin.Context) {

// 	questionID := c.Params.ByName("question")

// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 	defer cancel()

// 	type Answer struct {
// 		Server *string `json:"server"`
// 	}

// 	var answer Answer

// 	if err := c.BindJSON(&answer); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 		return
// 	}

// 	result, err := questionCollection.UpdateOne(ctx, bson.M{"_id": docID},
// 		bson.D{
// 			{Key: "$set", Value: bson.D{{Key: "server", Value: answer.Server}}},
// 		},
// 	)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, result.ModifiedCount)

// }

// delete a question given the user input
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
