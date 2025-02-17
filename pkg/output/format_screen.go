package output

import (
	"bytes"
	"fmt"
	"github.com/4dogs-cn/TXPortMap/pkg/Ginfo/Ghttp"
	"github.com/4dogs-cn/TXPortMap/pkg/conversion"
	"github.com/fatih/color"
	"strings"
)

// formatScreen formats the output for showing on screen.
func (w *StandardWriter) formatScreen(output *ResultEvent) []byte {
	builder := &bytes.Buffer{}
	builder.WriteString(color.RedString(fmt.Sprintf("%-20s", output.Target)))
	builder.WriteString(" ")
	if output.Info.Service != "unknown" {
		builder.WriteString(color.YellowString(output.Info.Service))
	}

	if output.Info.Service == "ssl/tls" || output.Info.Service == "http" {
		if len(output.Info.Cert) > 0 {
			builder.WriteString(" [")
			builder.WriteString(color.YellowString(output.Info.Cert))
			builder.WriteString("]")
		}
	}
	if output.WorkingEvent != nil {
		switch tmp := output.WorkingEvent.(type) {
		case Ghttp.Result:
			httpBanner := tmp.ToString()
			if len(httpBanner) > 0 {
				builder.WriteString(" | ")
				builder.WriteString(color.GreenString(httpBanner))
			}
		default:
			result := conversion.ToString(tmp)
			if strings.HasPrefix(result, "\\x") == false && len(result) > 0 {
				builder.WriteString(" | ")
				builder.WriteString(color.GreenString(result))
			}
		}
	} else {
		if strings.HasPrefix(output.Info.Banner, "\\x") == false && len(output.Info.Banner) > 0 {
			r := strings.Split(output.Info.Banner, "\\x0d\\x0a")
			builder.WriteString(" | ")
			builder.WriteString(color.GreenString(r[0]))
		}
	}
	return builder.Bytes()
}

// formatScreen formats the output for showing on screen.
func (w *StandardWriter) formatSuccessScreen(output *ResultSuccess) []byte {
	builder := &bytes.Buffer{}
	builder.WriteString(color.RedString(fmt.Sprintf("%-18s", output.Target)))
	builder.WriteString(" | ")
	builder.WriteString(color.GreenString(fmt.Sprintf("%-14s", output.StepIP)))
	builder.WriteString(" | ")
	builder.WriteString(color.YellowString(output.Country))
	builder.WriteString(" | ")
	builder.WriteString(color.BlueString(fmt.Sprintf("%dms", output.Ping)))
	builder.WriteString(" | ")
	if output.Domain != "" {
		builder.WriteString(color.MagentaString(output.Domain))
		builder.WriteString(" | ")
	}
	builder.WriteString(output.Time.Format("2006-01-02 15:04:05"))
	return builder.Bytes()
}
