package general

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type DataConfig struct {
	DBName       string
	DatabaseName string
	Username     string
	Password     string
	Prefix       string
	MaxConn      int
	MaxIdle      int
}
type DataObject struct {
	DataConfigParam []*DataConfig
}

//初始化数据库连接
func (this *DataObject) InitConnect() {

	//	beego.AppConfig.String(“dev::mysqluser”)
	for _, dataConfig := range this.DataConfigParam {
		// 提示数据库连接初始化成功
		orm.RegisterDataBase(dataConfig.DBName, "mysql", dataConfig.Username+":"+dataConfig.Password+"@/"+dataConfig.DatabaseName+"?charset=utf8", dataConfig.MaxIdle, dataConfig.MaxConn)

	}

}

//初始化所有的数据库
func InitDatabase() {
	//设置数据库调试模式
	orm.Debug = true
	//beego.LoadAppConfig("ini", "../common/conf/database.conf")
	orm.RegisterDriver("mysql", orm.DRMySQL)

	//dataNameArray := [...]string{"db_data", "db_admin", "db_user"}
	dataNameArray := beego.AppConfig.Strings("dbList")
	dataConfigParam := make([]*DataConfig, 0)
	for _, dataName := range dataNameArray {
		//加载数据库配置文件信息
		dConfig := new(DataConfig)
		dConfig.DBName = dataName
		//		beego.Info(dConfig.DBName)
		dConfig.DatabaseName = beego.AppConfig.String(dataName + "::databaseName")
		dConfig.Username = beego.AppConfig.String(dataName + "::username")
		dConfig.Password = beego.AppConfig.String(dataName + "::pwd")
		dConfig.Prefix = beego.AppConfig.String(dataName + "::prefix")
		maxConn, _ := beego.AppConfig.Int(dataName + "::maxConn")
		dConfig.MaxConn = maxConn
		maxIdle, _ := beego.AppConfig.Int(dataName + "::maxIdle")
		dConfig.MaxIdle = maxIdle
		dataConfigParam = append(dataConfigParam, dConfig)

	}

	//初始化数据库连接
	dataObject := new(DataObject)
	dataObject.DataConfigParam = dataConfigParam
	dataObject.InitConnect()

}
