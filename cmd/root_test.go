package cmd

import (
	"fmt"
	"math/big"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBigPercentOf(t *testing.T) {
	Convey("Given some big ints", t, func() {
		x := big.NewInt(860160)
		y := big.NewInt(14343298)
		z := bigPercentOf(x, y)
		perc := fmt.Sprintf("%0.2f%%", z)
		So(z, ShouldEqual, 5)
		So(perc, ShouldEqual, "5.00%")
	})
}
