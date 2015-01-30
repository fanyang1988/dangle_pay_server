package main

import (
    "github.com/fanyang1988/goconfig"
    "github.com/fanyang1988/gologger"
    _ "github.com/icattlecoder/godaemon"
)

var (
    ConfigManager = goconfig.New()
    LogManager    = gologger.New("logger", "./config/log.json", ConfigManager)
)

func init() {
}

func main() {
    defer ConfigManager.Close()
    defer LogManager.Close()

    mng := &HttpHanderManager{
        handers: make(map[string]*HttpHander),
    }

    mng.handers["pay"] = &HttpHander{
        Run: onPay,
    }

    mng.Start()
}
