package main

import (
        "flag"
        "fmt"
        config "github.com/gokyle/goconfig"
        "github.com/gokyle/webshell"
        "os"
)

const default_config = "urlshortenrc"

var (
        config_file string
        DEFAULT_HOST = ""
        DEFAULT_PORT = "8080"
)

func init() {
        config_server()
}

func main() {
        conf, err := config.ParseFile(config_file)
        if err != nil {
                fmt.Printf("[!] couldn't parse config file: %s\n", err.Error())
                os.Exit(1)
        }

        if conf["server"] == nil {
                webshell.SERVER_ADDR = DEFAULT_HOST
                webshell.SERVER_PORT = DEFAULT_PORT
        } else {
                if conf["server"]["port"] != "" {
                        webshell.SERVER_PORT = conf["server"]["port"]
                } else {
                        webshell.SERVER_PORT = DEFAULT_PORT
                }

                if conf["server"]["host"] != "" {
                        webshell.SERVER_ADDR = conf["server"]["host"]
                } else {
                        webshell.SERVER_ADDR = DEFAULT_HOST
                }
        }


        webshell.Serve(false, nil)
}

func config_server() {
        c_config_file := flag.String("c", default_config, "alternate config file")
        flag.Parse()

        config_file = *c_config_file
}
