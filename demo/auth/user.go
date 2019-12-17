package main

import "oauth-demo/demo/common"

func init() {
	common.CreateTalbe(User{})
}

type User struct {
	Username     string
	Password     string
	Actor        int64
	ClientId     string
	ClientSecret string
}

//将user 用kong client信息存入关系表
func createUser(user User) {
	msdb := common.GetMysqlDb()
	defer msdb.Close()

	msdb.Create(user)
}

func Registe(user *User) {
	//创建kong consumer
	consumer := createConsumer(user)
	//创建kong consumer oauth key
	consumerOAuth := createConsumerOAuth(consumer)

	user.ClientId = consumerOAuth.ClientId
	user.ClientSecret = consumerOAuth.ClientSecret

	//存入mysql
	createUser(*user)
}

func Login(user User) bool {
	msdb := common.GetMysqlDb()
	defer msdb.Close()

	find := msdb.Where("username = ? AND password = ?", user.Username, user.Password).Find(&User{}).Error

	if find != nil {
		return false
	}

	return true
}

func FindByUsername(username string) *User {
	msdb := common.GetMysqlDb()
	defer msdb.Close()

	user := User{}

	find := msdb.Where("username = ?", username).Find(&user).Error

	if find != nil {
		return nil
	}

	return &user
}

func ThirdOauth(thirdc ThirdClient) bool {
	msdb := common.GetMysqlDb()
	defer msdb.Close()

	find := msdb.Where("client_id = ? AND client_secret = ?", thirdc.ClientId, thirdc.ClientSecret).Find(&User{}).Error

	if find != nil {
		return false
	}

	return true
}

func FindByClientId(clientId string) *User {
	msdb := common.GetMysqlDb()
	defer msdb.Close()

	user := User{}

	find := msdb.Where("client_id = ?", clientId).Find(&user).Error

	if find != nil {
		return nil
	}

	return &user
}
