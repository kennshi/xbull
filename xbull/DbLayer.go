package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

type DBLayer struct {
	//db mysql.MySQLDriver
	hDb    *sql.DB
	inited bool
}

//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...¶mN=valueN]
// user@unix(/path/to/socket)/dbname
// root:pw@unix(/tmp/mysql.sock)/myDatabase?loc=Local
// user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
// user:password@/dbname?sql_mode=TRADITIONAL
// user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname?timeout=90s&collation=utf8mb4_unicode_ci
// id:password@tcp(your-amazonaws-uri.com:3306)/dbname
// user@cloudsql(project-id:instance-name)/dbname
// user@cloudsql(project-id:regionname:instance-name)/dbname
// user:password@tcp/dbname?charset=utf8mb4,utf8&sys_var=esc%40ped
// user:password@/dbname
// user:password@/

func (db *DBLayer) Inited() bool {
	return db.inited
}

func (db *DBLayer) Init(connString string) error {
	db.inited = false
	hDb, _ := sql.Open("mysql", connString)
	_, err := hDb.Exec("select value from config where item like 'id_colliery';")
	if err != nil {
		return err
	}
	db.inited = true
	db.hDb = hDb
	return nil
}

func (db *DBLayer) Insert(collieryId string, bridgeId string, tunNum float64, vehiPlate string) error {
	timestamp := time.Now().Format(standFormat)
	hDb := db.hDb
	sql := fmt.Sprintf("insert into weightdata (recId,collieryId, bridgeId,vehiNum,beginTime,valTime,endTime,weightValue) values (replace(uuid(),'-',''),'%s', '%s', '%s','%s','%s','%s',%f)", collieryId, bridgeId, vehiPlate, timestamp, timestamp, timestamp, tunNum)
	//sql := fmt.Sprintf("insert into weightdata (recId,weighttime,weightvalue,vehNum,collieryID,bridgeID) VALUES (replace(uuid(),'-',''),'%s',%f,'%s','%s','%s')",timestamp,tunNum,vehiPlate,collieryId,bridgeId)
	//sql := buildDataSql(collieryId, bridgeId, timestamp, timestamp, timestamp, tunNum)
	Log.Println(sql)
	_, err := hDb.Exec(sql)
	if err != nil {
		Err.Println("Insert Err! :", err)
		Err.Println("INSERT SQL:", sql)
		return err
	}
	return nil
}

func buildProcSql(collieryId string, bridgeId string, tunNum float64, vehiPlate string) string {
	timestamp := time.Now().Format(standFormat)
	sql := fmt.Sprintf("insert into weightproc (recId,weighttime,weightvalue,vehNum,collieryID,bridgeID) VALUES (replace(uuid(),'-',''),'%s',%f,'%s','%s','%s')", timestamp, tunNum, vehiPlate, collieryId, bridgeId)
	return sql
}

func buildEventSql(collieryId string, bridgeId string, eventId string, eventInfo string) string {
	timestamp := time.Now().Format(standFormat)
	sql := fmt.Sprintf("insert into eventlog (collieryid, bridgeid,eventtime,eventid,eventinfo) values ('%s', '%s', '%s','%s','%s')", collieryId, bridgeId, timestamp, eventId, eventInfo)
	return sql
}

func buildDataSql(collieryId string, bridgeId string, begTime string, valTime string, endTime string, tonNum float64) string {
	sql := fmt.Sprintf("insert into weightdata (recId,collieryId, bridgeId,vehiNum,beginTime,valTime,endTime,weightValue) values (replace(uuid(),'-',''),'%s', '%s', '','%s','%s','%s',%f)", collieryId, bridgeId, begTime, valTime, endTime, tonNum)
	return sql
}

func (db *DBLayer) InsertEvent(collieryId string, bridgeId string, eventId string, eventInfo string) error {
	timestamp := time.Now().Format(standFormat)
	hDb := db.hDb
	sql := fmt.Sprintf("insert into eventlog (collieryid, bridgeid,eventtime,eventid,eventinfo) values ('%s', '%s', '%s','%s','%s')", collieryId, bridgeId, timestamp, eventId, eventInfo)
	Log.Println(sql)
	_, err := hDb.Exec(sql)
	if err != nil {
		Err.Println("InsertEvent Err!:", err)
		Err.Println("INSERT SQL:", sql)
	}
	return nil
}

func (db *DBLayer) CheckLogin(name string, passwd string) bool {
	if !db.Inited() {
		return false
	}
	hDb := db.hDb

	sql := fmt.Sprintf("select passwd = md5(concat(salt,'%s')) as loginok from users where name= '%s';", passwd, name)
	r, err := hDb.Query(sql)
	if err != nil {
		return false
	}
	defer r.Close()
	r.Next()
	var loginok int = 0
	err = r.Scan(&loginok)
	if err != nil {
		return false
	}

	if loginok == 0 {
		return false
	}

	return true
}

