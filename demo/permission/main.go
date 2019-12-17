package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"oauth-demo/demo/common"
)

var RedisC *common.RedisClient

func init() {
	RedisC = new(common.RedisClient)
}

func main() {
	webServer()
}
func webServer() {
	http.HandleFunc("/competence", competence)

	http.ListenAndServe(":10002", nil)
}

//权限校验
func competence(writer http.ResponseWriter, request *http.Request) {
	msdb := common.GetMysqlDb()
	defer msdb.Close()

	token := request.URL.Query()["token"][0]
	interfaceName := request.URL.Query()["interface"][0] //获取Filter传来的接口名
	username, ok := RedisC.Get(token, 0)
	if ok == false {
		//登录失效
		writer.WriteHeader(451)
		return
	}
	iface := FindInterfaceByName(interfaceName)
	if iface == nil {
		//无效请求
		writer.WriteHeader(400)
		return
	}
	////查找接口信息
	//iface := Interface{}
	//find := msdb.Where(" name = ? ", interfaceName).Find(&iface).Error
	//if find != nil {
	//	//无效请求
	//	writer.WriteHeader(400);
	//	return;
	//}

	//查找用户信息
	user := FindByUsername(username)
	if user == nil {
		//没有该用户信息，无效登录
		writer.WriteHeader(451)
		return
	}

	//查找用户角色（用户类型）和权限是否符合权限
	actor_interface := InterfaceActor{}
	find := msdb.Where("interface_id = ? and actor_id = ?", iface.Id, user.Actor).Find(&actor_interface).Error
	if find != nil {
		//没有查到用户角色（用户类型）和接口的权限关系，无权访问
		writer.WriteHeader(401)
		return
	}

	//权限验证通过，返回200和用户信息
	data, _ := json.Marshal(&user)
	writer.Write(data)
}
