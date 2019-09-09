package errorReport

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"testing"
)

func TestEP(t *testing.T) {
	er := fmt.Errorf("nil return")
	erUid := Er.Mode("local").SaveError(errorx.Wrap(er), map[string]interface{}{
		"note": "this is a testing case for env local",
	})
	erUid2 := Er.Mode("dev").SaveError(errorx.Wrap(er), map[string]interface{}{
		"note": "this is a testing case for env dev",
	})
	erUid3 := Er.Mode("pro").SaveError(errorx.Wrap(er), map[string]interface{}{
		"note": "this is a testing case for env pro",
	})
	fmt.Println(erUid, erUid2, erUid3)
}
