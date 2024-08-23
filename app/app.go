package app

import (
	"fmt"
	"log"
	"net"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
)

func CreateApp(conf *config.AppConfig) (app *fiber.App) {
	var proxyHeader string
	var trustedProxies []string

	if conf.BEHIND_PROXY {
		if err := parseIPs(&conf.PROXY_IP_ADDRESSES); err != nil {
			log.Fatal(err)
		} else {
			proxyHeader = "X-Forwarded-For"
			trustedProxies = conf.PROXY_IP_ADDRESSES
		}
	}

	return fiber.New(fiber.Config{
		CaseSensitive:           true,
		EnableTrustedProxyCheck: conf.BEHIND_PROXY,
		ProxyHeader:             proxyHeader,
		TrustedProxies:          trustedProxies,
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		Prefork:                 true,
	})
}

func parseIPs(ipStrings *[]string) error {
	for _, s := range (*ipStrings) {
		if net.ParseIP(s) == nil {
			return fmt.Errorf("invalid IP: %s", s)
		}
	}

	return nil
}
