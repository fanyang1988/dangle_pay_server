package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    "strings"
)

const (
    AppKey = "cc5a0304d4efa71a07cde121d20fdb54"
)

func onPay(param map[string]string) (re string, err error) {
    fmt.Printf("onPay %s\n", param)
    ext, ext_ok := param["ext"]
    money, money_ok := param["money"]
    if ext_ok && money_ok {
        fmt.Printf("ext %s\n", ext)
        ext_info := strings.Split(ext, "|")
        if len(ext_info) < 4 {
            return "failed ext_info", nil
        }
        fmt.Printf("url %s\n", ext_info[0])
        fmt.Printf("uid %s\n", ext_info[1])
        fmt.Printf("gid %s\n", ext_info[2])
        fmt.Printf("money %s\n", money)

        money_float, money_float_err := strconv.ParseFloat(money, 10)
        if money_float_err != nil {
            return money_float_err.Error(), nil
        }

        u, _ := url.Parse(ext_info[0])
        q := u.Query()

        q.Set("app_key", AppKey)
        q.Set("product_id", ext_info[2])
        q.Set("amount", strconv.FormatInt(int64(money_float*100), 10))
        q.Set("app_uid", ext_info[1])
        q.Set("order_id", ext_info[3])
        q.Set("user_id", ext_info[1])
        q.Set("gateway_flag", "success")
        q.Set("app_order_id", ext_info[3])

        u.RawQuery = q.Encode()
        //TODO timeout
        url_toget := u.String()
        fmt.Printf("url %s\n", url_toget)
        res, err := http.Get(url_toget)
        if err != nil {
            return "get failed", nil
        }
        result, err := ioutil.ReadAll(res.Body)
        res.Body.Close()

        if err != nil {
            return "read failed", nil
        }
        fmt.Printf("result %s\n", result)

        return "success", nil
    }
    return "failed", nil
}
