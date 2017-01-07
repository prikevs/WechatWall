package config

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

const (
	CONFIGJSON = "crawler-config.json"
	COOKIE     = "crawler-cookie"
)

type Config struct {
	ImagePath     string `json:"image_path"`
	PoolSize      int    `json:"pool_size"`
	CrawlInterval int    `json:"crawl_interval"`

	Token       string
	UsersURL    string            `json:"users_url"`
	UsersParams map[string]string `json:"users_params"`
	ImageURL    string            `json:"image_url"`
	ImageParams map[string]string `json:"image_params"`
	ImageSuffix string            `json:"image_suffix"`

	Headers map[string]string
}

func MustGetConfigJson(dir string) []byte {
	p := path.Join(dir, CONFIGJSON)
	data, err := ioutil.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return data
}

func MustGetCookie(dir string) string {
	p := path.Join(dir, COOKIE)
	data, err := ioutil.ReadFile(p)
	if data[len(data)-1] == byte('\n') {
		data = data[:len(data)-1]
	}
	if err != nil {
		panic(err)
	}
	return string(data)
}

func New(dir string) *Config {
	jconfig := MustGetConfigJson(dir)
	cookie := MustGetCookie(dir)

	c := &Config{}
	if err := json.Unmarshal([]byte(jconfig), c); err != nil {
		panic(err)
	}
	c.UsersParams["token"] = c.Token
	c.ImageParams["token"] = c.Token
	c.Headers["Cookie"] = cookie
	return c
}

func NewForTest() *Config {
	c := &Config{}
	if err := json.Unmarshal([]byte(jconfig), c); err != nil {
		panic(err)
	}
	c.UsersParams["token"] = c.Token
	c.ImageParams["token"] = c.Token
	return c
}

var jconfig = `
{
    "image_path": "/Users/kevince/Documents/others/gocrawler/src/crawler/data/images",
	"pool_size" : 10,
	"crawl_interval": 10,
    "users_url": "https://mp.weixin.qq.com/cgi-bin/user_tag",
    "users_params": {
        "action": "get_user_list",
        "lang": "zh_CN",
        "f": "json",
        "ajax": "1",
        "limit": "300",
        "token": ""
    },
    "image_url": "https://mp.weixin.qq.com/misc/getheadimg",
    "image_params": {
        "fakeid":"",
        "token":"",
        "lang":"zh_CN"
    },
	"image_suffix": "jpg",
    "headers" : {
        "Host": "mp.weixin.qq.com", 
        "Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2", 
        "Accept-Encoding": "gzip, deflate, sdch, br", 
        "Upgrade-Insecure-Requests": "1", 
        "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36", 
        "Connection": "keep-alive", 
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Cookie": "noticeLoginFlag=1; ts_uid=9547811126; pac_uid=1_332748660; BDTUJIAID=0a89736a87a67560332288dbf71f40d3; RK=GWXiwKHXMb; tvfe_boss_uuid=f01a8a6ad8a369a8; pgv_pvi=6866613248; pgv_si=s8545646592; rv2=809E5FA674C91DB6CB61D848C29BC2EEA5F3D027CF34E740F3; property20=589F7F2ACC2A03830E6D55EED735770AAEF2089836FF6E5672D01198F2E841820E24C818A0748B5C; pt2gguin=o0332748660; uin=o0332748660; skey=MVzwPQNj67; ptisp=cnc; ptcz=05c635a03b8c270bd4a85f707adfb5f33a7cc2728020e0e0072d09cff333c950; pgv_info=ssid=s6761200060; pgv_pvid=5545100198; o_cookie=332748660; noticeLoginFlag=1; uuid=2d0ebd244f8e8e0588f2c513c0be4903; ticket=1f0d05d8b98a5cc0bbdc0aab5e3ff5afb113aefd; ticket_id=gh_d592b39f8508; account=bujie8660@qq.com; cert=rRS9fwlDaQPLBj1y23HdE8KM3ob_E2MT; data_bizuin=3091358366; data_ticket=W/54iCloXsI7mxahCCnU7wNAZiOWQPnGfvslmKhHhzwnHeAlOaVvrwu6uXYl+4XC; ua_id=jvC9UdN9AzIGK0G1AAAAAMjInhSrH8QhUinR9e4qz4A=; xid=be568d779498b7586f4ef86ae433ee40; openid2ticket_onmJCuOJwunbnxaG2E-a2zfLsjWU=NBNi9aeDB3E40LtpQe+1yfvf16BZmTGngmRSdCVkD6Q=; slave_user=gh_d592b39f8508; slave_sid=VG1yRkxfTkdOWUVPV2FjM2NEZmYxYU83MUdQYTBaRnZFeEZPUkhKTXJXV3pGOU5KVU85SVJITjFXaThHeDVKZW1ReWhKTkhMR29iendPbnhCY3UzMU9BSGVVZWJhcnRrX1ZrOXFRa2luX2lydUlKb2kyeUhhQmw2M0FDOWVyaHVSM28wVW5sZTRSUEVldlc1; bizuin=3011352575"
    },
    "token" : "1604976388"
}
`
