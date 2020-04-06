package collection

import (
	"github.com/PuerkitoBio/goquery"
	"ipProxy/utils/log"
	"ipProxy/utils/random"
	"net/http"
	"strconv"
	"strings"
)

func StartCollectionProxyIp() []proxyResult {
	parame := proxyParamet{
		ipIndex:           1,
		portIndex:         2,
		agreementIndex:    4,
		anonymousIndex:    3,
		regionIndex:       5,
		speedIndex:        6,
		verificationIndex: 7,
	}
	return CollectionResources("https://www.kuaidaili.com/free/", parame)
}
func CollectionResources(targetUrl string, parame proxyParamet) []proxyResult {
	proxyList := []proxyResult{}
	resp, err := http.Get(targetUrl)
	if resp.StatusCode != 200 || err != nil {
		log.Pr("spider", "请求出错", err)
	}
	resp.Header.Add("User-Agent", random.RandomUseragent())
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	doc.Find("tbody tr").Each(func(i int, selection *goquery.Selection) {
		ip := selection.Find("td:nth-child(" + strconv.Itoa(parame.ipIndex) + ")").Text()
		port, err := strconv.Atoi(selection.Find("td:nth-child(" + strconv.Itoa(parame.portIndex) + ")").Text())
		agreement := selection.Find("td:nth-child(" + strconv.Itoa(parame.agreementIndex) + ")").Text()
		anonymous := selection.Find("td:nth-child(" + strconv.Itoa(parame.anonymousIndex) + ")").Text()
		region := selection.Find("td:nth-child(" + strconv.Itoa(parame.regionIndex) + ")").Text()
		speedString := selection.Find("td:nth-child(" + strconv.Itoa(parame.speedIndex) + ")").Text()
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
			Verification: verification})
	})
	defer resp.Body.Close()
	return proxyList
}
