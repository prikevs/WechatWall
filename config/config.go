package config

import (
	"encoding/json"
)

type Config struct {
	ImagePath string `json:"image_path"`
	Token     string

	UsersURL    string            `json:"users_url"`
	UsersParams map[string]string `json:"users_params"`

	ImageURL    string            `json:"image_url"`
	ImageParams map[string]string `json:"image_params"`

	Cookie  string
	Headers map[string]string
}

func New() *Config {
	c := &Config{}
	if err := json.Unmarshal([]byte(jconfig), c); err != nil {
		panic(err)
	}
	c.UsersParams["token"] = c.Token
	c.ImageParams["token"] = c.Token
	c.Headers["Cookie"] = c.Cookie
	return c
}

var jconfig = `
{
    "image_path": "/Users/kevince/Documents/others/wechat/images",
    "users_url": "https://mp.weixin.qq.com/cgi-bin/user_tag",
    "users_params": {
        "action": "get_user_list",
        "lang": "zh_CN",
        "f": "json",
        "ajax": "1",
        "limit": "100",
        "token": ""
    },
    "image_url": "https://mp.weixin.qq.com/misc/getheadimg",
    "image_params": {
        "fakeid":"",
        "token":"",
        "lang":"zh_CN"
    },


    "headers" : {
        "Host": "mp.weixin.qq.com", 
        "Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2", 
        "Accept-Encoding": "gzip, deflate, sdch, br", 
        "Upgrade-Insecure-Requests": "1", 
        "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36", 
        "Connection": "keep-alive", 
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
        "Cookie": ""
    },

    "token" : "1080732709",
    "cookie" : "s_uid=9547811126; pac_uid=1_332748660; BDTUJIAID=0a89736a87a67560332288dbf71f40d3; RK=GWXiwKHXMb; tvfe_boss_uuid=f01a8a6ad8a369a8; pgv_pvi=6866613248; pgv_si=s8545646592; rv2=809E5FA674C91DB6CB61D848C29BC2EEA5F3D027CF34E740F3; property20=589F7F2ACC2A03830E6D55EED735770AAEF2089836FF6E5672D01198F2E841820E24C818A0748B5C; pt2gguin=o0332748660; uin=o0332748660; skey=MVzwPQNj67; ptisp=cnc; ptcz=05c635a03b8c270bd4a85f707adfb5f33a7cc2728020e0e0072d09cff333c950; pgv_info=ssid=s6761200060; pgv_pvid=5545100198; o_cookie=332748660; uuid=1b8ef83fd8d27f3cb5037a67cb1183dd; ticket=ef7bb0bf4c64787fc63a8870dd3a7e59818a4faa; ticket_id=gh_d592b39f8508; account=bujie8660@qq.com; cert=rRS9fwlDaQPLBj1y23HdE8KM3ob_E2MT; noticeLoginFlag=1; data_bizuin=3091358366; data_ticket=8R2amfl2EStKVtTwE7K88fvVfBY5fAXCh2niKy5msCK7hfF71L1rNnF9hBxz7Efv; ua_id=jvC9UdN9AzIGK0G1AAAAAMjInhSrH8QhUinR9e4qz4A=; xid=7b6b518b7a47a5b14fb30cc33c4dacb2; openid2ticket_onmJCuOJwunbnxaG2E-a2zfLsjWU=mw4lXLSSagONuNZlbmYAvA/drYa3jogrhu46tJdWWz8=; slave_user=gh_d592b39f8508; slave_sid=djdSVnZpdlNoQ0dCbFU2aDFZZGdMX3J3Tmd6bjVMYWNPUG5sZThLSnVQR2pUTWNKU0dXdTdPQmRKazlFZEM3em5nY0Z1eWRpZWRyQ25iMEdvQ1NKeTZpN29oeXNYM2xITlFtNzY0SWdDSDIza3dzZEF4OWNHZ2hlNGlQR3c3QXRwV1VNQUhoUFAxT09EU1hB; bizuin=3011352575"
}
`
