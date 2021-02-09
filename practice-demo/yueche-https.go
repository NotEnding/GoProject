package main

import (
  "../alipay"
  "./gomini/gocrypto"
  _ "./mysql"
  "fmt"
  // "regexp"
  "database/sql"
  "net/http"
  "strings"
  // "html"
  "bufio"
  "bytes"
  "crypto/hmac"
  "crypto/md5"
  "crypto/sha256"
  "encoding/base64"
  "encoding/hex"
  "encoding/json"
  "encoding/xml"
  "io/ioutil"
  "log"
  "math/rand"
  "os"
  "reflect"
  "strconv"
  "time"
)


//结构体
type Server struct {
  ServerName string
  ServerIP   string
}

type Serverslice struct {
  Servers []Server
  ServersID  string
}

//使用 tag 标记要返回的字段名。
type Userlogin struct {
 // Phone string 'json:"phone"'
 // Password string 'json:"password"'
  Phone string
  Password string
}

type UserloginRet struct {
  Ret int
  Message string
  Data struct {
    Token string
    Id    string
  } 
}


////////////////////////////////
// tonglar API json传输协议定义 --- begin
////////////////////////////////

type InputMessage struct {
  Body string
  Signature string
}

////////////////////////////////
// tonglar API json传输协议定义 --- end
////////////////////////////////

var signal string = "1234567890abcdef"
var db *sql.DB

func dbinit() {
  var err error
  db,err = sql.Open("mysql", "root:Yueche520!@#!@tcp(47.104.157.203:3306)/YUECHE_DB?charset=utf8&parseTime=true&loc=Local")
  if err != nil {
   log.Fatalf("Error on initializing database connection: %s", err.Error())
  }

  db.SetMaxOpenConns(200)
  db.SetMaxIdleConns(100)

  err = db.Ping() 
  if err != nil {
      log.Fatalf("Error on opening database connection: %s", err.Error())
  }
}

var alipay_client alipay.Client

func alipayinit() {
  /*
  //init alipay params
  alipay.AlipayPartner = "0000000000000000"
  alipay.AlipayKey = "000000000000000000000"    
  alipay.WebReturnUrl = base_url + "/api/ali/apaynotify"
  alipay.WebNotifyUrl = base_url + "/api/ali/apaynotify"
  alipay.WebSellerEmail = "2068714482@qq.com"
  */

    alipay_client = alipay.Client{
    Partner:   "2017011004960162", // 合作者ID
    Key:       "MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAMNsE6OwBO//uaamZF1wVVaO/IvnZULR1G+pgcnFdACvSKm0BBn+RN5UkvuD6Ax4ERS1HO4DGMyDYJnw+zBsRoqW5JvA84pboph4vhw3eBCpJIr5htzSLMqyqj0v3Jt0lyQ67xFG/ZqGqzLap9ZNBrlYm3q0tv6IsZREDvQYj4ZHAgMBAAECgYEAjUiUxlHar/zNJtvDUf6F5AeKNEd94Ro8oOIG0G5tmJUhTne0Q2qeNbMldKt/14vypWrvWHBqvGj7LTCZGgAd2o8YTz05lwsWVAZntmuDtBs6CKxesDrc/xAXYy46NWnRxrLxnEMpnCqML0gJ+HRoeSFckimimrWbjYWlKvfSUTkCQQDjjgs4b64617wUYnUjC8ScT1oY8NX1H2CdWT0zlm0XU6NOw/wcN/9iKdlUN2HNQq3ygWEBzlt8DVxaIAfEPxRDAkEA29nCN2lTAcFqdG2QiCN8+Ceag8bOrhUzVYdbr/nkfoSn266aBHtA4QqZDtTcqACdYNSKCYEIR24iMrKCQIsHrQJAYNaoK8JLUTtSDRLBasKtTx/t5cNIKmLKCOxbQUL49f5f9zssZQ3nnuzUUiSneGSyBgvNLqmVATvmW2xaIcf+ZQJAWm7ykvSCLoCvF4FSKI3gg/tWdco7jiQuX4o0TujN8rUCjzz9IcbJY0iGuTEaKwlFs2T5+vrWuvs0mgIPzhjiaQJBAJtyVqii2FZZhqWRiYO/WvQQBGe4gTsjAOKS1oeBQnRGbIgRZg0SIr7STzFEe1Gwaoa5xqhUWrMzfJkEnQNilKI=", // 合作者私钥
    ReturnUrl: "", // 同步返回地址
    //NotifyUrl: "http://120.25.251.146:9999/api/ali/alipaynotify", // 网站异步返回地址
    NotifyUrl: "", // 网站异步返回地址
    Email:     "", // 网站卖家邮箱地址
  }

  fmt.Println(alipay_client.NotifyUrl)
}

func init() {
  dbinit()
  alipayinit()
}

var ZERO_TIME string = "2018-01-01 00:00:00"

// 生成32位MD5
func MD5(text string) string{
   ctx := md5.New()
   ctx.Write([]byte(text))
   return hex.EncodeToString(ctx.Sum(nil))
}

// return len=8  salt
func GetRandomSalt() string {
   return GetRandomString(8)
}

//生成随机字符串
func GetRandomString(lens int64) string{
   str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
   bytes := []byte(str)
   result := []byte{}
   r := rand.New(rand.NewSource(time.Now().UnixNano()))
   var i int64
   for i = 0; i < lens; i++ {
      result = append(result, bytes[r.Intn(len(bytes))])
   }
   return string(result)
}

// 生成随机4位数字
func GetRandomNum4Str() string {
   return GetRandomNum(4)
}

// 生成随机6位数字
func GetRandomNum6Str() string {
   return GetRandomNum(6)
}

//生成随机数字
func GetRandomNum(lens int64) string{
   str := "0123456789"
   bytes := []byte(str)
   result := []byte{}
   r := rand.New(rand.NewSource(time.Now().UnixNano()))
   var i int64
   for i = 0; i < lens; i++ {
      result = append(result, bytes[r.Intn(len(bytes))])
   }
   return string(result)
}

func ComputeHmac256(message string, secret string) string {
  key := []byte(secret)
  h := hmac.New(sha256.New, key)
  h.Write([]byte(message))
  return hex.EncodeToString(h.Sum(nil))
}

// 如果messageMAC是message的合法HMAC标签，函数返回真
func CheckMAC(message, messageMAC, key []byte) bool {
  mac := hmac.New(sha256.New, key)  
  mac.Write(message)
  expectedMAC := mac.Sum(nil)
  //sha := base64.StdEncoding.EncodeToString(expectedMAC)
  //fmt.Println(sha)
  //fmt.Println("hmac256:")
  //fmt.Println(hex.EncodeToString(expectedMAC))
  return hmac.Equal(messageMAC, []byte(hex.EncodeToString(expectedMAC)))
}


func main() {
  
  // 测试接口
  http.HandleFunc("/api/upload",uploadHandler) 
  //http.HandleFunc("/api/getversion", handler) 

  // API接口
  http.HandleFunc("/api/wechat_createuser", wechat_createuserHandle)   // 创建用户 - 测试通过
  http.HandleFunc("/api/wechat_binduser", wechat_binduserHandle)   // 绑定用户手机号 - 测试通过
  http.HandleFunc("/api/wechat_createdriver", wechat_createdriverHandle)   // 创建司机用户 - 测试通过
  http.HandleFunc("/api/wechat_binddriver", wechat_binddriverHandle)   // 绑定司机用户信息 - 测试通过
  http.HandleFunc("/api/wechat_getbanner", wechat_getbannerHandle)   // 获取首页BANNER图-测试通过
  http.HandleFunc("/api/wechat_getcoupon", wechat_getcouponHandle)   // 获取所有优惠券-测试通过

  http.HandleFunc("/api/wechat_setusercoupon", wechat_setusercouponHandle)   // 添加优惠券到个人账户-测试通过
  http.HandleFunc("/api/wechat_getusercoupon", wechat_getusercouponHandle)   // 获取个人账户已经领取的优惠券-测试通过
  http.HandleFunc("/api/wechat_getqa", wechat_getqaHandle)   // 获取出行贴士-测试通过 - 427
  http.HandleFunc("/api/wechat_getuserorder", wechat_getuserorderHandle)   // 获取订单（已完成，未完成，已退票，已改签）-测试通过 - 427
  http.HandleFunc("/api/wechat_getlineinfo", wechat_getlineinfoHandle)   // 获取车辆路线信息-测试通过
  http.HandleFunc("/api/wechat_gebuslineinfo", wechat_getbuslineinfoHandle)   // 获取长阳大巴车辆路线信息
  http.HandleFunc("/api/wechat_createbigorder", wechat_createbigorderHandle)   // 创建大巴车购票订单
  http.HandleFunc("/api/wechat_createsoloorder", wechat_createsoloorderHandle)   // 创建购票订单-测试通过 - 427
  http.HandleFunc("/api/wechat_delsoloorder", wechat_delsoloorderHandle)   // 删除购票订单-测试通过
  http.HandleFunc("/api/wechat_createmultiorder", wechat_createmultiorderHandle)   // 创建包车订单-测试通过 - 427
  http.HandleFunc("/api/wechat_searchline", wechat_searchlineHandle)   // 车次查询
  http.HandleFunc("/api/wechat_getlinestartend", wechat_getlinestartendHandle)   // 获取出所有发地和目的地-测试通过 - 427
  http.HandleFunc("/api/wechat_getlineendbystart", wechat_getlineendbystartHandle)

  http.HandleFunc("/api/wechat_getairstartend", wechat_getairstartendHandle)   // 获取机场专线出所有发地和目的地-测试通过 - 427
  http.HandleFunc("/api/wechat_getairlineendbystart", wechat_getairlineendbystartHandle)   // 根据起点城市，获取机场专线终点城市列表

  //http.HandleFunc("/api/wechat_getlineendbystart", wechat_getlineendbystartHandle)  // 交换
  http.HandleFunc("/api/wechat_getselfbill", wechat_getselfbillHandle)   // 获取个人账单-测试通过

  http.HandleFunc("/api/wechat_getselfinfo", wechat_getselfinfoHandle)   // 获取个人出行信息,包括在线天数，个人订单总数，个人出行总公里数-测试通过

  http.HandleFunc("/api/wechat_ticketcheck", wechat_ticketcheckHandle)  // 乘客扫码后，将相关信息，个人用户信息绑定发送给服务器端，验证是否购买了车票-测试通过

  http.HandleFunc("/api/wechat_todayticketget", wechat_todayticketgetHandle)  // 查询个人当天订单，并返回订单信息，-测试通过 -427


  http.HandleFunc("/api/wechat_drivergetticket", wechat_drivergetticketHandle)  // 查询司机个人点票数，返回点票信息-测试通过

  http.HandleFunc("/api/wechat_drivergetinfo", wechat_drivergetinfoHandle)  // 获取司机综合信息

  http.HandleFunc("/api/wechat_drivergettop", wechat_drivergettopHandle)  // 获取司机排行榜

  http.HandleFunc("/api/wechat_driversetgps", wechat_driversetgpsHandle)  // 司机将GPS设置进数据库
  http.HandleFunc("/api/wechat_getgps", wechat_getgpsHandle)  // 获取司机的GPS

  http.HandleFunc("/api/wechat_getsmscode", wechat_getsmscodeHandle)  // 请求获取短信验证码

  http.HandleFunc("/api/wechat_areabycity", wechat_areabystartcityHandle)   // 根据城市获取城市区域名称,乘客使用 -427
  http.HandleFunc("/api/wechat_stationbyarea", wechat_stationbyareaHandle)   // 根据区域名称获取站点信息,乘客使用 -427

  http.HandleFunc("/api/wechat_endorseticket", wechat_endorseticketHandle)   // 改签车票
  http.HandleFunc("/api/wechat_endorseticket_ex", wechat_endorseticketexHandle)   // 改签车票扩展接口
  http.HandleFunc("/api/wechat_refundticket", wechat_refundticketHandle)   // 退票

  http.HandleFunc("/api/wechat_getbustime", wechat_getbustimeHandle)   // 获取发车时间
  http.HandleFunc("/api/wechat_getbigbustime", wechat_getbigbustimeHandle)   // 获取大巴发车时间

  // 司机上车打开，绑定车辆到司机表中，并且将当前绑定信息写入到司机出车表中
  http.HandleFunc("/api/wechat_aboarddriver", wechat_aboarddriverHandle)
  // 司机下车
  http.HandleFunc("/api/wechat_offdriver", wechat_offdriverHandle)

  // 司机扫码或者通过手动输入查验乘客电子票
  http.HandleFunc("/api/wechat_drivercheckticket", wechat_drivercheckticketHandle)

  // 获取包车出发点
  http.HandleFunc("/api/wechat_getmultistartcity", wechat_getmultistartcityHandle)

  // 获取包车公共信息
  http.HandleFunc("/api/wechat_getmultimodel", wechat_getmultimodelHandle)

  // 乘客建议
  http.HandleFunc("/api/wechat_setsuggest", wechat_setsuggestHandle)

  // 乘客评价订单
  http.HandleFunc("/api/wechat_evaluateorder", wechat_evaluateorderHandle)

  // 小程序微信支付接口
  http.HandleFunc("/api/weixin/JSAPIpay",JSAPI_weixinpayHandle) // - 427
  http.HandleFunc("/api/weixin/JSAPIpaynotify",JSAPI_weixinpaynotifyHandle)

  http.HandleFunc("/api/weixin/wechat_ordercomplete", wechat_ordercompleteHandle)  // 后台退钱完成后，修改订单状态为"已完成"

  http.HandleFunc("/api/weixin/wechat_ordercomplete_ext", wechat_ordercompleteextHandle)  // 后台退钱完成后，修改订单状态为"已完成"  


  ////////////////////////////   管理端API  ////////////////////////////////////////
  // 获取所有订单（已完成，未完成，已退票，已改签，已取消）
  http.HandleFunc("/api/manage_getuserorder", manage_getuserorderHandle)
  // 获取所有用户信息
  http.HandleFunc("/api/manage_getuserinfo", manage_getuserinfoHandle)

  // 根据用户手机号获取对应用户信息
  http.HandleFunc("/api/manage_getuserinfobyphone", manage_getuserinfobyphoneHandle)
  // 获取所有司机信息
  http.HandleFunc("/api/manage_getdriverinfo", manage_getdriverinfoHandle)

  // 司机上下线
  http.HandleFunc("/api/manage_updriverstatus", manage_updriverstatusHandle)

  // 获取所有累计收入，累计订票数量，当前在线车辆数
  http.HandleFunc("/api/manage_getsystemdata", manage_getsystemdataHandle)

  // 根据时间查询当天的任意站点的出票信息（订单号，票号，用户名称，用户电话等）
  http.HandleFunc("/api/manage_getorderbydaytime", manage_getorderbydaytimeHandle)

  // 后台管理员登录
  http.HandleFunc("/api/manage_adminlogin", manage_adminloginHandle)

  // 发送乘客票务信息给司机
  http.HandleFunc("/api/manage_sendsmstodriver", manage_sendsmstodriverHandle)

  // 获取所有司机相关信息，id，phone，name
  http.HandleFunc("/api/manage_getdriverlist", manage_getdriverlistHandle)

  // 设置发车时间
  http.HandleFunc("/api/manage_setbustime", manage_setbustimeHandle)
  // 获取发车时间
  http.HandleFunc("/api/manage_getbustime", manage_getbustimeHandle)   
  // 根据id号删除某一条发车时间
  http.HandleFunc("/api/manage_delbustime", manage_delbustimeHandle)   

  // 获取当天已分配司机的订单
  http.HandleFunc("/api/manage_getallotorder", manage_getallotorderHandle)   

  // 获取line_type获取区域名称
  http.HandleFunc("/api/manage_getareabylinetype", manage_getareabylinetypeHandle)   

  // 站点上下线
  http.HandleFunc("/api/manage_uplocal", manage_uplocalHandle)   

  // 新增站点
  http.HandleFunc("/api/manage_insertlocal", manage_insertlocalHandle)   

  // 获取当前所有站点
  http.HandleFunc("/api/manage_getalllocal", manage_getalllocalHandle)  


  // 新增区域
  http.HandleFunc("/api/manage_insertarea", manage_insertareaHandle)   

  // 获取当前所有区域
  http.HandleFunc("/api/manage_getallarea", manage_getallareaHandle)  

  // 获取所有城市
  http.HandleFunc("/api/manage_getallcity", manage_getallcityHandle)  

  // 添加城市
  http.HandleFunc("/api/manage_insertcity", manage_insertcityHandle) 

  // 获取所有城市线路
  http.HandleFunc("/api/manage_getline", manage_getlineHandle)  

  // 修改线路价格
  http.HandleFunc("/api/manage_updatelineprice", manage_updatelinepriceHandle)  


  // 根据用户手机号，获取当前用户历史上买了多少张票
  http.HandleFunc("/api/manage_getticketbyphone", manage_getticketbyphoneHandle) 




 // http.ListenAndServe(":9999", nil)
 // http.ListenAndServe(":9996", nil)
  err := http.ListenAndServeTLS(":9996","1929239_api.yueche520.com.pem","1929239_api.yueche520.com.key",nil)
  if err != nil {
    log.Fatal("ListenAndServeTLS:", err.Error())
  }
}

func get_type(i interface{}) {
  fmt.Println("reflect:", reflect.TypeOf(i))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()

  file, _, err := r.FormFile("file")
    if err != nil {
        log.Fatal("FormFile: ", err.Error())
        return
    }
    defer func() {
        if err := file.Close(); err != nil {
            log.Fatal("Close: ", err.Error())
            return
        }
    }()
 
    localFile, _ := os.Create("1.png")
    defer localFile.Close()
    writer := bufio.NewWriter(localFile)
    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatal("ReadAll: ", err.Error())
        return
    }
    writer.Write(bytes)
    writer.Flush()
}


type Wechat_CreateUserInput struct {
  Unionid string
  Openid  string
  Headimgurl string
  Nickname string
  Language string
  Sex  int
  Province string
  City string
  Country string
}

type Wechat_CreateUserRet struct {
  Ret string
  Message string
  Data struct {
    Id string
    Token string
  }
}

