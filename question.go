package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

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

	// questionIDでchoiceを検索し、全ての正解を選択できているかを判定する。
	_, err := controller.dbmap.Select(&choices,
		"SELECT * FROM t_choice where question_id = $1", answer.QuestionID)
	if err != nil {
		c.Logger().Error("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}
	c.Logger().Info("answer:", answer, "choices:", choices)
	// 正解判定
	answer.Correct = true

	// 回答の方法によって正解判定を切り分け
	if answer.AnswerType == "01" {
		// 単一選択形式(ラジオボタン)
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
	} else if answer.AnswerType == "02" {
		// 複数選択形式(チェックボックス)
		// 選択肢分ループ
		for _, c := range choices {
			// 選択肢を選ぶのが正解の場合
			if c.Correct {
				// ユーザの回答分ループ
				var choiceFlg bool = false

				for _, choiceID := range answer.ChoiceIDs {
					if c.ChoiceID == choiceID {
						choiceFlg = true // 正解の選択肢を選択されていた。
					}
				}
				if !choiceFlg {
					answer.Correct = false
				}
			} else {
				// 選択肢を選ばないのが正解の場合
				for _, choiceID := range answer.ChoiceIDs {
					if c.ChoiceID == choiceID {
						// 選択されていたので不正解
						answer.Correct = false
					}
				}
			}
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

	// オーナグループIDに紐づく問題を全て取得する。
	_, err := controller.dbmap.Select(&questions, "SELECT * FROM t_question where owner_group_id = $1", quiz.OwnerGroupID)
	if err != nil {
		c.Logger().Error("getQuestion(select t_question): ", err)
		return c.String(http.StatusBadRequest, "getQuestion(select t_question): "+err.Error())
	}

	// 取得した問題リストの中からランダムに決定した1問を出題用のquiz構造体に格納する。
	rand.Seed(time.Now().UnixNano())
	quiz.Question = questions[rand.Intn(len(questions))]

	// quiestion_idに紐づく選択肢を取得する。
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
