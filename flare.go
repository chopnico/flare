package flare

import (
	"fmt"
	"os"
	"strings"

	"github.com/chopnico/output"
	"github.com/cloudflare/cloudflare-go"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog"
)

// a wrapper for a logger and the cloudflare api client
type Api struct {
	Client *cloudflare.API
	Logger *zerolog.Logger
}

func PrintList(i *[]interface{}, properties string) {
	var o string

	if properties == "" {
		o = output.FormatList(i, nil)
	} else {
		b := strings.Split(properties, ",")
		o = output.FormatList(i, b)
	}

	fmt.Print(o)
}

func PrintTable(data [][]string, headers []string) {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader(headers)
	t.SetAutoWrapText(false)
	t.SetAutoFormatHeaders(true)
	t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetCenterSeparator("")
	t.SetColumnSeparator("")
	t.SetRowSeparator("")
	t.SetHeaderLine(false)
	t.SetBorder(false)
	t.SetTablePadding("\t")
	t.SetNoWhiteSpace(true)
	t.AppendBulk(data)
	t.Render()
}

func PrintJson(i interface{}) {
	var a []interface{}
	a = append(a, i)
	o := output.FormatJson(&a)

	fmt.Print(o)
}