//http.HandleFunc("/api/wechat_createuser", wechat_createuserHandle)   // 创建用户
func wechat_createuserHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_createuserHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_CreateUserRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_CreateUserInput
    json.Unmarshal([]byte(result), &input)

    // 判断Unionid是否存在
    querystr := fmt.Sprintf("select openid,user_no from YUECHE_USERS_BIND_WECHAT where openid = '%s'", input.Openid)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var openid string;
    var uid string
    for rows.Next() {
      err = rows.Scan(&openid, &uid)
      if err != nil {
        panic(err.Error())
        return
      }
    }

    if openid == input.Openid {
      retmessage.Ret = "0"
      retmessage.Message = "该微信号已注册"
      retmessage.Data.Id = uid
      b, _ := json.Marshal(retmessage)
      w.Header().Set("Access-Control-Allow-Origin","*")
      fmt.Fprintf(w, "%s", b)
      return   
    }

    //将数据写入 YUECHE_USERS_BIND_WECHAT
    token := GetRandomSalt()
    t := time.Now().Unix()
    user_no := fmt.Sprintf("WX_%d", t)

    stmtIns, err := db.Prepare("INSERT INTO YUECHE_USERS_BIND_WECHAT VALUES( ?,?,?,?,?,?,?,?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(user_no, input.Unionid, input.Openid, input.Headimgurl, input.Nickname, input.Language, input.Sex, input.Province, input.City, input.Country)
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    // 如果没有注册，则写入YUECHE_USER记录
    var phone string = ""
    var password string = ""
    var name string = ""
    var online_tot string = "0"
    var km_tot string = "0"
    var order_tot string = "0"

    stmtIns, err = db.Prepare("INSERT INTO YUECHE_USER VALUES( ?,?,?,?,?,?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(user_no, token, name, phone, password, online_tot, km_tot, order_tot)
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    // 首次注册，将优惠券写入用户优惠券表
    coupon_no := "CON_001"
    stmtIns, err = db.Prepare("INSERT INTO YUECHE_USER_COUPON VALUES( ?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(user_no, coupon_no, time.Now())
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    stmtIns, err = db.Prepare("INSERT INTO YUEHCE_TICKET_COUNT VALUES( ?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(user_no, 0)
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }


    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Id = user_no
    retmessage.Data.Token = token
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type Wechat_BindUserInput struct {
  User_no string
  Phone  string
}

type Wechat_BindUserRet struct {
  Ret string
  Message string
  Data struct {
    Id string
    Token string
  }
}

//http.HandleFunc("/api/wechat_binduser", wechat_binduserHandle)   // 绑定用户手机号
func wechat_binduserHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_binduserHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_BindUserRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_BindUserInput
    json.Unmarshal([]byte(result), &input)

    var phone string = ""
    querystr := fmt.Sprintf("select phone from YUECHE_USER where user_no = '%s'", input.User_no)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
      err = rows.Scan(&phone)
      if err != nil {
        panic(err.Error())
        return
      }
    }

    if phone == input.Phone {
      retmessage.Ret = "0"
      retmessage.Message = "该手机号已绑定"
      b, _ := json.Marshal(retmessage)
      w.Header().Set("Access-Control-Allow-Origin","*")
      fmt.Fprintf(w, "%s", b)
      return   
    }

    stmtIns, err := db.Prepare("update YUECHE_USER set phone=? where user_no=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(input.Phone, input.User_no)
    if err != nil {
      panic(err.Error())
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type Wechat_CreateDriverInput struct {
  Unionid string
  Openid  string
  Headimgurl string
  Nickname string
  Language string
  Sex  int
  Province string
  City string
  Country string
}

type Wechat_CreateDriverRet struct {
  Ret string
  Message string
  Data struct {
    Id string
    Token string
  }
}

//http.HandleFunc("/api/wechat_createdriver", wechat_createdriverHandle)   // 创建司机用户
func wechat_createdriverHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_createdriverHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_CreateDriverRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_CreateDriverInput
    json.Unmarshal([]byte(result), &input)

    // 判断Unionid是否存在
    querystr := fmt.Sprintf("select openid,driver_no from YUECHE_DRIVER_BIND_WECHAT where openid = '%s'", input.Openid)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var openid string;
    var uid string
    for rows.Next() {
      err = rows.Scan(&openid, &uid)
      if err != nil {
        panic(err.Error())
        return
      }
    }

    if openid == input.Openid {
      retmessage.Ret = "0"
      retmessage.Message = "该微信号已注册"
      retmessage.Data.Id = uid
      b, _ := json.Marshal(retmessage)
      w.Header().Set("Access-Control-Allow-Origin","*")
      fmt.Fprintf(w, "%s", b)
      return   
    }

    //将数据写入 YUECHE_DRIVER_BIND_WECHAT
    token := GetRandomSalt()
    t := time.Now().Unix()
    driver_no := fmt.Sprintf("WX_%d", t)

    stmtIns, err := db.Prepare("INSERT INTO YUECHE_DRIVER_BIND_WECHAT VALUES( ?,?,?,?,?,?,?,?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(driver_no, input.Unionid, input.Openid, input.Headimgurl, input.Nickname, input.Language, input.Sex, input.Province, input.City, input.Country)
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    // 如果没有注册，则写入YUECHE_DRIVER记录
    var phone string = ""
    var password string = ""
    var name string = ""
    var idcard string = ""
    var car_license string = ""
    var km_tot string = "0"
    var order_tot string = "0"
    var appraise string = "0"

    stmtIns, err = db.Prepare("INSERT INTO YUECHE_DRIVER VALUES( ?,?,?,?,?,?,?,?,?,?,?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(driver_no, token, name, phone, idcard, car_license, password, km_tot, order_tot, appraise,"0","0", "on")
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }


    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Id = driver_no
    retmessage.Data.Token = token
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

type Wechat_BindDriverInput struct {
  Driver_no string
  Name string
  Phone  string
  Idcard string
  Car_license string
}

type Wechat_BindDriverRet struct {
  Ret string
  Message string
  Data struct {
    Id string
    Token string
  }
}

//http.HandleFunc("/api/wechat_binddriver", wechat_binddriverHandle)   // 绑定司机用户信息
func wechat_binddriverHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_binddriverHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_BindDriverRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_BindDriverInput
    json.Unmarshal([]byte(result), &input)

    var phone string = ""
    querystr := fmt.Sprintf("select phone from YUECHE_DRIVER where driver_no = '%s'", input.Driver_no)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
      err = rows.Scan(&phone)
      if err != nil {
        panic(err.Error())
        return
      }
    }

    if phone == input.Phone {
      retmessage.Ret = "1"
      retmessage.Message = "该手机号已绑定"
      b, _ := json.Marshal(retmessage)
      w.Header().Set("Access-Control-Allow-Origin","*")
      fmt.Fprintf(w, "%s", b)
      return   
    }

    stmtIns, err := db.Prepare("update YUECHE_DRIVER set name=?,phone=?,idcard=? where driver_no=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(input.Name, input.Phone, input.Idcard, input.Driver_no)
    if err != nil {
      panic(err.Error())
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type Wechat_GetBannerRet struct {
  Ret string
  Message string
  Data struct {
    Banner []string
  }
}

//http.HandleFunc("/api/wechat_getbanner", wechat_getbannerHandle)   // 获取首页BANNER图
func wechat_getbannerHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getbannerHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetBannerRet

    var banner_pic string = ""
    var pic_slice []string 
    querystr := "select banner_pic from YUECHE_BANNER"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
      err = rows.Scan(&banner_pic)
      if err != nil {
        panic(err.Error())
        return
      }
      pic_slice = append(pic_slice, banner_pic)
    }

    // test
    //send_sms_order("18607127129", "武汉", "黄石", "2018-05-21 21:00:45")

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Banner = pic_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type CouponInfo struct {
  Coupon_no string
  Name string
  Line_no  string
  Tot string
  Residue_tot string
  Discount_money string
  Limit_money string
  Past_time string
  S_city string
  E_city string
}

type Wechat_GetCouponRet struct {
  Ret string
  Message string
  Data struct {
    Coupon []CouponInfo
  }
}

//http.HandleFunc("/api/wechat_getcoupon", wechat_getcouponHandle)   // 获取所有优惠券
func wechat_getcouponHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getcouponHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetCouponRet
    querystr := "select * from YUECHE_COUPON"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var cinfo CouponInfo
    var cinfo_slice []CouponInfo
    for rows.Next() {
      err = rows.Scan(&cinfo.Coupon_no, &cinfo.Name, &cinfo.Line_no, &cinfo.Tot, &cinfo.Residue_tot, &cinfo.Discount_money, &cinfo.Limit_money, &cinfo.Past_time, &cinfo.S_city, &cinfo.E_city)
      if err != nil {
        panic(err.Error())
        return
      }
      cinfo_slice = append(cinfo_slice, cinfo)
    }
    
    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Coupon = cinfo_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

type Wechat_SetCouponInput struct {
  User_no string
  Coupon_no string
}

type Wechat_SetCouponRet struct {
  Ret string
  Message string
}

//http.HandleFunc("/api/wechat_setusercoupon", wechat_setusercouponHandle)   // 添加优惠券到个人账户
func wechat_setusercouponHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_setusercouponHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_SetCouponRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_SetCouponInput
    json.Unmarshal([]byte(result), &input)

    var coupon_no string = ""
    querystr := fmt.Sprintf("select coupon_no from YUECHE_USER_COUPON where user_no = '%s'", input.User_no)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
      err = rows.Scan(&coupon_no)
      if err != nil {
        panic(err.Error())
        return
      }
    }

    if coupon_no == input.Coupon_no {
      retmessage.Ret = "1"
      retmessage.Message = "该优惠券已经领取"
      b, _ := json.Marshal(retmessage)
      w.Header().Set("Access-Control-Allow-Origin","*")
      fmt.Fprintf(w, "%s", b)
      return   
    }

    stmtIns, err := db.Prepare("INSERT INTO YUECHE_USER_COUPON VALUES(?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(input.User_no, input.Coupon_no, time.Now())
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type Wechat_GetUserCouponInput struct {
  User_no string
}

type CouponText struct {
  Money_sum string        // 优惠券金额
  Limit_txt string        // 例如：坐车就送 或者 满150元立减
  Coupon_type string      // 例如：现金立减
  Coupon_explain string   // 例如：适用于九州约车所有专线使用 或者 适用于从宜昌到武汉专线使用
  Coupon_date string // 例如：有效日期： 截止于2018-05-24
}

type Wechat_GetUserCouponRet struct {
  Ret string
  Message string
    Data struct {
    // Coupon []CouponInfo
      Coupon []CouponText
  }
}

//http.HandleFunc("/api/wechat_getusercoupon", wechat_getusercouponHandle)   // 获取个人账户已经领取的优惠券
func wechat_getusercouponHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getusercouponHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetUserCouponRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_GetUserCouponInput
    json.Unmarshal([]byte(result), &input)

    var coupon_no string = ""
    var ctext_slice []CouponText
    querystr := fmt.Sprintf("select coupon_no from YUECHE_USER_COUPON where user_no = '%s'", input.User_no)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
      err = rows.Scan(&coupon_no)
      if err != nil {
        panic(err.Error())
        return
      }

      querystr1 := fmt.Sprintf("select * from YUECHE_COUPON where coupon_no = '%s'", coupon_no)
      rows1,err1 := db.Query(querystr1)
      if err1 != nil {
        log.Fatal(err1)
        return
      }
      defer rows1.Close()

      var cinfo CouponInfo
      var ctext CouponText
      
      for rows1.Next() {
        err1 = rows1.Scan(&cinfo.Coupon_no, &cinfo.Name, &cinfo.Line_no, &cinfo.Tot, &cinfo.Residue_tot, &cinfo.Discount_money, &cinfo.Limit_money, &cinfo.Past_time, &cinfo.S_city, &cinfo.E_city)
        if err1 != nil {
          panic(err1.Error())
          return
        }

        // 根据信息格式化小程序需要显示的文本
        ctext.Money_sum = cinfo.Discount_money
        ctext.Coupon_type = "现金立减"

        if cinfo.Limit_money == "0" {
          ctext.Limit_txt = "坐车立减"
        } else {
          ctext.Limit_txt = "满"+cinfo.Limit_money+"元可用"
        }

        if cinfo.S_city == "通用" && cinfo.E_city == "通用" {
          ctext.Coupon_explain = "适用于九州约车所有专线使用"
        } else {
          ctext.Coupon_explain = "适用于从" + cinfo.S_city + "到" + cinfo.E_city + "专线使用"
        }
        
        s1 := cinfo.Past_time[0:10]
        s2 := cinfo.Past_time[11:19]
        s3 := s1 + " " + s2
        t1, _ := time.Parse("2006-01-02 15:04:05", s3)
        timeStr1 := t1.Format("2006-01-02")

        ctext.Coupon_date = "截止于"+timeStr1
        

        ctext_slice = append(ctext_slice, ctext)
      }
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    //retmessage.Data.Coupon = cinfo_slice
    retmessage.Data.Coupon = ctext_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type QaInfo struct {
  Qa_no string
  Qa_type string
  Question string
  Answer  string
}

type Wechat_GetQaRet struct {
  Ret string
  Message string
  Data struct {
    Qa []QaInfo
  }
}

//http.HandleFunc("/api/wechat_getqa", wechat_getqaHandle)   // 获取出行贴士
func wechat_getqaHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getqaHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetQaRet
    querystr := "select * from YUECHE_QA"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var cinfo QaInfo
    var cinfo_slice []QaInfo
    for rows.Next() {
      err = rows.Scan(&cinfo.Qa_no, &cinfo.Qa_type, &cinfo.Question, &cinfo.Answer)
      if err != nil {
        panic(err.Error())
        return
      }
      cinfo_slice = append(cinfo_slice, cinfo)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Qa = cinfo_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

func getBigLineNobycity(Line_no *string, s_city string, e_city string) {

  querystr := fmt.Sprintf("select line_no from YUECHE_BIG_BUS_LINE where s_city = '%s' and e_city = '%s'", s_city, e_city)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
  err = rows.Scan(Line_no)
    if err != nil {
      panic(err.Error())
      return
    }
  }
}

func getLineNobycity(Line_no *string, s_city string, e_city string) {

  querystr := fmt.Sprintf("select line_no from YUECHE_BUS_LINE where s_city = '%s' and e_city = '%s'", s_city, e_city)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
  err = rows.Scan(Line_no)
    if err != nil {
      panic(err.Error())
      return
    }
  }
}


func getLineinfobyid(Line_no string, s_city *string, e_city *string, price *string, km_tot *string) {

  querystr := fmt.Sprintf("select s_city,e_city,price,km_tot from YUECHE_BUS_LINE where line_no = '%s'", Line_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
  err = rows.Scan(s_city,e_city,price,km_tot)
    if err != nil {
      panic(err.Error())
      return
    }
  }
}

func getDriverinfobyid(driver_no string, name *string, phone *string, car_license *string) {

  querystr := fmt.Sprintf("select name,phone,car_license from YUECHE_DRIVER where driver_no = '%s'", driver_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
  err = rows.Scan(name,phone,car_license)
    if err != nil {
      panic(err.Error())
      return
    }
  }
}

func getUserinfobyid(user_no string, name *string, phone *string) {
  querystr := fmt.Sprintf("select name,phone from YUECHE_USER where user_no = '%s'", user_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
  err = rows.Scan(name,phone)
    if err != nil {
      panic(err.Error())
      return
    }
  }
}


type OrderInfo struct {
  User_name string
  User_phone string
  Order_no string
  Ticket_no string  // 车票号
  User_no string
  Line_no string
  S_city string   // 起点
  E_city string   // 终点
  Km_tot string   // 全程公里数
  Price_one string // 单张票价
  Amount string   // 购票数量
 // Ride_time_s string   // 乘车起始时间
 // Ride_time_e string   // 乘车终止时间
  Aboard_time string   // 上车时间
  Aboard_local_name string // 上车站点（途径点）
  End_local_name string // 下车地点
  Car_license string   // 车牌号
  Driver_no string     // 司机编号
  Driver_name string   // 司机姓名
  Driver_phone string  // 司机手机号
  Createtime string    // 订单创建时间
  Payprice string      // 支付金额
  Orderstatus string   // 已完成，未完成
  Evaluate string // 评价留言
}

type Wechat_GetUserOrderHInput struct {
  User_no string
  Orderstatus string
}

type Wechat_GetUserOrderRet struct {
  Ret string
  Message string
  Data struct {
    Order []OrderInfo
  }
}

//http.HandleFunc("/api/wechat_getuserorder", wechat_getuserorderHandle)   // 获取订单（已完成，未完成）
func wechat_getuserorderHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getuserorderHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetUserOrderRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_GetUserOrderHInput
    json.Unmarshal([]byte(result), &input)

    var cinfo_slice []OrderInfo

    // 取得符合条件的订单号
    querystr := fmt.Sprintf("select order_no,payprice from YUECHE_PAYMENTS where user_no = '%s' and orderstatus = '%s' order by order_no desc", input.User_no, input.Orderstatus)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var order_no string
    var payprice string
    for rows.Next() {
    err = rows.Scan(&order_no, &payprice)
      if err != nil {
        panic(err.Error())
        return
      }

      // 根据订单号取得对应的订单详细信息
      querystr1 := fmt.Sprintf("select order_no,ticket_no,user_no,line_no,amount,aboard_time,aboard_local_name,end_local_name,driver_no,createtime from YUECHE_ORDER where order_no = '%s'", order_no)
      rows1,err1 := db.Query(querystr1)
      if err1 != nil {
        log.Fatal(err1)
        return
      }
      defer rows1.Close()

      var cinfo OrderInfo
      
      for rows1.Next() {
        err1 = rows1.Scan(&cinfo.Order_no, &cinfo.Ticket_no, &cinfo.User_no, &cinfo.Line_no, &cinfo.Amount, &cinfo.Aboard_time, &cinfo.Aboard_local_name, &cinfo.End_local_name, &cinfo.Driver_no, &cinfo.Createtime)
        if err1 != nil {
          panic(err1.Error())
          return
        }

        getLineinfobyid(cinfo.Line_no, &cinfo.S_city, &cinfo.E_city, &cinfo.Price_one, &cinfo.Km_tot)
        getDriverinfobyid(cinfo.Driver_no, &cinfo.Driver_name, &cinfo.Driver_phone, &cinfo.Car_license)
        cinfo.Payprice = payprice
        cinfo.Orderstatus = input.Orderstatus

        var aboard_time string
        var start_city string
        var end_city string
        //var phone string
        var ticket_no string
        get_start_end_by_Order(cinfo.Order_no, &start_city, &end_city, &aboard_time, &cinfo.User_phone, &ticket_no)

        cinfo_slice = append(cinfo_slice, cinfo)
      }
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Order = cinfo_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type Wechat_GetLineInfoInput struct {
  Line_type string  // "common","air"
}

type Via_s_info struct {
  Via_name string
  Longitude string
  Latitude string
}

type LineInfo struct {
  Line_no string
  S_city string
  E_city string
  Price string
  Km_tot string
  Description string
  Line_via struct {
    Via_local_name []Via_s_info
  }
}

type Wechat_GetLineInfoRet struct {
  Ret string
  Message string
  Data struct {
    Line []LineInfo
  }
}

//http.HandleFunc("/api/wechat_getlineinfo", wechat_getlineinfoHandle)   // 获取车辆路线信息
func wechat_getlineinfoHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getlineinfoHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetLineInfoRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_GetLineInfoInput
    json.Unmarshal([]byte(result), &input)

    querystr := "select line_no,s_city,e_city,price,km_tot,description from YUECHE_BUS_LINE"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var cinfo LineInfo
    var cinfo_slice []LineInfo
    for rows.Next() {
    err = rows.Scan(&cinfo.Line_no, &cinfo.S_city, &cinfo.E_city, &cinfo.Price, &cinfo.Km_tot, &cinfo.Description)
      if err != nil {
        panic(err.Error())
        return
      }

/*
      if input.Line_type == "common" {
        if cinfo.Line_no != "LINE_001" {
          continue
        }
        
      } else if input.Line_type == "air" {
        if cinfo.Line_no != "LINE_010" {
          continue
        }
      } else {

      }
*/
      querystr1 := fmt.Sprintf("select via_local_name,longitude,latitude from YUECHE_BUS_LINE_VIA where line_type = '%s' and local_status='on'", input.Line_type)
      rows1,err1 := db.Query(querystr1)
      if err1 != nil {
        log.Fatal(err1)
        return
      }
      defer rows1.Close()

      var via_local_in Via_s_info
      // 这里有个BUG，所以暂时屏蔽
      //cinfo.Line_via.Via_local_name = append(cinfo.Line_via.Via_local_name[:0],cinfo.Line_via.Via_local_name[len(cinfo.Line_via.Via_local_name):]...)  // 将slice清空
      for rows1.Next() {
        err1 = rows1.Scan(&via_local_in.Via_name,&via_local_in.Longitude,&via_local_in.Latitude)
        if err1 != nil {
          panic(err1.Error())
          return
        }
        cinfo.Line_via.Via_local_name = append(cinfo.Line_via.Via_local_name, via_local_in)
      }
      cinfo_slice = append(cinfo_slice, cinfo)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Line = cinfo_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

type BigLineInfo struct {
  Line_no string
  S_city string
  E_city string
  Price string
  Description string
  Bus_time []string
}

type Wechat_GetBigLineInfoRet struct {
  Ret string
  Message string
  Data struct {
    Line []BigLineInfo
  }
}
 // 获取长阳大巴车辆路线信息
// http.HandleFunc("/api/wechat_gebuslineinfo", wechat_getbuslineinfoHandle)  
func wechat_getbuslineinfoHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getbuslineinfoHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetBigLineInfoRet
    //result, _:= ioutil.ReadAll(r.Body)

    querystr := "select line_no,s_city,e_city,price,description from YUECHE_BIG_BUS_LINE"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var cinfo_slice []BigLineInfo
    for rows.Next() {
      var cinfo BigLineInfo

      err = rows.Scan(&cinfo.Line_no, &cinfo.S_city, &cinfo.E_city, &cinfo.Price, &cinfo.Description)
      if err != nil {
        panic(err.Error())
        return
      }

      var bus_time string
      querystr1 := fmt.Sprintf("select bus_time from YUECHE_BIG_BUS_TIME where line_no = '%s'", cinfo.Line_no)
      rows1,err1 := db.Query(querystr1)
      if err1 != nil {
        log.Fatal(err1)
        return
      }
      defer rows1.Close()

      for rows1.Next() {
        err1 = rows1.Scan(&bus_time)
        if err1 != nil {
          panic(err1.Error())
          return
        }
        cinfo.Bus_time = append(cinfo.Bus_time, bus_time)
      }

      cinfo_slice = append(cinfo_slice, cinfo)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Line = cinfo_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

func get_usable_coupon_by_user(User_no string, Price string, Coupon_no *string, Discount_money *string) {

  var coupon_no string = ""
  querystr := fmt.Sprintf("select coupon_no from YUECHE_USER_COUPON where user_no = '%s'", User_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(&coupon_no)
    if err != nil {
      panic(err.Error())
      return
    }

    timestr := time.Now()
    querystr1 := fmt.Sprintf("select * from YUECHE_COUPON where coupon_no = '%s' and past_time > '%s'", coupon_no, timestr)
    rows1,err1 := db.Query(querystr1)
    if err1 != nil {
      log.Fatal(err1)
      return
    }
    defer rows1.Close()

    var cinfo CouponInfo
    for rows1.Next() {
      err1 = rows1.Scan(&cinfo.Coupon_no, &cinfo.Name, &cinfo.Line_no, &cinfo.Tot, &cinfo.Residue_tot, &cinfo.Discount_money, &cinfo.Limit_money, &cinfo.Past_time, &cinfo.S_city, &cinfo.E_city)
      if err1 != nil {
        panic(err1.Error())
        return
      }

      num_price, _ := strconv.Atoi(Price)
      num_Discount_money, _ := strconv.Atoi(cinfo.Discount_money)
      if num_price < num_Discount_money {
        continue
      } else {
        *Coupon_no = cinfo.Coupon_no
        *Discount_money = cinfo.Discount_money
      }

      // 如果后续添加针对某种特定线路的优惠券，这里还需要添加线路判断的代码
      // .......



    }
  }
}

type Wechat_CreateSoloOrderInput struct {
  User_no string
  S_city string
  E_city string
  Aboard_time string
  Aboard_local_name string
  End_local_name string
  Amount string
}

type Wechat_CreateSoloOrderRet struct {
  Ret string
  Message string
  Data struct {
    Order_no string
    Original_price string
    Payprice string
    S_city string
    E_city string
    Km_tot string
    Aboard_time string
    Aboard_local_name string
    End_local_name string
    Amount string
    Coupon_no string
    Discount_money string
  }
}

func is_user_ok(User_no string) string{
    var rcount string = "0"
    querystr := fmt.Sprintf("SELECT COUNT(*) FROM YUECHE_USER where user_no = '%s'", User_no)
 //   fmt.Println(string(querystr))
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return "error"
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&rcount)
      if err != nil {
        panic(err.Error())
        return "error"
      }
    }

    if rcount == "0" {
      return "error"
    } else {
      return "ok"
    }
}

//http.HandleFunc("/api/wechat_createsoloorder", wechat_createsoloorderHandle)   // 创建购票订单，这里的乘车起始时间务必为车票时间，否则验票无法通过
func wechat_createsoloorderHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_createsoloorderHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_CreateSoloOrderRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_CreateSoloOrderInput
    json.Unmarshal([]byte(result), &input)

    var line_no string
    getLineNobycity(&line_no, input.S_city, input.E_city)
   // fmt.Println(string(input.S_city))
   // fmt.Println(string(input.E_city))
   // fmt.Println(string(line_no))

    isok := is_user_ok(input.User_no)
    if isok == "error" {
      retmessage.Ret = "1"
      retmessage.Message = "您的账户存在异常，请重新安装小程序并通过验证！"
      b, _ := json.Marshal(retmessage)
      w.Header().Set("Access-Control-Allow-Origin","*")
      fmt.Fprintf(w, "%s", b)
      return   
    }

    t := time.Now().Unix()
    order_no := fmt.Sprintf("ORDER_%d", t)

    //ticket_no := fmt.Sprintf("%d", t)
    ticket_no := GetRandomNum6Str()

    // 避免同步造成的订单号相同
    order_no = order_no + ticket_no

    stmtIns, err := db.Prepare("INSERT INTO YUECHE_ORDER VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    t1, _ := time.Parse("2006-01-02 15:04:05", input.Aboard_time)
    timeStr1 := t1.Format("2006-01-02")
    //fmt.Println(string(timeStr))
    timeStr2 := timeStr1 + " 23:59:59"
    //fmt.Println(string(timeStr))


    _, err = stmtIns.Exec(order_no, ticket_no, input.User_no, line_no, input.Amount, timeStr1, timeStr2, input.Aboard_time, input.Aboard_local_name, input.End_local_name,"","",time.Now(), "未付款")
    //_, err = stmtIns.Exec(multi_no, input.User_no, input.Line_no, "", timeStr, "确认中", "")
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    // 统计用户需要支付的钱数，以分为单位
    var price string
    var km_tot string
    querystr := fmt.Sprintf("select price,km_tot from YUECHE_BUS_LINE where line_no = '%s'", line_no)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&price,&km_tot)
      if err != nil {
        panic(err.Error())
        return
      }
    }

    num, _ := strconv.Atoi(price)
    amount, _ := strconv.Atoi(input.Amount)
    lastprice := num * amount


    // 获取用户的优惠券，并将可以使用的优惠券金额从付款总金额中去掉。
    var Coupon_no string = "0"
    var Discount_money string = "0"
    get_usable_coupon_by_user(input.User_no, strconv.Itoa(int(lastprice)), &Coupon_no, &Discount_money)
    num_Discount_money, _ := strconv.Atoi(Discount_money)

    retmessage.Data.Coupon_no = Coupon_no
    retmessage.Data.Discount_money = Discount_money


    // 同时将使用过的优惠券从用户优惠券列表中去掉 - 这个操作放在微信回调中更合适
    retmessage.Data.Original_price = strconv.Itoa(int(lastprice))
    if lastprice - num_Discount_money < 0 {
      retmessage.Data.Payprice = "0"
    } else {
      retmessage.Data.Payprice = strconv.Itoa(int(lastprice - num_Discount_money))  
    }
    
    fmt.Println(string(retmessage.Data.Payprice))


    order_type := "车票付款"
    coupon_no := Coupon_no  
    payid := ""
    tradenum := ""
    payprice := retmessage.Data.Payprice
    paystatus := "none"
    orderstatus := "待付款"  // 0：待付款，1：未完成，3:已完成，4：已退票，5：已取消，6：已删除
    payway := "weico"   // 到付（cash），支付宝（ali），微信（weico），银行卡（unionpay）

    stmtIns, err = db.Prepare("INSERT INTO YUECHE_PAYMENTS VALUES( ?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(input.User_no, order_no, order_type, payid, tradenum, payprice, coupon_no , payprice, paystatus, time.Now(), orderstatus, "无发票", "空", payway)
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }


    // 获取当天最早时间
    //timeStr := time.Now().Format("2006-01-02")
    // 获取当天最晚时间
    //timeStr = timeStr + " 23:59:59"
    // 将当天最早和最晚时间分别写入YUECHE_ORDER中的ride_time_s和ride_time_e

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Order_no = order_no
    retmessage.Data.S_city = input.S_city
    retmessage.Data.E_city = input.E_city
    retmessage.Data.Km_tot = km_tot
    retmessage.Data.Aboard_time = input.Aboard_time
    retmessage.Data.Aboard_local_name = input.Aboard_local_name
    retmessage.Data.End_local_name = input.End_local_name
    retmessage.Data.Amount = input.Amount
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

type Wechat_CreateBigOrderInput struct {
  User_no string
  S_city string
  E_city string
  Aboard_time string
//  Aboard_local_name string
//  End_local_name string
  Amount string
}

type Wechat_CreateBigOrderRet struct {
  Ret string
  Message string
  Data struct {
    Order_no string
    Original_price string
    Payprice string
    S_city string
    E_city string
    Km_tot string
    Aboard_time string
    Aboard_local_name string
    End_local_name string
    Amount string
    Coupon_no string
    Discount_money string
  }
}


// 创建大巴车购票订单
// http.HandleFunc("/api/wechat_createbigorder", wechat_createbigorderHandle)   
func wechat_createbigorderHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_createbigorderHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_CreateBigOrderRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_CreateBigOrderInput
    json.Unmarshal([]byte(result), &input)

    var line_no string
    getBigLineNobycity(&line_no, input.S_city, input.E_city)
   // fmt.Println(string(input.S_city))
   // fmt.Println(string(input.E_city))
   // fmt.Println(string(line_no))

    isok := is_user_ok(input.User_no)
    if isok == "error" {
      retmessage.Ret = "1"
      retmessage.Message = "您的账户存在异常，请重新安装小程序并通过验证！"
      b, _ := json.Marshal(retmessage)
      w.Header().Set("Access-Control-Allow-Origin","*")
      fmt.Fprintf(w, "%s", b)
      return   
    }

    t := time.Now().Unix()
    order_no := fmt.Sprintf("ORDER_B_%d", t)

    //ticket_no := fmt.Sprintf("%d", t)
    ticket_no := GetRandomNum6Str()

    // 避免同步造成的订单号相同
    order_no = order_no + ticket_no

    stmtIns, err := db.Prepare("INSERT INTO YUECHE_ORDER VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    t1, _ := time.Parse("2006-01-02 15:04:05", input.Aboard_time)
    timeStr1 := t1.Format("2006-01-02")
    //fmt.Println(string(timeStr))
    timeStr2 := timeStr1 + " 23:59:59"
    //fmt.Println(string(timeStr))

    Aboard_local_name := input.S_city
    End_local_name := input.E_city

    _, err = stmtIns.Exec(order_no, ticket_no, input.User_no, line_no, input.Amount, timeStr1, timeStr2, input.Aboard_time, Aboard_local_name, End_local_name,"","",time.Now(), "未付款")
    //_, err = stmtIns.Exec(multi_no, input.User_no, input.Line_no, "", timeStr, "确认中", "")
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    // 统计用户需要支付的钱数，以分为单位
    var price string
    var km_tot string = "0"
    querystr := fmt.Sprintf("select price from YUECHE_BIG_BUS_LINE where line_no = '%s'", line_no)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&price)
      if err != nil {
        panic(err.Error())
        return
      }
    }

    num, _ := strconv.Atoi(price)
    amount, _ := strconv.Atoi(input.Amount)
    lastprice := num * amount


    // 获取用户的优惠券，并将可以使用的优惠券金额从付款总金额中去掉。
    var Coupon_no string = "0"
    var Discount_money string = "0"
    get_usable_coupon_by_user(input.User_no, strconv.Itoa(int(lastprice)), &Coupon_no, &Discount_money)
    num_Discount_money, _ := strconv.Atoi(Discount_money)

    retmessage.Data.Coupon_no = Coupon_no
    retmessage.Data.Discount_money = Discount_money


    // 同时将使用过的优惠券从用户优惠券列表中去掉 - 这个操作放在微信回调中更合适
    retmessage.Data.Original_price = strconv.Itoa(int(lastprice))
    if lastprice - num_Discount_money < 0 {
      retmessage.Data.Payprice = "0"
    } else {
      retmessage.Data.Payprice = strconv.Itoa(int(lastprice - num_Discount_money))  
    }
    
    fmt.Println(string(retmessage.Data.Payprice))


    order_type := "车票付款"
    coupon_no := Coupon_no  
    payid := ""
    tradenum := ""
    payprice := retmessage.Data.Payprice
    paystatus := "none"
    orderstatus := "待付款"  // 0：待付款，1：未完成，3:已完成，4：已退票，5：已取消，6：已删除
    payway := "weico"   // 到付（cash），支付宝（ali），微信（weico），银行卡（unionpay）

    stmtIns, err = db.Prepare("INSERT INTO YUECHE_PAYMENTS VALUES( ?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(input.User_no, order_no, order_type, payid, tradenum, payprice, coupon_no , payprice, paystatus, time.Now(), orderstatus, "无发票", "空", payway)
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }


    // 获取当天最早时间
    //timeStr := time.Now().Format("2006-01-02")
    // 获取当天最晚时间
    //timeStr = timeStr + " 23:59:59"
    // 将当天最早和最晚时间分别写入YUECHE_ORDER中的ride_time_s和ride_time_e

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Order_no = order_no
    retmessage.Data.S_city = input.S_city
    retmessage.Data.E_city = input.E_city
    retmessage.Data.Km_tot = km_tot
    retmessage.Data.Aboard_time = input.Aboard_time
    retmessage.Data.Aboard_local_name = Aboard_local_name
    retmessage.Data.End_local_name = End_local_name
    retmessage.Data.Amount = input.Amount
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}



type Wechat_DelSoloOrderInput struct {
  User_no string
  Order_no string
}

type Wechat_DelSoloOrderRet struct {
  Ret string
  Message string
}
//http.HandleFunc("/api/wechat_delsoloorder", wechat_delsoloorderHandle)   // 删除购票订单
func wechat_delsoloorderHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_delsoloorderHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_DelSoloOrderRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_DelSoloOrderInput
    json.Unmarshal([]byte(result), &input)

    stmtIns, err := db.Prepare("update YUECHE_ORDER set description=? where order_no=? and user_no =?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec("已取消", input.Order_no, input.User_no)
    if err != nil {
      panic(err.Error())
    }

    stmtIns, err = db.Prepare("update YUECHE_PAYMENTS set orderstatus=? where order_no=? and user_no =?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec("已取消", input.Order_no, input.User_no)
    if err != nil {
      panic(err.Error())
    }
/*
    stmtIns, err = db.Prepare("DELETE FROM YUECHE_PAYMENTS WHERE user_no =? and order_no =?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(input.User_no, input.Order_no)
    if err != nil {
      panic(err.Error())
    }
    */
    // 获取当天最早时间
    //timeStr := time.Now().Format("2006-01-02")
    // 获取当天最晚时间
    //timeStr = timeStr + " 23:59:59"
    // 将当天最早和最晚时间分别写入YUECHE_ORDER中的ride_time_s和ride_time_e

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type Wechat_CreateMultiOrderInput struct {
  User_no string
  S_city string
  E_city string
  Go_time string
  Travel_model string
}

type Wechat_CreateMultiOrderRet struct {
  Ret string
  Message string
  Data struct {
    Price string
    Order_no string
  }
}

//http.HandleFunc("/api/wechat_createmultiorder", wechat_createmultiorderHandle)   // 创建包车订单
func wechat_createmultiorderHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_createmultiorderHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_CreateMultiOrderRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_CreateMultiOrderInput
    json.Unmarshal([]byte(result), &input)

    //var line_no string
    //getLineNobycity(&line_no, input.S_city, input.E_city)

    t := time.Now().Unix()
    multi_no := fmt.Sprintf("MULTI_%d", t)

    var price string = "0"

    if input.Travel_model == "同城" {
      price = "10000"
    } else if input.Travel_model == "省内" {
      price = "30000"
    } else if input.Travel_model == "跨省" {
      price = "50000"
    } else {
      retmessage.Ret = "1"
      retmessage.Message = "出行模式参数错误"
      b, _ := json.Marshal(retmessage)
      w.Header().Set("Access-Control-Allow-Origin","*")
      fmt.Fprintf(w, "%s", b)
      return   
    }

    stmtIns, err := db.Prepare("INSERT INTO YUECHE_MULTI VALUES(?,?,?,?,?,?,?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(multi_no, input.User_no, input.S_city, "", input.Go_time, time.Now(), "确认中", input.Travel_model, price)
    //_, err = stmtIns.Exec(multi_no, input.User_no, input.Line_no, "", timeStr, "确认中", "")
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    order_type := "包车订金"
    payprice := price
    coupon_no := ""  
    payid := ""
    tradenum := ""
    //payprice := retmessage.Data.Payprice
    paystatus := "none"
    orderstatus := "待付款"  // 0：待付款，1：待发货，3:待收货，4：退换货，5：已取消，6：已删除
    payway := "weico"   // 到付（cash），支付宝（ali），微信（weico），银行卡（unionpay）

    stmtIns, err = db.Prepare("INSERT INTO YUECHE_PAYMENTS VALUES( ?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(input.User_no, multi_no, order_type, payid, tradenum, payprice, coupon_no , payprice, paystatus, time.Now(), orderstatus, "无发票", "空", payway)
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Price = price
    retmessage.Data.Order_no = multi_no
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type Wechat_SearchLineInput struct {
  User_no string
  Order_type string
}

type Wechat_SearchLineRet struct {
  Ret string
  Message string
}

//http.HandleFunc("/api/wechat_searchline", wechat_searchlineHandle)   // 车次查询
func wechat_searchlineHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_searchlineHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_SearchLineRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_SearchLineInput
    json.Unmarshal([]byte(result), &input)

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

type Wechat_GetLineStartEndRet struct {
  Ret string
  Message string
  Data struct {
    Start_City []string
    End_City []string
  }
}

//http.HandleFunc("/api/wechat_getlinestartend", wechat_getlinestartendHandle)   // 获取出所有发地和目的地
func wechat_getlinestartendHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getlinestartendHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetLineStartEndRet

    var s_city string = ""
    var s_city_slice []string 
    querystr := "select distinct s_city from YUECHE_BUS_LINE"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
      err = rows.Scan(&s_city)
      if err != nil {
        panic(err.Error())
        return
      }

      if s_city == "武汉机场" {
        continue
      }
      s_city_slice = append(s_city_slice, s_city)
    }

    var e_city string = ""
    var e_city_slice []string 
    querystr1 := "select distinct e_city from YUECHE_BUS_LINE"
    rows1,err1 := db.Query(querystr1)
    if err1 != nil {
      log.Fatal(err1)
      return
    }
    defer rows1.Close()

    for rows1.Next() {
      err1 = rows1.Scan(&e_city)
      if err1 != nil {
        panic(err1.Error())
        return
      }
      if e_city == "武汉机场" {
        continue
      }
      e_city_slice = append(e_city_slice, e_city)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Start_City = s_city_slice
    retmessage.Data.End_City = e_city_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type Wechat_GetLineEndbyStartInput struct {
  Start_City string
}

type Wechat_GetLineEndbyStartRet struct {
  Ret string
  Message string
  Data struct {
    End_City []string
  }
}

// http.HandleFunc("/api/wechat_getlineendbystart", wechat_getlineendbystartHandle)   // 根据起点城市，获取终点城市列表
func wechat_getlineendbystartHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getlineendbystartHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetLineEndbyStartRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_GetLineEndbyStartInput
    json.Unmarshal([]byte(result), &input)

    var e_city string = ""
    var e_city_slice []string 
    querystr := fmt.Sprintf("select e_city from YUECHE_BUS_LINE where s_city = '%s'", input.Start_City)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
      err = rows.Scan(&e_city)
      if err != nil {
        panic(err.Error())
        return
      }
      if e_city == "武汉机场" {
        continue;
      }

      e_city_slice = append(e_city_slice, e_city)
    }


    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.End_City = e_city_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}




type Wechat_GetAirRet struct {
  Ret string
  Message string
  Data struct {
    Start_City []string
    End_City []string
  }
}

// http.HandleFunc("/api/wechat_getairstartend", wechat_getairstartendHandle)   // 获取机场专线出所有发地和目的地
func wechat_getairstartendHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getairstartendHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetAirRet
    var s_city_slice []string
    var e_city_slice []string

    var s1 string = "武汉机场"
    var s2 string = "宜昌"
    var s3 string = "宜都"
    var s4 string = "枝江"

    s_city_slice = append(s_city_slice, s1)
    s_city_slice = append(s_city_slice, s2)
    s_city_slice = append(s_city_slice, s3)
    s_city_slice = append(s_city_slice, s4)

    e_city_slice = append(e_city_slice, s1)
    e_city_slice = append(e_city_slice, s2)
    e_city_slice = append(e_city_slice, s3)
    e_city_slice = append(e_city_slice, s4)


    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Start_City = s_city_slice
    retmessage.Data.End_City = e_city_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type Wechat_GetAirLineEndbyStartInput struct {
  Start_City string
}

type Wechat_GetAirLineEndbyStartRet struct {
  Ret string
  Message string
  Data struct {
    End_City []string
  }
}

// http.HandleFunc("/api/wechat_getairlineendbystart", wechat_getairlineendbystartHandle)   // 根据起点城市，获取机场专线终点城市列表
func wechat_getairlineendbystartHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getairlineendbystartHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetAirLineEndbyStartRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_GetAirLineEndbyStartInput
    json.Unmarshal([]byte(result), &input)

    var e_city string = ""
    var e_city_slice []string 

    if input.Start_City == "枝江" || input.Start_City == "宜昌" || input.Start_City == "宜都" {
      e_city = "武汉机场"
      e_city_slice = append(e_city_slice, e_city)
    } else if input.Start_City == "武汉机场" {
      e1 := "枝江"
      e2 := "宜昌"
      e3 := "宜都"

      e_city_slice = append(e_city_slice, e1)
      e_city_slice = append(e_city_slice, e2)
      e_city_slice = append(e_city_slice, e3)
    } else {
      retmessage.Ret = "1"
      retmessage.Message = "start city is failure"
      b, _ := json.Marshal(retmessage)
      w.Header().Set("Access-Control-Allow-Origin","*")
      fmt.Fprintf(w, "%s", b)
      return  
    }


    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.End_City = e_city_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type UserBillInfo struct {
  Payprice string
  Createtime string
}

type Wechat_GetSelfBillInput struct {
  User_no string
}

type Wechat_GetSelfBillRet struct {
  Ret string
  Message string
    Data struct {
    BillInfo []UserBillInfo
  }
}

//http.HandleFunc("/api/wechat_getselfbill", wechat_getselfbillHandle)   // 获取个人账单
func wechat_getselfbillHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getselfbillHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetSelfBillRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_GetSelfBillInput
    json.Unmarshal([]byte(result), &input)

    var billinfo UserBillInfo
    var billinfo_slice []UserBillInfo
    querystr := fmt.Sprintf("select payprice,createtime from YUECHE_PAYMENTS where user_no = '%s' and paystatus='success'", input.User_no)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&billinfo.Payprice, &billinfo.Createtime)
      if err != nil {
        panic(err.Error())
        return
      }

      billinfo_slice = append(billinfo_slice, billinfo)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.BillInfo = billinfo_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


func getLineMilebyid(Line_no string, km_tot *string) {

  querystr := fmt.Sprintf("select km_tot from YUECHE_BUS_LINE where line_no = '%s'", Line_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
  err = rows.Scan(km_tot)
    if err != nil {
      panic(err.Error())
      return
    }
  }
}

func milecountbyuserno(user_no string, milecount *string) {
    var mileint int = 0
    var km_tot string
    var line_no string
    querystr := fmt.Sprintf("SELECT line_no FROM YUECHE_ORDER where user_no = '%s'", user_no)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&line_no)
      if err != nil {
        panic(err.Error())
        return
      }

      getLineMilebyid(line_no, &km_tot)
      num, _ := strconv.Atoi(km_tot)
      mileint = num + mileint
    }

    *milecount = strconv.Itoa(int(mileint))
}


type SelfInfo struct {
  Name string
  Phone string
  OnlineDay string
  MileCount string
  OrderCount string
  Unionid string
  Openid  string
  Headimgurl string
  Nickname string
}

type Wechat_GetSelfInfoInput struct {
  User_no string
}

type Wechat_GetSelfInfoRet struct {
  Ret string
  Message string
  Data struct {
    SelfInfoMess SelfInfo
  }
}

//http.HandleFunc("/api/wechat_getselfinfo", wechat_getselfinfoHandle)   // 获取个人出行信息
func wechat_getselfinfoHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getselfinfoHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetSelfInfoRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_GetSelfInfoInput
    json.Unmarshal([]byte(result), &input)

    var rcount string = "0"
    var self_info SelfInfo
    querystr := fmt.Sprintf("SELECT COUNT(*) FROM YUECHE_ORDER where user_no = '%s'", input.User_no)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&rcount)
      if err != nil {
        panic(err.Error())
        return
      }
    }
    self_info.OnlineDay = rcount
    self_info.OrderCount = rcount

    milecountbyuserno(input.User_no, &self_info.MileCount)


    querystr1 := fmt.Sprintf("SELECT name,phone FROM YUECHE_USER where user_no = '%s'", input.User_no)
    rows1,err1 := db.Query(querystr1)
    if err1 != nil {
      log.Fatal(err1)
      return
    }
    defer rows1.Close()

    for rows1.Next() {
    err1 = rows1.Scan(&self_info.Name, &self_info.Phone)
      if err1 != nil {
        panic(err1.Error())
        return
      }
    }

    querystr2 := fmt.Sprintf("select unionid,openid,headimgurl,nickname from YUECHE_USERS_BIND_WECHAT where user_no = '%s'", input.User_no)
    rows2,err2 := db.Query(querystr2)
    if err2 != nil {
      log.Fatal(err2)
      return
    }
    defer rows2.Close()

    for rows2.Next() {
      err2 = rows2.Scan(&self_info.Unionid, &self_info.Openid, &self_info.Headimgurl, &self_info.Nickname)
      if err2 != nil {
        panic(err2.Error())
        return
      }
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.SelfInfoMess = self_info
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}



