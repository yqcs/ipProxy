# ipProxy
使用goQuery根据模板自动采集IP，验证采集完毕验证代理的访问速度，使用筛选合格的IP作为请求的代理
 
 ##### 测试demo，如有问题提交issues
 
## 一、goQuery库：

#### **1.简述**
据官方描述，goQuery实现了与jQuery相似的DOM操作功能。与jQuery不同的是，jQuery返回的是完整的DOM树，而goQuery返回的是DOM节点。goQuery底层由golang标准库net/html实现，解析器要求文档必须是 UTF-8 编码，使用者应按需转换文档编码。[goQuery-readme](https://github.com/PuerkitoBio/goquery/blob/master/README.md)
#### **2.主要方法**
**2.1 Document：** 返回要被操作的HTML文档

```go
// Document represents an HTML document to be manipulated. Unlike jQuery, which
// is loaded as part of a DOM document, and thus acts upon its containing
// document, GoQuery doesn't know which HTML document to act upon. So it needs
// to be told, and that's what the Document class is for. It holds the root
// document node to manipulate, and can make selections on this document.
type Document struct {
	*Selection
	Url      *url.URL
	rootNode *html.Node
}
```
**2.2 Selection：** 符合指定条件的节点。
```go
// Selection represents a collection of nodes matching some criteria. The
// initial Selection can be created by using Document.Find, and then
// manipulated using the jQuery-like chainable syntax and methods.
type Selection struct {
	Nodes    []*html.Node
	document *Document
	prevSel  *Selection
}
```
**2.3文档操作函数：** 

```go
Eq()
Index()
Last()
Slice()
Get()
······
```
## 二、采集代理IP：
#### **1.代理IP池**
 1. [高可用全球免费代理IP库](https://ip.jiangxianli.com/)
 2. [西刺免费代理IP](https://www.xicidaili.com/nn)
 3. [全网IP代理](http://www.goubanjia.com/)
 4.  [快代理](https://www.kuaidaili.com/free/)
#### **2.goQuery采集：**
不多说，上代码：

```go
import "github.com/PuerkitoBio/goquery"
```
导入goquery库
```go
//采集代理返回的参数
type proxyResult struct {
	Ip           string `json:"ip"`           //ip
	Port         int    `json:port`           //端口
	Agreement    string `json:agreement`      //请求协议
	Anonymous    string `json:anonymous`      //透明度
	Region       string `json:region`         //地区
	Speed        string `json:"speed"`        //响应速度
	Source       string `json:"source"`       //来源（采集资源站）
	Verification string `json:"verification"` //验证时间
}

//采集代理所需的参数
type proxyParamet struct {
	ipIndex           int `json:"ipIndex"`           //ip下标
	portIndex         int `json:"portIndex"`         //端口下标
	agreementIndex    int `json:"agreementIndex"`    //请求协议下标
	anonymousIndex    int `json:"anonymousIndex"`    //透明度下标
	regionIndex       int `json:"regionIndex"`       //地区下标
	speedIndex        int `json:"speedIndex"`        //响应速度下标
	sourceIndex       int `json:"sourceIndex"`       //来源（采集资源站）下标
	verificationIndex int `json:"verificationIndex"` //验证时间下标
}
```
创建两个结构体，`proxyResult`保存返回的数据，`proxyParamet` 用来携带传入参数
```go
func CollectionResources(targetUrl string, parame proxyParamet)[]proxyResult  {
	//用来存储采集结果
	  proxyList := []proxyResult{}
	//请求目标站点方式（GET/POST）
    resp, err := http.Get(targetUrl)
    //请求失败则输出日志
    if resp.StatusCode != 200 || err != nil {
        log.Pr("spider", "请求出错", err)
    }
    //添加随机User-Agent
    resp.Header.Add("User-Agent", random.RandomUseragent())
    //返回目标站html文档
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    //因为大部分都是以表格方式展示，所以这里就直接抓取tbody的内容
    doc.Find("tbody tr").Each(func(i int, selection *goquery.Selection) {
    	//golang不支持自动类型转换，这里手动转换拼接
        ip := selection.Find("td:nth-child(" + strconv.Itoa(parame.ipIndex) + ")").Text()
        port, err := strconv.Atoi(selection.Find("td:nth-child(" + strconv.Itoa(parame.portIndex) + ")").Text())
        agreement := selection.Find("td:nth-child(" + strconv.Itoa(parame.agreementIndex) + ")").Text()
        anonymous := selection.Find("td:nth-child(" + strconv.Itoa(parame.anonymousIndex) + ")").Text()
        region := selection.Find("td:nth-child(" + strconv.Itoa(parame.regionIndex) + ")").Text()
        speedString := selection.Find("td:nth-child(" + strconv.Itoa(parame.speedIndex) + ")").Text()
        //有的代理池会携带上单位：秒，加上这段代码可以去掉。可加可不加，按需设置
        speed := strings.TrimRight(speedString, "秒")
        verification := selection.Find("td:nth-child(" + strconv.Itoa(parame.verificationIndex) + ")").Text()
        if err != nil {
            log.Pr("spider", "数据转换出错", err)
        }
		  proxyList = append(proxyList, proxyResult{Ip: ip,
            Port:         port,
            Agreement:    agreement,
            Anonymous:    anonymous,
            Region:       region,
            Speed:        speed,
            Source:       targetUrl,
            Verification: verification}, )
    })
    //采集完毕之后清理
    defer resp.Body.Close()
    return proxyList 
}
```
**调用：**

```go
 parame := proxyParamet{
        ipIndex:           1,
        portIndex:         2,
        agreementIndex:    4,
        anonymousIndex:    3,
        regionIndex:       5,
        speedIndex:        6,
        verificationIndex: 7,
    }
    
	CollectionResources("https://www.kuaidaili.com/free/", parame)
```
parame的参数对应着表格标题的下标（从1开始）
![在这里插入图片描述](https://img-blog.csdnimg.cn/20200406180759810.png)

设置代理

```go
func StartRequestProxy(address string, proxyIpInfo [] proxyResult) string {
	proxyAddr := "协议+地址+端口 /如：https://127.0.0.1:8080"
	url := address
	cli := newHttpClient(proxyAddr)
	data, _ := httpGET(cli, url)
	return string(data)
}

func newHttpClient(proxyAddr string) *http.Client {
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		return nil
	}
	netTransport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(10))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   10,                             //每个host最大空闲连接
		ResponseHeaderTimeout: time.Second * time.Duration(5), //数据收发5秒超时
	}

	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

func httpGET(client *http.Client, url string) (body []byte, err error) {
	rsp, err := client.Get(url)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK || err != nil {
		err = fmt.Errorf("HTTP GET Code=%v, URI=%v, err=%v", rsp.StatusCode, url, err)
		log.Pr("HttpGet", "Request error", err)
		return
	}
	return ioutil.ReadAll(rsp.Body)
}
```
运行结果：
 ![Alt text](https://raw.githubusercontent.com/yqcs/ipProxy/master/describePictures/all.png)

 
注：goquery的部分介绍引用了[go语言中文网](http://blog.studygolang.com/2015/04/go-jquery-goquery/)
