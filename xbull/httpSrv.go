package main

import (
	"encoding/json"
	"fmt"
	//"time"
	//"os"
	"io/ioutil"
	"log"
	"net/http"
	//"strconv"
)

/*
IdColliery string	`json:"id_colliery"`
Dsn    string         `json:"dsn"`
IntervalHeartbeatApp int  `json:"interval_heartbeat_app"`
IntervalReConnect int `json:"interval_reconnect"`
TimeoutEndPointData int `json:"timeout_endpoint_data"`
DurationEndPointData int `json:"duration_endpoint_data"`
MinEndPointData float64 `json:"min_endpoint_data"`
NumEndPoint int `json:"num_endpoint"`
EventEPDataChange bool `json:"event_epdata_change"`
SaveEPDataZero bool `json:"save_epdata_zero"`
EndPoints []EndPointAttr `json:"endpoints"`

Ipaddr   string `json:"ipaddr"`
Tcpport int    `json:"tcpport"`
Unit    int    `json:"unit"` // 1,2,3
ProtNb     int    `json:"protnb"`   //   -5 ï½ž 5
EpId   string    `json:"epid"`
*/

var htmlbody = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<title>采集软件配置与连通性测试</title>
<style>
	input,div{padding:1px 1px 1px 1px;color:blue;}
	div.table{width:100%;display:table;background-color:#FDFDFD;}
	div.table>div{display:table-row;}
	div.table>div>div{display:table-cell;vertical-align:middle;}
	.lb{width:210px;text-align:right;color:gray;font-size:0.9em;}
	
	button.funcbtn{ 
		width:120px; height:36px; 
		text-align:center; vertical-align: middle;
		font-size:14px; border:0px solid #FFF; 
		margin:15px auto; padding:0 auto; 
		border-radius:20px; color:#FFFF;
		background-color:yellowgreen;
		/*display:block;*/
	}
	
	.area{padding:50px 50px; font-size:1em;width:700px;}
	.title{padding:20px 20px 20px 0px; font-size:1.3em; text-align:center}
	
	#EventEPDataChange,#SaveEPDataZero{width:22px;height:22px;}
	#dsn{width:400px;}
	
	.connected {background-color:lightgreen;width:25px; height:15px;margin:2px 1px 2px 10px;}
	.connecting {background-color:lightyellow;width:25px; height:15px;margin:2px 1px 2px 10px;}
	.connectnk {background-color:tomato;width:25px; height:15px;margin:2px 1px 2px 10px;}
	
	.epid{width:120px;}
	.ipaddr{width:120px;}
	.tcpport{width:120px;}
	.unit{width:120px;}
	.protNb{width:120px;}
</style>
</head>
<body>
<div class="area">
	<div class ="title" style="text-align:center"> 采集软件配置与连通性测试 </div>
	<div class="table" id="tbconfig">
		<div> <div class="lb">煤矿编号:</div><div><input id="id_colliery" type="text" value= "140000000" /></div> </div>
		<div> <div class="lb">Dsn:</div><div><input id="dsn" type="text" value= "root:13811237916sS@tcp(localhost:3306)/bridge" /> </div> </div>
		<div> <div class="lb">心跳消息间隔(秒):</div><div><input id="interval_heartbeat_app" type="tel" value= "60" /> </div> </div>
		<div> <div class="lb">网络重连时间间隔(秒):</div><div><input id="interval_reconnect" type="tel" value= "60" /> </div> </div>
		<div> <div class="lb">数据超时间隔(秒):</div><div><input id="timeout_endpoint_data" type="tel" value= "60" /> </div> </div>
		<div> <div class="lb">数据持续时间(秒):</div><div><input id="duration_endpoint_data" type="tel" value= "10" /> </div> </div>
		<div> <div class="lb">数据最小值(吨):</div><div><input id="min_endpoint_data" type="tel" value= "1.8" /> </div> </div>
		<div> <div class="lb">设备数量:</div><div><input id="num_endpoint" type="tel" value= "1" /> <button type="button" class=“numbtn” onclick="checkNumEP()"> 确认</button></div> </div>
		<div> <div class="lb">数据值变化事件是否写入数据库:</div><div>   <input id="event_epdata_change" type="checkbox" checked /> </div> </div>
		<div> <div class="lb">零值数据是否写入数据库:</div><div>   <input id="save_epdata_zero" type="checkbox" checked /> </div> </div>
	</div>

	<div class="title">设备列表</div>
	<div class="table"  id="endpoints">
		<div style="font-size:0.9em;color:gray;">
			<div>设备编号</div>
			<div>设备IP地址</div>
			<div>TCP端口</div>
			<div>数据单位</div>
			<div>协议编号</div>
			<div></div>
		</div>
		<div id="tempEP" style="display:none;">
			<div  class="epid" >  <input class ="epid" value="000"/> </div>
			<div class="ipaddr">  <input class="ipaddr" value="192.168.10.1"/> </div>
			<div class="tcpport">  <input class="tcpport" value="6000"/> </div>
			<div  class="unit">  <select class="unit">
				<option value="0">吨</option>
				<option value="1">公斤</option>
				</select> 
			</div>
			<div class="protnb">  <select class="protnb">
				<option value="1">1号协议</option>
				<option value="2">2号协议</option>
				<option value="3">3号协议</option>
				</select>
			</div>
			<div> <div class="testconn"> </div></div>
		</div>
	</div>

	<div class="title"> 
		<button type="button" class="funcbtn" onclick="PostJsonData()" id="submitBTN"> 提交</button>
		<button type="button" class="funcbtn" onclick="TestConn()" id="testBTN"> 连通性测试</button>
		<button type="button" class="funcbtn" onclick="Logout()" >退出</button>
	</div>

</div>
<script>

	var EndPoints;
	window.addEventListener("load", function(event) {
		console.log("All resources finished loading!");
		var xhr = new XMLHttpRequest();
		xhr.addEventListener("load", function(){
			var cfg = JSON.parse(this.responseText);
			var tbconfig = document.querySelector("#tbconfig");
			for(var key in cfg){
				if(key == "endpoints"){
					EndPoints = cfg[key];
					var endPS = document.querySelector("#endpoints");
					var endP = endPS.querySelector("#tempEP");
					for(var i=1,l=EndPoints.length;i<=l;i++){
						var NewEndP = document.createElement('div');
						NewEndP.id= "ep" + i;
						NewEndP.innerHTML = endP.innerHTML;
						var epattr = EndPoints[i-1];
						NewEndP.querySelector("input.epid").value=epattr.epid;
						NewEndP.querySelector("input.ipaddr").value=epattr.ipaddr;
						NewEndP.querySelector("input.tcpport").value=epattr.tcpport;
						NewEndP.querySelector("select.unit").value=epattr.unit;
						NewEndP.querySelector("select.protnb").value=epattr.protnb;
						NewEndP.querySelector("div.testconn").id="T"+epattr.epid;
						NewEndP.classList.add("endpoint");
						endPS.appendChild(NewEndP);
					}
				}else {
					var obj = tbconfig.querySelector("#"+key);
					if(!obj){
						continue
					}
					if(key == "event_epdata_change" || key == "save_epdata_zero"){
						obj.checked = cfg[key];
					}else  {
						obj.value = cfg[key];
					}
				}
			}
			
		});
		xhr.open("GET", "/GetConfig");
		xhr.send();
	});

	function PostJsonData(){
		var cfg = {};
		var tbconfig = document.querySelector("#tbconfig");
		var cfgs = tbconfig.querySelectorAll("input");
		for(i=0;i<cfgs.length;i++){
			var e = cfgs[i];
			if(e.id == "event_epdata_change" || e.id == "save_epdata_zero"){
				cfg[e.id] = e.checked;
			}else {
				cfg[e.id] = e.value;
			}
		}
		cfg["endpoints"] = [];
		var endPS = document.querySelector("#endpoints");
		var eps = endPS.querySelectorAll("div.endpoint");
		for(i=0; i< eps.length;i++){
			ep = eps[i];
			var epattr = {};
			epattr["epid"] = ep.querySelector("input.epid").value;
			epattr["ipaddr"] = ep.querySelector("input.ipaddr").value;
			epattr["tcpport"] = ep.querySelector("input.tcpport").value;
			epattr["unit"] = ep.querySelector("select.unit").value;
			epattr["protnb"] = ep.querySelector("select.protnb").value;
			cfg["endpoints"].push(epattr);
		}

		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/SetConfig",true);
		xhr.addEventListener("load", function(){
			var cfg = JSON.parse(this.responseText);
			window.location.href = "/"
		});
		xhr.send(JSON.stringify(cfg));
	}

	function TestConn(){
		var xhr = new XMLHttpRequest();
		xhr.addEventListener("load", function(){
			var cfg = JSON.parse(this.responseText);
			var endPS = document.querySelector("#endpoints");
			for(var key in cfg){
				var testconn = endPS.querySelector("#T"+key)
				testconn.classList.add(cfg[key]);
			}
		});
		xhr.open("GET", "/TestConn");
		xhr.send();		
	}
	function checkNumEP(){
		var NumEP = document.querySelector("#num_endpoint");
		if(NumEP.value<1){
			alert("must >=1");
			NumEP.focus();
			return;
		}
		if(NumEP.value>32){
			alert("too large!! must <=8");
			NumEP.focus();
			return;
		}

		var endPS = document.querySelector("#endpoints");
		var eps = endPS.querySelectorAll("div.endpoint");
		for (i = 0; i < eps.length; ++i) {
			endPS.removeChild(eps[i]);
		}

		var numEP = NumEP.value; 
		var endPS = document.querySelector("#endpoints");
		var endP = endPS.querySelector("#tempEP");
		for(var i=1,l=numEP;i<=l;i++){
			
			var NewEndP = document.createElement('div');
			NewEndP.id= "ep" + i;
			NewEndP.innerHTML = endP.innerHTML;

			if (EndPoints &&  i <= EndPoints.length) {
				var epattr = EndPoints[i-1];
				NewEndP.querySelector("input.epid").value=epattr.epid;
				NewEndP.querySelector("input.ipaddr").value=epattr.ipaddr;
				NewEndP.querySelector("input.tcpport").value=epattr.tcpport;
				NewEndP.querySelector("select.unit").value=epattr.unit;
				NewEndP.querySelector("select.protnb").value=epattr.protnb;
				NewEndP.querySelector("div.testconn").id="T"+epattr.epid;
			} else {
				var num = ""+i;
				var epid = ('0').repeat(3-num.length) + num;
				NewEndP.querySelector("input.epid").value=epid;
			}
			NewEndP.classList.add("endpoint");
			endPS.appendChild(NewEndP);
		}

	 }
	 function Logout(){
		var loginData = {"uname":"","upass":""};
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/Login",true);
		xhr.addEventListener("load", function(){
			var cfg = JSON.parse(this.responseText);
			window.location.href = "/"
		});
		xhr.send(JSON.stringify(loginData));
	 }
</script>
</body>
</html>
`

var loginForm = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<title>采集软件配置与连通性测试登录</title>
<style>
	input,div{padding:1px 1px 1px 1px;color:blue;}

	button{ 
		width:100px; height:32px; 
		text-align:center; vertical-align: middle;
		font-size:14px; border:0px solid #FFF; 
		margin:15px auto; padding:0 auto; 
		border-radius:20px; color:#FFFF;
		background-color:yellowgreen;
		/*display:block;*/
	}
	
	.area{padding:50px 50px; font-size:1em;width:500px;background-color:#FDFDFD;}
	.title{padding:20px 20px 20px 0px; font-size:1.3em; }
	
	label {
		display: block;
		font: 0.9rem 'Fira Sans', sans-serif;
	}

	input[type='submit'],
	label {
		font:1rem;
		margin-top: 1rem;
	}
	

</style>
</head>
<body>
<div class="area">
	<div class ="title"> 您需要登录！ </div>
	<div>
		<label for="username">用户名:</label>
		<input type="text" id="username" name="username">
	</div>

	<div>
    <label for="pass">密码:</label>
    <input type="password" id="pass" name="password"
           minlength="8" required>
	</div>
	<div class = "title"></div>
	<button onclick="PostJsonData()"> 登录</button>
	
</div>
<script>
	function PostJsonData(){
		var loginData = {};
		var uname = document.querySelector("#username");
		var upass = document.querySelector("#pass");
		loginData["uname"]=uname.value;
		loginData["upass"]=upass.value;
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/Login",true);
		xhr.addEventListener("load", function(){
			var cfg = JSON.parse(this.responseText);
			window.location.href = "/"
		});
		xhr.send(JSON.stringify(loginData));
	}
</script>
</body>
</html>
`

var logined = false

func httpsrv(sqlChannel chan string) {
	Log.Println("Start HttpServer for UIConfig Manager....")
	handler := func(w http.ResponseWriter, r *http.Request) {
		if !logined {
			fmt.Fprintf(w, loginForm)
			return
		}
		fmt.Fprintf(w, htmlbody)
	}
	handlerGetConfig := func(w http.ResponseWriter, r *http.Request) {
		if !logined {
			fmt.Fprintf(w, loginForm)
			return
		}
		ra, err := ioutil.ReadFile("xbull.json")
		if err != nil {
			//log.Println(err)
			fmt.Fprintf(w, `{"msg":"error"}`)
		}
		fmt.Fprintf(w, string(ra))
	}
	handlerSetConfig := func(w http.ResponseWriter, r *http.Request) {
		if !logined {
			fmt.Fprintf(w, loginForm)
			return
		}
		ra := `{"result":"Ok"}`
		body, _ := ioutil.ReadAll(r.Body)
		bodystr := string(body)
		DbMgr.saveConfig(bodystr)
		ioutil.WriteFile("xbull.json", []byte(bodystr), 0666) //os.ModeAppend)
		fmt.Fprintf(w, ra)
		sqlChannel <- "quit......."
	}
	handlerTestConn := func(w http.ResponseWriter, r *http.Request) {
		if !logined {
			fmt.Fprintf(w, loginForm)
			return
		}
		b, err := json.Marshal(EPStatusLst)
		//fmt.Println("b:", string(b), "err:",err)
		if err != nil {
			fmt.Fprintf(w, `{"err":"1"}`)
			return
		}
		fmt.Fprintf(w, string(b))
	}
	handlerLogin := func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		txt := string(body)
		//fmt.Println("bodyData:",txt)

		type loginData struct {
			UName string `json:"uname"`
			UPass string `json:"upass"`
		}
		m := loginData{}
		err := json.Unmarshal([]byte(txt), &m)
		//fmt.Println("Uname:",m.UName, "UPass:" ,m.UPass)
		if err == nil && DbMgr.CheckLogin(m.UName, m.UPass) {
			logined = true
			fmt.Fprintf(w, `{"result":"Ok"}`)
			return
		}
		logined = false
		fmt.Fprintf(w, `{"result":"Err"}`)
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/Login", handlerLogin)
	http.HandleFunc("/GetConfig", handlerGetConfig)
	http.HandleFunc("/SetConfig", handlerSetConfig)
	http.HandleFunc("/TestConn", handlerTestConn)
	log.Fatal(http.ListenAndServe(":8864", nil))
}

func startHttpSrv(sqlChannel chan string) {
	go httpsrv(sqlChannel)
}
