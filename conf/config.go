package conf

import (
	"gin_mall/dao"
	"strings"

	"gopkg.in/ini.v1"
)

var (
	AppModel    string
	HttpPort    string
	DB          string
	DbHost      string
	DbPort      string
	DbUser      string
	DbPassword  string
	DbName      string
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
	VaildEmail  string
	SmtpHost    string
	SmtpEmail   string
	SmtpPass    string
	Host        string
	ProductPath string
	AvatarPath  string
)

func Init() {
	//本地读取环境变量
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		panic(err)
	}
	LoadServer(file)
	LoadMysql(file)
	LoadRedis(file)
	LoadEmail(file)
	LoadPhotoPath(file)
	//mysql read 主
	pathRead := strings.Join([]string{DbUser, ":", DbPassword, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=true"}, "")
	//mysql write 从
	pathWrite := strings.Join([]string{DbUser, ":", DbPassword, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=true"}, "")
	dao.Database(pathRead, pathWrite)
}
func LoadServer(file *ini.File) {
	AppModel = file.Section("service").Key("AppModel").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}
func LoadMysql(file *ini.File) {
	DB = file.Section("mysql").Key("DB").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassword = file.Section("mysql").Key("DbPassword").String()
	DbName = file.Section("mysql").Key("DbName").String()
}
func LoadRedis(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDb = file.Section("redis").Key("RedisDb").String()
}
func LoadEmail(file *ini.File) {
	VaildEmail = file.Section("email").Key("VaildEmail").String()
	SmtpHost = file.Section("email").Key("SmtpHost").String()
	SmtpEmail = file.Section("email").Key("SmtpEmail").String()
	SmtpPass = file.Section("email").Key("SmtpPass").String()
}
func LoadPhotoPath(file *ini.File) {
	Host = file.Section("path").Key("Host").String()
	ProductPath = file.Section("path").Key("ProductPath").String()
	AvatarPath = file.Section("path").Key("AvatarPath").String()
}