type Wechat_TicketCheckInput struct {
  User_no string
  Car_license string
 // Driver_no string
  Order_no string
}

type Wechat_TicketCheckRet struct {
  Ret string
  Message string
}

func  get_driver_by_carlicense(Car_license string, Driver_no *string) {
  querystr := fmt.Sprintf("select driver_no from YUECHE_DRIVER where car_license = '%s'", Car_license)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(Driver_no)
    if err != nil {
      panic(err.Error())
      return
    }
  }
}

func get_residueticket_by_orderno(Order_no string, Residue *int) {

  var amount string = "0"
  querystr := fmt.Sprintf("select amount from YUECHE_ORDER where order_no = '%s'", Order_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(&amount)
    if err != nil {
      panic(err.Error())
      return
    }
  }
  num_amount, _ := strconv.Atoi(amount)

  // 统计已经上车的人数
  var rcount string = "0"
  querystr = fmt.Sprintf("SELECT COUNT(*) FROM YUECHE_GET_ON_CAR where order_no = '%s'", Order_no)
  rows,err = db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
  err = rows.Scan(&rcount)
    if err != nil {
      panic(err.Error())
      return
    }
  }
  num_on, _ := strconv.Atoi(rcount)

  Residue_i := num_amount - num_on

  if Residue_i > 0 {
    *Residue = Residue_i
  } else {
    *Residue = 0
  }
}

