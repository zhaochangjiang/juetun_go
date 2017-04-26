package models

//初始化后台的Models
func InitFrontendModels() {
	//注册带前缀的Model
	//orm.RegisterModelWithPrefix("prefix_", new(User))

	//注册不带前缀的Model
	//orm.RegisterModel(new(User))
}
