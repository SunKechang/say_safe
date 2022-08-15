package safe

import (
	"bytes"
	"fmt"
	"gin-test/database/safe"
	"gin-test/util/flag"
	"gin-test/util/log"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

const JSESSIONID = "JSESSIONID"

type SafeService struct {
}

func NewSafeService() SafeService {
	return SafeService{}
}

func (p *SafeService) SendSafe(username, password string) (string, error) {
	safeJobDao := safe.NewSafeJobDao()
	safeJob, err := safeJobDao.GetJobByUserID(username)
	if err != nil {
		return "", err
	}

	safeLogDao := safe.NewSafeLogDao()
	safeLog := &safe.SafeLog{
		UserId: safeJob.UserId,
		JobId:  safeJob.ID,
	}
	defer func() {
		err := safeLogDao.CreateLog(safeLog)
		if err != nil {
			log.Log(fmt.Sprintf("sendSafe failed: %s\n", err))
			return
		}
	}()
	//todo 1.获取登录界面，拿到lt
	ltCode, sessionId, err := p.getLtSession()
	if err != nil {
		safeLog.Success = 0
		log.Log(fmt.Sprintf("sendSafe failed: %s\n", err))
		return "", err
	}
	//todo 2.通过lt，user，password，基于des加密计算得到rsa
	rsa, err := p.getRsa(username, password, ltCode)
	if err != nil {
		safeLog.Success = 0
		log.Log(fmt.Sprintf("sendSafe failed: %s\n", err))
		return "", err
	}
	//todo 3.登录，获取sessionID，route
	route, sessionId, err := p.getRoute(sessionId, rsa, ltCode, len(username), len(password))
	if err != nil {
		safeLog.Success = 0
		log.Log(fmt.Sprintf("sendSafe failed: %s\n", err))
		return "", err
	}
	//todo 4.读取要发送的表单数据到file
	log.Log(flag.SafeRoot + safeJob.Path + "\n")
	file, err := ioutil.ReadFile(path.Join(flag.SafeRoot, safeJob.Path))
	if err != nil {
		safeLog.Success = 0
		log.Log(fmt.Sprintf("sendSafe failed: %s\n", err))
		return "", err
	}
	//todo 5.报平安
	res, err := p.saySafe(sessionId, route, string(file))
	if err != nil {
		safeLog.Success = 0
		log.Log(fmt.Sprintf("sendSafe failed: %s\n", err))
		return "", err
	}
	safeLog.Success = 1
	safeLog.Result = string(res)
	return string(res), nil
}

func (p *SafeService) getLtSession() (string, string, error) { //获取lt与sessionId
	connectUrl := "http://cas.bjfu.edu.cn/cas/login?service=https%3A%2F%2Fs.bjfu.edu.cn%2Ftp_fp%2Findex.jsp"
	request, err := http.NewRequest("GET", connectUrl, nil)
	if err != nil {
		return "", "", err
	}
	resp, err := http.DefaultClient.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return "", "", err
	}
	//从html中找到对应节点获取lt，其实lt只是用于后续计算rsa的一个随机值，lt值大概长这样：LT-753722-hCsoUZ4f4QmSmeqFzPDyCFtAUAMHnu-cas
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	ltNode := doc.Find("#lt").Nodes[0]
	return ltNode.Attr[3].Val, getSingleCookie(resp, JSESSIONID), nil
}

func (p *SafeService) getRsa(user, password, ltCode string) (string, error) { //计算rsa值，该值在登录时登录用得到
	//rsa是通过用户名，密码，lt进行DES计算得到的值
	//具体算法在strconv.js文件中，此处通过goja包跨语言调用了该方法
	file, err := ioutil.ReadFile(path.Join(flag.SafeRoot + "/strconv.js"))
	if err != nil {
		return "", err
	}
	vm := goja.New()
	_, err = vm.RunString(string(file))
	if err != nil {
		log.Log("JS 有问题")
		return "", err
	}
	var fn func(string) string
	err = vm.ExportTo(vm.Get("main"), &fn)
	if err != nil {
		log.Log("Js函数映射到 Go 函数失败！")
		return "", err
	}
	word := user + password + ltCode
	return fn(word), nil
}

