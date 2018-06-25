package request

import (
	"fmt"
	"github.com/sclevine/agouti"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Proxy struct {
	Ip           string
	Port         string
	Protocol     string
	ResponseTime int
}

var Proxies []Proxy

func fetchPageByAgouti(u string) {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{"--headless", "--disable-gpu", "--no-sandbox"}),
	)
	//driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		panic(err)
	}
	page, err := driver.NewPage()
	defer page.CloseWindow()
	if err != nil {
		panic(err)
	}
	if err := page.Navigate(u); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 3)
	pageSource, _ := page.HTML()
	fmt.Println(pageSource)
}

func Get(_url string, c chan string, cc chan int) {
	fetchStart := time.Now()

	httpClient := &http.Client{
		Timeout: time.Duration(time.Second * 10),
	}
	req, err := http.NewRequest("GET", _url, nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:60.0) Gecko/20100101 Firefox/60.0`)
	resp, err := httpClient.Do(req)
	if err != nil {
		c <- fmt.Sprintf(" >> Fetch url : %v with error (%v) \n", _url, err.Error())
		return
	}
	defer resp.Body.Close()
	pageSource, err := ioutil.ReadAll(resp.Body)

	fetchPageByAgouti(_url)

	r, _ := regexp.Compile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}((<[^>]+>)|(\s)|(\:))+\d{2,5}`)
	result := r.FindAllString(string(pageSource), -1)
	for _, item := range result {
		var re = regexp.MustCompile(`((<[^>]+>)|(\s))+`)
		p := re.ReplaceAllString(item, `:`)
		pi := strings.Split(p, `:`)
		if len(pi) < 2 {
			continue
		}
		Proxies = append(Proxies, Proxy{
			pi[0],
			pi[1],
			"",
			0,
		},
		)
	}
	c <- fmt.Sprintf(" >> Fetch %v spent : %v seconds \n", _url, time.Since(fetchStart).Seconds())
	cc <- <-cc - 1
	return
}

func (p *Proxy) CheckProxy() (result bool, err error) {
	if p.Protocol == "" {
		p.Protocol = "http://"
	}
	proxy := p.Protocol + p.Ip + ":" + p.Port
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}
	httpClient := &http.Client{
		Timeout: time.Duration(time.Second * 5),
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	req, err := http.NewRequest("GET", "http://weixin.sogou.com/robots.txt", nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:60.0) Gecko/20100101 Firefox/60.0`)
	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	fmt.Println(string(body))
	return true, nil
}
