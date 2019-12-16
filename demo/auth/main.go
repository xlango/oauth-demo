package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type RespCode struct {
	RedirectUri string `json:"redirect_uri"`
}

type RespCodeEntity struct {
	RedirectUrl string
	Code        string
	State       string
}

type TokenEntity struct {
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	State        string `json:"state"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
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

	http.ListenAndServe(":10001", nil)
}

func stringToBase64(key string) string {
	input := []byte(key)
	// base64编码
	encodeString := base64.StdEncoding.EncodeToString(input)
	return encodeString
}

func getOAuthCode(username string, password string) *RespCodeEntity {
	apiUrl := "https://192.168.10.33:8443"
	resource := "/apigwtest/oauth2/authorize"
	data := url.Values{}
	data.Set("client_id", "BxbjIJwOEHdarSjsfadjlw3whezCPTGn")
	data.Set("response_type", "code")
	data.Set("scope", "email address")
	//data.Set("provision_key", "PzFa0aSZm06KfaMzYlkOuQyWdyeuyV7T")
	data.Set("provision_key", "hGg6tZ5OwiqftvjZJ09Z1n9LptXJ8aAl")
	data.Set("authenticated_userid", username)
	data.Set("state", "1")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "http://127.0.0.1/tpost"

	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 5 * time.Second, Transport: tr}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	fmt.Println(stringToBase64(fmt.Sprintf("%v:%v", username, password)))
	r.Header.Add("Authorization", stringToBase64(fmt.Sprintf("%v:%v", username, password)))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)

	rc := RespCode{}
	json.Unmarshal(result, &rc)

	parse, _ := url.Parse(rc.RedirectUri)
	values, _ := url.ParseQuery(parse.RawQuery)

	entity := &RespCodeEntity{}
	entity.RedirectUrl = parse.Host
	if len(values["code"]) > 0 {
		entity.Code = values["code"][0]
	}
	if len(values["state"]) > 0 {
		entity.State = values["state"][0]
	}
	return entity
}

func getToken(code string) *TokenEntity {
	apiUrl := "https://192.168.10.33:8443"
	resource := "/apigwtest/oauth2/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", "BxbjIJwOEHdarSjsfadjlw3whezCPTGn")
	data.Set("client_secret", "hl9EAX1z620qYKGMd1C149Tng4MIpYvo")
	data.Set("code", code)
	data.Set("state", "1")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "http://127.0.0.1/tpost"

	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 5 * time.Second, Transport: tr}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)

	rc := TokenEntity{}
	json.Unmarshal(result, &rc)

	return &rc
}

//通过refresh_token重新获取token
func refreshToken(refreshToken string) *TokenEntity {
	apiUrl := "https://192.168.10.33:8443"
	resource := "/apigwtest/oauth2/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", "BxbjIJwOEHdarSjsfadjlw3whezCPTGn")
	data.Set("client_secret", "hl9EAX1z620qYKGMd1C149Tng4MIpYvo")
	data.Set("refresh_token", refreshToken)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "http://127.0.0.1/tpost"

	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 5 * time.Second, Transport: tr}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)

	rc := TokenEntity{}
	json.Unmarshal(result, &rc)

	return &rc
}

func loginOAuth(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	con, _ := ioutil.ReadAll(request.Body) //获取post的数据

	person := Person{}
	json.Unmarshal(con, &person)

	code := getOAuthCode(person.Username, person.Password)
	token := getToken(code.Code)

	bytes, _ := json.Marshal(token)
	writer.Write(bytes)

	//if person.Username == "xhyl" && person.Password == "123456" {
	//	code := getOAuthCode("xhyl", "123456")
	//	token := getToken(code.Code)
	//
	//	bytes, _ := json.Marshal(token)
	//	writer.Write(bytes)
	//} else {
	//	writer.Write([]byte("Auth failed!"))
	//}

}

//注册
func handleRegist(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	con, _ := ioutil.ReadAll(request.Body) //获取post的数据

	user := User{}
	json.Unmarshal(con, &user)

	Registe(user)
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

func handlthird(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(1)
	person := Person{
		Name:    "a",
		Age:     22,
		Content: "aaaaaaaaaa",
	}
	bytes, _ := json.Marshal(person)
	writer.Write(bytes)
}
