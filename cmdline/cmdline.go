package cmdline

import (
	"flag"
	"fmt"
	"os"

	nirilof "github.com/azr4e1/niri-lof"
)

func Main() int {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <appID> <cmd>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, `%s is a simple utility for Niri; it will
focus an open app identified by <appID>; if no app
with that ID is open, it will launch the corresponding <cmd>`, os.Args[0])
		fmt.Fprintf(os.Stderr, "\n\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  -h, -help\n")
		fmt.Fprintf(os.Stderr, "        Show this help message\n")
		fmt.Fprintf(os.Stderr, "\nPositional arguments:\n")
		fmt.Fprintf(os.Stderr, "  appID    Application identifier\n")
		fmt.Fprintf(os.Stderr, "  cmd      Command to execute\n")
	}
	flag.Parse()

	args := flag.Args()

	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "you need to provide exactly two arguments: <appID> and <cmd>\n")
		flag.Usage()
		return 1
	}

	appID := args[0]
	cmd := args[1]

	err := nirilof.LaunchOrFocus(appID, cmd)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return 2
	}

	return 0
}
