package cmdline

import (
	"flag"
	"fmt"
	"os"

	nirilof "github.com/azr4e1/niri-lof"
)

const VERSION = "v0.4.0"

func Main() int {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <appID> [<cmd>]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, `%s is a simple utility for Niri; it will
focus an open app identified by <appID>; if no app
with that ID is open, it will launch the corresponding <cmd>`, os.Args[0])
		fmt.Fprintf(os.Stderr, "\n\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  -h, -help\n")
		fmt.Fprintf(os.Stderr, "        Show this help message\n")
		fmt.Fprintf(os.Stderr, "  -s, -show\n")
		fmt.Fprintf(os.Stderr, "        Show all windows in tabular format\n")
		fmt.Fprintf(os.Stderr, "  -t, -format FORMAT\n")
		fmt.Fprintf(os.Stderr, "        Select which header to show for the window table\n")
		fmt.Fprintf(os.Stderr, "  -v, -version\n")
		fmt.Fprintf(os.Stderr, "        Show version\n")
		fmt.Fprintf(os.Stderr, "\nPositional arguments:\n")
		fmt.Fprintf(os.Stderr, "  appID    Application identifier\n")
		fmt.Fprintf(os.Stderr, "  cmd      Command to execute\n")
	}
	var version bool
	var show bool
	var format string
	flag.BoolVar(&version, "version", false, "show version")
	flag.BoolVar(&version, "v", false, "show version")
	flag.BoolVar(&show, "show", false, "show windows")
	flag.BoolVar(&show, "s", false, "show windows")
	flag.StringVar(&format, "format", "", "format windows table")
	flag.StringVar(&format, "t", "", "format windows table")
	flag.Parse()

	args := flag.Args()

	if version {
		fmt.Fprintf(os.Stdout, "version %s\n", VERSION)
		return 0
	}

	runner := NewRunner()
	if format != "" && !show {
		fmt.Fprintln(os.Stderr, "must provide -s")
		flag.Usage()
		return 1
	}
	if show {
		options, err := ValidateFormat(format)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			flag.Usage()
			return 1
		}
		windows, err := nirilof.GetWindows(runner)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			flag.Usage()
			return 1
		}
		windowsTable := FormatWindows(windows, options)
		fmt.Fprintln(os.Stdout, windowsTable)

		return 0
	}

	var appID string
	var cmd string
	if len(args) == 1 {
		appID = args[0]
		cmd = ""
	} else if len(args) == 2 {
		appID = args[0]
		cmd = args[1]
	} else {
		errMsg := "you need to provide one or two arguments: <appID> and optionally <cmd>\n"
		fmt.Fprintln(os.Stderr, errMsg)
		flag.Usage()
		return 1
	}

	err := nirilof.LaunchOrFocus(runner, appID, cmd)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return 2
	}

	return 0
}