// select * from YUECHE_ORDER where ride_time_s >'2018-02-21 00:00:00' and ride_time_s <'2018-02-21 17:00:00'
//http.HandleFunc("/api/wechat_ticketcheck", wechat_ticketcheckHandle)  // 扫码后，将司机信息，个人用户信息绑定发送给服务器端，验证是否购买了该司机的车票
func wechat_ticketcheckHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_ticketcheckHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_TicketCheckRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_TicketCheckInput
    json.Unmarshal([]byte(result), &input)

    retmessage.Ret = "1"
    retmessage.Message = "failure"
    timeStr := time.Now()
    var rcount string = "0"
    querystr := fmt.Sprintf("SELECT COUNT(*) FROM YUECHE_ORDER where user_no = '%s' and order_no = '%s' and description <> '已取消' and description <> '已上车' and ride_time_s < '%s' and ride_time_e > '%s'", input.User_no, input.Order_no, timeStr, timeStr)
    fmt.Println(querystr)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&rcount)
      if err != nil {
        panic(err.Error())
        return
      }
    }

    if rcount != "0" {
      retmessage.Ret = "0"
      retmessage.Message = "success"
    }

    // 根据当前车牌号获取司机信息
    var Driver_no string
    get_driver_by_carlicense(input.Car_license, &Driver_no)

    // 将当前司机号写入到当前订单中
    stmtIns, err := db.Prepare("update YUECHE_ORDER set driver_no=?,description=? where order_no=? and description <> '已上车'")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(Driver_no, "已上车",input.Order_no)
    if err != nil {
      panic(err.Error())
    }


    // 将订单表的状态修改为“已完成”
    stmtIns, err = db.Prepare("update YUECHE_PAYMENTS set orderstatus=? where order_no=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec("已完成",input.Order_no)
    if err != nil {
      panic(err.Error())
    }

    // 将验票记录插入YUECHE_GET_ON_CAR中
    var Residue int = 0
    var Longitude string = "0"
    var Latitude string = "0"
    var car_license string = ""
    get_carlicense_by_driver(Driver_no, &car_license)
    get_residueticket_by_orderno(input.Order_no, &Residue)
    for j := Residue; j > 0; j-- {

      stmtIns, err := db.Prepare("INSERT INTO YUECHE_GET_ON_CAR VALUES( ?,?,?,?,?,?)")
      if err != nil {
          panic(err.Error()) // proper error handling instead of panic in your app
      }
      defer stmtIns.Close() 

      _, err = stmtIns.Exec(Driver_no, car_license, input.Order_no, time.Now(), Longitude, Latitude)
      if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
      }
    }

    //retmessage.Ret = "0"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type Wechat_TodayTicketInfo struct {
  User_no string
  Order_no string
  Ticket_no string
  Line_no string
  Amount string
  Aboard_local_name string
  End_local_name string
  Aboard_time string
}

type Wechat_TodayTicketGetInput struct {
  User_no string
}

type Wechat_TodayTicketGetRet struct {
  Ret string
  Message string
  Data struct {
    TodayTicket []Wechat_TodayTicketInfo
  }
}
//http.HandleFunc("/api/wechat_todayticketget", wechat_todayticketgetHandle)  // 查询个人当天未完成订单，并返回订单信息，-测试通过
func wechat_todayticketgetHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_todayticketgetHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_TodayTicketGetRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_TodayTicketGetInput
    json.Unmarshal([]byte(result), &input)

    timeStr := time.Now()
 //   var user_no string = "0"
 //   var line_no string = "0"
 //   var amount string = "0"
 //   var aboard_local_name string = "0"
 //   var aboard_time string = "0"
    querystr := fmt.Sprintf("SELECT user_no,order_no,ticket_no, line_no,amount,aboard_local_name,end_local_name,aboard_time FROM YUECHE_ORDER where user_no = '%s' and ride_time_s < '%s' and ride_time_e > '%s' and description <> '已取消' and description <> '已上车' and description <> '未付款' order by order_no desc", input.User_no, timeStr, timeStr)
    //fmt.Println(querystr)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var ticket Wechat_TodayTicketInfo
    var ticket_slice []Wechat_TodayTicketInfo
    for rows.Next() {
    err = rows.Scan(&ticket.User_no,&ticket.Order_no,&ticket.Ticket_no,&ticket.Line_no,&ticket.Amount,&ticket.Aboard_local_name,&ticket.End_local_name,&ticket.Aboard_time)
      if err != nil {
        panic(err.Error())
        return
      }
      ticket_slice = append(ticket_slice, ticket)

      ss := ticket.Order_no[0:5]
      fmt.Println(ss)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.TodayTicket = ticket_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


type Wechat_DriverTicketInfo struct {
  User_no string
  Order_no string
  Line_no string
  S_city string
  E_city string
  Price_one string
  Price_amount string
  Km_tot string
  Amount string
  Aboard_local_name string
  End_local_name string
  Aboard_time string
}

type Wechat_DriverGetTicketInput struct {
  Driver_no string
}

type Wechat_DriverGetTicketRet struct {
  Ret string
  Message string
  Data struct {
    Ticket []Wechat_DriverTicketInfo
  }
}
//http.HandleFunc("/api/wechat_drivergetticket", wechat_drivergetticketHandle)  // 查询司机个人点票数，返回点票信息
func wechat_drivergetticketHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_drivergetticketHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_DriverGetTicketRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_DriverGetTicketInput
    json.Unmarshal([]byte(result), &input)

    querystr := fmt.Sprintf("SELECT user_no,order_no,line_no,amount,aboard_local_name,end_local_name,aboard_time FROM YUECHE_ORDER where driver_no = '%s' and description <> '已取消' ", input.Driver_no)
    //fmt.Println(querystr)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var ticket Wechat_DriverTicketInfo
    var ticket_slice []Wechat_DriverTicketInfo
    for rows.Next() {
    err = rows.Scan(&ticket.User_no,&ticket.Order_no,&ticket.Line_no,&ticket.Amount,&ticket.Aboard_local_name,&ticket.End_local_name,&ticket.Aboard_time)
      if err != nil {
        panic(err.Error())
        return
      }
      getLineinfobyid(ticket.Line_no, &ticket.S_city, &ticket.E_city, &ticket.Price_one, &ticket.Km_tot)

      num, _ := strconv.Atoi(ticket.Price_one)
      amount, _ := strconv.Atoi(ticket.Amount)
      lastprice := num * amount
      ticket.Price_amount = strconv.Itoa(int(lastprice))

      ticket_slice = append(ticket_slice, ticket)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Ticket = ticket_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

// 获取司机基本信息
func get_driver_info(Driver_no string, Name *string, Phone *string, Car_license *string, Headimgurl *string, Nickname *string) {
  querystr := fmt.Sprintf("select name,phone,car_license from YUECHE_DRIVER where driver_no = '%s'", Driver_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(Name, Phone, Car_license)
    if err != nil {
      panic(err.Error())
      return
    }
  }

  querystr = fmt.Sprintf("select headimgurl,nickname from YUECHE_DRIVER_BIND_WECHAT where driver_no = '%s'", Driver_no)
  rows,err = db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(Headimgurl, Nickname)
    if err != nil {
      panic(err.Error())
      return
    }
  }
}


type Wechat_DriverGetInfoInput struct {
  Driver_no string
}

type Wechat_DriverGetInfoRet struct {
  Ret string
  Message string
  Data struct {
    Driver_no string
    Name string
    Phone string
    Headimgurl string
    Nickname string
    Car_license string    // 车牌号
    Amount_people string  // 总人数
    Amount_Km_tot string  // 累计公里数
    Amount_star string    // 累计评级
  }
}

// http.HandleFunc("/api/wechat_drivergetinfo", wechat_drivergetinfoHandle)  // 获取司机综合信息
func wechat_drivergetinfoHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_drivergetinfoHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_DriverGetInfoRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_DriverGetInfoInput
    json.Unmarshal([]byte(result), &input)

    var order_no string
    var line_no string
    var amount string
    var amount_Km_tot int
    var amount_people int
    querystr := fmt.Sprintf("SELECT order_no,line_no,amount FROM YUECHE_ORDER where driver_no = '%s' and description <> '已取消' ", input.Driver_no)
    //fmt.Println(querystr)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&order_no,&line_no,&amount)
      if err != nil {
        panic(err.Error())
        return
      }
      var s_city string
      var e_city string
      var price_one string
      var km_tot string
      getLineinfobyid(line_no, &s_city, &e_city, &price_one, &km_tot)

      // 统计总公里数
      km_tot_tmp, _ := strconv.Atoi(km_tot)
      amount_Km_tot = amount_Km_tot + km_tot_tmp

      // 统计总人数
      amounttmp, _ := strconv.Atoi(amount)
      amount_people = amount_people + amounttmp
    }

    get_driver_info(input.Driver_no, &retmessage.Data.Name, &retmessage.Data.Phone, &retmessage.Data.Car_license, &retmessage.Data.Headimgurl, &retmessage.Data.Nickname)

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Driver_no = input.Driver_no
    retmessage.Data.Amount_people = strconv.Itoa(int(amount_people))
    retmessage.Data.Amount_Km_tot = strconv.Itoa(int(amount_Km_tot))
    retmessage.Data.Amount_star = "4.9"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

type Wechat_DriverKmtRankInfo struct {
  Driver_no string
  Driver_name string
  Ranking string
  km_tot string
}

type Wechat_DriverStarRankInfo struct {
  Driver_no string
  Driver_name string
  Ranking string
  Star string
}

type Wechat_DriverGetTopInput struct {
  Driver_no string
}

type Wechat_DriverGetTopRet struct {
  Ret string
  Message string
  Data struct {
    Ranking string
    Kmt_rank []Wechat_DriverKmtRankInfo
    Star_rank []Wechat_DriverStarRankInfo
  }
}

//http.HandleFunc("/api/wechat_drivergettop", wechat_drivergettopHandle)  // 获取司机排行榜
func wechat_drivergettopHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_drivergettopHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_DriverGetTopRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_DriverGetTopInput
    json.Unmarshal([]byte(result), &input)
/*
    var order_no string
    var line_no string
    var amount string
    var amount_Km_tot int
    var amount_people int
    querystr := fmt.Sprintf("SELECT order_no,line_no,amount FROM YUECHE_ORDER where driver_no = '%s' and description <> '已取消' ", input.Driver_no)
    //fmt.Println(querystr)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&order_no,&line_no,&amount)
      if err != nil {
        panic(err.Error())
        return
      }
      var s_city string
      var e_city string
      var price_one string
      var km_tot string
      getLineinfobyid(line_no, &s_city, &e_city, &price_one, &km_tot)

      // 统计总公里数
      km_tot_tmp, _ := strconv.Atoi(km_tot)
      amount_Km_tot = amount_Km_tot + km_tot_tmp

      // 统计总人数
      amounttmp, _ := strconv.Atoi(amount)
      amount_people = amount_people + amounttmp
    }
*/
    retmessage.Ret = "0"
    retmessage.Message = "success"
   // retmessage.Data.Driver_no = input.Driver_no
   // retmessage.Data.Amount_people = strconv.Itoa(int(amount_people))
   // retmessage.Data.Amount_Km_tot = strconv.Itoa(int(amount_Km_tot))
   // retmessage.Data.Amount_star = "4.9"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

type Wechat_DriverSetGpsInput struct {
  Driver_no string
  Longitude string
  Latitude  string
}

type Wechat_DriverSetGpsRet struct {
  Ret string
  Message string
}
//http.HandleFunc("/api/wechat_driversetgps", wechat_driversetgpsHandle)  // 司机将GPS设置进数据库
func wechat_driversetgpsHandle(w http.ResponseWriter, r *http.Request) { 
 // log.Println("performance wechat_driversetgpsHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_DriverSetGpsRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_DriverSetGpsInput
    json.Unmarshal([]byte(result), &input)

    stmtIns, err := db.Prepare("update YUECHE_DRIVER set longitude=?,latitude=? where driver_no=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(input.Longitude, input.Latitude, input.Driver_no)
    if err != nil {
      panic(err.Error())
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

type Wechat_GpsInfo struct {
  Driver_no string
  Longitude string
  Latitude  string
  Name string
  Phone string
  Car_license string
}
/*
type Wechat_GetGpsInput struct {
  Driver_no string
}*/

type Wechat_GetGpsRet struct {
  Ret string
  Message string
  Data struct {
    Driver_gps []Wechat_GpsInfo
  }
}
//http.HandleFunc("/api/wechat_getgps", wechat_getgpsHandle)  // 获取司机的GPS
func wechat_getgpsHandle(w http.ResponseWriter, r *http.Request) { 
//  log.Println("performance wechat_getgpsHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetGpsRet
    //result, _:= ioutil.ReadAll(r.Body)
    //var input Wechat_GetGpsInput
    //json.Unmarshal([]byte(result), &input)

    //querystr := fmt.Sprintf("SELECT driver_no,longitude,latitude FROM YUECHE_DRIVER where driver_no = '%s'", input.Driver_no)
    querystr := "SELECT driver_no,longitude,latitude,name,phone,car_license FROM YUECHE_DRIVER where d_status='on'"
    //fmt.Println(querystr)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var gps_info Wechat_GpsInfo
    var gps_info_slice []Wechat_GpsInfo
    for rows.Next() {
    err = rows.Scan(&gps_info.Driver_no,&gps_info.Longitude,&gps_info.Latitude,&gps_info.Name,&gps_info.Phone,&gps_info.Car_license)
      if err != nil {
        panic(err.Error())
        return
      }
      //gps_info.Driver_no = input.Driver_no
      gps_info_slice = append(gps_info_slice, gps_info)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Driver_gps = gps_info_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}


type Wechat_GetSmsCodeInput struct {
  User_no string
  Phone string
}

type Wechat_GetSmsCodeRet struct {
  Ret string
  Message string
  Data struct {
    Sms_code string
  }
}

//http.HandleFunc("/api/wechat_getsmscode", wechat_getsmscodeHandle)  // 请求获取短信验证码
func wechat_getsmscodeHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getsmscodeHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetSmsCodeRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_GetSmsCodeInput
    json.Unmarshal([]byte(result), &input)

    sms_code := GetRandomNum4Str()

    //messagexsend
    messageconfig := make(map[string]string)
    messageconfig["appid"] = "21084"
    messageconfig["appkey"] = "323cdfa1b38815aaecacea37b3369938"
    messageconfig["signtype"] = "md5"

    messagexsend := submail.CreateMessageXSend()
    submail.MessageXSendAddTo(messagexsend, input.Phone)
    submail.MessageXSendSetProject(messagexsend, "29r2V3")
    submail.MessageXSendAddVar(messagexsend, "code", sms_code)
    fmt.Println("MessageXSend ", submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig))

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Sms_code = sms_code
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}




type Wechat_AreabyStartcityInput struct {
  City string
}

type Wechat_AreabyStartcityRet struct {
  Ret string
  Message string
  Data struct {
    Area []string
  }
}

//http.HandleFunc("/api/wechat_areabycity", wechat_areabystartcityHandle)   // 根据城市获取城市区域名称
func wechat_areabystartcityHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_areabystartcityHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_AreabyStartcityRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_AreabyStartcityInput
    json.Unmarshal([]byte(result), &input)

    var area string
    var area_slice []string
    querystr := fmt.Sprintf("select area from YUECHE_BUS_AREA where city = '%s'", input.City)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&area)
      if err != nil {
        panic(err.Error())
        return
      }

      area_slice = append(area_slice, area)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Area = area_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}



type Wechat_StationbyAreaInput struct {
  Area string
}

type Wechat_StationbyAreaRet struct {
  Ret string
  Message string
  Data struct {
    Station []string
  }
}

// http.HandleFunc("/api/wechat_stationbyarea", wechat_stationbyareaHandle)   // 根据区域名称获取站点信息,乘客使用 -427
func wechat_stationbyareaHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_stationbyareaHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_StationbyAreaRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_StationbyAreaInput
    json.Unmarshal([]byte(result), &input)

    var station string
    var station_slice []string
    querystr := fmt.Sprintf("select via_local_name from YUECHE_BUS_LINE_VIA where description = '%s' and local_status='on'", input.Area)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&station)
      if err != nil {
        panic(err.Error())
        return
      }

      station_slice = append(station_slice, station)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Station = station_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}




