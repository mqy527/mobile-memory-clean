package decide

import (
	"testing"

	"github.com/mqy527/mobile-memory-clean/command"
	"github.com/mqy527/mobile-memory-clean/log"
	"github.com/stretchr/testify/require"
)

func TestSortByRecentUsePkgs(t *testing.T) {
	// USER       PID  PPID    VSZ   RSS NAME                                                                           STIME
	out := `
u0_a138  13916   756 412119 47970 com.autonavi.minimap                                             2022-11-21 16:17:53
u0_a109  22406   756 282679 26295 com.sina.weibo                                                   2022-11-21 17:41:07
u0_a124  17802   756 701719 16462 com.xiaomi.smarthome                                             2022-11-21 17:20:26
u0_a109  22875   756 253260 97068 com.sina.weibo.image                                             2022-11-21 17:41:16
u0_a37    4788   756 176714 94456 com.android.systemui                                             2022-11-21 16:12:05
u0_a109  22532   756 254296 94040 com.sina.weibo:remote                                            2022-11-21 17:41:09
u0_a140  21247   757 148772 86324 com.tencent.tim                                                  2022-11-21 17:40:51
u0_a138  19106   756 247085 73152 com.autonavi.minimap:locationservice                             2022-11-21 17:21:58
u0_a16    5438   756 171972 66768 com.meizu.flyme.launcher                                         2022-11-20 16:23:50
u0_a37    4851   756 170978 63852 com.android.systemui:recents                                     2022-11-21 16:12:05
u0_a140  22181   757 151234 55064 com.tencent.tim:mail                                             2022-11-21 17:41:04
u0_a140  21305   757 134029 44568 com.tencent.tim:MSF                                              2022-11-21 17:40:53
u0_a28   24987   757 117683 41124 com.meizu.flyme.input                                            2022-11-21 15:03:07
u0_a138  15788   756 244662 28952 com.autonavi.minimap:sandboxed_privilege_process0                2022-11-21 16:18:58
u0_a19    2849   756 166220 21364 com.meizu.location                                               2022-10-27 19:39:55
u0_a37   17400   756 232063 20908 com.flyme.systemuitools                                          2022-11-20 19:11:05
u0_a64   23427   756 167807 20904 com.meizu.account                                                2022-11-19 22:25:50
u0_a9     3030   756 232182 16596 android.process.media                                            2022-10-27 19:39:55
u0_a136  16490   756 163492 11944 com.example.cpuactive                                            2022-11-20 19:07:21
u0_a12     956   756 162310 11604 android.ext.services                                             2022-11-20 23:53:33
u0_a1     2786   756 163555 10768 com.android.incallui                                             2022-10-27 19:39:55
u0_a76    2867   756 163352 10360 com.oma.drm.server                                               2022-10-27 19:39:55
u0_a124  18680 17802 378735  7444 /data/app/com.xiaomi.smarthome-1/lib/arm64//libweexjsb.so        2022-11-21 17:20:53`
	ps, err := command.ToToyBoxPsResult(out)
	require.NoError(t, err)

	sortByStartTimeAsc(ps)

	recentPackages := []string{"com.sina.weibo", "com.meizu.flyme.launcher", "com.xiaomi.smarthome", "com.tencent.tim", "com.autonavi.minimap", "com.android.systemui", "com.clcw.exejia", "com.android.settings", "com.galaxy.stock", "com.meizu.net.search", "com.coolapk.market", "cn.xuexi.android", "com.beisen.italent", "com.finhub.fenbeitong", "com.meizu.media.reader", "com.ss.android.article.lite", "com.meizu.flyme.flymebbs", "com.android.browser", "com.meizu.battery", "com.android.chrome", "com.example.cpuactive"}
	res := sortByRecentUsePkgs(ps, recentPackages)

	logger := log.GetLogger("TestSortByRecentUsePkgs")
	logger.Info("res:")
	for i := 0; i < len(res); i++ {
		logger.Info(res[i].NAME, ",", res[i].STIME)
	}
}
