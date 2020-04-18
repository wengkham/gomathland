package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var questionList = []question{}
var answerList = []question{}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		questionList = []question{}
		answerList = []question{}

		qno := c.DefaultQuery("q", "10")
		iqno, _ := strconv.Atoi(qno)
		questions := genQuestions(iqno)

		fmt.Println(answerList)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":        "Math Land",
			"payload":      questions,
			"has_error":    true,
			"just_started": true,
		})
	})

	r.POST("/submit", func(c *gin.Context) {
		var returnList = []question{}
		has_error := false
		answers := c.PostFormArray("answers")

		for i, s := range answers {
			ss, _ := strconv.Atoi(s)
			is_correct := correctAnswer(ss, answerList[i].Numbers, answerList[i].QID)
			if is_correct {
				returnList = append(returnList, answerList[i])
			} else {
				has_error = true
				returnList = append(returnList, questionList[i])
			}
		}

		rl, _ := json.Marshal(returnList)
		saveResult(string(rl))

		fmt.Println(returnList)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":        "Math Land",
			"payload":      returnList,
			"has_error":    has_error,
			"just_started": false,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func genQuestions(size int) []question {
	per_q := 6
	for j := 0; j < size; j++ {
		q := make([]int, per_q)
		ans := make([]int, per_q)
		start := genRandomInt(randomIntMinMax{min: 300, max: 1000})
		incredby := genRandomInt(randomIntMinMax{min: 10, max: 50})
		removeby := genRandomInt(randomIntMinMax{min: 0, max: per_q - 1})
		q[0] = start

		fmt.Println(start, incredby, start%2)
		for i := 0; i < per_q; i++ {
			if start%2 == 0 {
				q[i] = start + incredby*i
			} else {
				q[i] = start - incredby*i
			}
			ans[i] = q[i]
			if i == removeby {
				q[i] = 0
				continue
			}
		}

		qq := question{ID: j + 1, Numbers: q, QID: removeby}
		anss := question{ID: j + 1, Numbers: ans, QID: removeby}
		questionList = append(questionList, qq)
		answerList = append(answerList, anss)
	}

	return questionList
}

func genRandomInt(p randomIntMinMax) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(p.max-p.min+1) + p.min
}

//todo return index
func correctAnswer(a int, list []int, qid int) bool {
	for i, b := range list {
		if b == a && qid == i {
			return true
		}
	}
	return false
}

func saveResult(content string) {
	file, err := os.OpenFile("results.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		log.Fatal(err)
	}
}

type randomIntMinMax struct {
	min, max int
}

type question struct {
	ID      int
	Numbers []int
	QID     int
}