type Wechat_EndorseTicketInput struct {
  User_no string
  Order_no string
}

type Wechat_EndorseTicketRet struct {
  Ret string
  Message string
}

//http.HandleFunc("/api/wechat_endorseticket", wechat_endorseticketHandle)   // 改签车票
func wechat_endorseticketHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_endorseticketHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_EndorseTicketRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_EndorseTicketInput
    json.Unmarshal([]byte(result), &input)

    stmtIns, err := db.Prepare("update YUECHE_PAYMENTS set orderstatus=? where order_no=? and user_no =?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec("已改签", input.Order_no, input.User_no)
    if err != nil {
      panic(err.Error())
    }
    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}


type Wechat_EndorseTicketExInput struct {
  User_no string
  Order_no string
  New_time string
}

type Wechat_EndorseTicketExRet struct {
  Ret string
  Message string
}

//http.HandleFunc("/api/wechat_endorseticket_ex", wechat_endorseticketexHandle)   // 改签车票
func wechat_endorseticketexHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_endorseticketexHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_EndorseTicketExRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_EndorseTicketExInput
    json.Unmarshal([]byte(result), &input)

    fmt.Println(string(input.New_time))

    t1, _ := time.Parse("2006-01-02 15:04:05", input.New_time)
    timeStr1 := t1.Format("2006-01-02")
    timeStr2 := timeStr1 + " 23:59:59"

    stmtIns, err := db.Prepare("update YUECHE_ORDER set aboard_time=?,ride_time_s=?,ride_time_e=? where order_no=? and user_no =?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(input.New_time, timeStr1, timeStr2, input.Order_no, input.User_no)
    if err != nil {
      panic(err.Error())
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}



type Wechat_RefundTicketInput struct {
  User_no string
  Order_no string
}

type Wechat_RefundTicketRet struct {
  Ret string
  Message string
}

//http.HandleFunc("/api/wechat_refundticket", wechat_refundticketHandle)   // 退票
func wechat_refundticketHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_refundticketHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_RefundTicketRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_RefundTicketInput
    json.Unmarshal([]byte(result), &input)

    stmtIns, err := db.Prepare("update YUECHE_PAYMENTS set orderstatus=? where order_no=? and user_no =?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec("已退票", input.Order_no, input.User_no)
    if err != nil {
      panic(err.Error())
    }
    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}



// http.HandleFunc("/api/wechat_getbustime", wechat_getbustimeHandle)   // 获取发车时间
type Wechat_GetbustimeRet struct {
  Ret string
  Message string
  Data struct {
    Air_time []string
    Common_time []string
  }
}

func wechat_getbustimeHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getbustimeHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetbustimeRet

    var bus_time string
    querystr := "select bus_time from YUECHE_BUS_TIME where line_name = '机场专线' order by bus_time"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
      err = rows.Scan(&bus_time)
      if err != nil {
        panic(err.Error())
        return
      }

      retmessage.Data.Air_time = append(retmessage.Data.Air_time, bus_time)
    }


    querystr = "select bus_time from YUECHE_BUS_TIME where line_name = '一般线路' order by bus_time"
    rows,err = db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
      err = rows.Scan(&bus_time)
      if err != nil {
        panic(err.Error())
        return
      }

      retmessage.Data.Common_time = append(retmessage.Data.Common_time, bus_time)
    }


    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}


// http.HandleFunc("/api/wechat_getbigbustime", wechat_getbigbustimeHandle)   // 获取大巴发车时间
type Bigbustime struct {
  Line_no string
  Bus_time []string
}

type Wechat_GetbigbustimeRet struct {
  Ret string
  Message string
  Data struct {
    Big_bus_time []Bigbustime
  }
}

func wechat_getbigbustimeHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getbigbustimeHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetbigbustimeRet

    querystr := "select distinct line_no from YUECHE_BIG_BUS_TIME"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var line_no string
    var line_no_slice []string
    for rows.Next() {
      err = rows.Scan(&line_no)
      if err != nil {
        panic(err.Error())
        return
      }
      line_no_slice = append(line_no_slice, line_no)
    }

    
    var bigbustime_slice []Bigbustime
    for _, value := range line_no_slice {
      var bigbustime Bigbustime
      var bus_time string
      querystr := fmt.Sprintf("select bus_time from YUECHE_BIG_BUS_TIME where line_no = '%s'", value)
      rows,err := db.Query(querystr)
      if err != nil {
        log.Fatal(err)
        return
      }
      defer rows.Close()

      bigbustime.Line_no = value

      for rows.Next() {
        err = rows.Scan(&bus_time)
        if err != nil {
          panic(err.Error())
          return
        }
        bigbustime.Bus_time = append(bigbustime.Bus_time, bus_time)
      }

      bigbustime_slice = append(bigbustime_slice, bigbustime)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Big_bus_time = bigbustime_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}

// 司机上车打开，绑定车辆到司机表中，并且将当前绑定信息写入到司机出车表中
type Wechat_AboardDriverInput struct {
  Driver_no string
  Car_license string
}

type Wechat_AboardDriverRet struct {
  Ret string
  Message string
  Data struct {
    Diapatch_no string
  }
}

//http.HandleFunc("/api/wechat_aboarddriver", wechat_aboarddriverHandle)  
func wechat_aboarddriverHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_aboarddriverHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_AboardDriverRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_AboardDriverInput
    json.Unmarshal([]byte(result), &input)

    stmtIns, err := db.Prepare("update YUECHE_DRIVER set car_license=? where driver_no=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(input.Car_license, input.Driver_no)
    if err != nil {
      panic(err.Error())
    }

    // 将当前信息写入到司机出车表中
    t := time.Now().Unix()
    diapatch_no := fmt.Sprintf("DIS%d", t)

    stmtIns, err = db.Prepare("INSERT INTO YUECHE_DISPATCH VALUES( ?,?,?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(diapatch_no, input.Driver_no, input.Car_license, time.Now(), ZERO_TIME)
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Diapatch_no = diapatch_no
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}

// 司机下车
type Wechat_OffDriverInput struct {
  Driver_no string
  Car_license string
}

type Wechat_OffDriverRet struct {
  Ret string
  Message string
}

//http.HandleFunc("/api/wechat_offdriver", wechat_offdriverHandle)
func wechat_offdriverHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_offdriverHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_AboardDriverRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_AboardDriverInput
    json.Unmarshal([]byte(result), &input)

    stmtIns, err := db.Prepare("update YUECHE_DISPATCH set off_time=? where driver_no=? and car_license=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(time.Now(), input.Driver_no, input.Car_license)
    if err != nil {
      panic(err.Error())
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}


// 司机查验乘客电子票,每扫描一次，则YUECHE_GET_ON_CAR中的记录加1
type Wechat_DriverCheckTicketInput struct {
  Driver_no string
  Ticket_no string
  Longitude string
  Latitude string
}

type Wechat_DriverCheckTicketRet struct {
  Ret string
  Message string
  Data struct {
    Residue int
  }
}

func get_residueticket_by_ticketno(Ticket_no string, Residue *int, Order_no *string) {

  var amount string = "0"
  var order_no string = "0"
  querystr := fmt.Sprintf("select amount,order_no from YUECHE_ORDER where ticket_no = '%s'", Ticket_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(&amount, &order_no)
    if err != nil {
      panic(err.Error())
      return
    }
  }
  num_amount, _ := strconv.Atoi(amount)
  *Order_no = order_no

  // 统计已经上车的人数
  var rcount string = "0"
  querystr = fmt.Sprintf("SELECT COUNT(*) FROM YUECHE_GET_ON_CAR where order_no = '%s'", order_no)
  rows,err = db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
  err = rows.Scan(&rcount)
    if err != nil {
      panic(err.Error())
      return
    }
  }
  num_on, _ := strconv.Atoi(rcount)

  Residue_i := num_amount - num_on

  if Residue_i > 0 {
    *Residue = Residue_i
  } else {
    *Residue = 0
  }
}

func get_carlicense_by_driver(Driver_no string, Car_license *string) {
  querystr := fmt.Sprintf("select car_license from YUECHE_DRIVER where driver_no = '%s'", Driver_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(Car_license)
    if err != nil {
      panic(err.Error())
      return
    }
  }
}

//http.HandleFunc("/api/wechat_drivercheckticket", wechat_drivercheckticketHandle)
func wechat_drivercheckticketHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_drivercheckticketHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_DriverCheckTicketRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_DriverCheckTicketInput
    json.Unmarshal([]byte(result), &input)

    // 通过Ticket_no获取订单no,然后通过订单no返回当前的余票数
    var Residue int = 0
    var Order_no string = "0"
    get_residueticket_by_ticketno(input.Ticket_no, &Residue, &Order_no)
    if Order_no == "0" {
      retmessage.Ret = "1"
      retmessage.Message = "failure"
      b, _ := json.Marshal(retmessage)
      w.Header().Set("Access-Control-Allow-Origin","*")
      fmt.Fprintf(w, "%s", b)
      return
    }

    timeStr := time.Now()
    var rcount string = "0"
    querystr := fmt.Sprintf("SELECT COUNT(*) FROM YUECHE_ORDER where order_no = '%s' and (description <> '已取消' and description <> '已上车' and description <> '未付款') and (ride_time_s < '%s' and ride_time_e > '%s')", Order_no, timeStr, timeStr)
    fmt.Println(querystr)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&rcount)
      if err != nil {
        panic(err.Error())
        return
      }
    }

    if rcount == "0" {
      retmessage.Ret = "1"
      retmessage.Message = "failure"
      b, _ := json.Marshal(retmessage)
      w.Header().Set("Access-Control-Allow-Origin","*")
      fmt.Fprintf(w, "%s", b)
      return
    }


    if Residue > 0{
      var car_license string = ""
      get_carlicense_by_driver(input.Driver_no, &car_license)

      stmtIns, err := db.Prepare("INSERT INTO YUECHE_GET_ON_CAR VALUES( ?,?,?,?,?,?)")
      if err != nil {
          panic(err.Error()) // proper error handling instead of panic in your app
      }
      defer stmtIns.Close() 

      _, err = stmtIns.Exec(input.Driver_no, car_license, Order_no, time.Now(), input.Longitude, input.Latitude)
      if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
      }
    } else {
      retmessage.Ret = "1"
      retmessage.Message = "failure"
      b, _ := json.Marshal(retmessage)
      w.Header().Set("Access-Control-Allow-Origin","*")
      fmt.Fprintf(w, "%s", b)
      return
    }

    // 当前票为最后一张的时候
    if Residue == 1 {
      // 将订单表的状态修改为“已完成”
      stmtIns, err := db.Prepare("update YUECHE_PAYMENTS set orderstatus=? where order_no=?")
      if err != nil {
          panic(err.Error())
      }
      defer stmtIns.Close() 

      _, err = stmtIns.Exec("已完成",Order_no)
      if err != nil {
        panic(err.Error())
      }

      stmtIns, err = db.Prepare("update YUECHE_ORDER set description=?,driver_no=? where order_no=?")
      if err != nil {
          panic(err.Error())
      }
      defer stmtIns.Close() 

      _, err = stmtIns.Exec("已上车", input.Driver_no, Order_no)
      if err != nil {
        panic(err.Error())
      }
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Residue = Residue - 1
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}




// 获取包车出发点
type Wechat_GetmultiStartcityRet struct {
  Ret string
  Message string
  Data struct {
    S_city []string
  }
}

//http.HandleFunc("/api/wechat_getmultistartcity", wechat_getmultistartcityHandle)
func wechat_getmultistartcityHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getmultistartcityHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetmultiStartcityRet

    querystr := "select start_city from YUECHE_MULTI_START_CITY"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var city string
    var city_slice []string
    for rows.Next() {
      err = rows.Scan(&city)
      if err != nil {
        panic(err.Error())
        return
      }

      city_slice = append(city_slice, city)
    }


    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.S_city = city_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}


// 获取包车公共信息
type Multi_model struct {
  Travel_model string
  Price string
}


type Wechat_GetmultiModelRet struct {
  Ret string
  Message string
  Data struct {
    Travel_info []Multi_model
  }
}

//http.HandleFunc("/api/wechat_getmultimodel", wechat_getmultimodelHandle)
func wechat_getmultimodelHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_getmultimodelHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_GetmultiModelRet

    querystr := "select travel_model,price from YUECHE_MULTI_MODEL"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var t_model Multi_model
    var t_model_slice []Multi_model
    for rows.Next() {
      err = rows.Scan(&t_model.Travel_model, &t_model.Price)
      if err != nil {
        panic(err.Error())
        return
      }

      t_model_slice = append(t_model_slice, t_model)
    }


    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Travel_info = t_model_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}



// 乘客建议
type Wechat_SetSuggestInput struct {
  User_no string
  Suggest string   // 乘客输入的建议信息
}


type Wechat_SetSuggestRet struct {
  Ret string
  Message string
}

// http.HandleFunc("/api/wechat_setsuggest", wechat_setsuggestHandle)
func wechat_setsuggestHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_setsuggestHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_SetSuggestRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_SetSuggestInput
    json.Unmarshal([]byte(result), &input)

    stmtIns, err := db.Prepare("INSERT INTO YUECHE_SUGGEST VALUES( ?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(input.User_no, input.Suggest, time.Now())
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}


// 乘客评价订单-- 只有已完成订单才允许乘客评价
type Wechat_EvaluateOrderInput struct {
  User_no string
  Order_no string  
  Star_rank string  // 星级
  Desp string    // 文字评价
}


type Wechat_EvaluateOrderRet struct {
  Ret string
  Message string
}

//http.HandleFunc("/api/wechat_evaluateorder", wechat_evaluateorderHandle)
func wechat_evaluateorderHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_evaluateorderHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Wechat_EvaluateOrderRet
    result, _:= ioutil.ReadAll(r.Body)
    var input Wechat_EvaluateOrderInput
    json.Unmarshal([]byte(result), &input)

    stmtIns, err := db.Prepare("INSERT INTO YUECHE_ORDER_EVALUATE VALUES( ?,?,?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(input.Star_rank, input.Order_no, input.User_no, input.Desp, time.Now())
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }


    stmtIns, err = db.Prepare("update YUECHE_PAYMENTS set orderstatus=? where order_no=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec("已评价", input.Order_no)
    if err != nil {
      panic(err.Error())
    }


    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return
  }  
}


type JSAPIWeixinpayInput struct {
  Appid string
  Openid string
  Out_trade_no string
  Spbill_create_ip  string
  Total_fee  string
}

type JSAPIWeixinpayRet struct {
  Ret string
  Message string
  Data struct {
    /*
    Return_code string
    Return_msg string
    Appid string
    Mch_id string
    Nonce_str string
    Sign string
    Result_code string
    Prepay_id string
    Trade_type string
    */
    Appid string
    Partnerid string
    Prepayid string
    Package string
    Noncestr string
    Timestamp string
    Sign string
  }
}

type JSAPIXmlDataWeixinPay struct {
  Appid string `xml:"appid"` 
  Attach string `xml:"attach"` 
  Body string `xml:"body"` 
  Mch_id string `xml:"mch_id"` 
  Nonce_str string `xml:"nonce_str"` 
  Notify_url string `xml:"notify_url"` 
  Openid string `xml:"openid"` 
  Out_trade_no string `xml:"out_trade_no"` 
  Spbill_create_ip string `xml:"spbill_create_ip"` 
  Total_fee string `xml:"total_fee"` 
  Trade_type string `xml:"trade_type"` 
  Sign string `xml:"sign"` 
}

type JSAPIXmlDataWeixinPayRet struct {
  Return_code string `xml:"return_code"` 
  Return_msg string `xml:"return_msg"` 
  Appid string `xml:"appid"` 
  Mch_id string `xml:"mch_id"` 
  Nonce_str string `xml:"nonce_str"` 
  Sign string `xml:"sign"` 
  Result_code string `xml:"result_code"` 
  Prepay_id string `xml:"prepay_id"` 
  Trade_type string `xml:"trade_type"`  
}

var JSAPIappid string = "wxb994de5246988709"
var JSAPImch_id string = "1498107782"
var JSAPInotify_url string = "/api/weixin/JSAPIpaynotify"
var JSAPIbase_url string = "https://api.yueche520.com"
var JSAPImach_key string = "fadfadfaf38943jlnvmakmfakfu78343"

// 小程序微信支付接口
//http.HandleFunc("/api/weixin/JSAPIpay",JSAPI_weixinpayHandle)
func JSAPI_weixinpayHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance JSAPI_weixinpayHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage JSAPIWeixinpayRet
    result, _:= ioutil.ReadAll(r.Body)
    var input JSAPIWeixinpayInput
    json.Unmarshal([]byte(result), &input)

    fmt.Println(result)
    fmt.Println(input.Appid)

    // 对传入的支付数与数据库进行对比-重要环节，防篡改
/*
    // 通过Out_trade_no【订单号】查找出应付金额，与传递进来的金额进行比较，如果不相同，则直接返回错误
    querystr := fmt.Sprintf("select payprice from TONGLAR_PAYMENTS where orderid = '%s'", input.Out_trade_no)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var payprice string
    for rows.Next() {
      err = rows.Scan(&payprice)
      if err != nil {
        panic(err.Error())
        return
      }
    }

    if payprice != input.Total_fee {
      retmessage.Ret = "1"
      retmessage.Message = "fuck you"
      b, _ := json.Marshal(retmessage)
      fmt.Fprintf(w, "%s", b)
      return  
    }
*/


/*
   <appid>wx2421b1c4370ec43b</appid>
   <attach>支付测试</attach>
   <body>APP支付测试</body>
   <mch_id>10000100</mch_id>
   <nonce_str>1add1a30ac87aa2db72f57a2375d8fec</nonce_str>
   <notify_url>http://wxpay.wxutil.com/pub_v2/pay/notify.v2.php</notify_url>
   <openid>oUpF8uMuAJO_M2pxb1Q9zNjWeS6o</openid>
   <out_trade_no>1415659990</out_trade_no>
   <spbill_create_ip>14.23.150.211</spbill_create_ip>
   <total_fee>1</total_fee>
   <trade_type>JSAPI</trade_type>
   <sign>0CB01533B8C1EF103065174F50BCA001</sign>
*/
    nonce_str := GetRandomString(32)
    last_url := JSAPIbase_url + JSAPInotify_url
    stringA := fmt.Sprintf("appid=%s&attach=九州约车&body=九州约车&mch_id=%s&nonce_str=%s&notify_url=%s&openid=%s&out_trade_no=%s&spbill_create_ip=%s&total_fee=%s&trade_type=JSAPI&key=%s", 
                            input.Appid, JSAPImch_id, nonce_str, last_url, input.Openid, input.Out_trade_no, input.Spbill_create_ip, input.Total_fee, JSAPImach_key)

    sign := strings.ToUpper(MD5(stringA))

    var v JSAPIXmlDataWeixinPay
    v.Appid = input.Appid
    v.Attach = "九州约车"
    v.Body = "九州约车"
    v.Mch_id = JSAPImch_id
    v.Nonce_str = nonce_str
    v.Notify_url = last_url
    v.Openid = input.Openid
    v.Out_trade_no = input.Out_trade_no
    v.Spbill_create_ip = input.Spbill_create_ip
    v.Total_fee = input.Total_fee
    v.Trade_type = "JSAPI"
    v.Sign = sign

    output, err := xml.MarshalIndent(&v, "", "\t")
    if err != nil {
      fmt.Printf("error: %v\n", err)
    }

    // 发送到微信官方接口
    var neturl string = "https://api.mch.weixin.qq.com/pay/unifiedorder"
    body := bytes.NewBuffer([]byte(output))
    res,err := http.Post(neturl, "application/xml;charset=utf-8", body)
    if err != nil {
      log.Fatal(err)
      return
    }
    result, err = ioutil.ReadAll(res.Body)
    if err != nil {
      log.Fatal(err)
      return 
    }

    fmt.Println(string(result))

