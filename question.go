package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// Choice is a strucct to hold unit of request and response
type Choice struct {
	ChoiceID    int64  `json:"choiceId" form:"choiceId" db:"choice_id,primarykey,autoincrement"`
	QuestionID  int64  `json:"questionId" form:"questionId" db:"question_id,notnull"`
	ChoiceLabel string `json:"choiceLabel" form:"choiceLabel" db:"choice_label,notnull"`
	Correct     bool   `json:"correct" form:"correct" db:"correct,notnull"`
}

// Question is a struct to hold unit of request and response
type Question struct {
	QuestionID   int64  `json:"questionId" form:"questionId" db:"question_id,primarykey,autoincrement"`
	Question     string `json:"question" form:"question" db:"question,notnull"`
	AnswerType   string `json:"answerType" form:answerType" db:"answer_type,notnull,size:2"`
	ChoiceNum    int    `json:"choiceNum" form:choiceNum" db:"choice_num,notnull"`
	OwnerGroupID int64  `json:"ownerGroupId" form:ownerGroupId" db:"owner_group_id,notnull"`
	LessonID     int64  `json:"lessonId" form:lessonId" db:"lesson_id,notnull"`
}

// Quiz 問題とその選択肢を格納する構造体
type Quiz struct {
	OwnerGroupID int64 `json:"ownerGroupId" form:ownerGroupId"`
	Question     Question
	Choices      []Choice
}

// Answer ユーザーの回答
type Answer struct {
	QuestionID int64 `json:"questionId"`
	Choice     int64 `json:"choice"`
	Correct    bool  `json:"correct"`
}

func printQuestion(q Question) {
	fmt.Println("QuestionID:", q.QuestionID)
	fmt.Println("Question:", q.Question)
	fmt.Println("AnswerType:", q.AnswerType)
	fmt.Println("ChoiceNum:", q.ChoiceNum)
	fmt.Println("OwnerGroupID:", q.OwnerGroupID)
	fmt.Println("LessonID:", q.LessonID)

}

// InsertQuestion is POST handler to insert record.
func (controller *Controller) InsertQuestion(c echo.Context) error {
	var question Question
	var choices []Choice

	attachTable(controller)

	// bind request to question struct
	if err := c.Bind(&question); err != nil {
		c.Logger().Error("Bind question: ", err)
		return c.String(http.StatusBadRequest, "Bind question: "+err.Error())
	}
	printQuestion(question)

	if err := c.Bind(&choices); err != nil {
		c.Logger().Error("Bind Choices: ", err)
		return c.String(http.StatusBadRequest, "Bind Choices: "+err.Error())
	}

	// validate request
	if err := c.Validate(&question); err != nil {
		c.Logger().Error("Validate: ", err)
		return c.JSON(http.StatusBadRequest, &Error{Error: err.Error()})
	}
	trans, err := controller.dbmap.Begin()
	if err != nil {
		c.Logger().Error("Insert: ", err)
		return c.String(http.StatusBadRequest, "Insert: "+err.Error())
	}

	// insert t_question record
	if err := trans.Insert(&question); err != nil {
		c.Logger().Error("Insert t_question: ", err)
		return c.String(http.StatusBadRequest, "Insert t_question: "+err.Error())
	}

	// insert t_choice record
	for choice := range choices {
		if err := trans.Insert((&choice)); err != nil {
			c.Logger().Error("Insert t_choice:", err)
			return c.String(http.StatusBadRequest, "Insert t_choice: "+err.Error())
		}
	}

	trans.Commit()

	c.Logger().Infof("inserted quesion: %v", question.QuestionID)
	return c.NoContent(http.StatusCreated)
}

