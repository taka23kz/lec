package main

import (
	"net/http"

	"github.com/labstack/echo"
)

/*
Signup ...
 サインアップ画面の入力内容を保持するための構造体
*/
type Signup struct {
	UserName  string `json:"userName"`
	Password  string `json:"password"`
	Password2 string `json:"password2"`
	Message   string `json:"message"`
}

/*
Signup ...
 サインアップ画面でユーザ登録ボタンが押下された際のAPI
*/
func (controller *Controller) Signup(c echo.Context) error {
	var signup Signup

	// ユーザ登録画面の入力値をsignupにバインド
	if err := c.Bind(&signup); err != nil {
		c.Logger().Error("Bind signup: ", err)
		return c.String(http.StatusBadRequest, "Bind signup: "+err.Error())
	}
	// 入力チェック
	if signup.UserName == "" {
		signup.Message = "ユーザ名は必須です。"
		c.Logger().Error("userName check: ", signup.Message)
		return c.JSON(http.StatusOK, signup)
	}
	if signup.Password == "" {
		signup.Message = "パスワードは必須です。"
		c.Logger().Error("password check: ", signup.Message)
		return c.JSON(http.StatusOK, signup)
	}
	if signup.Password2 == "" {
		signup.Message = "パスワード(確認用)は必須です。"
		c.Logger().Error("password2 check: ", signup.Message)
		return c.JSON(http.StatusOK, signup)
	}
	if signup.Password != signup.Password2 {
		signup.Message = "パスワードとパスワード(確認用)が一致していません。"
		c.Logger().Error("password check: ", signup.Message)
		return c.JSON(http.StatusOK, signup)
	}

	// 同一ユーザ存在チェック

	return c.JSON(http.StatusOK, signup)
}