/*
    <xml>
   <return_code><![CDATA[SUCCESS]]></return_code>
   <return_msg><![CDATA[OK]]></return_msg>
   <appid><![CDATA[wx2421b1c4370ec43b]]></appid>
   <mch_id><![CDATA[10000100]]></mch_id>
   <nonce_str><![CDATA[IITRi8Iabbblz1Jc]]></nonce_str>
   <sign><![CDATA[7921E432F65EB8ED0CE9755F0E86D72F]]></sign>
   <result_code><![CDATA[SUCCESS]]></result_code>
   <prepay_id><![CDATA[wx201411101639507cbf6ffd8b0779950874]]></prepay_id>
   <trade_type><![CDATA[APP]]></trade_type>
  </xml> 
*/
    var weixinret JSAPIXmlDataWeixinPayRet
    err = xml.Unmarshal(result, &weixinret)
    if err != nil {
      log.Fatal(err)
      return
    }

    fmt.Println(weixinret.Prepay_id)

    retmessage.Ret = "0"
    retmessage.Message = "success"
    /*
    retmessage.Data.Return_code = weixinret.Return_code
    retmessage.Data.Return_msg = weixinret.Return_msg
    retmessage.Data.Appid = weixinret.Appid
    retmessage.Data.Mch_id = weixinret.Mch_id
    retmessage.Data.Nonce_str = weixinret.Nonce_str
    retmessage.Data.Sign = weixinret.Sign
    retmessage.Data.Result_code = weixinret.Result_code
    retmessage.Data.Prepay_id = weixinret.Prepay_id
    retmessage.Data.Trade_type = weixinret.Trade_type
*/

    t := time.Now().Unix()
    stringB := fmt.Sprintf("appId=%s&nonceStr=%s&package=prepay_id=%s&signType=MD5&timeStamp=%d&key=%s", 
                            weixinret.Appid, weixinret.Nonce_str, weixinret.Prepay_id, t, JSAPImach_key)


    stringC := fmt.Sprintf("prepay_id=%s", weixinret.Prepay_id)
    signB := strings.ToUpper(MD5(stringB))
    fmt.Println(stringB)
    fmt.Println(signB)
    retmessage.Data.Appid = weixinret.Appid
    retmessage.Data.Partnerid = weixinret.Mch_id
    retmessage.Data.Prepayid = weixinret.Prepay_id
    retmessage.Data.Package = stringC
    retmessage.Data.Noncestr = weixinret.Nonce_str
    retmessage.Data.Timestamp = strconv.FormatInt(t, 10)
    retmessage.Data.Sign = signB


    b, _ := json.Marshal(retmessage)
    fmt.Println(b)
    fmt.Fprintf(w, "%s", b)
    return   
  }
}

func get_start_end_by_Order(Order_no string, Start *string, End *string, Aboard_time *string, Phone *string, Ticket_no *string) {
  var line_no string 
  var user_no string
  var querystr string = ""


  querystr = fmt.Sprintf("select line_no,user_no,ticket_no,aboard_time from YUECHE_ORDER where order_no = '%s'", Order_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(&line_no, &user_no, Ticket_no, Aboard_time)
    if err != nil {
      panic(err.Error())
      return
    }
  }

  order_prefix := Order_no[6:7]
  if order_prefix == "B" {
    querystr = fmt.Sprintf("select s_city,e_city from YUECHE_BIG_BUS_LINE where line_no = '%s'", line_no)
  } else {
    querystr = fmt.Sprintf("select s_city,e_city from YUECHE_BUS_LINE where line_no = '%s'", line_no)
  }

  //querystr = fmt.Sprintf("select s_city,e_city from YUECHE_BUS_LINE where line_no = '%s'", line_no)
  rows,err = db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(Start, End)
    if err != nil {
      panic(err.Error())
      return
    }
  }

  querystr = fmt.Sprintf("select phone from YUECHE_USER where user_no = '%s'", user_no)
  rows,err = db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(Phone)
    if err != nil {
      panic(err.Error())
      return
    }
  }
}

func get_start_end_by_Multi(Multi_no string, Phone *string, Desp *string) {

  var user_no string
  querystr := fmt.Sprintf("select description,user_no from YUECHE_MULTI where multi_no = '%s'", Multi_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(Desp, &user_no)
    if err != nil {
      panic(err.Error())
      return
    }
  }


  querystr = fmt.Sprintf("select phone from YUECHE_USER where user_no = '%s'", user_no)
  rows,err = db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(Phone)
    if err != nil {
      panic(err.Error())
      return
    }
  }
}


func send_sms_order(phone string, start_city string, end_city string, aboard_time string, ticket_no string) { 
  log.Println("performance send_sms_order");

  //messagexsend
  messageconfig := make(map[string]string)
  messageconfig["appid"] = "21084"
  messageconfig["appkey"] = "323cdfa1b38815aaecacea37b3369938"
  messageconfig["signtype"] = "md5"

  messagexsend := submail.CreateMessageXSend()
  submail.MessageXSendAddTo(messagexsend, phone)
  submail.MessageXSendSetProject(messagexsend, "y2kKY1")
  submail.MessageXSendAddVar(messagexsend, "start_city", start_city)
  submail.MessageXSendAddVar(messagexsend, "end_city", end_city)
  submail.MessageXSendAddVar(messagexsend, "aboard_time", aboard_time)
  submail.MessageXSendAddVar(messagexsend, "ticket_no", ticket_no)
  fmt.Println("MessageXSend ", submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig))
}

func send_sms_order_tuikuan(phone string, order_no string) { 
  log.Println("performance send_sms_order_tuikuan");

  //messagexsend
  messageconfig := make(map[string]string)
  messageconfig["appid"] = "21084"
  messageconfig["appkey"] = "323cdfa1b38815aaecacea37b3369938"
  messageconfig["signtype"] = "md5"

  messagexsend := submail.CreateMessageXSend()
  submail.MessageXSendAddTo(messagexsend, phone)
  submail.MessageXSendSetProject(messagexsend, "7h3ez1")
  submail.MessageXSendAddVar(messagexsend, "order_no", order_no)
  fmt.Println("MessageXSend ", submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig))
}

func send_sms_order_multi(phone string, multi_no string, desp string) { 
  log.Println("performance send_sms_order_multi");

  //messagexsend
  messageconfig := make(map[string]string)
  messageconfig["appid"] = "21084"
  messageconfig["appkey"] = "323cdfa1b38815aaecacea37b3369938"
  messageconfig["signtype"] = "md5"

  messagexsend := submail.CreateMessageXSend()
  submail.MessageXSendAddTo(messagexsend, phone)
  submail.MessageXSendSetProject(messagexsend, "Ysl4K4")
  submail.MessageXSendAddVar(messagexsend, "desp", desp)
  submail.MessageXSendAddVar(messagexsend, "multi_no", multi_no)
  fmt.Println("MessageXSend ", submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig))
}

func statistics_ticket_by_user(Order_no string) {
  var user_no string
  querystr := fmt.Sprintf("select user_no from YUECHE_PAYMENTS where order_no = '%s'", Order_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
  err = rows.Scan(&user_no)
    if err != nil {
      panic(err.Error())
      return
    }
  }

  // 通过 user_no 查看统计表中，该用户购买了多少张票，如果本次达到10张，总送一次价值180元的优惠券，然后将购票次数改为0
  var ticket_count int = 0
  querystr = fmt.Sprintf("select ticket_count from YUEHCE_TICKET_COUNT where user_no = '%s'", user_no)
  rows,err = db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
  err = rows.Scan(&ticket_count)
    if err != nil {
      panic(err.Error())
      return
    }
  }

  if ticket_count == 9 {
    coupon_no := "CON_002"
    stmtIns, err := db.Prepare("INSERT INTO YUECHE_USER_COUPON VALUES( ?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(user_no, coupon_no, time.Now())
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }
  }

  ticket_count = ticket_count + 1
  stmtIns, err := db.Prepare("update YUEHCE_TICKET_COUNT set ticket_count=? where user_no=?")
  if err != nil {
      panic(err.Error())
  }
  defer stmtIns.Close() 

  _, err = stmtIns.Exec(ticket_count, user_no)
  if err != nil {
    panic(err.Error())
  }

/*
  if ticket_count < 9 {
    ticket_count = ticket_count + 1
    stmtIns, err := db.Prepare("update YUEHCE_TICKET_COUNT set ticket_count=? where user_no=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(ticket_count, user_no)
    if err != nil {
      panic(err.Error())
    }
  } else {
    // 赠与用户180元优惠券，并将计数清0
    coupon_no := "CON_002"
    stmtIns, err := db.Prepare("INSERT INTO YUECHE_USER_COUPON VALUES( ?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(user_no, coupon_no, time.Now())
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    // 计数清0
    stmtIns, err = db.Prepare("update YUEHCE_TICKET_COUNT set ticket_count=? where user_no=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(0, user_no)
    if err != nil {
      panic(err.Error())
    }
  }
  */
}


func delete_counpon_by_order(Order_no string) {

  var coupon_no,user_no string
  querystr := fmt.Sprintf("select coupon_no,user_no from YUECHE_PAYMENTS where order_no = '%s'", Order_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
  err = rows.Scan(&coupon_no,&user_no)
    if err != nil {
      panic(err.Error())
      return
    }
  }

  stmtIns, err := db.Prepare("DELETE FROM YUECHE_USER_COUPON WHERE user_no =? and coupon_no =?")
  if err != nil {
      panic(err.Error())
  }
  defer stmtIns.Close() 

  _, err = stmtIns.Exec(user_no, coupon_no)
  if err != nil {
    panic(err.Error())
  }
}


type JSAPIWeixinInputMessage struct {
    Appid string `xml:"appid"` 
    Attach string `xml:"attach"` 
    Bank_type string `xml:"bank_type"` 
    Fee_type string `xml:"fee_type"` 
    Is_subscribe string `xml:"is_subscribe"` 
    Mch_id string `xml:"mch_id"` 
    Nonce_str string `xml:"nonce_str"` 
    Openid string `xml:"openid"` 
    Out_trade_no string `xml:"out_trade_no"` 
    Result_code string `xml:"result_code"` 
    Return_code string `xml:"return_code"` 
    Sign string `xml:"sign"` 
    Sub_mch_id string `xml:"sub_mch_id"` 
    Time_end string `xml:"time_end"` 
    Total_fee string `xml:"total_fee"` 
    Trade_type string `xml:"trade_type"` 
    Transaction_id string `xml:"transaction_id"` 
}

type JSAPIWeixinNotifyRet struct {
    Return_code string `xml:"return_code"` 
  //  Return_msg string `xml:"return_msg"` 
}

//http.HandleFunc("/api/weixin/JSAPIpaynotify",JSAPI_weixinpaynotifyHandle)
func JSAPI_weixinpaynotifyHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance JSAPI_weixinpaynotifyHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var v JSAPIWeixinNotifyRet
/*
  <xml>
  <appid><![CDATA[wx2421b1c4370ec43b]]></appid>
  <attach><![CDATA[支付测试]]></attach>
  <bank_type><![CDATA[CFT]]></bank_type>
  <fee_type><![CDATA[CNY]]></fee_type>
  <is_subscribe><![CDATA[Y]]></is_subscribe>
  <mch_id><![CDATA[10000100]]></mch_id>
  <nonce_str><![CDATA[5d2b6c2a8db53831f7eda20af46e531c]]></nonce_str>
  <openid><![CDATA[oUpF8uMEb4qRXf22hE3X68TekukE]]></openid>
  <out_trade_no><![CDATA[1409811653]]></out_trade_no>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <sign><![CDATA[B552ED6B279343CB493C5DD0D78AB241]]></sign>
  <sub_mch_id><![CDATA[10000100]]></sub_mch_id>
  <time_end><![CDATA[20140903131540]]></time_end>
  <total_fee>1</total_fee>
  <trade_type><![CDATA[JSAPI]]></trade_type>
  <transaction_id><![CDATA[1004400740201409030005092168]]></transaction_id>
  </xml> 
 */
    // 接受微信服务器信息
    result, _:= ioutil.ReadAll(r.Body)
    var imes JSAPIWeixinInputMessage
    xml.Unmarshal([]byte(result), &imes)

    fmt.Println(string(result))

    // 中间做一些业务逻辑处理
    // 当收到通知进行处理时，首先检查对应业务数据的状态，判断该通知是否已经处理过
    
    var tradenum string 
    querystr := fmt.Sprintf("select tradenum from YUECHE_PAYMENTS where tradenum = '%s'", imes.Transaction_id)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
      err = rows.Scan(&tradenum)
      if err != nil {
        panic(err.Error())
        return
      }
    }
    if tradenum == imes.Transaction_id {
      v.Return_code = "SUCCESS"

      output, err := xml.MarshalIndent(&v, "", "\t")
      if err != nil {
        fmt.Printf("error: %v\n", err)
      }
      fmt.Fprintf(w, "%s", output)
      return
    }

    // 将收到的微信订单号[transaction_id],和商户订单号【out_trade_no】一起绑定到TONGLAR_PAYMENTS表去
    orderstatus := "未完成"  // 未完成
    payway := "weico"
    paystatus := "success"
    t := time.Now().Unix()
    payid := fmt.Sprintf("PAY_%d", t)
    stmtIns, err := db.Prepare("update YUECHE_PAYMENTS set tradenum=?,orderstatus=?,payway=?,paystatus=?,payid=? where order_no=?")
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(imes.Transaction_id, orderstatus, payway, paystatus, payid, imes.Out_trade_no)
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    order_prefix := imes.Out_trade_no[0:5]
    if order_prefix == "ORDER" {
      // 个人订单支付
      stmtIns, err := db.Prepare("update YUECHE_ORDER set description=? where order_no=?")
      if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
      }
      defer stmtIns.Close() 

      _, err = stmtIns.Exec("已付款", imes.Out_trade_no)
      if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
      }

      // 根据订单号，获取相关信息，并发送短信
      var aboard_time string
      var start_city string
      var end_city string
      var phone string
      var ticket_no string
      get_start_end_by_Order(imes.Out_trade_no, &start_city, &end_city, &aboard_time, &phone, &ticket_no)

      s1 := aboard_time[0:10]
      s2 := aboard_time[11:19]
      s3 := s1 + " " + s2
      send_sms_order(phone, start_city, end_city, s3, ticket_no)


      // 将order_no对应的用户优惠券删除掉
      delete_counpon_by_order(imes.Out_trade_no)

      // 对用户购票计数，每十次送一次优惠
      statistics_ticket_by_user(imes.Out_trade_no)
      
    } else if order_prefix == "MULTI" {
      // 包车租金支付
      stmtIns, err := db.Prepare("update YUECHE_MULTI set multi_status=? where multi_no=?")
      if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
      }
      defer stmtIns.Close() 

      _, err = stmtIns.Exec("已付款", imes.Out_trade_no)
      if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
      }

      var phone string
      var desp string
      get_start_end_by_Multi(imes.Out_trade_no, &phone, &desp)
      send_sms_order_multi(phone, imes.Out_trade_no, desp)
    }


    // 返回成功信息给微信服务器
    
    v.Return_code = "SUCCESS"

    output, err := xml.MarshalIndent(&v, "", "\t")
    if err != nil {
      fmt.Printf("error: %v\n", err)
    }

    fmt.Fprintf(w, "%s", output)
    return   
  }
}


// 获取所有订单（已完成，未完成，已退票，已改签，已取消）
type Manage_GetUserOrderHInput struct {
  Orderstatus string
}

type Manage_GetUserOrderRet struct {
  Ret string
  Message string
  Data struct {
    Order []OrderInfo
    Limit string
    Offset string
    LCount string
  }
}

func getevaluatebyorder(Order_no string, Evaluate *string) {
  querystr := fmt.Sprintf("select desp from YUECHE_ORDER_EVALUATE where order_no = '%s'", Order_no)
  rows,err := db.Query(querystr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(Evaluate)
    if err != nil {
      panic(err.Error())
      return
    }
  }
}

//http.HandleFunc("/api/manage_getuserorder", manage_getuserorderHandle)
func manage_getuserorderHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getuserorderHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetUserOrderRet
    var orderstatus_in string = "0"
    var limit_in string = "0"
    var offset_in string = "0"

    if len(r.Form["Orderstatus"]) > 0 {
      orderstatus_in = r.Form["Orderstatus"][0]
      fmt.Println(r.Form["Orderstatus"][0])
    }

    if len(r.Form["Limit"]) > 0 {
      limit_in = r.Form["Limit"][0]
      fmt.Println(r.Form["Limit"][0])
    }

    if len(r.Form["Offset"]) > 0 {
      offset_in = r.Form["Offset"][0]
      fmt.Println(r.Form["Offset"][0])
    }

    var cinfo_slice []OrderInfo

    // 取得符合条件的订单号
    querystr := fmt.Sprintf("select order_no,payprice from YUECHE_PAYMENTS where orderstatus = '%s' order by order_no desc limit %s,%s", orderstatus_in, offset_in, limit_in)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var order_no string
    var payprice string
    for rows.Next() {
    err = rows.Scan(&order_no, &payprice)
      if err != nil {
        panic(err.Error())
        return
      }

      // 根据订单号取得对应的订单详细信息
      querystr1 := fmt.Sprintf("select order_no,ticket_no,user_no,line_no,amount,aboard_time,aboard_local_name,end_local_name,driver_no,createtime from YUECHE_ORDER where order_no = '%s'", order_no)
      rows1,err1 := db.Query(querystr1)
      if err1 != nil {
        log.Fatal(err1)
        return
      }
      defer rows1.Close()

      var cinfo OrderInfo
      
      for rows1.Next() {
        err1 = rows1.Scan(&cinfo.Order_no, &cinfo.Ticket_no, &cinfo.User_no, &cinfo.Line_no, &cinfo.Amount, &cinfo.Aboard_time, &cinfo.Aboard_local_name, &cinfo.End_local_name, &cinfo.Driver_no, &cinfo.Createtime)
        if err1 != nil {
          panic(err1.Error())
          return
        }

        getUserinfobyid(cinfo.User_no, &cinfo.User_name, &cinfo.User_phone)
        getLineinfobyid(cinfo.Line_no, &cinfo.S_city, &cinfo.E_city, &cinfo.Price_one, &cinfo.Km_tot)
        getDriverinfobyid(cinfo.Driver_no, &cinfo.Driver_name, &cinfo.Driver_phone, &cinfo.Car_license)
        cinfo.Payprice = payprice
        cinfo.Orderstatus = orderstatus_in

        cinfo.Evaluate = "无评价"
        if orderstatus_in == "已评价" {
          getevaluatebyorder(cinfo.Order_no, &cinfo.Evaluate)
        }

        cinfo_slice = append(cinfo_slice, cinfo)
      }
    }



    var rcount string = "0"
    querystr_count := fmt.Sprintf("select COUNT(*) from YUECHE_PAYMENTS where orderstatus = '%s'", orderstatus_in)
    rows,err = db.Query(querystr_count)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&rcount)
      if err != nil {
        panic(err.Error())
        return
      }
    }

    retmessage.Data.LCount = rcount
    retmessage.Data.Limit = limit_in
    retmessage.Data.Offset = offset_in

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Order = cinfo_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}



// 获取所有用户信息
type Userinfo struct {
  User_no string
  Name string
  Phone string
  Openid string
  Headimgurl string
  Nickname string
  Sex string
  Province string
  City string
  Country string
  Ticket_num string
 //Coupon_no []string
}

type Manage_GetUserinfoRet struct {
  Ret string
  Message string
  Data struct {
    Info_slice []Userinfo
    Limit string
    Offset string
    LCount string
  }
}

func get_weico_info_by_id(User_no string, Openid *string, Headimgurl *string, Nickname *string, Sex *string, Province *string, City *string, Country *string) {
    
    querystr := fmt.Sprintf("select openid,headimgurl,nickname,sex,province,city,country from YUECHE_USERS_BIND_WECHAT where user_no = '%s'", User_no)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(Openid, Headimgurl, Nickname, Sex, Province, City, Country)
      if err != nil {
        panic(err.Error())
        return
      }
    }
}

func get_weico_info_by_driver(Driver_no string, Openid *string, Headimgurl *string, Nickname *string, Sex *string, Province *string, City *string, Country *string) {
    
    querystr := fmt.Sprintf("select openid,headimgurl,nickname,sex,province,city,country from YUECHE_DRIVER_BIND_WECHAT where driver_no = '%s'", Driver_no)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(Openid, Headimgurl, Nickname, Sex, Province, City, Country)
      if err != nil {
        panic(err.Error())
        return
      }
    }
}

//http.HandleFunc("/api/manage_getuserinfo", manage_getuserinfoHandle)
func manage_getuserinfoHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getuserinfoHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetUserinfoRet

    var info Userinfo
    var info_slice []Userinfo

    var limit_in string = "0"
    var offset_in string = "0"

    if len(r.Form["Limit"]) > 0 {
      limit_in = r.Form["Limit"][0]
      fmt.Println(r.Form["Limit"][0])
    }

    if len(r.Form["Offset"]) > 0 {
      offset_in = r.Form["Offset"][0]
      fmt.Println(r.Form["Offset"][0])
    }

    //querystr := "select user_no,name,phone from YUECHE_USER"
    querystr := fmt.Sprintf("select user_no,name,phone from YUECHE_USER limit %s,%s", offset_in, limit_in)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
      err = rows.Scan(&info.User_no, &info.Name, &info.Phone)
      if err != nil {
        panic(err.Error())
        return
      }
      get_weico_info_by_id(info.User_no, &info.Openid, &info.Headimgurl, &info.Nickname, &info.Sex, &info.Province, &info.City, &info.Country)

      getticketbyuserno(info.User_no, &info.Ticket_num)  // 通过手机号获取用户购买了多少张票

      info_slice = append(info_slice, info)
    }


    var rcount string = "0"
    querystr_count := "select COUNT(*) from YUECHE_USER"
    rows,err = db.Query(querystr_count)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&rcount)
      if err != nil {
        panic(err.Error())
        return
      }
    }

    retmessage.Data.LCount = rcount
    retmessage.Data.Limit = limit_in
    retmessage.Data.Offset = offset_in

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Info_slice = info_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


