package sip

import (
	"io"
	"strconv"
	"strings"
)

// Uri is parsed form of
// sip:user:password@host:port;uri-parameters?headers
// In case of `sips:“ Encrypted is set to true
type Uri struct {
	// True if and only if the URI is a SIPS URI.
	Encrypted bool
	Wildcard  bool

	// The user part of the URI: the 'joe' in sip:joe@bloggs.com
	// This is a pointer, so that URIs without a user part can have 'nil'.
	User string

	// The password field of the URI. This is represented in the URI as joe:hunter2@bloggs.com.
	// Note that if a URI has a password field, it *must* have a user field as well.
	// Note that RFC 3261 strongly recommends against the use of password fields in SIP URIs,
	// as they are fundamentally insecure.
	Password string

	// The host part of the URI. This can be a domain, or a string representation of an IP address.
	Host string

	// The port part of the URI. This is optional, and can be empty.
	Port int

	// Any parameters associated with the URI.
	// These are used to provide information about requests that may be constructed from the URI.
	// (For more details, see RFC 3261 section 19.1.1).
	// These appear as a semicolon-separated list of key=value pairs following the host[:port] part.
	UriParams HeaderParams

	// Any headers to be included on requests constructed from this URI.
	// These appear as a '&'-separated list at the end of the URI, introduced by '?'.
	Headers HeaderParams
}

// Generates the string representation of a SipUri struct.
func (uri *Uri) String() string {
	var buffer strings.Builder
	uri.StringWrite(&buffer)

	return buffer.String()
}

// StringWrite writes uri string to buffer
func (uri *Uri) StringWrite(buffer io.StringWriter) {
	// Compulsory protocol identifier.
	if uri.IsEncrypted() {
		buffer.WriteString("sips")
		buffer.WriteString(":")
	} else {
		buffer.WriteString("sip")
		buffer.WriteString(":")
	}

	// Optional userinfo part.
	if uri.User != "" {
		buffer.WriteString(uri.User)
		if uri.Password != "" {
			buffer.WriteString(":")
			buffer.WriteString(uri.Password)
		}
		buffer.WriteString("@")
	}

	// Compulsory hostname.
	buffer.WriteString(uri.Host)

	// Optional port number.
	if uri.Port > 0 {
		buffer.WriteString(":")
		buffer.WriteString(strconv.Itoa(uri.Port))
	}

	if (uri.UriParams != nil) && uri.UriParams.Length() > 0 {
		buffer.WriteString(";")
		buffer.WriteString(uri.UriParams.ToString(';'))
	}

	if (uri.Headers != nil) && uri.Headers.Length() > 0 {
		buffer.WriteString("?")
		buffer.WriteString(uri.Headers.ToString('&'))
	}
}

// Clone
func (uri *Uri) Clone() *Uri {
	c := *uri
	if uri.UriParams != nil {
		c.UriParams = uri.UriParams.clone()
	}
	if uri.Headers != nil {
		c.Headers = uri.Headers.clone()
	}
	return &c
}

// IsEncrypted returns true if uri is SIPS uri
func (uri *Uri) IsEncrypted() bool {
	return uri.Encrypted
}

// Endpoint is uri user identifier. user@host[:port]
func (uri *Uri) Endpoint() string {
	addr := uri.User + "@" + uri.Host
	if uri.Port > 0 {
		addr += ":" + strconv.Itoa(uri.Port)
	}
	return addr
}

// Addr is uri part without headers and params. sip[s]:user@host[:port]
func (uri *Uri) Addr() string {
	addr := uri.User + "@" + uri.Host
	if uri.Port > 0 {
		addr += ":" + strconv.Itoa(uri.Port)
	}

	if uri.Encrypted {
		return "sips:" + addr
	}
	return "sip:" + addr
}

// HostPort represents host:port part
func (uri *Uri) HostPort() string {
	p := strconv.Itoa(uri.Port)
	// Default the SIP port if it is not provided
	if p == "0" {
		p = "5060"
	}
	return uri.Host + ":" + p
}
