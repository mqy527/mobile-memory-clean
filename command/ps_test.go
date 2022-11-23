package command

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPsResult(t *testing.T) {
	// USER      PID   PPID  VSIZE  RSS   WCHAN              PC  NAME
	out := `
	system    1693  756   2670596 145120 SyS_epoll_ 0000000000 S system_server
	u0_a25    4494  756   2432324 67292 SyS_epoll_ 0000000000 S com.meizu.media.music
	u0_a25    6281  756   2330812 32876 SyS_epoll_ 0000000000 S com.meizu.media.music:receiver
	u0_a37    2631  756   2328468 37436 SyS_epoll_ 0000000000 S com.flyme.systemuitools
	u0_a9     3030  756   2321820 21180 SyS_epoll_ 0000000000 S android.process.media
	root      756   1     2207328 13004 poll_sched 0000000000 S zygote64
	system    2601  756   1860544 129960 SyS_epoll_ 0000000000 S com.meizu.safe
	u0_a37    2015  756   1813004 143400 SyS_epoll_ 0000000000 S com.android.systemui
	system    2748  756   1784220 100816 SyS_epoll_ 0000000000 S com.android.settings
	system    29936 756   1763400 57592 SyS_epoll_ 0000000000 S com.meizu.battery`
	res, err := ToPsResult(out)
	require.NoError(t, err)
	fmt.Println("res:\n", res)
}

func TestToyBoxPsResult(t *testing.T) {
	// USER       PID  PPID    VSZ   RSS NAME                                                                           STIME
	out := `
system    1693   756 269193 16102 system_server                                                    2022-10-27 19:39:40
u0_a37    2015   756 181773 15093 com.android.systemui                                             2022-10-27 19:39:49
u0_a16    3063   756 178579 84632 com.meizu.flyme.launcher                                         2022-10-27 19:39:55
system    2601   756 173701 61372 com.meizu.safe                                                   2022-11-14 09:42:50
u0_a37    2120   756 171062 58608 com.android.systemui:recents                                     2022-10-27 19:39:49
u0_a25   18812   756 242572 56272 com.meizu.media.music                                            2022-11-14 15:59:11
system   29936   756 176340 53228 com.meizu.battery                                                2022-10-29 22:47:15
system   20655   756 166621 52752 com.meizu.flyme.update                                           2022-11-14 16:08:20
u0_a15   20372   757 112421 50444 com.android.email                                                2022-11-14 15:59:42
u0_a88   20729   757 115676 47768 com.meizu.net.search                                             2022-11-14 16:08:21
u0_a64   16302   756 166478 47748 com.meizu.account                                                2022-11-14 15:44:58
system   12020   756 165704 40764 com.meizu.connectivitysettings                                   2022-11-14 14:32:42
radio     2219   756 171948 40712 com.android.phone                                                2022-10-27 19:39:49
system   16159   756 168325 38400 com.meizu.safe:MzSecService                                      2022-11-14 15:44:28
u0_a95   20678   756 164472 37892 com.meizu.flyme.weather                                          2022-11-14 16:08:21
u0_a25   20710   756 233081 37128 com.meizu.media.music:receiver                                   2022-11-14 16:08:21
system    2834   756 167916 36276 com.meizu.dataservice                                            2022-10-27 19:39:55
u0_a28   30991   757 118482 35784 com.meizu.flyme.input                                            2022-10-29 23:00:31`
	res, err := ToToyBoxPsResult(out)
	require.NoError(t, err)
	fmt.Println("now: ", time.Now())
	fmt.Println("res:\n", res)
}

func TestPsResultEq(t *testing.T) {
	tmp := PsResult{NAME: "123"}
	ps1 := PsResult{NAME: "abc"}
	ps2 := PsResult{NAME: "abc"}
	require.False(t, ps1 == ps2)

	ps1 = tmp
	ps2 = tmp
	require.True(t, ps1 == ps2)
}
