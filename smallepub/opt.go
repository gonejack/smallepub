package smallepub

import (
	"path/filepath"

	"github.com/alecthomas/kong"
)

type Options struct {
	About   bool `help:"About."`
	Verbose bool `short:"v" help:"Verbose printing."`
	Quality int  `short:"q" default:"40" help:"Picture compress rate/quality (1-100)"`

	EPUB []string `arg:"" optional:""`
}

func MustParseOptions() (opts Options) {
	kong.Parse(&opts,
		kong.Name("html-to-epub"),
		kong.Description("This command line converts .html to .epub with images embed"),
		kong.UsageOnError(),
	)
	if len(opts.EPUB) == 0 || opts.EPUB[0] == "*.epub" {
		opts.EPUB, _ = filepath.Glob("*.epub")
	}
	return
}
