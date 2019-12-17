package main

import "oauth_demo/demo/common"

func init() {
	common.CreateTalbe(Interface{})
	common.CreateTalbe(InterfaceActor{})
}

type Interface struct {
	Id int64 `gorm:"AUTO_INCREMENT"`
	//Id		int64
	Name string
}

type InterfaceActor struct {
	InterfaceId int64
	ActorId     int64
}

func FindInterfaceByName(name string) *Interface {
	msdb := common.GetMysqlDb()
	defer msdb.Close()
	//查找接口信息
	iface := Interface{}
	find := msdb.Where(" name = ? ", name).Find(&iface).Error
	if find != nil {
		//无效请求
		return nil
	}
	return &iface
}
