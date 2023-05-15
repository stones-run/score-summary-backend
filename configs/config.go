package configs

var Config *config

type config struct {
	Mysql MysqlConf
}

type MysqlConf struct {
	UserName string
	Password string
	IpHost   string
	DbName   string
}

func init() {
	Config = new(config)
	Config.Mysql.UserName = "root"
	Config.Mysql.Password = ""
	Config.Mysql.IpHost = "127.0.0.1:3306"
	Config.Mysql.DbName = "score-summary"

}
