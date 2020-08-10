package tls

// TLS is a data structure for Kubernetes TLS configuration with command line
// flags.
type TLS struct {
	CAFile  string
	CrtFile string
	KeyFile string
}
