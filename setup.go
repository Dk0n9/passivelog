package passivelog

import (
	"log"
	"os"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/rs/zerolog"
)

const passiveLogPluginName = "passivelog"

func init() { plugin.Register(passiveLogPluginName, setupPassiveLog) }

func setupPassiveLog(c *caddy.Controller) error {
	var fileName string
	var logger zerolog.Logger

	timeFormat := "2006-01-02 15:04:05"
	zerolog.TimeFieldFormat = timeFormat
	c.Next()
	if c.NextArg() {
		fileName = c.Val()
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		logger = zerolog.New(file).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return PluginPassive{
			Next:   next,
			Logger: logger,
		}
	})

	return nil
}
