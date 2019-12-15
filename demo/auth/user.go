package main

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

//创建kong consumer
func createConsumer() {

}

func Registe(user User) {

}
