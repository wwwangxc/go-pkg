package etcd

// ClientOption etcd client proxy option
type ClientOption func(*clientConfig)

// WithTarget set target
//
// Target is a list of URLs, multiple URL split by ','
func WithTarget(target string) ClientOption {
	return func(sc *clientConfig) {
		sc.Target = target
	}
}

// WithTimeout set timeout
//
// Unit millisecond, default 3000
func WithTimeout(timeout int) ClientOption {
	return func(sc *clientConfig) {
		sc.Timeout = timeout
	}
}

// WithUsername set user name for authentication
func WithUsername(username string) ClientOption {
	return func(sc *clientConfig) {
		sc.Username = username
	}
}

// WithPassword set password for authentication
func WithPassword(password string) ClientOption {
	return func(sc *clientConfig) {
		sc.Password = password
	}
}

// WithTLSKeyPath set tls key file path
func WithTLSKeyPath(tlsKeyPath string) ClientOption {
	return func(sc *clientConfig) {
		sc.TLSKeyPath = tlsKeyPath
	}
}

// WithTLSKeyPath set tls cert file path
func WithTLSCertPath(tlsCertPath string) ClientOption {
	return func(sc *clientConfig) {
		sc.TLSCertPath = tlsCertPath
	}
}

// WithCACertPath set ca cert file path
func WithCACertPath(caCertPath string) ClientOption {
	return func(sc *clientConfig) {
		sc.CACertPath = caCertPath
	}
}
