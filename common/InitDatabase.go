package common

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type DataConfig struct {
	DBName       string
	DatabaseName string
	Username     string
	Password     string
	MaxConn      int
	MaxIdle      int
}
type DataObject struct {
	DataConfigParam []*DataConfig
}

func (this *DataObject) InitConnect() {
	//	beego.AppConfig.String(“dev::mysqluser”)
	for _, dataConfig := range this.DataConfigParam {
		orm.RegisterDataBase(dataConfig.DBName, "mysql", dataConfig.Username+":"+dataConfig.Password+"@/"+dataConfig.DatabaseName+"?charset=utf8", dataConfig.MaxIdle, dataConfig.MaxConn)

	}

}
func InitDatabase() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	//dataNameArray := [...]string{"db_data", "db_admin", "db_user"}
	dataNameArray := beego.AppConfig.Strings("dbList")
	dataConfigParam := make([]*DataConfig, 0)
	for _, dataName := range dataNameArray {
		dConfig := new(DataConfig)
		dConfig.DBName = dataName
		dConfig.DatabaseName = beego.AppConfig.String(dataName + "::databaseName")
		dConfig.Username = beego.AppConfig.String(dataName + "::username")
		dConfig.Password = beego.AppConfig.String(dataName + "::pwd")
		maxConn, _ := beego.AppConfig.Int(dataName + "::maxConn")
		dConfig.MaxConn = maxConn
		maxIdle, _ := beego.AppConfig.Int(dataName + "::maxIdle")
		fmt.Println(maxIdle)
		dConfig.MaxIdle = maxIdle

		dataConfigParam = append(dataConfigParam, dConfig)
		fmt.Println(dataConfigParam)
	}

	dataObject := new(DataObject)
	dataObject.DataConfigParam = dataConfigParam
	dataObject.InitConnect()

}
