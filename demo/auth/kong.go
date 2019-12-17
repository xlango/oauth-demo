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

func stringToBase64(key string) string {
	input := []byte(key)
	// base64编码
	encodeString := base64.StdEncoding.EncodeToString(input)
	return encodeString
}

func getOAuthCode(username string, password string, clientId string) *RespCodeEntity {
	apiUrl := "https://192.168.10.234:8443"
	resource := "/permission/oauth2/authorize"
	data := url.Values{}
	//data.Set("client_id", "BxbjIJwOEHdarSjsfadjlw3whezCPTGn")
	data.Set("client_id", clientId)
	data.Set("response_type", "code")
	data.Set("scope", "email address")
	//data.Set("provision_key", "PzFa0aSZm06KfaMzYlkOuQyWdyeuyV7T")
	//data.Set("provision_key", "hGg6tZ5OwiqftvjZJ09Z1n9LptXJ8aAl")
	data.Set("provision_key", "neREK21ph1FPkFr9YrbKVQN54pTrZ1q4")
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
	str:=string(result)
	fmt.Println(str)
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

func getToken(code string, clientId string, clientSecret string) *TokenEntity {
	apiUrl := "https://192.168.10.234:8443"
	resource := "/permission/oauth2/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	//data.Set("client_id", "BxbjIJwOEHdarSjsfadjlw3whezCPTGn")
	//data.Set("client_secret", "hl9EAX1z620qYKGMd1C149Tng4MIpYvo")
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
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
	apiUrl := "https://192.168.10.234:8443"
	resource := "/permission/oauth2/token"
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

type ConsumerResp struct {
	CustomId  string `json:"custom_id"`
	CreatedAt string `json:"created_at"`
	Id        string `json:"id"`
	Tags      string `json:"tags"`
	Username  string `json:"username"`
}

type ConsumerOAuthResp struct {
	ClientSecret string `json:"client_secret"`
	ClientId     string `json:"client_id"`
	Id           string `json:"id"`
	Tags         string `json:"tags"`
	name         string `json:"name"`
}

//创建consumer
func createConsumer(user *User) *ConsumerResp {
	apiUrl := "http://192.168.10.234:8001"
	resource := "/consumers"
	data := url.Values{}
	data.Set("username", user.Username)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{Timeout: 5 * time.Second}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)

	cr := ConsumerResp{}
	json.Unmarshal(result, &cr)

	return &cr
}

//创建consumer
func createConsumerOAuth(consumer *ConsumerResp) *ConsumerOAuthResp {
	apiUrl := "http://192.168.10.234:8001"
	resource := fmt.Sprintf("/consumers/%v/oauth2", consumer.Id)
	data := url.Values{}
	data.Set("name", fmt.Sprintf("%vOauth", consumer.Username))
	data.Set("redirect_uris", "http://getkong.org/")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{Timeout: 5 * time.Second}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)

	cr := ConsumerOAuthResp{}
	json.Unmarshal(result, &cr)

	return &cr
}
