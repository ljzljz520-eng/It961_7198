package transport

import (
	"flag"
	"fmt"
)

type Options struct {
	Database string
	Format   string
	Campus   string
	Subject  string
	Version  string
}

func ParseOptions(args []string) (Options, error) {
	flags := flag.NewFlagSet("drivehub", flag.ContinueOnError)
	flags.SetOutput(nil)
	options := Options{Database: "drivehub.db", Format: "table"}
	flags.StringVar(&options.Database, "db", options.Database, "database")
	flags.StringVar(&options.Format, "format", options.Format, "output format")
	flags.StringVar(&options.Campus, "campus", "", "campus")
	flags.StringVar(&options.Subject, "subject", "", "subject")
	flags.StringVar(&options.Version, "version", "", "version")
	if err := flags.Parse(args); err != nil {
		return options, err
	}
	if options.Format != "table" && options.Format != "json" && options.Format != "csv" {
		return options, fmt.Errorf("unsupported format %s", options.Format)
	}
	return options, nil
}

func (o Options) Filter() map[string]string {
	return map[string]string{"campus": o.Campus, "subject": o.Subject, "version": o.Version}
}
