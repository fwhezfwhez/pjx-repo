package errorReport

import (
	"github.com/fwhezfwhez/errorx"
)

var Er *errorx.Reporter

const (
	mode = "dev"
)

func init() {
	Er = errorx.NewReporter(mode)
	Er.SetContextName("request")
	Er.AddURL("dev", "https://xyx.zonst.com/dev/err/error/").
		AddURL("pro", "https://xyx.zonst.com/err/error/")

	Er.AddModeHandler("dev", Er.Mode("dev").ReportURLHandler)
	Er.AddModeHandler("pro", Er.Mode("pro").ReportURLHandler)
	Er.AddModeHandler("local", errorx.DefaultHandler)
}

// If you don't want to add extra dependency,this allow caller to call it in outer place, and it will replace init().
//
//func Init(mode string) {
//	Er = errorx.NewReporter(mode)
//	Er.SetContextName("request")
//	Er.AddURL("dev", "https://xyx.zonst.com/dev/err/error/").
//		AddURL("pro", "https://xyx.zonst.com/err/error/")
//
//	Er.AddModeHandler("dev", Er.Mode("dev").ReportURLHandler)
//	Er.AddModeHandler("pro", Er.Mode("pro").ReportURLHandler)
//	Er.AddModeHandler("local", errorx.DefaultHandler)
//}
