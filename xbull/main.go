package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	//"math/rand"
	"flag"
	"io"
	"log"
	"os"
	"time"
	//"encoding/hex"
)

const standFormat = "2006-01-02 15:04:05.000"
const VERSION string = "2.0.16"
const EVENT_APP_START = "A0"
const EVENT_APP_HEART = "A1"
const EVENT_EP_Connected = "E0"
const EVENT_EP_Connecting = "E1"
const EVENT_EP_ConnBreak = "E2"
const EVENT_EP_HEART = "E3"
const EVENT_EP_TIMEOUT = "E4"
const EVENT_EP_DATA = "E5"

//var dsn = flag.String("dsn", "root:13811237916sS@tcp(localhost:3306)/bridge", "DataSourceName Of Mysql")
var dsn = flag.String("dsn", "", "DataSourceName Of Mysql")
var doLog = flag.Bool("log", false, "enable log")
var help = flag.Bool("h", false, "print this help")
var testDB = flag.Bool("testDB", false, "test mysql")

var EPStatusLst = make(map[string]string)

var (
	Log   *log.Logger
	Err   *log.Logger
	DbMgr = DBLayer{}
)
var startTime = 

func init() {
	errFile, err := os.OpenFile("xbullErr.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Open Err LogFile Error!", err)
	}
	Err = log.New(io.MultiWriter(os.Stdout, errFile), "", log.Ldate|log.Ltime) // | log.Lshortfile)
	Log = log.New(ioutil.Discard, "", log.LstdFlags)
}

func main() {

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}

	if *doLog {
		logFile, err := os.OpenFile("xbull.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Open LogFile Error!", err)
		}
		Log = log.New(io.MultiWriter(os.Stdout, logFile), "", log.Ldate|log.Ltime) // | log.Lshortfile)
	}
	Err.Println("xbull start VERSION", VERSION)

	if *testDB {
		db := DBLayer{}
		err := db.Init(*dsn)
		if err != nil {
			Err.Fatalln("Can't Init MySQL Database! Please check DSNString err:", err)
		}

		checkLogin := func(uname string, passwd string) {
			loginOk := db.CheckLogin(uname, passwd)
			fmt.Println("UName:", uname, "UPasswd:", passwd, "LoginOk:", loginOk)
		}

		checkLogin("admin", "13910580009")
		checkLogin("admin", "13910580008")
		checkLogin("admin2", "13910580009")

		//db.Insert("140000000", "001", 8888.88, "京JR8899")
		//db.InsertEvent("140000000", "000", EVENT_APP_START, VERSION)

		json, err := db.loadConfig()
		if err != nil {
			fmt.Println("db.loadConfig() Err:", err)
		}
		fmt.Println("db.loadConfig Ok!!")

		fmt.Println()
		fmt.Println()
		fmt.Println("Test Config from jsonString")
		cfg := Config{}
		cfg.fromJsonString(json)
		fmt.Println(cfg.toJsonString())
		return
	}

	cfg := Config{}
	if _, err := os.Stat("xbull.json"); err == nil && *dsn == "" {
		cfg.LoadFromFile("xbull.json")
		*dsn = cfg.Dsn
	}

	sqlChannel := make(chan string, 64)
	defer close(sqlChannel)

	//hDb := DBLayer{}
	//fmt.Println(*dsn)
	err := DbMgr.Init(*dsn)
	if err == nil {
		//Err.Fatalln("Can't Init MySQL Database! Please check DSNString err:", err)
		json, err := DbMgr.loadConfig()
		if err != nil {
			Err.Fatalln("DB.LoadConfig() Err:", err)
		}

		cfg.fromJsonString(json)
		cfg.Dsn = *dsn
		cfg.SaveToFile("xbull.json")

		for _, epcfg := range cfg.EndPoints {
			EPoint := NewEndPoint()
			EPoint.Init(cfg, epcfg, sqlChannel)
			//defer EPoint.Close()
		}
	}

	startHttpSrv(sqlChannel)
	sqlChannel <- buildEventSql(cfg.IdColliery, "000", EVENT_APP_START, VERSION)

	Log.Println("Start RunLoop.....")
ExitFor:
	for {
		select {
		case sql := <-sqlChannel:
			{
				//Log.Println("mainLoop, sql:",sql)
				if strings.HasPrefix(sql, "status") {
					a := strings.Split(sql, ",")
					if len(a) == 3 {
						EPStatusLst[a[1]] = a[2]
					}
				} else if strings.HasPrefix(sql, "quit") {

					Log.Println("Recv a Quit sign , will quit......")
					time.Sleep(5000)
					break ExitFor

				} else {
					if DbMgr.Inited() {
						err := DbMgr.ExecSql(sql)
						if err != nil {
							Err.Fatalln("EXECSQL Err!!!!!!!!!")
						}
					}
				}
			}
		case <-time.After(time.Duration(cfg.IntervalHeartbeatApp) * time.Second):
			{
				sqlChannel <- buildEventSql(cfg.IdColliery, "000", EVENT_APP_HEART, VERSION)
			}
		}
	}

	Log.Println("xbull quited......")
}
