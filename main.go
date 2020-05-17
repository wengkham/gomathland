package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var questionMemTable = map[int64][]question{}
var answerMemTable = map[int64][]question{}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		qno := c.DefaultQuery("q", "10")
		iqno, _ := strconv.Atoi(qno)
		sid := genQuestions(iqno)
		questions := questionMemTable[sid]

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":        "Math Land",
			"payload":      questions,
			"sid":          sid,
			"hasError":     true,
			"just_started": true,
		})
	})

	r.POST("/submit", func(c *gin.Context) {
		var returnList = []question{}
		hasError := false
		answers := c.PostFormArray("answers")
		sid, _ := strconv.Atoi(c.PostForm("sid"))
		id := int64(sid)

		answerList := answerMemTable[id]
		questionList := questionMemTable[id]

		for i, s := range answers {
			ss, _ := strconv.Atoi(s)
			isCorrect := correctAnswer(ss, answerList[i].Numbers, answerList[i].QID)
			if isCorrect {
				returnList = append(returnList, answerList[i])
			} else {
				hasError = true
				returnList = append(returnList, questionList[i])
			}
		}

		rl, _ := json.Marshal(returnList)
		saveResult(string(rl))

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":        "Math Land",
			"payload":      returnList,
			"hasError":     hasError,
			"sid":          sid,
			"just_started": false,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func genQuestions(size int) int64 {
	questionList := []question{}
	answerList := []question{}
	perQ := 6
	for j := 0; j < size; j++ {
		q := make([]int, perQ)
		ans := make([]int, perQ)
		start := genRandomInt(randomIntMinMax{min: 300, max: 1000})
		incredby := genRandomInt(randomIntMinMax{min: 10, max: 50})
		removeby := genRandomInt(randomIntMinMax{min: 0, max: perQ - 1})
		q[0] = start

		for i := 0; i < perQ; i++ {
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

	id := createMemTable(questionList, answerList)

	return id
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

func createMemTable(questions []question, answer []question) int64 {
	now := time.Now()
	sec := now.Unix()

	questionMemTable[sec] = questions
	answerMemTable[sec] = answer

	return sec
}

type randomIntMinMax struct {
	min, max int
}

type question struct {
	ID      int
	Numbers []int
	QID     int
}
