package sipgo

import (
	"crypto/tls"
	"log/slog"
	"net"

	"github.com/emiago/sipgo/sip"
)

type UserAgent struct {
	name        string
	hostname    string
	dnsResolver *net.Resolver
	tlsConfig   *tls.Config
	parser      *sip.Parser
	txLog       *slog.Logger
	tpLog       *slog.Logger
	tp          *sip.TransportLayer
	tx          *sip.TransactionLayer
}

type UserAgentOption func(s *UserAgent) error

// WithUserAgent changes user agent name
// Default: sipgo
func WithUserAgent(ua string) UserAgentOption {
	return func(s *UserAgent) error {
		s.name = ua
		return nil
	}
}

// WithUserAgentHostname represents FQDN of user that can be presented in From header
func WithUserAgentHostname(hostname string) UserAgentOption {
	return func(s *UserAgent) error {
		s.hostname = hostname
		return nil
	}
}

// WithUserAgentDNSResolver allows customizing default DNS resolver for transport layer
func WithUserAgentDNSResolver(r *net.Resolver) UserAgentOption {
	return func(s *UserAgent) error {
		s.dnsResolver = r
		return nil
	}
}

// WithUserAgenTLSConfig allows customizing default tls config.
func WithUserAgenTLSConfig(c *tls.Config) UserAgentOption {
	return func(s *UserAgent) error {
		s.tlsConfig = c
		return nil
	}
}

// WithUserAgentParser allows removing default behavior of parser
// You can define and remove default headers parser map and pass here.
// Only use if your benchmarks are better than default
func WithUserAgentParser(p *sip.Parser) UserAgentOption {
	return func(s *UserAgent) error {
		s.parser = p
		return nil
	}
}

// WithUserAgentTransactionLogger allows setting the logger for the transaction
func WithUserAgentTransactionLogger(logger *slog.Logger) UserAgentOption {
	return func(s *UserAgent) error {
		s.txLog = logger
		return nil
	}
}

// WithUserAgentTransactionLogger allows setting the logger for the transaction
func WithUserAgentTransportLogger(logger *slog.Logger) UserAgentOption {
	return func(s *UserAgent) error {
		s.tpLog = logger
		return nil
	}
}

// NewUA creates User Agent
// User Agent will create transport and transaction layer
// Check options for customizing user agent
func NewUA(options ...UserAgentOption) (*UserAgent, error) {
	ua := &UserAgent{
		name:        "sipgo",
		hostname:    "localhost",
		dnsResolver: net.DefaultResolver,
		parser:      sip.NewParser(),
	}

	for _, o := range options {
		if err := o(ua); err != nil {
			return nil, err
		}
	}

	ua.tp = sip.NewTransportLayer(ua.dnsResolver, ua.parser, ua.tlsConfig, sip.WithTransportLayerLogger(ua.tpLog))
	ua.tx = sip.NewTransactionLayer(ua.tp, sip.WithTransactionLayerLogger(ua.txLog))
	return ua, nil
}

func (ua *UserAgent) Close() error {
	// stop transaction layer
	ua.tx.Close()

	// stop transport layer
	return ua.tp.Close()
}

func (ua *UserAgent) Name() string {
	return ua.name
}

func (ua *UserAgent) Hostname() string {
	return ua.hostname
}

func (ua *UserAgent) TransportLayer() *sip.TransportLayer {
	return ua.tp
}

func (ua *UserAgent) TransactionLayer() *sip.TransactionLayer {
	return ua.tx
}
