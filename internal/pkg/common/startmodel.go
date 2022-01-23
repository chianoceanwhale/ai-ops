package common

type (
	Mode string
)

const (
	ModeDev   Mode = "dev"    //开发模式
	ModeTest  Mode = "test"   //测试模式
	ModeProd  Mode = "prod"   //生产模式
	ModeDebug Mode = "debug"  //调试模式
	Mysql          = "mysql"  //mysql数据库标识
	Sqlite         = "sqlite" //sqlite
)
