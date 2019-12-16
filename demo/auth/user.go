package main

func init() {
	CreateTalbe(User{})
}

type User struct {
	Username     string
	Password     string
	ClientId     string
	ClientSecret string
}

//将user 用kong client信息存入关系表
func createUser(user User) {
	msdb := GetMysqlDb()
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
	msdb := GetMysqlDb()
	defer msdb.Close()

	find := msdb.Where("username = ? AND password = ?", user.Username, user.Password).Find(&User{}).Error

	if find != nil {
		return false
	}

	return true
}

func FindByUsername(username string) *User {
	msdb := GetMysqlDb()
	defer msdb.Close()

	user := User{}

	find := msdb.Where("username = ?", username).Find(&user).Error

	if find != nil {
		return nil
	}

	return &user
}
