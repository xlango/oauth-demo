package main

func init() {
	CreateTalbe(User{})
}

type User struct {
	Username     string
	Password     string
	Actor		 int64
	ClientId     string
	ClientSecret string
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
