package models

//初始化前台的Models
func InitBackendModels() {
	//注册带前缀的Model
	//orm.RegisterModelWithPrefix("prefix_", new(User))

	//注册不带前缀的Model
	//orm.RegisterModel(new(User))
}