// 根据用户手机号获取对应用户信息

// http.HandleFunc("/api/manage_getuserinfobyphone", manage_getuserinfobyphoneHandle)
func manage_getuserinfobyphoneHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getuserinfobyphoneHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetUserinfoRet

    var info Userinfo
    var info_slice []Userinfo

    var Phone_in string = "0"

     if len(r.Form["Phone"]) > 0 {
      Phone_in = r.Form["Phone"][0]
      fmt.Println(r.Form["Phone"][0])
    }

    //querystr := "select user_no,name,phone from YUECHE_USER"
    querystr := fmt.Sprintf("select user_no,name,phone from YUECHE_USER where phone='%s'", Phone_in)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&info.User_no, &info.Name, &info.Phone)
      if err != nil {
        panic(err.Error())
        return
      }
      get_weico_info_by_id(info.User_no, &info.Openid, &info.Headimgurl, &info.Nickname, &info.Sex, &info.Province, &info.City, &info.Country)

      getticketbyuserno(info.User_no, &info.Ticket_num)  // 通过手机号获取用户购买了多少张票

      info_slice = append(info_slice, info)
    }

    retmessage.Data.LCount = "0"
    retmessage.Data.Limit = "0"
    retmessage.Data.Offset = "0"

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Info_slice = info_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}



// 获取所有司机信息
type Driverinfo struct {
  Driver_no string
  Name string
  Phone string
  Idcard string
  Car_license string
  Longitude string
  Latitude string
  D_status string
  Openid string
  Headimgurl string
  Nickname string
  Sex string
  Province string
  City string
  Country string
}

type Manage_GetDriverinfoRet struct {
  Ret string
  Message string
  Data struct {
    Info_slice []Driverinfo
    Limit string
    Offset string
    LCount string
  }
}

//http.HandleFunc("/api/manage_getdriverinfo", manage_getdriverinfoHandle)
func manage_getdriverinfoHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getdriverinfoHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetDriverinfoRet

    var info Driverinfo
    var info_slice []Driverinfo

    var limit_in string = "0"
    var offset_in string = "0"

    if len(r.Form["Limit"]) > 0 {
      limit_in = r.Form["Limit"][0]
      fmt.Println(r.Form["Limit"][0])
    }

    if len(r.Form["Offset"]) > 0 {
      offset_in = r.Form["Offset"][0]
      fmt.Println(r.Form["Offset"][0])
    }

    //querystr := "select driver_no,name,phone,idcard,car_license,longitude,latitude from YUECHE_DRIVER"
    querystr := fmt.Sprintf("select driver_no,name,phone,idcard,car_license,longitude,latitude,d_status from YUECHE_DRIVER limit %s,%s", offset_in, limit_in)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&info.Driver_no, &info.Name, &info.Phone, &info.Idcard, &info.Car_license, &info.Longitude, &info.Latitude, &info.D_status)
      if err != nil {
        panic(err.Error())
        return
      }
      get_weico_info_by_driver(info.Driver_no, &info.Openid, &info.Headimgurl, &info.Nickname, &info.Sex, &info.Province, &info.City, &info.Country)
      info_slice = append(info_slice, info)
    }

    var rcount string = "0"
    querystr_count := "select COUNT(*) from YUECHE_DRIVER"
    rows,err = db.Query(querystr_count)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&rcount)
      if err != nil {
        panic(err.Error())
        return
      }
    }

    retmessage.Data.LCount = rcount
    retmessage.Data.Limit = limit_in
    retmessage.Data.Offset = offset_in

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Info_slice = info_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}



// 司机上下线
type Manage_UpdriverstatusRet struct {
  Ret string
  Message string
}

//http.HandleFunc("/api/manage_updriverstatus", manage_updriverstatusHandle)
func manage_updriverstatusHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_updriverstatusHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_UpdriverstatusRet


    var Status_in string = "0"
    var Driver_no_in string = "0"

    if len(r.Form["Driver_no"]) > 0 {
      Driver_no_in = r.Form["Driver_no"][0]
      //fmt.Println(r.Form["Driver_no"][0])
    }

    if len(r.Form["Status"]) > 0 {
      Status_in = r.Form["Status"][0]
     // fmt.Println(r.Form["Status"][0])
    }

    stmtIns, err := db.Prepare("update YUECHE_DRIVER set d_status=? where driver_no=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(Status_in, Driver_no_in)
    if err != nil {
      panic(err.Error())
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


// 获取所有累计收入，累计订票数量，当前在线车辆数
type Manage_GetsystemdataRet struct {
  Ret string
  Message string
  Data struct {
    Income_money string   // 累计收入（单位分）
    Ticket_amount string  // 累计订票数量
    Online_car string  // 当前在线车辆数
  }
}

func get_income_money(money *string) {
    querystr := "select price from YUECHE_PAYMENTS where paystatus='success'"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var price string = "0"
    var amount int = 0
    for rows.Next() {
      err = rows.Scan(&price)
      if err != nil {
        panic(err.Error())
        return
      }

      num, _ := strconv.Atoi(price)
      amount = amount + num
    }

    *money = strconv.Itoa(int(amount))
}

func get_ticket_amount(Ticket_amount *string) {
    querystr := "select amount from YUECHE_ORDER where description='已付款'"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var amount int = 0
    var ticket_amount int = 0
    for rows.Next() {
      err = rows.Scan(&amount)
      if err != nil {
        panic(err.Error())
        return
      }

      ticket_amount = ticket_amount + amount
    }

    *Ticket_amount = strconv.Itoa(int(ticket_amount))
}

func online_car(on_car *string) {

    timeStr:=time.Now().Format("2006-01-02 15:04:05")
    t1, _ := time.Parse("2006-01-02 15:04:05", timeStr)
    timeStr1 := t1.Format("2006-01-02")
    timeStr2 := timeStr1 + " 23:59:59"

    querystr := fmt.Sprintf("SELECT COUNT(*) FROM YUECHE_DISPATCH where (on_time < '%s' and on_time > '%s') and off_time='2018-01-01 00:00:00'", timeStr2, timeStr1)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(on_car)
      if err != nil {
        panic(err.Error())
        return
      }
    }
}

//http.HandleFunc("/api/manage_getsystemdata", manage_getsystemdataHandle)
func manage_getsystemdataHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getsystemdataHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetsystemdataRet

    retmessage.Data.Income_money = "0"
    retmessage.Data.Ticket_amount = "0"
    retmessage.Data.Online_car = "0"
    get_income_money(&retmessage.Data.Income_money)
    get_ticket_amount(&retmessage.Data.Ticket_amount)
    online_car(&retmessage.Data.Online_car)

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

// 根据时间查询当天的任意站点的出票信息（订单号，票号，用户名称，用户电话等）
type Manage_GetorderbydaytimeRet struct {
  Ret string
  Message string
  Data struct {
    Order_slice []OrderInfo
  }
}

//http.HandleFunc("/api/manage_getorderbydaytime", manage_getorderbydaytimeHandle)
func manage_getorderbydaytimeHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getorderbydaytimeHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetorderbydaytimeRet

/*
    var limit_in string = "0"
    var offset_in string = "0"

    if len(r.Form["Limit"]) > 0 {
      limit_in = r.Form["Limit"][0]
      fmt.Println(r.Form["Limit"][0])
    }

    if len(r.Form["Offset"]) > 0 {
      offset_in = r.Form["Offset"][0]
      fmt.Println(r.Form["Offset"][0])
    }
*/

    var Allot_time_in string = "0"
    if len(r.Form["Allot_time"]) > 0 {
      Allot_time_in = r.Form["Allot_time"][0]
    }

    t := time.Now()
    timeStr := t.Format("2006-01-02 15:04:05")
    if Allot_time_in != "0" {
      timeStr = Allot_time_in
    }

    //timeStr := time.Now()

   // querystr := fmt.Sprintf("SELECT user_no,order_no,ticket_no, line_no,amount,aboard_local_name,end_local_name,aboard_time FROM YUECHE_ORDER where user_no = '%s' and ride_time_s < '%s' and ride_time_e > '%s' and description <> '已取消' and description <> '已上车' and description <> '未付款' order by order_no desc", input.User_no, timeStr, timeStr)

    var cinfo_slice []OrderInfo



    // 根据订单号取得对应的订单详细信息
    querystr1 := fmt.Sprintf("select order_no,ticket_no,user_no,line_no,amount,aboard_time,aboard_local_name,end_local_name,driver_no,createtime from YUECHE_ORDER where  (ride_time_s < '%s' and ride_time_e > '%s') and description <> '已取消' and description <> '已上车' and description <> '未付款' and description <> '已分配司机'",  timeStr, timeStr)
    rows1,err1 := db.Query(querystr1)
    if err1 != nil {
      log.Fatal(err1)
      return
    }
    defer rows1.Close()

    var cinfo OrderInfo
    
    for rows1.Next() {
      err1 = rows1.Scan(&cinfo.Order_no, &cinfo.Ticket_no, &cinfo.User_no, &cinfo.Line_no, &cinfo.Amount, &cinfo.Aboard_time, &cinfo.Aboard_local_name, &cinfo.End_local_name, &cinfo.Driver_no, &cinfo.Createtime)
      if err1 != nil {
        panic(err1.Error())
        return
      }

      // 取得符合条件的订单号
      querystr := fmt.Sprintf("select order_no,payprice from YUECHE_PAYMENTS where (orderstatus = '未完成' or orderstatus = '已改签') and paystatus='success' and order_no='%s' order by order_no desc", cinfo.Order_no)
      rows,err := db.Query(querystr)
      if err != nil {
        log.Fatal(err)
        return
      }
      defer rows.Close()

      var order_no string
      var payprice string
      for rows.Next() {
        err = rows.Scan(&order_no, &payprice)
        if err != nil {
          panic(err.Error())
          return
        }

        getLineinfobyid(cinfo.Line_no, &cinfo.S_city, &cinfo.E_city, &cinfo.Price_one, &cinfo.Km_tot)
        getDriverinfobyid(cinfo.Driver_no, &cinfo.Driver_name, &cinfo.Driver_phone, &cinfo.Car_license)
        cinfo.Payprice = payprice
        cinfo.Orderstatus = "未完成"

       // var aboard_time string
      //  var start_city string
      //  var end_city string
        //var phone string
       // var ticket_no string
        //get_start_end_by_Order(cinfo.Order_no, &cinfo.S_city,&cinfo.E_city, &aboard_time, &cinfo.User_phone, &ticket_no)


        querystr2 := fmt.Sprintf("SELECT phone FROM YUECHE_USER where user_no = '%s'", cinfo.User_no)
        rows2,err2 := db.Query(querystr2)
        if err2 != nil {
          log.Fatal(err2)
          return
        }
        defer rows2.Close()

        for rows2.Next() {
        err2 = rows2.Scan(&cinfo.User_phone)
          if err2 != nil {
            panic(err2.Error())
            return
          }
        }
        cinfo_slice = append(cinfo_slice, cinfo)
      }
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Order_slice = cinfo_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


// 后台管理员登录
type Manage_AdminloginRet struct {
  Ret string
  Message string
  Data struct {
    Admin_no string
    Token string
  }
}

//http.HandleFunc("/api/manage_adminlogin", manage_adminloginHandle)
func manage_adminloginHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_adminloginHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_AdminloginRet


    var name_in string = "0"
    var pwssword_in string = "0"

    if len(r.Form["Name"]) > 0 {
      name_in = r.Form["Name"][0]
      //fmt.Println(r.Form["Name"][0])
    }

    if len(r.Form["Password"]) > 0 {
      pwssword_in = r.Form["Password"][0]
      //fmt.Println(r.Form["Password"][0])
    }

    var rcount string = "0"
    querystr := fmt.Sprintf("SELECT COUNT(*) FROM YUECHE_ADMIN where admin_name = '%s' and admin_password = '%s'", name_in, pwssword_in)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&rcount)
      if err != nil {
        panic(err.Error())
        return
      }
    }


    if rcount == "0" {
      retmessage.Ret = "1"
      retmessage.Message = "用户名或者密码错误，请重新输入！"  
    } else {

      // 将token改变
      token := GetRandomSalt()

      var admin_no string = "0"
      querystr := fmt.Sprintf("SELECT admin_id FROM YUECHE_ADMIN where admin_name = '%s' and admin_password = '%s'", name_in, pwssword_in)
      rows,err := db.Query(querystr)
      if err != nil {
        log.Fatal(err)
        return
      }
      defer rows.Close()

      for rows.Next() {
      err = rows.Scan(&admin_no)
        if err != nil {
          panic(err.Error())
          return
        }
      }

      stmtIns, err := db.Prepare("update YUECHE_ADMIN set token=?,login_time=? where admin_id=?")
      if err != nil {
          panic(err.Error())
      }
      defer stmtIns.Close() 

      _, err = stmtIns.Exec(token, time.Now(), admin_no)
      if err != nil {
        panic(err.Error())
      }

      retmessage.Ret = "0"
      retmessage.Message = "success" 
      retmessage.Data.Admin_no = admin_no
      retmessage.Data.Token = token
    }

    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

// 发送乘客票务信息给司机
type Manage_SendsmstodriverRet struct {
  Ret string
  Message string
}

func get_orderinfo_by_id(Order_no string, User_phone *string, Ticket_no *string, S_city *string, E_city *string, Aboard_local_name *string, End_local_name *string, Aboard_time *string, Amount *string) {

  // 根据订单号取得对应的订单详细信息
  querystr1 := fmt.Sprintf("select amount,aboard_local_name,end_local_name from YUECHE_ORDER where order_no = '%s'", Order_no)
  rows1,err1 := db.Query(querystr1)
  if err1 != nil {
    log.Fatal(err1)
    return
  }
  defer rows1.Close()

  
  for rows1.Next() {
    err1 = rows1.Scan(Amount, Aboard_local_name, End_local_name)
    if err1 != nil {
      panic(err1.Error())
      return
    }

    get_start_end_by_Order(Order_no, S_city, E_city, Aboard_time, User_phone, Ticket_no)
  }
}

// http.HandleFunc("/api/manage_sendsmstodriver", manage_sendsmstodriverHandle)
func manage_sendsmstodriverHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_sendsmstodriverHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_SendsmstodriverRet

    var orderno_in string = "0"
    if len(r.Form["Order_no"]) > 0 {
      orderno_in = r.Form["Order_no"][0]
      //fmt.Println(r.Form["Name"][0])
    }

    var driverphone_in string = "0"
    if len(r.Form["Driverphone"]) > 0 {
      driverphone_in = r.Form["Driverphone"][0]
      //fmt.Println(r.Form["Name"][0])
    }

    var driverno_in string = "0"
    if len(r.Form["Driver_no"]) > 0 {
      driverno_in = r.Form["Driver_no"][0]
      //fmt.Println(r.Form["Name"][0])
    }

    orderno_slice := strings.Split(orderno_in, ",")
    for _, value := range orderno_slice {
      var User_phone string
      var Ticket_no string
      var S_city string
      var E_city string
      var Aboard_local_name string
      var End_local_name string
      var Aboard_time string
      var Amount string
      get_orderinfo_by_id(value, &User_phone, &Ticket_no, &S_city, &E_city, &Aboard_local_name, &End_local_name, &Aboard_time, &Amount)

      s1 := Aboard_time[0:10]
      s2 := Aboard_time[11:19]
      s3 := s1 + " " + s2

      //messagexsend
      messageconfig := make(map[string]string)
      messageconfig["appid"] = "21084"
      messageconfig["appkey"] = "323cdfa1b38815aaecacea37b3369938"
      messageconfig["signtype"] = "md5"

      messagexsend := submail.CreateMessageXSend()
      submail.MessageXSendAddTo(messagexsend, driverphone_in)
      submail.MessageXSendSetProject(messagexsend, "eq7g13")
      submail.MessageXSendAddVar(messagexsend, "Orderno", value)
      submail.MessageXSendAddVar(messagexsend, "User_phone", User_phone)
      submail.MessageXSendAddVar(messagexsend, "Ticket_no", Ticket_no)
      submail.MessageXSendAddVar(messagexsend, "S_city", S_city)
      submail.MessageXSendAddVar(messagexsend, "E_city", E_city)
      submail.MessageXSendAddVar(messagexsend, "Aboard_local_name", Aboard_local_name)
      //submail.MessageXSendAddVar(messagexsend, "End_local_name", End_local_name)
      submail.MessageXSendAddVar(messagexsend, "Aboard_time", s3)
      submail.MessageXSendAddVar(messagexsend, "Amount", Amount)
      fmt.Println("MessageXSend ", submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig))


      // 将订单表中的 description 设置为"已分配司机"状态
      stmtIns, err := db.Prepare("update YUECHE_ORDER set driver_no=?,description=? where order_no=?")
      if err != nil {
          panic(err.Error())
      }
      defer stmtIns.Close() 

      _, err = stmtIns.Exec(driverno_in, "已分配司机",value)
      if err != nil {
        panic(err.Error())
      }
    }


    // 订单@var(Orderno)需要您的处理，乘客电话[@var(User_phone)],车票号[@var(Ticket_no)],由@var(S_city)到@var(E_city)，上车点为@var(Aboard_local_name)，上车时间为@var(Aboard_time)，共购票@var(Amount)张。请及时联系乘客，谢谢！
    retmessage.Ret = "0"
    retmessage.Message = "success" 
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}



// 获取所有司机相关信息，id，phone，name
type DriverList struct {
  Driver_no string
  Name string
  Phone string
}

type Manage_GetDriverListRet struct {
  Ret string
  Message string
  Data struct {
    Info_slice []DriverList
  }
}

// http.HandleFunc("/api/manage_getdriverlist", manage_getdriverlistHandle)
func manage_getdriverlistHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getdriverlistHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetDriverListRet

    var info DriverList
    var info_slice []DriverList


    //querystr := "select driver_no,name,phone,idcard,car_license,longitude,latitude from YUECHE_DRIVER"
    querystr := "select driver_no,name,phone from YUECHE_DRIVER where d_status = 'on'"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
    err = rows.Scan(&info.Driver_no, &info.Name, &info.Phone)
      if err != nil {
        panic(err.Error())
        return
      }
      info_slice = append(info_slice, info)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Info_slice = info_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

// 设置发车时间
type Manage_setbustimeRet struct {
  Ret string
  Message string
}

//http.HandleFunc("/api/manage_setbustime", manage_setbustimeHandle)
func manage_setbustimeHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_setbustimeHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_setbustimeRet

    var Line_type_in string = "0"
    if len(r.Form["Line_type"]) > 0 {
      Line_type_in = r.Form["Line_type"][0]
    }

    var Bus_time_in string = "0"
    if len(r.Form["Bus_time"]) > 0 {
      Bus_time_in = r.Form["Bus_time"][0]
    }

    t := time.Now().Unix()
    //num, _ := strconv.Atoi(t)

    if (Line_type_in == "air" || Line_type_in == "common") && Bus_time_in != "0" {
      Bus_time_slice := strings.Split(Bus_time_in, ",")
      for _, value := range Bus_time_slice {
        t = t + 1
        timeid := fmt.Sprintf("TIME_%d", t)

        if Line_type_in == "air" {
          stmtIns, err := db.Prepare("INSERT INTO YUECHE_BUS_TIME VALUES( ?,?,?)")
          if err != nil {
              panic(err.Error()) // proper error handling instead of panic in your app
          }
          defer stmtIns.Close() 

          _, err = stmtIns.Exec(timeid, value, "机场专线")
          if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
          }
        } else {
          stmtIns, err := db.Prepare("INSERT INTO YUECHE_BUS_TIME VALUES( ?,?,?)")
          if err != nil {
              panic(err.Error()) // proper error handling instead of panic in your app
          }
          defer stmtIns.Close() 

          _, err = stmtIns.Exec(timeid, value, "一般线路")
          if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
          }
        }
      }
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


// 获取发车时间
type Bustime_info struct {
  Timeid string
  Bus_time string
  Line_name string
}

type Manage_getbustimeRet struct {
  Ret string
  Message string
  Data struct {
    Bus_time_slice []Bustime_info
  }
}

