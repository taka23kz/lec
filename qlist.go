package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/x-motemen/ghq/logger"
)

/*
QListInit ...
 問題検索画面の初期表示用API
*/
func (controller *Controller) InitQList(c echo.Context) error {
	var searchCondition SearchCondition
	var qlist Qlist

	// 問題検索画面の入力値をsearchConditionにバインド
	if err := c.Bind(&searchCondition); err != nil {
		c.Logger().Error("Bind searchCondition: ", err)
		return c.String(http.StatusBadRequest, "Bind searchCondition: "+err.Error())
	}

	qlist.Lessons = controller.selectLesson(searchCondition.OwnerGroupID)
	qlist.SearchCondition = searchCondition

	return c.JSON(http.StatusOK, qlist)
}

func (controller *Controller) selectLesson(ownerGroupID int64) []Lesson {
	var lessons []Lesson

	_, err := controller.dbmap.Select(&lessons, "SELECT lesson_id, lesson_name FROM t_lesson where owner_group_id = $1", ownerGroupID)
	if err != nil {
		logger.Logf("selectLesson(select t_lesson): ", err.Error())
		return nil
	}
	return lessons
}

func (controller *Controller) searchQuestion(c echo.Context) error {
	var questions []Question
	var searchCondition SearchCondition
	var query string
	var queryParam []interface{}

	// 問題検索画面の入力値をsearchConditionにバインド
	if err := c.Bind(&searchCondition); err != nil {
		c.Logger().Error("Bind searchCondition: ", err)
		return c.String(http.StatusBadRequest, "Bind searchCondition: "+err.Error())
	}

	// 検索条件に合致する問題を全て取得する。
	query = "SELECT * FROM t_question where owner_group_id = $1"
	queryParam = append(queryParam, searchCondition.OwnerGroupID)
	if searchCondition.LessonID != -1 {
		query += " and lesson_id = $2"
		queryParam = append(queryParam, searchCondition.LessonID)
	}

	_, err := controller.dbmap.Select(&questions, query, queryParam)
	if err != nil {
		c.Logger().Error("GetQuestionList(select t_question): ", err)
		return c.String(http.StatusBadRequest, "GetQuestionList(select t_question): "+err.Error())
	}
	return c.JSON(http.StatusOK, questions)
}
