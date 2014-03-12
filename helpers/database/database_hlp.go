package database

import (
	"github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/thrsafe"
	"log"
	"os"
)

var (
	// MySQL Connection Handler
	CurtDevDb = autorc.New("tcp", "", "127.0.0.1:3306", "root", "", "CurtDev")
	AdminDb   = autorc.New("tcp", "", "127.0.0.1:3306", "root", "", "admin")
	Pcdb      = autorc.New("tcp", "", "127.0.0.1:3306", "root", "", "pcdb")
	Vcdb      = autorc.New("tcp", "", "127.0.0.1:3306", "root", "", "vcdb")
)

func BindDatabase() {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		curtdev_name := os.Getenv("CURT_DEV_NAME")
		pcdb_name := os.Getenv("PCDB_NAME")
		vcdb_name := os.Getenv("VCDB_NAME")
		admin_name := os.Getenv("ADMIN_NAME")
		CurtDevDb = autorc.New(proto, "", addr, user, pass, curtdev_name)
		AdminDb = autorc.New(proto, "", addr, user, pass, admin_name)
		Vcdb = autorc.New(proto, "", addr, user, pass, vcdb_name)
		Pcdb = autorc.New(proto, "", addr, user, pass, pcdb_name)
	}
}

func MysqlError(err error) (ret bool) {
	ret = (err != nil)
	if ret {
		log.Println("MySQL error: ", err)
	}
	return
}

func MysqlErrExit(err error) {
	if MysqlError(err) {
		os.Exit(1)
	}
}
