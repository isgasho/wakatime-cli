package legacy

import (
	"errors"
	"os"

	"github.com/wakatime/wakatime-cli/cmd/legacy/configread"
	"github.com/wakatime/wakatime-cli/cmd/legacy/today"
	"github.com/wakatime/wakatime-cli/pkg/exitcode"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// Run executes legacy commands following the interface of the old python implementation of the WakaTime script.
func Run(v *viper.Viper) {
	setVerbose(v)

	if v.GetBool("version") {
		runVersion()

		os.Exit(exitcode.Success)
	}

	if err := ReadInConfig(v, ConfigFilePath); err != nil {
		jww.CRITICAL.Printf("err: %s", err)

		var cfperr ErrConfigFileParse
		if errors.As(err, &cfperr) {
			os.Exit(exitcode.ErrConfigFileParse)
		}

		os.Exit(exitcode.ErrDefault)
	}

	if v.GetString("config-read") != "" {
		configread.Run(v)
	}

	if v.GetBool("today") {
		today.Run(v)
	}

	os.Exit(exitcode.Success)
}

func setVerbose(v *viper.Viper) {
	if v.GetBool("verbose") {
		jww.SetStdoutThreshold(jww.LevelDebug)
	}
}
