package app

import (
	"fmt"
	"log"
	"net"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/routes"
)

func CreateApp(conf *config.AppConfig) (app *fiber.App, dbs *databases.Databases) {
	var proxyHeader string
	var trustedProxies []string

	if conf.GO_FIBER_BEHIND_PROXY {
		if err := parseIPs(&conf.GO_FIBER_PROXY_IP_ADDRESSES); err != nil {
			log.Fatal(err)
		} else {
			proxyHeader = "X-Forwarded-For"
			trustedProxies = conf.GO_FIBER_PROXY_IP_ADDRESSES
		}
	}

	app = fiber.New(fiber.Config{
		CaseSensitive:           true,
		EnableTrustedProxyCheck: conf.GO_FIBER_BEHIND_PROXY,
		ProxyHeader:             proxyHeader,
		TrustedProxies:          trustedProxies,
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		Prefork:                 true,
	})

	dbs = databases.Init(conf)
	routes.Register(app, dbs, conf)

	return app, dbs
}

func parseIPs(ipStrings *[]string) error {
	for _, s := range (*ipStrings) {
		if net.ParseIP(s) == nil {
			return fmt.Errorf("invalid IP: %s", s)
		}
	}

	return nil
}
