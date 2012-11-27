package main

import (
	"flag"
	"fmt"
	config "github.com/gokyle/goconfig"
	"github.com/gokyle/webshell"
	"os"
)

const (
	default_config = "urlshortenrc"
	DEFAULT_HOST   = ""
	DEFAULT_PORT   = "8080"
	DEFAULT_TITLE  = "urlshorten.go"
)

var (
	config_file string
	check_auth  = true
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

		if conf["server"]["development"] == "false" {
			server_dev = false
			server_secure = true
		}

		if conf["server"]["dbfile"] != "" {
			dbFile = conf["server"]["dbfile"]
		}

		if conf["server"]["authenticate"] == "false" {
			check_auth = false
		}

		if conf["server"]["admin_user"] != "" {
			admin_user = conf["server"]["admin_user"]
			ok, err := userExists(admin_user)
			if err != nil {
				panic(err)
			} else if !ok {
				panic("User does not exists.")
			}
		}

                if conf["server"]["access_log"] != "" {
                        access_logfile = conf["server"]["access_log"]
                }
	        error_logfile = access_logfile
                if conf["server"]["error_log"] != "" {
                        error_logfile = conf["server"]["error_log"]
                }
	}

	if conf["page"] == nil {
		page_title = DEFAULT_TITLE
		server_host = "localhost"
	} else {
		if conf["page"]["title"] != "" {
			page_title = conf["page"]["title"]
		} else {
			page_title = DEFAULT_TITLE
		}

		if conf["page"]["host"] != "" {
			server_host = conf["page"]["title"]
		} else {
			server_host = "localhost"
		}
	}

	if server_dev {
		server_host = fmt.Sprintf("%s:%s", server_host, webshell.SERVER_PORT)
	}

	webshell.StaticRoute("/assets/", "assets/")
	webshell.AddRoute("/", topRoute)
	if check_auth {
		webshell.AddRoute("/change", changePass)
		webshell.AddRoute("/add", addUser)
	}
	webshell.AddRoute("/views/", getViews)
	webshell.Serve(false, nil)
}

func config_server() {
	c_config_file := flag.String("c", default_config, "alternate config file")
	flag.Parse()

	config_file = *c_config_file
}
