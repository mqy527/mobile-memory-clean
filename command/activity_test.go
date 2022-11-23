package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitTopActivityPackage(t *testing.T) {
	actInfo := "   ACTIVITY com.clcw.exejia/.activity.MainActivity a351b99 pid=1845"
	pkg := splitTopActivityPackage(actInfo)
	require.Equal(t, pkg, "com.clcw.exejia")
}
