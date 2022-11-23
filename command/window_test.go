package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitTopWindowPackage(t *testing.T) {
	windowInfo := "mCurrentFocus=Window{878e21e u0 com.clcw.exejia/com.clcw.exejia.activity.MainActivity}"
	pkg := splitTopWindowPackage(windowInfo)
	require.Equal(t, pkg, "com.clcw.exejia")
}
