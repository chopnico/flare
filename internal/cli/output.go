package cli

import (
	"fmt"
	"strings"

	"github.com/chopnico/output"
)

func writeOutput(i *[]interface{}, p, o string) {
	switch o {
	case "list":
		if p == "" {
			fmt.Printf(output.FormatList(i, nil))
		} else {
			properties := strings.Split(p, ",")
			fmt.Printf(output.FormatList(i, properties))
		}
	case "json":
		fmt.Printf(output.FormatJson(i))
	case "pretty-json":
		fmt.Printf(output.FormatPrettyJson(i))
	}
}
