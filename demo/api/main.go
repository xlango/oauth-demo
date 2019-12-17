package main

import (
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/test", test)
	http.HandleFunc("/check", handlCheck)

	http.ListenAndServe(":10003", nil)
}

func test(writer http.ResponseWriter, request *http.Request) {
	//模拟远程权限校验

	token := request.Header.Get("Authorization")
	token= token[7:]

	ifaceName := "Add"
	client := &http.Client{}
	url := "http://127.0.0.1:10002/competence?token=" + token + "&interface=" + ifaceName
	reqest, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}

	response, _ := client.Do(reqest)
	status := response.StatusCode

	//权限未通过
	if status != 200 {
		if status == 400 {
			writer.Write([]byte("无效请求"))
		}
		if status == 451 {
			writer.Write([]byte("无效登录"))
		}
		if status == 401 {
			writer.Write([]byte("无访问权限"))
		}
		return
	}

	//var user User
	body, _ := ioutil.ReadAll(response.Body)
	//err = json.Unmarshal(body, &user)
	writer.Write(body)

	//response.Body.Read(data)
	//权限已通过，输出用户信息
	//writer.Write((byte[])user)
}

func handlCheck(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("api gateway check"))
}
