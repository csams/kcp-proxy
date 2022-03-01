package cmd

import (
	"github.com/csams/kcp-proxy/pkg/kcpproxy"
	"github.com/spf13/cobra"
)

var (
	server  = kcpproxy.Server{}
	rootCmd = &cobra.Command{
		Use: "kcp-proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.Serve()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&server.ListenAddress, "listen-address", ":8083", "Address and port for the proxy to listen on")

	rootCmd.PersistentFlags().StringVar(&server.ClientCACert, "client-ca-cert", "certs/ca-cert.pem", "CA cert used to validate client certs")
	rootCmd.MarkFlagRequired("client-ca-cert")

	rootCmd.PersistentFlags().StringVar(&server.ServerCertFile, "server-cert-file", "certs/server-cert.pem", "The proxy's serving cert file")
	rootCmd.MarkFlagRequired("server-cert-file")

	rootCmd.PersistentFlags().StringVar(&server.ServerKeyFile, "server-key-file", "certs/server-key.pem", "The proxy's serving private key file")
	rootCmd.MarkFlagRequired("server-key-file")

	rootCmd.PersistentFlags().StringVar(&server.MappingFile, "mapping-file", "", "Config file mapping paths to backends")
	rootCmd.MarkFlagRequired("mapping-file")
}