// UpdateQuestion is POST handler to insert record.
func (controller *Controller) UpdateQuestion(c echo.Context) error {
	var question Question

	if err := c.Bind(&question); err != nil {
		c.Logger().Error("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}
	// validate request
	if err := c.Validate(&question); err != nil {
		c.Logger().Error("Validate: ", err)
		return c.JSON(http.StatusBadRequest, &Error{Error: err.Error()})
	}
	// insert record
	trans, err := controller.dbmap.Begin()
	if err != nil {
		c.Logger().Error("Update: ", err)
		return c.String(http.StatusBadRequest, "Update: "+err.Error())
	}
	if _, err := trans.Update(&question); err != nil {
		c.Logger().Error("Update: ", err)
		return c.String(http.StatusBadRequest, "Update: "+err.Error())
	}
	trans.Commit()

	c.Logger().Infof("updated quesion: %v", question.QuestionID)
	return c.NoContent(http.StatusCreated)
}

// SelectQuestion is GET handler to return records.
func (controller *Controller) SelectQuestion(c echo.Context) error {
	var questions []Question

	attachTable(controller)

	_, err := controller.dbmap.Select(&questions,
		"SELECT * FROM t_question where owner_group_id = $1 ", 1)
	if err != nil {
		c.Logger().Error("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}
	return c.JSON(http.StatusOK, questions)
}

// JudgeAnswer is judge user answer
func (controller *Controller) JudgeAnswer(c echo.Context) error {
	var answer Answer
	var choices []Choice
	if err := c.Bind(&answer); err != nil {
		c.Logger().Error("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}
	fmt.Println(answer.Choice)

	// questionIDでchoiceを検索し、全ての正解を選択できているかを判定する。
	_, err := controller.dbmap.Select(&choices,
		"SELECT * FROM t_choice where question_id = $1", answer.QuestionID)
	if err != nil {
		c.Logger().Error("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}

	// 正解判定
	answer.Correct = true
	for _, c := range choices {
		if c.Correct && answer.Choice != c.ChoiceID {
			// 正解の選択肢を選択していなければ不正解
			answer.Correct = false
			break

		} else if !c.Correct && answer.Choice == c.ChoiceID {
			// 不正解の選択肢を選択していれば不正解
			answer.Correct = false
			break
		}
	}
	return c.JSON(http.StatusOK, answer)
}

// GetQuestion is GET question to return question
func (controller *Controller) GetQuestion(c echo.Context) error {
	var questions []Question
	var choices []Choice
	var quiz Quiz

	attachTable(controller)
	// bind request to question struct
	if err := c.Bind(&quiz); err != nil {
		c.Logger().Error("Bind quiz: ", err)
		return c.String(http.StatusBadRequest, "Bind quiz: "+err.Error())
	}
	_, err := controller.dbmap.Select(&questions, "SELECT * FROM t_question where owner_group_id = $1", quiz.OwnerGroupID)
	if err != nil {
		c.Logger().Error("getQuestion(select t_question): ", err)
		return c.String(http.StatusBadRequest, "getQuestion(select t_question): "+err.Error())
	}
	rand.Seed(time.Now().UnixNano())
	quiz.Question = questions[rand.Intn(len(questions))]

	_, err = controller.dbmap.Select(&choices, "SELECT * FROM t_choice where question_id = $1", quiz.Question.QuestionID)
	if err != nil {
		c.Logger().Error("getQuestion(select t_choice): ", err)
		return c.String(http.StatusBadRequest, "getQuestion(select t_choice): "+err.Error())
	}

	quiz.Choices = choices

	return c.JSON(http.StatusOK, quiz)
}

func attachTable(controller *Controller) {
	question := controller.dbmap.AddTableWithName(Question{}, "t_question")
	question.ColMap("AnswerType").Rename("answer_type")
	question.ColMap("ChoiceNum").Rename("choice_num")
	question.ColMap("OwnerGroupID").Rename("owner_group_id")
	question.ColMap("LessonID").Rename("lesson_id")

	choice := controller.dbmap.AddTableWithName(Choice{}, "t_choice")
	choice.ColMap("ChoiceID").Rename("choice_id")
	choice.ColMap("QuestionID").Rename("question_id")
	choice.ColMap("ChoiceLabel").Rename("choice_label")
}
