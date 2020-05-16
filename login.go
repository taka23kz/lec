package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

/*
Login ...
 ログイン画面の入力内容を保持するための構造体
*/
type Login struct {
	// input
	MailAddress string `json:"mailAddress"`
	Passwd      string `json:"passwd"`

	// output
	Message string `json:"message"`
}

/*
Login ...
 ログイン画面でログインボタンが押下された際のAPI
*/
func (controller *Controller) Login(c echo.Context) error {
	var login Login
	var queryParam []interface{}
	var userRows []User
	var chk bool

	// ログイン画面の入力値をloginにバインド
	if err := c.Bind(&login); err != nil {
		c.Logger().Error("Bind login: ", err)
		return c.String(http.StatusBadRequest, "Bind login: "+err.Error())
	}
	// 入力チェック
	chk, login.Message = checkLoginInput(login)

	if !chk {
		c.Logger().Error("checkLoginInput: ", login.Message)
		return c.JSON(http.StatusOK, login)
	}

	// 同一ユーザ存在チェック
	hashedPassword := encode(login.Passwd)
	query := "SELECT * FROM t_user where mail_address = $1 and passwd = $2 and delete_flag = '0'"
	queryParam = append(queryParam, login.MailAddress)
	queryParam = append(queryParam, hashedPassword)

	_, err := controller.dbmap.Select(&userRows, query, queryParam...)
	if err != nil {
		c.Logger().Error("ユーザ検索に失敗しました。: ", err)
		return c.String(http.StatusBadRequest, "ユーザ検索に失敗しました。: "+err.Error())
	}

	if len(userRows) == 0 {
		// ユーザが見つからない→メールアドレスまたはパスワードが誤っている
		login.Message = "メールアドレスまたはパスワードが誤っています。"
		c.Logger().Warn("login check: ", login.Message)
		return c.JSON(http.StatusOK, login)
	}

	return c.JSON(http.StatusOK, login)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Println("bcrypt.GenerateFromPassword :", err)
		return "", err
	}
	return string(hash), nil
}

func checkLoginInput(login Login) (bool, string) {

	if login.MailAddress == "" {
		return false, "メールアドレスは必須です。"
	}
	if login.Passwd == "" {
		return false, "パスワードは必須です。"
	}
	return true, ""
}
