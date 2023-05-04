package controller

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"time"

	Algorithm "BE/String-Matching-Algorithm"
	"BE/server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Get response by user input using KMP
func GetResponseKMP(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Retrieve user input
	question := c.Params.ByName("question")
	var result []bson.M
	var questions []bson.M

	// Revoke all data in database
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

	// Get all listed question in database: "list pertanyaan"
	regexGetAll := regexp.MustCompile(`^(?:\s+)?list(?:\s+)?pertanyaan(?:\s+)?(.+?)?(?:\s*|\b)$`)
	matchGetAll := regexGetAll.MatchString(question)
	if matchGetAll {

		// No questions message
		if len(questions) == 0 {
			flag := bson.M{"answer": "Belum ada pertanyaan yang terdaftar!"}
			result = append(result, flag)

		} else { // Get all questions
			flag := bson.M{"answer": "List pertanyaan yang telah terdaftar:"}
			result = append(result, flag)
			result = append(result, questions...)
		}
		c.JSON(http.StatusOK, result)
		return
	}

	// Add Question: "tambah pertanyaan .... jawaban ...." or "tambah pertanyaan .... dengan jawaban ...."
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
		AddQuestionKMP(c, questionAdded, answerAdded, questions)
		return
	}

	// Delete Question prompt: "hapus pertanyaan ...." or "hapus ...."
	regexDelete := regexp.MustCompile(`^(?:\s+)?(?:meng)?hapus(?:lah)?(?:kan)?(?:\s+)?(?:(?:pertanyaan(?:\s+)?)?(.+?)(?:\s*|\b)$)`)
	matchDelete := regexDelete.MatchString(question)
	parseDelete := regexDelete.FindStringSubmatch(question)

	// Match with regex
	if matchDelete {
		questionDeleted := parseDelete[1]
		DeleteQuestionKMP(c, questionDeleted, questions)
		return
	}

	if questions != nil {
		for _, elmt := range questions {

			// Match using KMP, add into result
			// fmt.Println(Algorithm.KMPSearch(question, string(elmt["question"].(string))))
			if Algorithm.KMPSearch(question, string(elmt["question"].(string))) != -1 {
				result = append(result, elmt)
				break

			} else {
				// Not match KMP but LCS >= 90%, add into result
				// fmt.Println(Algorithm.LongestCommonSubstring(question, string(elmt["question"].(string))))
				if Algorithm.LongestCommonSubstring(question, string(elmt["question"].(string))) >= 90.0 {
					result = append(result, elmt)
					break
				}
			}
		}

		// Question not match KMP and LCS < 90%
		if len(result) != 1 {
			flag := bson.M{"answer": "Pertanyaan tidak ditemukan, mungkin maksudnya:"}
			result = append(result, flag)
			rank := []float64{}

			// Get rank of the similiarity using LCS and sort them descending
			for i := 0; i < len(questions); i++ {
				rank = append(rank, Algorithm.LongestCommonSubstring(question, string(questions[i]["question"].(string))))
			}
			sort.Sort(sort.Reverse(sort.Float64Slice(rank)))

			numOfRecommendation := len(rank)
			if len(rank) >= 3 {
				numOfRecommendation = 3
			}
			for i := 0; i < numOfRecommendation; i++ {
				for j := 0; j < len(questions); j++ {

					// Add 3 biggest LCS to result based on rank
					if Algorithm.LongestCommonSubstring(question, string(questions[j]["question"].(string))) == rank[i] {
						result = append(result, questions[j])
					}
				}
			}
		}
	}
	c.JSON(http.StatusOK, result)
}

// Add a question or update the answer if question already exists using KMP
func AddQuestionKMP(c *gin.Context, questionAdded string, answerAdded string, questions []bson.M) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var result []bson.M
	var question models.Question
	question.Question = &questionAdded
	question.Answer = &answerAdded

	// Validate the question struct
	validationErr := validate.Struct(question)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}

	for _, elmt := range questions {

		// If matching with KMP or LCS >= 90%
		// fmt.Println(Algorithm.LongestCommonSubstring(questionAdded, string(elmt["question"].(string))))
		if Algorithm.KMPSearch(questionAdded, string(elmt["question"].(string))) != -1 ||
			Algorithm.LongestCommonSubstring(questionAdded,
				string(elmt["question"].(string))) >= 90.0 {

			if answerAdded == "" { // Case where user doesn't input the answer yet

				// Question exists but the answer doesn't
				if string(elmt["answer"].(string)) == "" {
					flag := bson.M{"answer": "Pertanyaan \"" + elmt["question"].(string) + "\" sudah ada, namun belum tersimpan jawabannya. Silakan update jawaban"}
					result = append(result, flag)
					c.JSON(http.StatusOK, result)
					return

				} else {
					// Question exists and the answer as well
					flag := bson.M{"answer": "Pertanyaan \"" + elmt["question"].(string) + "\" sudah ada dan telah tersimpan jawaban: \"" + string(elmt["answer"].(string)) + "\""}
					result = append(result, flag)
					c.JSON(http.StatusOK, result)
					return
				}

			} else { // Update question with new answer
				flag := bson.M{"answer": "Pertanyaan \"" + elmt["question"].(string) + "\" sudah ada! Jawaban diupdate menjadi: \"" + answerAdded + "\""}
				result = append(result, flag)
				// fmt.Println(result[0]["answer"])
				c.JSON(http.StatusOK, result)
			}

			// Delete recent question
			DeleteQuestionKMP(c, string(elmt["question"].(string)), questions)

			// Re-Add question with new answer
			temp := string(elmt["question"].(string))
			question.Question = &temp
			_, insertErr := questionCollection.InsertOne(ctx, question)
			if insertErr != nil {
				msg := "question was not added"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				fmt.Println(insertErr)
				return
			}
			return
		}
	}

	// Add question to database
	_, insertErr := questionCollection.InsertOne(ctx, question)
	if insertErr != nil {
		msg := "question was not added"
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}

	// Add question succeed message
	if answerAdded == "" {
		flag := bson.M{"answer": "Berhasil menambahkan pertanyaan \"" + questionAdded + "\" tanpa jawaban. Update untuk menambahkan jawaban"}
		result = append(result, flag)
	} else {
		flag := bson.M{"answer": "Berhasil menambahkan pertanyaan \"" + questionAdded + "\" dengan jawaban \"" + answerAdded + "\""}
		result = append(result, flag)
	}
	c.JSON(http.StatusOK, result)
}

// Delete a question given the user input using KMP
func DeleteQuestionKMP(c *gin.Context, questionDeleted string, questions []bson.M) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var result []bson.M

	for _, elmt := range questions {

		// If matching with KMP algorithm or LCS >= 90.0
		if Algorithm.KMPSearch(questionDeleted, string(elmt["question"].(string))) != -1 ||
			Algorithm.LongestCommonSubstring(questionDeleted, string(elmt["question"].(string))) >= 90.0 {

			// Delete in database
			_, err := questionCollection.DeleteOne(ctx, bson.M{"question": string(elmt["question"].(string))})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				fmt.Println(err)
				return
			}

			// Delete succeed response
			flag := bson.M{"answer": "Pertanyaan \"" + questionDeleted + "\" berhasil dihapus!"}
			result = append(result, flag)
			c.JSON(http.StatusOK, result)
			return
		}
	}

	// Cannot find the question in database
	flag := bson.M{"answer": "Pertanyaan \"" + questionDeleted + "\" tidak ditemukan sehingga tidak bisa dihapus!"}
	result = append(result, flag)
	c.JSON(http.StatusOK, result)
}
