package strings

// NewInitialism a factory function.
func NewInitialism(options ...InitialismOption) *Initialism {
	init := &Initialism{
		list: map[string]struct{}{
			"acl": {}, "api": {}, "ascii": {},
			"cpu": {}, "css": {},
			"dns": {},
			"eod": {}, "eof": {}, "eol": {},
			"guid": {},
			"html": {}, "http": {}, "https": {},
			"id": {}, "ip": {},
			"json": {},
			"ok":   {}, "otp": {},
			"ram": {}, "rpc": {},
			"sdk": {}, "sms": {}, "smtp": {}, "sku": {}, "sql": {}, "ssh": {},
			"tcp": {}, "tls": {}, "ttl": {},
			"udp": {}, "ui": {}, "uid": {}, "uuid": {}, "uri": {}, "url": {}, "utf": {},
			"vm":  {},
			"xml": {}, "xmpp": {}, "xsrf": {}, "xss": {},
		},
	}

	for i := range options {
		options[i](init)
	}

	return init
}
