package main

import (
	"github.com/go-gorp/gorp"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

// Controller is a controller for this application.
type Controller struct {
	dbmap *gorp.DbMap
}

// Error indicate response erorr
type Error struct {
	Error string `json:"error"`
}

// Validator is implementation of validation of rquest values.
type Validator struct {
	trans     ut.Translator
	validator *validator.Validate
}

/*
Choice ...
 t_choiceの内容を保持するための構造体
*/
type Choice struct {
	ChoiceID    int64  `json:"choiceId" form:"choiceId" db:"choice_id,primarykey,autoincrement"` // 問題に対する選択肢を一意に特定するためのID
	QuestionID  int64  `json:"questionId" form:"questionId" db:"question_id,notnull"`            // 問題を一意に特定するためのキー
	ChoiceLabel string `json:"choiceLabel" form:"choiceLabel" db:"choice_label,notnull"`         // 選択肢
	Correct     bool   `json:"correct" form:"correct" db:"correct,notnull"`                      // 選択肢が正解かを判別するためのフラグ
}

/*
Question ...
 t_questionの内容を保持するための構造体
*/
type Question struct {
	QuestionID   int64  `json:"questionId" form:"questionId" db:"question_id,primarykey,autoincrement"` // 問題を一意に特定するためのキー
	Question     string `json:"question" form:"question" db:"question,notnull"`                         // 問題文
	AnswerType   string `json:"answerType" form:"answerType" db:"answer_type,notnull,size:2"`           // 回答の形式 ラジオボタン形式、チェックボックス形式、etc..詳細はddl.sql参照
	ChoiceNum    int    `json:"choiceNum" form:"choiceNum" db:"choice_num,notnull"`                     // 回答として表示する選択肢の数
	OwnerGroupID int64  `json:"ownerGroupId" form:"ownerGroupId" db:"owner_group_id,notnull"`           // 問題のオーナーグループのID
	LessonID     int64  `json:"lessonId" form:"lessonId" db:"lesson_id,notnull"`                        // 問題のLESSONのID
}

/*
Lesson ...
 t_lessonテーブルの内容を保持するための構造体
*/
type Lesson struct {
	LessonID       int64  `json:"lessonId" form:"lessonId" db:"lesson_id, primaryKey, autoincrement"`
	LessonName     string `json:"lessonName" form:"lessonName" db:"lesson_name"`
	SubNumber      int64  `json:"subNumber" form:"subNumber" db:"sub_number"`
	ContentType    string `json:"contentType" form:"contentType" db:"content_type"`
	ParentLessonID int64  `json:"parentLessonId" form:"parentLessonId" db:"parentLessonId"`
	OwnerGroupID   int64  `json:"ownerGroupId" form:"ownerGroupId" db:"ownerGroupId"`
}

/*
Quiz ...
 quiz.htmlで出題する問題とその選択肢を保持するための構造体
*/
type Quiz struct {
	OwnerGroupID int64    `json:"ownerGroupId" form:"ownerGroupId"` // 問題のオーナグループのID
	Question     Question // 問題を格納している構造体
	Choices      []Choice // 選択肢構造体のリスト
}

/*
Answer ...
 quiz.htmlで出題された問題に対するユーザの回答
*/
type Answer struct {
	// input field
	QuestionID int64   `json:"questionId"` // quiz.htmlで出題された問題のID
	Choice     int64   `json:"choice"`     // radioボタン形式の場合に選択した回答の選択ID
	ChoiceIDs  []int64 `json:"choiceIds"`  // チェックボックス形式の場合に選択した回答の選択IDのリスト
	AnswerType string  `json:"answerType"` // 回答の形式 ラジオボタン形式、チェックボックス形式、etc..詳細はddl.sql参照

	// response field
	Correct bool `json:"correct"` // ユーザの回答結果が正解か不正解かのフラグ
}

/*
Qlist ...
 qlist.htmlの入力出力を纏めた構造体
*/
type Qlist struct {

	// input
	SearchCondition SearchCondition `json:"searchCondition"` // 検索条件

	// output field
	Lessons   []Lesson   `json:"lessons"`   // 検索条件-Lessonのリスト
	Questions []Question `json:"questions"` // 検索結果
}

/*
SearchCondition ...
 qlist.htmlで設定された検索情報を保持する構造体
*/
type SearchCondition struct {
	OwnerGroupID int64
	LessonID     int64
}
