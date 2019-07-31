package config

type Config struct {
	RemoteAddr string
	LocalPort  int
	SSHVia     string
	NameSuffix string

	Purge  bool
	DryRun bool
}
