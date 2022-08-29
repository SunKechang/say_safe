package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin-test/database"
	"gin-test/database/user"
	"gin-test/handler/response"
	"gin-test/util/log"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

const (
	Admin = 2
	User  = 1

	Role     = "role"
	UserName = "username"
	Password = "password"

	JSESSIONID = "JSESSIONID"
)

func Login() func(*gin.Context) {
	return func(c *gin.Context) {
		login(c)
	}
}
func Logout() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		logout(ctx)
	}
}
func SignUp() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		signup(ctx)
	}
}
func login(c *gin.Context) {
	res := response.NewResponse()
	defer c.JSON(res["code"].(int), res)
	session := sessions.Default(c)
	session.Set(Role, nil)

	//读取request body中数据
	temp, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Log(fmt.Sprintf("[LOGIN] failed: %s\n", err.Error()))
	}
	//转化为LoginRequest结构
	body := LoginRequest{}
	err = json.Unmarshal(temp, &body)
	if err != nil {
		res["code"] = http.StatusBadRequest
		res[Message] = "传入参数有误"
		return
	}
	id := body.UserName
	password := body.Password

	userDao := user.NewUserDao()
	userInfo, err := userDao.GetUserByID(id)
	if err != nil {
		res["code"] = http.StatusBadRequest
		if database.IsRecordNotFound(err) {
			res[Message] = "账号未注册"
		} else {
			res[Message] = err.Error()
		}
		log.Log(fmt.Sprintf("Login failed: %s\n", err))
		return
	}

	if userInfo.Password == password {
		res[Message] = "login success"
		session.Set(Role, User)
		session.Set(UserName, id)
		session.Set(Password, password)
	} else {
		res[Message] = "login fail"
	}

	err = session.Save()
	if err != nil {
		res["code"] = http.StatusForbidden
		res[Message] = "Login failed"
		log.Log(fmt.Sprintf("[LOGIN] failed:%s\n", err))
	}
}

func logout(c *gin.Context) {
	res := response.NewResponse()
	defer c.JSON(res["code"].(int), res)
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		res["code"] = http.StatusInternalServerError
		res[Message] = "退出失败"
	}
}

func signup(c *gin.Context) {
	res := response.NewResponse()
	defer c.JSON(res["code"].(int), res)
	temp, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Log(fmt.Sprintf("[SIGNUP] failed: %s\n", err.Error()))
	}
	//转化为LoginRequest结构
	body := SignUpRequest{}
	err = json.Unmarshal(temp, &body)
	if err != nil {
		log.Logger("unmarshal failed: %s\n", err.Error())
	}

	log.Logger("signup body %v\n", body)
	userDao := user.NewUserDao()
	_, err = userDao.GetUserByID(body.UserName)
	log.Logger("database %v\n", err)
	if database.IsError(err) {
		res["code"] = http.StatusInternalServerError
		res[Message] = "数据库内部错误"
		return
	}
	if !database.IsRecordNotFound(err) {
		res["code"] = http.StatusOK
		res[Message] = "用户已注册"
		return
	}

	ltCode, sessionId := getLtSession()
	log.Logger("ltcode %s\n", ltCode)
	//todo 2.通过lt，user，password，基于des加密计算得到rsa
	rsa := getRsa(body.UserName, body.Password, ltCode)
	log.Logger("rsa: %s\n", rsa)
	//todo 3.登录，获取sessionID，route
	route, sessionId := getRoute(sessionId, rsa, ltCode, len(body.UserName), len(body.Password))
	if len(route) == 0 {
		res[Message] = "账号或密码错误"
		return
	}
	log.Logger("route: %s, session: %s\n", route, sessionId)
	newUser := user.User{
		ID:       body.UserName,
		UserName: "",
		Salt:     "",
		Password: body.Password,
		Class:    "",
		IsMan:    0,
	}
	log.Logger("create user %v\n", newUser)
	err = userDao.CreateUser(&newUser)
	if err != nil {
		log.Logger("create user %v\n", newUser)
		res[Message] = err.Error()
		return
	}
	session := sessions.Default(c)
	session.Set(Role, User)
	session.Set(UserName, body.UserName)
	session.Set(Password, body.Password)
	res[Message] = "注册成功"
}

