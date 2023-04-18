package wechat

import (
	"fmt"
	"github.com/medivhzhan/weapp/v2"
)

type Service struct {
	AppID     string
	AppSecret string
}

// Resolve wechat专管用户微信登录这块
func (s *Service) Resolve(code string) (string, error) {
	res, err := weapp.Login(s.AppID, s.AppSecret, code)
	if err != nil {
		return "", fmt.Errorf("weapp.Login:%v", err)
	}
	if err := res.GetResponseError(); err != nil {
		return "", fmt.Errorf("weapp resopnse error:%v", err)
	}
	return res.OpenID, nil
}
