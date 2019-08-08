package config

type Config struct {
	RemoteUser string
	RemoteIP   string

	LocalPort  int
	SSHVia     string
	NameSuffix string

	Purge  bool
	DryRun bool
	Force  bool

	PrintSSHForwarding bool
}
