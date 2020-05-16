package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

/*
Signup ...
 サインアップ画面の入力内容を保持するための構造体
*/
type Signup struct {
	// input
	MailAddress string `json:"mailAddress"`
	UserName    string `json:"userName"`
	Passwd      string `json:"passwd"`
	Passwd2     string `json:"passwd2"`

	// output
	Message string `json:"message"`
}

/*
Signup ...
 サインアップ画面でユーザ登録ボタンが押下された際のAPI
*/
func (controller *Controller) Signup(c echo.Context) error {
	var signup Signup
	var queryParam []interface{}
	var userRows []User
	var chk bool

	// ユーザ登録画面の入力値をsignupにバインド
	if err := c.Bind(&signup); err != nil {
		c.Logger().Error("Bind signup: ", err)
		return c.String(http.StatusBadRequest, "Bind signup: "+err.Error())
	}
	// 入力チェック
	chk, signup.Message = checkInput(signup)

	if !chk {
		c.Logger().Error("checkInput: ", signup.Message)
		return c.JSON(http.StatusOK, signup)
	}

	// 同一ユーザ存在チェック
	query := "SELECT * FROM t_user where mail_address = $1 and delete_flag = '0'"
	queryParam = append(queryParam, signup.MailAddress)

	_, err := controller.dbmap.Select(&userRows, query, queryParam...)
	if err != nil {
		c.Logger().Error("ユーザ検索に失敗しました。: ", err)
		return c.String(http.StatusBadRequest, "ユーザ検索に失敗しました。: "+err.Error())
	}

	if len(userRows) != 0 {
		// 同一メールアドレスのユーザがいたら登録できない。
		signup.Message = "既にユーザが存在するため登録できません。"
		c.Logger().Warn("userName check: ", fmt.Sprintf("同じメールアドレスのユーザが存在するため登録できません。メールアドレス=[%s]", signup.MailAddress))
		return c.JSON(http.StatusOK, signup)
	}

	// 同名のユーザがいないので登録可能
	// 登録情報を作成してユーザ登録
	// サインアップ時はサインアップ画面の入力内容＋
	// 「00:仮登録状態,true:有料機能制限あり」で登録する。
	user, err := createUser(signup, "00", true)
	if err != nil {
		c.Logger().Error("userName check: ", signup.Message, "detail: ", err)
	}

	// ユーザ情報をDB(t_user)に登録する。
	err = controller.dbmap.Insert(user)
	if err != nil {
		signup.Message = "ユーザ登録でエラーが発生しました。"
		c.Logger().Error("userName check: ", signup.Message, "detail: ", err)
		return c.JSON(http.StatusOK, signup)
	}

	return c.JSON(http.StatusOK, signup)
}

/*
 createUser
 サインアップ画面で登録するユーザ情報を作成する。
*/
func createUser(signup Signup, userStatus string, limitFlag bool) (*User, error) {
	var user User
	user.MailAddress = signup.MailAddress
	user.UserName = signup.UserName
	user.UserStatus = userStatus
	user.LimitFlag = limitFlag

	// パスワードをハッシュ化
	user.Passwd = encode(signup.Passwd)
	fmt.Println("パスワード: ", signup.Passwd)
	fmt.Println("コンバート後のパスワード: ", user.Passwd)

	user.Created = time.Now()
	user.Updated = time.Now()

	return &user, nil
}

func checkInput(signup Signup) (bool, string) {

	if signup.MailAddress == "" {
		return false, "メールアドレスは必須です。"
	}
	if signup.UserName == "" {
		return false, "ユーザ名は必須です。"
	}
	if signup.Passwd == "" {
		return false, "パスワードは必須です。"
	}
	if signup.Passwd2 == "" {
		return false, "パスワード(確認用)は必須です。"
	}
	if signup.Passwd != signup.Passwd2 {
		return false, "パスワードとパスワード(確認用)が一致していません。"
	}
	return true, ""
}
