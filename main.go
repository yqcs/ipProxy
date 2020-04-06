package main

import (
	"fmt"
	"ipProxy/collection"
	"ipProxy/proxy"
	"net/http"
	"strconv"
	"time"
)

func main() {
	fmt.Println("GitHub：https://github.com/yqcs")
	result := collection.StartCollectionProxyIp()
	fmt.Println(">>>>采集完成，信息：")
	fmt.Println(result)
	fmt.Println(">>>>开始筛选IP")
	for _, proxyInfo := range result {
		fmt.Println(">>IP:" + proxyInfo.Ip)
		url := proxyInfo.Agreement + "://" + proxyInfo.Ip + ":" + strconv.Itoa(proxyInfo.Port)
		timeout := time.Duration(3 * time.Second)
		//设置3秒超时
		client := http.Client{
			Timeout: timeout,
		}
		_, err := client.Get(url)

		//请求出错或者超时，就更换下一个
		if err != nil {
			fmt.Println("请求出错，或者是响应超时,目标：" + url)
			continue
		}
		fmt.Println(">>>开始匿名请求")
		proxy.StartRequestProxy("https://www.baidu.com/s?wd=ip", url)

	}
}