func (db *DBLayer) ExecSql(sql string) error {
	hDb := db.hDb
	Log.Println(sql)
	_, err := hDb.Exec(sql)
	if err != nil {
		Err.Println("SQL:", sql, " ERR:", err)
		return err
	}
	return nil
}

func (db *DBLayer) saveConfig(cfg string) error {
	if !db.Inited() {
		return errors.New("mysql not yet inited")
	}
	var rr interface{}
	json.Unmarshal([]byte(cfg), &rr)
	conf, ok := rr.(map[string]interface{})
	if ok {
		for k, v := range conf {
			if k != "endpoints" {
				val, ok := v.(string)
				if ok {
					sql := fmt.Sprintf("UPDATE config SET value = '%s' WHERE item like '%s';", val, k)
					db.ExecSql(sql)
				} else {
					val, _ := v.(bool)
					sql := fmt.Sprintf("UPDATE config SET value = '%s' WHERE item like '%s';", strconv.FormatBool(val), k)
					db.ExecSql(sql)
				}

			} else {
				eps, _ := v.([]interface{})
				db.ExecSql("delete from endpoints;")
				for _, epr := range eps {
					ep, _ := epr.(map[string]interface{})
					ipaddr, _ := ep["ipaddr"].(string)
					tcpport, _ := ep["tcpport"].(string)
					unit, _ := ep["unit"].(string)
					protnb, _ := ep["protnb"].(string)
					epid, _ := ep["epid"].(string)
					sql := fmt.Sprintf("INSERT INTO endpoints (ipaddr, tcpport, unit, protnb, epid) VALUES ('%s', %s, %s, %s, '%s');",
						ipaddr, tcpport, unit, protnb, epid)
					db.ExecSql(sql)
				}

			}
		}

	}
	return nil
}

func (db *DBLayer) loadConfig() (string, error) {
	const rn = "\r\n"
	hDb := db.hDb
	json := "{\r\n\t\"dsn\":\"\",\r\n"
	numEndPoint := 0
	{
		sql := "select item, value, itemtype from config"
		r, err := hDb.Query(sql)
		if err != nil {
			Err.Fatalln("DbLayer_loadConfig err:", err)
		}
		defer r.Close()
		for r.Next() {
			var (
				item     string
				value    string
				itemtype string
			)
			if err := r.Scan(&item, &value, &itemtype); err != nil {
				Log.Fatalln("DbLayer_loadConfig r.Scan err:", err)
			}
			json += "\t\"" + item + "\":"
			if itemtype == "S" {
				json += "\""
			}

			if itemtype == "T" {
				if value == "1" || value == "true" {
					value = "true"
				} else {
					value = "false"
				}
			}
			json += value
			if itemtype == "S" {
				json += "\""
			}
			json += ",\r\n"

			if item == "num_endpoint" {
				numEndPoint, err = strconv.Atoi(value)
				if err != nil {
					Err.Fatalln("DbLayer_loadConfig get  NumOfEndPoint err:", err)
				}
			}
		}
	}
	{
		sql := "select ipaddr, tcpport, unit, protnb, epid from endpoints order by  epid limit " + strconv.Itoa(numEndPoint)
		r, err := hDb.Query(sql)
		if err != nil {
			Err.Fatalln("DbLayer_loadConfig err:", err)
		}
		defer r.Close()

		json += "\t\"endpoints\": [\r\n"
		for r.Next() {
			var (
				ipaddr  string
				tcpport int
				unit    int
				portnb  int
				epid    string
			)
			if err := r.Scan(&ipaddr, &tcpport, &unit, &portnb, &epid); err != nil {
				Log.Fatalln("DbLayer_loadConfig r.Scan err:", err)
			}
			json += "\t\t{\r\n"
			json += "\t\t\t\"ipaddr\":\"" + ipaddr + "\",\r\n"
			json += "\t\t\t\"tcpport\":" + strconv.Itoa(tcpport) + ",\r\n"
			json += "\t\t\t\"unit\":" + strconv.Itoa(unit) + ",\r\n"
			json += "\t\t\t\"protnb\":" + strconv.Itoa(portnb) + ",\r\n"
			json += "\t\t\t\"epid\":\"" + epid + "\"\r\n"
			json += "\t\t},\r\n"
		}
		json = json[:len(json)-3]
		json += "\t\r\n]\r\n"
	}
	json += "}\r\n"

	return json, nil
}

func (db *DBLayer) Close() error {
	db.hDb.Close()
	return nil
}
