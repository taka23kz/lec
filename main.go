package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/go-playground/locales/ja_JP"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-yaml/yaml"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/go-playground/validator.v9/translations/ja"
)

type db struct {
	DbDriver string `json:"dbDriver" yaml:"dbDriver"`
	Dsn      string `json:"dsn" yaml:"dsn"`
}

type config struct {
	Db db `json:"db" yaml:"db"`
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("loadConfig error:", err)
	}

	dbmap, err := setupDB(cfg.Db.DbDriver, cfg.Db.Dsn)
	if err != nil {
		log.Fatal("setupDB error:", err)
	}
	controller := &Controller{dbmap: dbmap}

	e := setupEcho()

	e.POST("/api/question/regist", controller.InsertQuestion) // 管理画面用:問題登録
	e.POST("/api/question/update", controller.UpdateQuestion) // 管理画面用:問題更新
	e.POST("/api/question/select", controller.SelectQuestion) // 管理画面用:問題検索

	// 検索画面(qlist.html)で使用するapi
	e.POST("/api/qlist/init", controller.InitQList)
	e.POST("/api/qlist/list", controller.SearchQuestion) // 問題の一覧取得

	// 出題＆回答画面(quiz.html)で使用するapi
	e.POST("/api/question/quiz", controller.GetQuestion)   // 出題
	e.POST("/api/question/answer", controller.JudgeAnswer) // 回答判定

	// ユーザ登録画面(signup.html)で使用するapi
	e.POST("/api/signup/signup", controller.Signup) // ユーザ登録

	e.Static("/", "static/")
	e.Logger.Fatal(e.Start(":8080"))
}

// Validate do validation for request value.
func (v *Validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)
	msg := ""
	for _, v := range errs.Translate(v.trans) {
		if msg != "" {
			msg += ", "
		}
		msg += v
	}
	return errors.New(msg)
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Debug = true
	e.Logger.SetOutput(os.Stderr)

	// setup japanese translation
	japanese := ja_JP.New()
	uni := ut.New(japanese, japanese)
	trans, _ := uni.GetTranslator("ja")
	validate := validator.New()
	err := ja.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		log.Fatal(err)
	}

	// register japanese translation for input field
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		switch name {
		case "question":
			return "質問"
		case "choice":
			return "選択肢"
		case "-":
			return ""
		}
		return name
	})

	e.Validator = &Validator{validator: validate, trans: trans}
	return e
}

func setupDB(dbDriver string, dsn string) (*gorp.DbMap, error) {
	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10 * time.Second)

	var diarect gorp.Dialect = gorp.PostgresDialect{}
	dbmap := &gorp.DbMap{Db: db, Dialect: diarect}
	return dbmap, nil
}

func loadConfig() (*config, error) {
	f, err := os.Open("config.yml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg config
	err = yaml.NewDecoder(f).Decode(&cfg)
	return &cfg, err
}
