package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"oauth-demo/demo/common"
)

var RedisC *common.RedisClient

func init() {
	RedisC = new(common.RedisClient)
}

type Person struct {
	Name     string
	Age      int
	Content  string
	Username string
	Password string
}

func main() {

	webServer()

}

func webServer() {
	http.HandleFunc("/login", loginOAuth)
	http.HandleFunc("/test/a", handla)
	http.HandleFunc("/third/oauth", handlthird)
	http.HandleFunc("/check", handlCheck)
	http.HandleFunc("/register", handleRegiste)

	http.ListenAndServe(":10001", nil)
}

func loginOAuth(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	con, _ := ioutil.ReadAll(request.Body) //获取post的数据

	user := User{}
	json.Unmarshal(con, &user)

	//code := getOAuthCode(user.Username, user.Password)
	//token := getToken(code.Code)
	//
	//bytes, _ := json.Marshal(token)
	//writer.Write(bytes)

	if Login(user) {
		user := FindByUsername(user.Username)
		code := getOAuthCode(user.Username, user.Password, user.ClientId)
		token := getToken(code.Code, user.ClientId, user.ClientSecret)

		RedisC.SetExpTime(token.AccessToken, user.Username, 7200)

		bytes, _ := json.Marshal(token)

		writer.Write(bytes)
	} else {
		writer.Write([]byte("Auth failed!"))
	}

}

//注册
func handleRegiste(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	con, _ := ioutil.ReadAll(request.Body) //获取post的数据

	user := User{}
	json.Unmarshal(con, &user)

	Registe(&user)
}

func handla(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(1)
	person := Person{
		Name:    "a",
		Age:     22,
		Content: "aaaaaaaaaa",
	}
	bytes, _ := json.Marshal(person)
	writer.Write(bytes)
}

func handlCheck(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("api gateway check")
	writer.Write([]byte("api gateway check"))
}

type ThirdClient struct {
	ClientSecret string `json:"client_secret"`
	ClientId     string `json:"client_id"`
}

func handlthird(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	con, _ := ioutil.ReadAll(request.Body) //获取post的数据

	third := ThirdClient{}
	json.Unmarshal(con, &third)

	if ThirdOauth(third) {
		user := FindByClientId(third.ClientId)
		code := getOAuthCode(user.Username, user.Password, user.ClientId)
		token := getToken(code.Code, user.ClientId, user.ClientSecret)

		RedisC.SetExpTime(token.AccessToken, user.Username, 7200)

		bytes, _ := json.Marshal(token)

		writer.Write(bytes)
	} else {
		writer.Write([]byte("Auth failed!"))
	}
}
