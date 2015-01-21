package main

import (
	"fmt"
	"github.com/fanyang1988/goconfig"
	gologger "github.com/fanyang1988/gologger"
	_ "github.com/icattlecoder/godaemon"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func onPay(param map[string]string) (re string, err error) {
	fmt.Printf("onPay %s\n", param)
	ext, ext_ok := param["ext"]
	money, money_ok := param["money"]
	if ext_ok && money_ok {
		fmt.Printf("ext %s\n", ext)
		ext_info := strings.Split(ext, "|")
		if len(ext_info) < 2 {
			return "failed", nil
		}
		fmt.Printf("url %s\n", ext_info[0])
		fmt.Printf("uid %s\n", ext_info[1])
		fmt.Printf("money %s\n", money)

		u, _ := url.Parse(ext_info[0])
		q := u.Query()

		q.Set("app_key", "cc5a0304d4efa71a07cde121d20fdb54")
		q.Set("product_id", "TODO")
		q.Set("amount", money)
		q.Set("app_uid", "1111")
		q.Set("order_id", "1111")
		q.Set("user_id", ext_info[1])
		q.Set("gateway_flag", "success")
		q.Set("app_order_id", "1111")

		u.RawQuery = q.Encode()
		//TODO timeout
		res, err := http.Get(u.String())
		if err != nil {
			return "failed", nil
		}
		result, err := ioutil.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			return "failed", nil
		}
		fmt.Printf("result %s\n", result)

		return "success", nil
	}
	return "failed", nil
}

func main() {
	fmt.Printf("dangle_pay_server start\n")
	configMng := goconfig.NewConfig()
	defer configMng.Close()

	logMng := gologger.NewLog("logger", "./config/log_config.json", configMng)

	logMng.Init()
	defer logMng.Close()
	fmt.Printf("dangle_pay_server logMng Init\n")

	mng := &HttpHanderManager{
		handers: make(map[string]*HttpHander),
		config:  configMng,
		logger:  logMng.GetLogger("http"),
	}

	mng.handers["pay"] = &HttpHander{
		Run: onPay,
	}

	mng.Start()
}
