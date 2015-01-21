package main

import (
    "fmt"
    "github.com/fanyang1988/goconfig"
    "github.com/fanyang1988/gologger"
    "net/http"
)

type HttpHander struct {
    Run func(param map[string]string) (re string, err error)
}

type HttpHanderManager struct {
    config  *goconfig.Config
    logger  *log.Logger
    handers map[string]*HttpHander
}

func (self *HttpHanderManager) OnHttp(w http.ResponseWriter, r *http.Request) {
    //
    // http://127.0.0.1:3324/KpReq?p1=p1Info&p2=p2Info
    //
    path := r.URL.Path
    method := r.Method

    if method != "GET" {
        return
    }

    switch path {
    case "/":
        r.ParseForm()
        self.logger.Info("New Http Request")
        self.logger.Info(r.Host)
        self.logger.Info(r.Form.Encode())

        params := make(map[string]string)
        for k, v := range r.Form {
            if len(v) > 0 {
                params[k] = v[0]
            }

        }
        re, err := self.OnGet("pay", params)
        if err == nil {
            fmt.Fprintf(w, re)
        }

    default:
        fmt.Fprintf(w, "NULL")
    }

    return
}

func (self *HttpHanderManager) OnGet(typ string, param map[string]string) (re string, err error) {
    handler := self.handers[typ]
    if handler == nil {
        err = fmt.Errorf("OnGet : no handler")
        self.logger.Error("OnGet : no handler")
        return
    }

    re, err = handler.Run(param)
    return
}

func (self *HttpHanderManager) Start() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        self.OnHttp(w, r)
    })

    self.config.Reg("http", "./config/http.json", false)
    http_config := self.config.Get("http")

    port, err := http_config.Get("http_port").String()
    if err != nil {
        self.logger.Error("Start : no http_port")
        return
    }

    listen_str := ":" + port
    self.logger.Error("Start " + listen_str)
    http.ListenAndServe(listen_str, nil)
}
