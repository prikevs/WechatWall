package ucrawler

import (
	"encoding/json"
)

type AllData struct {
	UserList struct {
		UserInfoList []User `json:"user_info_list"`
	} `json:"user_list"`
}

type User struct {
	UserOpenid     string `json:"user_openid"`
	UserName       string `json:"user_name"`
	UserRemark     string `json:"user_remark"`
	UserCreateTime int64  `json:"user_create_time"`
}

func parseAllData(data []byte) (*AllData, error) {
	all := &AllData{}
	if err := json.Unmarshal(data, all); err != nil {
		return nil, err
	}
	return all, nil
}

func Parse(data []byte) ([]User, error) {
	all := &AllData{}
	if err := json.Unmarshal(data, all); err != nil {
		return nil, err
	}
	return all.UserList.UserInfoList, nil
}
