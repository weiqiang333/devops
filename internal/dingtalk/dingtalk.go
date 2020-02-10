package dingtalk

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)


//Dingtalk 发送钉钉消息至群.(接收参数 token , '发送信息内容', true/false)
func Dingtalk(access_token, message, isAtAll string) error {
	data := fmt.Sprintf(`{
        "msgtype": "text",
            "text": {
            "content": "%s",
        },
        "at": {
            "isAtAll": "%s",
        },
    }`, message, isAtAll)
	bodys := strings.NewReader(data)
	resp, err := http.Post(fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", access_token),
		"application/json", bodys)
	if err != nil {
		log.Println(http.StatusInternalServerError)
		return err
	}
	log.Println(resp.StatusCode)
	return nil
}
