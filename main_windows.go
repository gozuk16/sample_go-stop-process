import (
	"os"
)

func init() {
	stopSignal = os.Interrupt
}
