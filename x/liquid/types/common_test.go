package types_test

import (
	"os"
	"testing"

	"github.com/mokitanetwork/aether/app"
)

func TestMain(m *testing.M) {
	app.SetSDKConfig()
	os.Exit(m.Run())
}
