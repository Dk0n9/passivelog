package passivelog

import (
	"context"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
	"github.com/rs/zerolog"
)

type ResponseWriter struct {
	queryTime time.Time
	query     *dns.Msg
	ctx       context.Context
	dns.ResponseWriter
	PluginPassive
}

type PluginPassive struct {
	Next   plugin.Handler
	Logger zerolog.Logger
}

func (w ResponseWriter) WriteMsg(resp *dns.Msg) error {
	err := w.ResponseWriter.WriteMsg(resp)
	if err != nil {
		return err
	}

	if len(resp.Answer) > 0 {
		for _, ans := range resp.Answer {
			go w.logRecord(ans)
		}
	}

	return nil
}

func (p PluginPassive) Name() string { return passiveLogPluginName }

func (p PluginPassive) ServeDNS(ctx context.Context, w dns.ResponseWriter, req *dns.Msg) (int, error) {
	rw := &ResponseWriter{
		PluginPassive:  p,
		ResponseWriter: w,
		query:          req,
		ctx:            ctx,
		queryTime:      time.Now(),
	}
	return plugin.NextOrFailure(p.Name(), p.Next, ctx, rw, req)
}

func (p PluginPassive) logRecord(rr dns.RR) {
	switch rr.(type) {
	case *dns.A:
		rec := rr.(*dns.A)
		p.Logger.Info().
			Str("name", rec.Hdr.Name).
			Str("type", "A").
			Str("value", rec.A.String()).
			Msg(rec.String())
		break
	case *dns.MX:
		rec := rr.(*dns.MX)
		p.Logger.Info().
			Str("name", rec.Hdr.Name).
			Str("type", "MX").
			Str("value", rec.Mx).
			Msg(rec.String())
		break
	case *dns.CNAME:
		rec := rr.(*dns.CNAME)
		p.Logger.Info().
			Str("name", rec.Hdr.Name).
			Str("type", "CNAME").
			Str("value", rec.Target).
			Msg(rec.String())
		break
	case *dns.AAAA:
		rec := rr.(*dns.AAAA)
		p.Logger.Info().
			Str("name", rec.Hdr.Name).
			Str("type", "AAAA").
			Str("value", rec.AAAA.String()).
			Msg(rec.String())
		break
	}
}
