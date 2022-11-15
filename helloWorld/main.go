package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func userLogin(writer http.ResponseWriter, request *http.Request) {

	//io.WriteString(writer, "Hello World !")
	//
	//str := `{"code":0, "data":{"id":1, "token":"test"}}`
	//if !loginok {
	//	// 登录失败
	//	str := `{"code":-1, "msg":"密码不正确"`
	//	// 设置 header 为 json 默认为 text/html, 所以特别指出返回的为 application/json
	//	writer.Header().Set("Content-Type", "application/json")
	//	writer.WriteHeader(http.StatusOK) // 设置状态 200
	//	_, _ = writer.Write([]byte(str))
	//} else {
	//	// 登录成功
	//	// 设置 header 为 json 默认为 text/html, 所以特别指出返回的为 application/json
	//	writer.Header().Set("Content-Type", "application/json")
	//	writer.WriteHeader(http.StatusOK) // 设置状态 200
	//	_, _ = writer.Write([]byte(str))
	//}

	// 1. 获取前端传递的参数，mobile, password
	request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	password := request.PostForm.Get("password")
	// 2. 解析前端传来的参数
	loginok := false
	if mobile == "18800138580" && password == "0000" {
		loginok = true
	}
	if loginok {
		data := make(map[string]interface{})
		data["id"] = 1
		data["token"] = "test"
		Resp(writer, 0, data, "")
	} else {
		Resp(writer, -1, nil, "密码不正确")
	}

}

var DbEngin *xorm.Engine

func init() {
	driveName := "mysql"
	DsName := "root:Yizhili80@(localhost:3306)/chat?charset=utf8"
	DbEngine, err := xorm.NewEngine(driveName, DsName)
	if err != nil {
		log.Fatal(err.Error())
	}
	// 是否显示 SQL 语句
	DbEngine.ShowSQL(true)
	// 设置数据库最大打开的连接数
	DbEngine.SetMaxOpenConns(2)

	// 自动创建 User 表结构
	// DbEngine.Sync2(new(User))
	fmt.Println("Done Init DataBase ")
}

// 定义一个结构体
type H struct {
	Code int         `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"` // omitempty 为空的时候不显示
}

func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 设置状态 200
	//_, _ = writer.Write([]byte(str))

	h := H{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	// 将结构体转换成 json 字符串
	res, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, _ = w.Write(res)
}

func RegisterView() {
	// 模版解析
	tpl, err := template.ParseGlob("view/**/*")
	if err != nil {
		log.Fatalln(err.Error()) // 打印并直接退出
	}
	for _, v := range tpl.Templates() {
		tplName := v.Name()
		http.HandleFunc(tplName, func(w http.ResponseWriter, r *http.Request) {
			err = tpl.ExecuteTemplate(w, tplName, nil)
			if err != nil {
				log.Fatalln(err.Error()) // 打印并直接退出
			}
		})
	}

}

func main() {
	// 绑定请求和处理函数
	http.HandleFunc("/user/login", userLogin)
	// 1. 提供静态资源目录支持
	//http.Handle("/", http.FileServer(http.Dir("."))) // 代码暴露了，可访问

	// 2. 提供指定目录的静态文件支持
	http.Handle("/asset/", http.FileServer(http.Dir(".")))
	//➜  ~ curl http://localhost:8080/helloWorld/asset/js/util.js
	//404 page not found
	//➜  ~ curl http://localhost:8080/helloWorld/main.go
	//404 page not found

	//// 登录
	//http.HandleFunc("/user/login.shtml", func(w http.ResponseWriter, r *http.Request) {
	//	// 模版解析
	//	tpl, err := template.ParseFiles("view/user/login.html")
	//	if err != nil {
	//		// 打印并直接退出
	//		log.Fatalln(err.Error())
	//	}
	//	err = tpl.ExecuteTemplate(w, "/user/login.shtml", nil)
	//	if err != nil {
	//		// 打印并直接退出
	//		log.Fatalln(err.Error())
	//	}
	//})
	//
	//// 注册
	//http.HandleFunc("/user/register.shtml", func(w http.ResponseWriter, r *http.Request) {
	//	// 模版解析
	//	tpl, err := template.ParseFiles("view/user/register.html")
	//	if err != nil {
	//		// 打印并直接退出
	//		log.Fatalln(err.Error())
	//	}
	//	err = tpl.ExecuteTemplate(w, "/user/register.shtml", nil)
	//	if err != nil {
	//		// 打印并直接退出
	//		log.Fatalln(err.Error())
	//	}
	//})

	RegisterView()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
