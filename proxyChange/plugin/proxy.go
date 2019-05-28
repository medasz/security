package plugin

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"github.com/levigross/grequests"
	"net/url"
	"time"
)

type IpPort struct {
	Ip   string
	Port string
}

type Proxy interface {
	Get() []IpPort
	CheckUrlOk() bool
	CheckProxyOk(data []IpPort) IpPort
	GetProxy() IpPort
}

type CnProxy struct {
	url string
}

func (c *CnProxy) Load(url string) {
	c.url = url
}

func (c *CnProxy) Get() []IpPort {
	var proxyList []IpPort
	rq := &grequests.RequestOptions{
		RequestTimeout: time.Second * 5,
	}
	for j := 0; j < 5; j++ {
		res, err := grequests.Get(c.url, rq)
		if err != nil {
			logs.Error("发送请求失败，原因:%s", err)
			continue
		}
		dom, err := goquery.NewDocumentFromReader(res)
		if err != nil {
			logs.Error("解析html文件失败，原因:%s", err)
			continue
		}
		dom.Find("tbody").Each(func(i int, selection *goquery.Selection) {
			selection.Find("tr").Each(func(i int, selection *goquery.Selection) {
				data := IpPort{
					selection.Find("td").Nodes[0].FirstChild.Data,
					selection.Find("td").Nodes[1].FirstChild.Data,
				}
				proxyList = append(proxyList, data)
			})
		})
		return proxyList
	}
	return proxyList
}

func (c *CnProxy) CheckUrlOk() bool {
	rq := &grequests.RequestOptions{
		RequestTimeout: time.Second * 5,
	}
	res, err := grequests.Get(c.url, rq)
	if err != nil {
		logs.Error("请求网址%s出错，原因:%s", c.url, err)
		return false
	}
	return res.Ok
}

func (c *CnProxy) CheckProxyOk(data []IpPort) IpPort {
	for _, y := range data {
		proxyURL, err := url.Parse("http://" + y.Ip + ":" + y.Port) // Proxy URL
		if err != nil {
			logs.Error("解析代理出错，原因:%s", err)
			continue
		}
		resp, err := grequests.Get("http://www.baidu.com/",
			&grequests.RequestOptions{
				Proxies:        map[string]*url.URL{proxyURL.Scheme: proxyURL},
				RequestTimeout: time.Second * 5,
			})
		if err != nil {
			logs.Error("验证代理失败，原因:%s", err)
			continue
		}
		if resp.Ok {
			return y
		} else {
			continue
		}
	}
	return IpPort{}
}

func (c *CnProxy) GetProxy() IpPort {
	if c.CheckUrlOk() {
		result := c.CheckProxyOk(c.Get())
		return result
	} else {
		fmt.Println("无可用代理")
		return IpPort{}
	}
}