func getLtSession() (string, string) { //获取lt与sessionId
	connectUrl := "http://cas.bjfu.edu.cn/cas/login?service=https%3A%2F%2Fs.bjfu.edu.cn%2Ftp_fp%2Findex.jsp"
	log.Logger("connection url: %s\n", connectUrl)
	request, err := http.NewRequest("GET", connectUrl, nil)
	if err != nil {
		log.Log(err.Error())
	}
	resp, err := http.DefaultClient.Do(request)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("error: ", err)
	}
	//从html中找到对应节点获取lt，其实lt只是用于后续计算rsa的一个随机值，lt值大概长这样：LT-753722-hCsoUZ4f4QmSmeqFzPDyCFtAUAMHnu-cas
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	ltNode := doc.Find("#lt").Nodes[0]
	log.Logger("ltNode: %v\n", ltNode)
	return ltNode.Attr[3].Val, getSingleCookie(resp, JSESSIONID)
}

func getRsa(user, password, ltCode string) string { //计算rsa值，该值在登录时登录用得到
	//rsa是通过用户名，密码，lt进行DES计算得到的值
	//具体算法在strconv.js文件中，此处通过goja包跨语言调用了该方法
	file, err := ioutil.ReadFile(path.Join("./strconv.js"))
	if err != nil {
		fmt.Println(err)
	}
	vm := goja.New()
	_, err = vm.RunString(string(file))
	if err != nil {
		fmt.Println("JS代码有问题！")
		return ""
	}
	var fn func(string) string
	err = vm.ExportTo(vm.Get("main"), &fn)
	if err != nil {
		fmt.Println("Js函数映射到 Go 函数失败！")
		return ""
	}
	word := user + password + ltCode
	return fn(word)
}

func getRoute(originId, rsa, ltCode string, ul, pl int) (route, sessionId string) { //获取route与sessionId
	//由于报平安请求需要访问s.bjfu.cn，而在此之前都是在请求cas.bjfu.cn，因此需要重新获取sessionId，该sessionId是用于与s.bjfu.cn连接后产生的
	surl := "http://cas.bjfu.edu.cn/cas/login?service=https%3A%2F%2Fs.bjfu.edu.cn%2Ftp_fp%2Findex.jsp"
	// 用url.values方式构造form-data参数
	formValues := url.Values{}
	formValues.Set("rsa", rsa)
	formValues.Set("ul", strconv.Itoa(ul))
	formValues.Set("pl", strconv.Itoa(pl))
	formValues.Set("lt", ltCode)
	formValues.Set("execution", "e1s1")
	formValues.Set("_eventId", "submit")
	formDataStr := formValues.Encode()
	formDataBytes := []byte(formDataStr)
	formBytesReader := bytes.NewReader(formDataBytes)
	request, err := http.NewRequest("POST", surl, formBytesReader)
	if err != nil {
		fmt.Println("err: ", err)
	}

	refer := "http://cas.bjfu.edu.cn/cas/login?service=https%3A%2F%2Fs.bjfu.edu.cn%2Ftp_fp%2Findex.jsp"
	host := "cas.bjfu.edu.cn"
	setCommonHeader(request, host, refer)
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Cookie", "JSESSIONID="+originId+"; cas_hash=; Language=zh_CN")
	request.Header.Set("Origin", "http://cas.bjfu.edu.cn")

	casResp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer casResp.Body.Close()
	if casResp.Request.Response == nil {
		return "", ""
	}
	sessionId = getSingleCookie(casResp.Request.Response.Request.Response, JSESSIONID)
	route = getSingleCookie(casResp.Request.Response.Request.Response, "route")
	return route, sessionId
}

func setCommonHeader(r *http.Request, host, refer string) {
	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	r.Header.Set("Accept-Encoding", "gzip, deflate, br")
	r.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	r.Header.Set("Cache-Control", "max-age=0")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	r.Header.Set("Referer", refer)
	r.Host = host
}

func getSingleCookie(response *http.Response, goalName string) string { //从Cookie中查找想要的参数对应的值，JSESSIONID，route
	cookies := response.Cookies()
	for _, cookie := range cookies {
		if strings.EqualFold(cookie.Name, goalName) {
			return cookie.Value
		}
	}
	return ""
}

func writeFile(fileName string, data []byte) { //向文件中写入内容，并在前面添加系统当前时间
	fl, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	defer fl.Close()
	timeInfo := []byte(time.Now().String())
	temp := string(timeInfo) + string(data)
	write, err := fl.Write([]byte(temp))
	if err != nil {
		fmt.Println(err, write)
	}
}