//http.HandleFunc("/api/manage_getbustime", manage_getbustimeHandle) 
func manage_getbustimeHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getbustimeHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_getbustimeRet

    var Line_name_in string = "0"
    if len(r.Form["Line_name"]) > 0 {
      Line_name_in = r.Form["Line_name"][0]
    }

    var bus Bustime_info
    var bus_slice []Bustime_info 

    if Line_name_in != "0" {
      querystr := fmt.Sprintf("select timeid,bus_time,line_name from YUECHE_BUS_TIME where line_name='%s' order by bus_time", Line_name_in)
      rows,err := db.Query(querystr)
      if err != nil {
        log.Fatal(err)
        return
      }
      defer rows.Close()

      for rows.Next() {
        err = rows.Scan(&bus.Timeid, &bus.Bus_time, &bus.Line_name)
        if err != nil {
          panic(err.Error())
          return
        }
        bus_slice = append(bus_slice, bus)
      }
    }



    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Bus_time_slice = bus_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}



// 根据id号删除某一条发车时间
type Manage_delbustimeRet struct {
  Ret string
  Message string
}

// http.HandleFunc("/api/manage_delbustime", manage_delbustimeHandle) 
func manage_delbustimeHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_delbustimeHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_delbustimeRet

    var Timeid_in string = "0"
    if len(r.Form["Timeid"]) > 0 {
      Timeid_in = r.Form["Timeid"][0]
    }


    stmtIns, err := db.Prepare("DELETE FROM YUECHE_BUS_TIME WHERE timeid =?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(Timeid_in)
    if err != nil {
      panic(err.Error())
    }


    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


// 获取当天已分配司机的订单
type Manage_GetallotorderRet struct {
  Ret string
  Message string
  Data struct {
    Order_slice []OrderInfo
  }
}

// http.HandleFunc("/api/manage_getallotorder", manage_getallotorderHandle) 
func manage_getallotorderHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getallotorderHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetallotorderRet

    var Allot_time_in string = "0"
    if len(r.Form["Allot_time"]) > 0 {
      Allot_time_in = r.Form["Allot_time"][0]
    }

    t := time.Now()
    timeStr := t.Format("2006-01-02 15:04:05")
    if Allot_time_in != "0" {
      timeStr = Allot_time_in
    }

    var cinfo_slice []OrderInfo

    // 根据订单号取得对应的订单详细信息
    querystr1 := fmt.Sprintf("select order_no,ticket_no,user_no,line_no,amount,aboard_time,aboard_local_name,end_local_name,driver_no,createtime from YUECHE_ORDER where  (ride_time_s < '%s' and ride_time_e > '%s') and description = '已分配司机'",  timeStr, timeStr)
    rows1,err1 := db.Query(querystr1)
    if err1 != nil {
      log.Fatal(err1)
      return
    }
    defer rows1.Close()

    var cinfo OrderInfo
    
    for rows1.Next() {
      err1 = rows1.Scan(&cinfo.Order_no, &cinfo.Ticket_no, &cinfo.User_no, &cinfo.Line_no, &cinfo.Amount, &cinfo.Aboard_time, &cinfo.Aboard_local_name, &cinfo.End_local_name, &cinfo.Driver_no, &cinfo.Createtime)
      if err1 != nil {
        panic(err1.Error())
        return
      }

      // 取得符合条件的订单号
      querystr := fmt.Sprintf("select order_no,payprice from YUECHE_PAYMENTS where (orderstatus = '未完成' or orderstatus = '已改签') and paystatus='success' and order_no='%s' order by order_no desc", cinfo.Order_no)
      rows,err := db.Query(querystr)
      if err != nil {
        log.Fatal(err)
        return
      }
      defer rows.Close()

      var order_no string
      var payprice string
      for rows.Next() {
        err = rows.Scan(&order_no, &payprice)
        if err != nil {
          panic(err.Error())
          return
        }

        getLineinfobyid(cinfo.Line_no, &cinfo.S_city, &cinfo.E_city, &cinfo.Price_one, &cinfo.Km_tot)
        getDriverinfobyid(cinfo.Driver_no, &cinfo.Driver_name, &cinfo.Driver_phone, &cinfo.Car_license)
        cinfo.Payprice = payprice
        cinfo.Orderstatus = "未完成"

       // var aboard_time string
      //  var start_city string
      //  var end_city string
        //var phone string
       // var ticket_no string
        //get_start_end_by_Order(cinfo.Order_no, &cinfo.S_city,&cinfo.E_city, &aboard_time, &cinfo.User_phone, &ticket_no)


        querystr2 := fmt.Sprintf("SELECT phone FROM YUECHE_USER where user_no = '%s'", cinfo.User_no)
        rows2,err2 := db.Query(querystr2)
        if err2 != nil {
          log.Fatal(err2)
          return
        }
        defer rows2.Close()

        for rows2.Next() {
        err2 = rows2.Scan(&cinfo.User_phone)
          if err2 != nil {
            panic(err2.Error())
            return
          }
        }
        cinfo_slice = append(cinfo_slice, cinfo)
      }
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Order_slice = cinfo_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


// 根据line_type获取区域名称
type Manage_GetareabylinetypeRet struct {
  Ret string
  Message string
  Data struct {
    Area_name []string
  }
}

//http.HandleFunc("/api/manage_getareabylinetype", manage_getareabylinetypeHandle) 
func manage_getareabylinetypeHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getareabylinetypeHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetareabylinetypeRet

    var Line_type_in string = "0"
    if len(r.Form["Line_type"]) > 0 {
      Line_type_in = r.Form["Line_type"][0]
    }

    querystr := fmt.Sprintf("select distinct description from YUECHE_BUS_LINE_VIA WHERE line_type='%s'", Line_type_in)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var des string
    var des_slice []string
    for rows.Next() {
      err = rows.Scan(&des)
      if err != nil {
        panic(err.Error())
        return
      }
      des_slice = append(des_slice, des)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Area_name = des_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}



// 站点上下线
type Manage_UplocalRet struct {
  Ret string
  Message string
}

//http.HandleFunc("/api/manage_uplocal", manage_uplocalHandle) 
func manage_uplocalHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_uplocalHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_UplocalRet

    var Via_local_num_in string = "0"
    if len(r.Form["Via_local_num"]) > 0 {
      Via_local_num_in = r.Form["Via_local_num"][0]
    }

    var Local_status_in string = "0"
    if len(r.Form["Local_status"]) > 0 {
      Local_status_in = r.Form["Local_status"][0]
    }

    stmtIns, err := db.Prepare("update YUECHE_BUS_LINE_VIA set local_status=? where via_local_num=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(Local_status_in, Via_local_num_in)
    if err != nil {
      panic(err.Error())
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}



// 新增站点
type Manage_InsertlocalRet struct {
  Ret string
  Message string
  Data struct {
    Via_local_num string
  }
}

//http.HandleFunc("/api/manage_insertlocal", manage_insertlocalHandle)   
func manage_insertlocalHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_insertlocalHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_InsertlocalRet

    var Line_type_in string = "0"
    if len(r.Form["Line_type"]) > 0 {
      Line_type_in = r.Form["Line_type"][0]
    }

    var Via_local_name_in string = "0"
    if len(r.Form["Via_local_name"]) > 0 {
      Via_local_name_in = r.Form["Via_local_name"][0]
    }

    var Longitude_in string = "0"
    if len(r.Form["Longitude"]) > 0 {
      Longitude_in = r.Form["Longitude"][0]
    }

    var Latitude_in string = "0"
    if len(r.Form["Latitude"]) > 0 {
      Latitude_in = r.Form["Latitude"][0]
    }

    var Description_in string = "0"
    if len(r.Form["Description"]) > 0 {
      Description_in = r.Form["Description"][0]
    }

    var line_no string = "0"
    if Line_type_in == "common" {
      line_no = "LINE_001"
    } else {
      line_no = "LINE_010"
    }
    var via_local_num string = "0"
    t := time.Now().Unix()
    via_local_num = fmt.Sprintf("V%d", t)


    stmtIns, err := db.Prepare("INSERT INTO YUECHE_BUS_LINE_VIA VALUES( ?,?,?,?,?,?,?,?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(line_no, Line_type_in, Via_local_name_in, via_local_num, Longitude_in, Latitude_in, Description_in, "on")
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Via_local_num = via_local_num
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


// 获取当前所有站点
type AllLocalInfo struct {
  Line_type string    // 站点类型
  Via_local_name string  // 站点名称
  Via_local_num string   // 站点编号
  Longitude string
  Latitude string
  Description string    // 所属区域
  Local_status string   // 站点状态
}

type Manage_GetalllocalRet struct {
  Ret string
  Message string
  Data struct {
    Local_slice []AllLocalInfo
  }
}

//http.HandleFunc("/api/manage_getalllocal", manage_getalllocalHandle)   
func manage_getalllocalHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getalllocalHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetalllocalRet

    querystr := "select line_type,via_local_name,via_local_num,longitude,latitude,description,local_status from YUECHE_BUS_LINE_VIA"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var local AllLocalInfo
    var local_slice []AllLocalInfo
    for rows.Next() {
      err = rows.Scan(&local.Line_type, &local.Via_local_name, &local.Via_local_num, &local.Longitude, &local.Latitude, &local.Description, &local.Local_status)
      if err != nil {
        panic(err.Error())
        return
      }
      local_slice = append(local_slice, local)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Local_slice = local_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


// 新增区域
type Manage_InsertareaRet struct {
  Ret string
  Message string
}

//http.HandleFunc("/api/manage_insertarea", manage_insertareaHandle)   
func manage_insertareaHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_insertareaHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_InsertareaRet

    var City_in string = "0"
    if len(r.Form["City"]) > 0 {
      City_in = r.Form["City"][0]
    }

    var Area_in string = "0"
    if len(r.Form["Area"]) > 0 {
      Area_in = r.Form["Area"][0]
    }

    stmtIns, err := db.Prepare("INSERT INTO YUECHE_BUS_AREA VALUES( ?,?)")
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(City_in, Area_in)
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}



// 获取当前所有区域
type AllAreaInfo struct {
  City string
  Area string
}

type Manage_GetallareaRet struct {
  Ret string
  Message string
  Data struct {
    Area_slice []AllAreaInfo
  }
}

//http.HandleFunc("/api/manage_getallarea", manage_getallareaHandle)  
func manage_getallareaHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getallareaHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetallareaRet

    querystr := "select city,area from YUECHE_BUS_AREA"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var area AllAreaInfo
    var area_slice []AllAreaInfo
    for rows.Next() {
      err = rows.Scan(&area.City, &area.Area)
      if err != nil {
        panic(err.Error())
        return
      }
      area_slice = append(area_slice, area)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Area_slice = area_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


// 获取所有城市
type Manage_GetallcityRet struct {
  Ret string
  Message string
  Data struct {
    City_slice []string
  }
}

// http.HandleFunc("/api/manage_getallcity", manage_getallcityHandle)
func manage_getallcityHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getallcityHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetallcityRet
/*
    querystr := "select distinct city from YUECHE_BUS_AREA"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var city string
    var city_slice []string
    for rows.Next() {
      err = rows.Scan(&city)
      if err != nil {
        panic(err.Error())
        return
      }
      city_slice = append(city_slice, city)
    }
*/
    querystr := "select city from YUECHE_CITY"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var city string
    var city_slice []string
    for rows.Next() {
      err = rows.Scan(&city)
      if err != nil {
        panic(err.Error())
        return
      }
      city_slice = append(city_slice, city)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.City_slice = city_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


// 添加城市
type Manage_InsertcityRet struct {
  Ret string
  Message string
}

// http.HandleFunc("/api/manage_insertcity", manage_insertcityHandle) 
func manage_insertcityHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_insertcityHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_InsertcityRet

    var City_in string = "0"
    if len(r.Form["City"]) > 0 {
      City_in = r.Form["City"][0]
    }

    stmtIns, err := db.Prepare("INSERT INTO YUECHE_CITY VALUES(?)")
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(City_in)
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}

// 获取所有城市线路
type Lineinfomini struct {
  Line_no string
  S_city string  // 起点城市
  E_city string  // 终点城市
  Price string // 票价
}

type Manage_GetLineRet struct {
  Ret string
  Message string
  Data struct {
    Line_slice []Lineinfomini
  }
}

// http.HandleFunc("/api/manage_getline", manage_getlineHandle)
func manage_getlineHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getlineHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetLineRet

    querystr := "select line_no,s_city,e_city,price from YUECHE_BUS_LINE"
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    var cinfo Lineinfomini
    var cinfo_slice []Lineinfomini
    for rows.Next() {
    err = rows.Scan(&cinfo.Line_no, &cinfo.S_city, &cinfo.E_city, &cinfo.Price)
      if err != nil {
        panic(err.Error())
        return
      }
      cinfo_slice = append(cinfo_slice, cinfo)
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Line_slice = cinfo_slice
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}  


// 修改线路价格
type Manage_UpdatelinepriceRet struct {
  Ret string
  Message string
}

// http.HandleFunc("/api/manage_updatelineprice", manage_updatelinepriceHandle)  
func manage_updatelinepriceHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_updatelinepriceHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_UpdatelinepriceRet

    var Price_in string = "0"
    if len(r.Form["Price"]) > 0 {
      Price_in = r.Form["Price"][0]
    }

    var Line_no_in string = "0"
    if len(r.Form["Line_no"]) > 0 {
      Line_no_in = r.Form["Line_no"][0]
    }

    stmtIns, err := db.Prepare("update YUECHE_BUS_LINE set price=? where line_no=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec(Price_in, Line_no_in)
    if err != nil {
      panic(err.Error())
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
} 


// 根据用户手机号，获取当前用户历史上买了多少张票
type Manage_GetticketbyphoneRet struct {
  Ret string
  Message string
  Data struct {
    Ticket_amount string
  }
}

func getticketbyuserno(User_no string, Amount *string) {
    var last_num int = 0
    var amount_str string = "0"
    querystr := fmt.Sprintf("select amount from YUECHE_ORDER where user_no = '%s' and (description = '已付款' or description = '已上车')", User_no)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
      err = rows.Scan(&amount_str)
      if err != nil {
        panic(err.Error())
        return
      }
      amount_num, _ := strconv.Atoi(amount_str)
      last_num = last_num + amount_num
    }

  *Amount = strconv.Itoa(int(last_num))
}

// http.HandleFunc("/api/manage_getticketbyphone", manage_getticketbyphoneHandle) 
func manage_getticketbyphoneHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance manage_getticketbyphoneHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    var retmessage Manage_GetticketbyphoneRet
    
    var Phone_in string = "0"
    if len(r.Form["Phone"]) > 0 {
      Phone_in = r.Form["Phone"][0]
    }

    var user_no string = "0"
    var user_no_slice []string
    querystr := fmt.Sprintf("select user_no from YUECHE_USER where phone = '%s'", Phone_in)
    rows,err := db.Query(querystr)
    if err != nil {
      log.Fatal(err)
      return
    }
    defer rows.Close()

    for rows.Next() {
      err = rows.Scan(&user_no)
      if err != nil {
        panic(err.Error())
        return
      }
      user_no_slice = append(user_no_slice, user_no)
    }

    var last_num int = 0
    for i:=0;i<len(user_no_slice);i++ {

      var amount_str string = "0"
      querystr := fmt.Sprintf("select amount from YUECHE_ORDER where user_no = '%s' and (description = '已付款' or description = '已上车')", user_no_slice[i])
      rows,err := db.Query(querystr)
      if err != nil {
        log.Fatal(err)
        return
      }
      defer rows.Close()

      for rows.Next() {
        err = rows.Scan(&amount_str)
        if err != nil {
          panic(err.Error())
          return
        }
        amount_num, _ := strconv.Atoi(amount_str)
        last_num = last_num + amount_num
      }
    }

    retmessage.Ret = "0"
    retmessage.Message = "success"
    retmessage.Data.Ticket_amount = strconv.Itoa(int(last_num))
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}




// 后台退钱完成后，修改订单状态为"已完成"
type JSAPIWeixintuikuanInputMessage struct {
  Return_code string `xml:"return_code"` 
  Appid string `xml:"appid"` 
  Mch_id string `xml:"mch_id"` 
  Nonce_str string `xml:"nonce_str"`
  Req_info string `xml:"req_info"` 
}

type Req_info_Message struct {
  Out_refund_no string `xml:"out_refund_no"` 
  Out_trade_no string `xml:"out_trade_no"` 
  Refund_account string `xml:"refund_account"` 
  Refund_fee string `xml:"refund_fee"`
  Refund_id string `xml:"refund_id"` 
  Refund_recv_accout string `xml:"refund_recv_accout"` 
  Refund_request_source string `xml:"refund_request_source"` 
  Refund_status string `xml:"refund_status"` 
  Settlement_refund_fee string `xml:"settlement_refund_fee"` 
  Settlement_total_fee string `xml:"settlement_total_fee"` 
  Success_time string `xml:"success_time"` 
  Total_fee string `xml:"total_fee"` 
  Transaction_id string `xml:"transaction_id"` 
}

//http.HandleFunc("/api/weixin/wechat_ordercomplete", wechat_ordercompleteHandle)
func wechat_ordercompleteHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_ordercompleteHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {
    
    result, _:= ioutil.ReadAll(r.Body)
    var imes JSAPIWeixintuikuanInputMessage
    xml.Unmarshal([]byte(result), &imes)

    fmt.Println(string(result))

    if imes.Return_code == "SUCCESS" {
      // 解密
      // 1. 对加密串A做base64解码，得到加密串B
      // 2. 对商户key做md5，得到32位小写key* ( key设置路径：微信商户平台(pay.weixin.qq.com)–>账户设置–>API安全–>密钥设置 )
      // 3. 用key*对加密串B做AES-256-ECB解密
      b, err := base64.StdEncoding.DecodeString(imes.Req_info)
      if err != nil {
        fmt.Println(b, err)
        return
      }
      gocrypto.SetAesKey(strings.ToLower(gocrypto.Md5(JSAPImach_key)))

      plaintext, err := gocrypto.AesECBDecrypt(b)
      if err != nil {
        fmt.Println(err)
        return
      }

      fmt.Println(string(plaintext))

      // 转换XML
      var req_info Req_info_Message
      xml.Unmarshal([]byte(plaintext), &req_info)

      // req_info.out_trade_no 商户订单号,将订单状态修改为
      stmtIns, err := db.Prepare("update YUECHE_PAYMENTS set orderstatus=? where order_no=?")
      if err != nil {
          panic(err.Error())
      }
      defer stmtIns.Close() 

      _, err = stmtIns.Exec("退款完成",req_info.Out_trade_no)
      if err != nil {
        panic(err.Error())
      }

      var aboard_time string
      var start_city string
      var end_city string
      var phone string
      var ticket_no string
      get_start_end_by_Order(req_info.Out_trade_no, &start_city, &end_city, &aboard_time, &phone, &ticket_no)
      send_sms_order_tuikuan(phone, req_info.Out_trade_no)
    }


    var v JSAPIWeixinNotifyRet
    v.Return_code = "SUCCESS"

    output, err := xml.MarshalIndent(&v, "", "\t")
    if err != nil {
      fmt.Printf("error: %v\n", err)
    }

    fmt.Fprintf(w, "%s", output)
    return     
  }  
}


type Manage_ordercompleteextRet struct {
  Ret string
  Message string
}

// http.HandleFunc("/api/weixin/wechat_ordercomplete_ext", wechat_ordercompleteextHandle)  // 后台退钱完成后，修改订单状态为"已完成"  
func wechat_ordercompleteextHandle(w http.ResponseWriter, r *http.Request) { 
  log.Println("performance wechat_ordercompleteextHandle");
  r.ParseForm() //解析参数，默认是不会解析的 

  if r.Method == "GET" {

  } else if r.Method == "POST" {

    var retmessage Manage_ordercompleteextRet

    var Order_no_in string = "0"
    if len(r.Form["Order_no"]) > 0 {
      Order_no_in = r.Form["Order_no"][0]
    }
    
    // req_info.out_trade_no 商户订单号,将订单状态修改为
    stmtIns, err := db.Prepare("update YUECHE_PAYMENTS set orderstatus=? where order_no=?")
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close() 

    _, err = stmtIns.Exec("退款完成",Order_no_in)
    if err != nil {
      panic(err.Error())
    }

    var aboard_time string
    var start_city string
    var end_city string
    var phone string
    var ticket_no string
    get_start_end_by_Order(Order_no_in, &start_city, &end_city, &aboard_time, &phone, &ticket_no)
    send_sms_order_tuikuan(phone, Order_no_in)

    retmessage.Ret = "0"
    retmessage.Message = "success"
    b, _ := json.Marshal(retmessage)
    w.Header().Set("Access-Control-Allow-Origin","*")
    fmt.Fprintf(w, "%s", b)
    return   
  }  
}


