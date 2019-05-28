package main

import (
	"./plugin"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/levigross/grequests"
	"net/url"
	"time"
)

var(
	ChinaProxy plugin.CnProxy
	DlProxy plugin.DailiProxy
	Proxy2 plugin.DailiProxy
	Proxy3 plugin.XilaProxy
)

func init()  {
	//加载代理网站插件
	ChinaProxy=plugin.CnProxy{}
	ChinaProxy.Load("https://cn-proxy.com/")
	DlProxy=plugin.DailiProxy{}
	DlProxy.Load("https://www.kuaidaili.com/free/inha/1/")
	Proxy2=plugin.DailiProxy{}
	Proxy2.Load("https://www.kuaidaili.com/free/intr/")
	Proxy3=plugin.XilaProxy{}
	Proxy3.Load("http://www.xiladaili.com/")

}

func main(){
	asd:=Proxy3.GetProxy()
	fmt.Println(asd)
	proxyURL, err := url.Parse("http://" + asd.Ip+":"+asd.Port) // Proxy URL
	if err != nil {
		logs.Error("解析代理出错，原因:%s", err)
	}
	resp, err := grequests.Get("http://www.baidu.com/",
	&grequests.RequestOptions{
		Proxies:        map[string]*url.URL{proxyURL.Scheme: proxyURL},
		RequestTimeout: time.Second * 5,
	})
	println(resp.String())
}