func (p *SafeService) getRoute(originId, rsa, ltCode string, ul, pl int) (string, string, error) { //获取route与sessionId
	//由于报平安请求需要访问s.bjfu.cn，而在此之前都是在情趣cas.bjfu.cn，因此需要重新获取sessionId，该sessionId是用于与s.bjfu.cn连接后产生的
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
		return "", "", err
	}

	request.Host = "cas.bjfu.edu.cn"
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	request.Header.Set("Cache-Control", "max-age=0")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Origin", "http://cas.bjfu.edu.cn")
	request.Header.Set("Referer", "http://cas.bjfu.edu.cn/cas/login?service=https%3A%2F%2Fs.bjfu.edu.cn%2Ftp_fp%2Findex.jsp")
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	request.Header.Set("Cookie", "JSESSIONID="+originId+"; cas_hash=; Language=zh_CN")

	casResp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", "", err
	}
	defer casResp.Body.Close()
	sessionId := getSingleCookie(casResp.Request.Response.Request.Response, JSESSIONID)
	route := getSingleCookie(casResp.Request.Response.Request.Response, "route")
	return route, sessionId, nil
}

func (p *SafeService) saySafe(sessionId, route, sendInfo string) ([]byte, error) { //报平安请求

	actionUrl := "https://s.bjfu.edu.cn/tp_fp/formParser?status=update&formid=7394b770-ba93-4041-91b7-80198a68&workflowAction=startProcess&seqId=&unitId=&applyCode=&workitemid=&process=bae380db-7db4-4c7c-9458-d79188fa359a"
	//todo 将要发送的表单数据sendInfo放入request中
	reader := bytes.NewReader([]byte(sendInfo))
	request, err := http.NewRequest("POST", actionUrl, reader)
	if err != nil {
		return nil, err
	}
	//todo 设置request header，并在Cookie中添加route和sessionId
	request.Host = "s.bjfu.edu.cn"
	request.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Accept-Language", "zh-cn")
	request.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	request.Header.Set("Origin", "https://s.bjfu.edu.cn")
	request.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.23(0x1800172e) NetType/WIFI Language/zh_CN")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Referer", "https://s.bjfu.edu.cn/tp_fp/formParser?status=select&formid=7394b770-ba93-4041-91b7-80198a68&service_id=99f0cf19-3ca4-4786-badb-521f0f734cad&process=bae380db-7db4-4c7c-9458-d79188fa359a&seqId=&seqPid=&privilegeId=8467766035e3e965668c850086270762")
	request.Header.Set("Cookie", JSESSIONID+"="+sessionId+"; route="+route)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//todo 拿到response body写入到日志中
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
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

//6.28 无头绪状态，感觉需要模拟登陆状态，再模拟发送请求，但是登录连postman都无法模拟，很懵
//6.30 已经完成80%啦～之前遇到的postman无法模拟是因为302跳转了很多个请求，但是只要从尽头找需要什么参数，往前就会发现那些参数都是通过一次次请求拿来的。
//7.1 目前要做的是自动获取默认提交值，加油！
//7.5 已经部署到服务器了，自动获取默认值没有完成，就算是了解了大致过程了吧。

func (p *SafeService) AddSafe(username string, safeInfo []byte) error {
	// todo 将用户传递的报平安内容保存到root/学号/学号.txt
	relaPath := path.Join(username, username+".txt")
	filePath := path.Join(flag.SafeRoot, relaPath)
	var file *os.File
	if _, err := os.Stat(filePath); err != nil {
		file, err = os.Create(filePath)
		if err != nil {
			log.Log(fmt.Sprintf("AddSafe failed: %s\n", err))
			return err
		}
	} else {
		file, err = os.OpenFile(filePath, os.O_RDWR, 0644)
		if err != nil {
			log.Log(fmt.Sprintf("AddSafe failed: %s\n", err))
			return err
		}
	}

	_, err := file.Write(safeInfo)
	if err != nil {
		log.Log(fmt.Sprintf("AddSafe failed: %s\n", err))
		return err
	}

	safeJobDao := safe.NewSafeJobDao()
	err = safeJobDao.DeleteJobsByUserID(username)
	if err != nil {
		log.Log(fmt.Sprintf("AddSafe failed: %s\n", err))
		return err
	}

	newJob := safe.SafeJob{
		UserId: username,
		Path:   relaPath,
	}
	err = safeJobDao.CreateSafeJob(&newJob)
	if err != nil {
		log.Log(fmt.Sprintf("AddSafe failed: %s\n", err))
		return err
	}
	return nil
}

func (p *SafeService) GetSafe(username string) (string, error) {
	jobDao := safe.NewSafeJobDao()
	safeInfo, err := jobDao.GetJobByUserID(username)
	if err != nil {
		log.Log(fmt.Sprintf("GetSafe failed: %s\n", err))
		return "", err
	}
	filePath := path.Join(flag.SafeRoot, safeInfo.Path)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Log(fmt.Sprintf("GetSafe failed: %s\n", err))
		return "", err
	}
	return string(content), nil
}
