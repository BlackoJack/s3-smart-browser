package version

import (
	"fmt"
	"runtime"
	"strings"
)

// Эти переменные будут установлены во время сборки через ldflags
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
)

// GetVersion возвращает информацию о версии
func GetVersion() map[string]string {
	return map[string]string{
		"version":    Version,
		"git_commit": GitCommit,
		"build_time": BuildTime,
		"go_version": GoVersion,
	}
}

// GetVersionString возвращает строку версии для отображения
func GetVersionString() string {
	if Version == "dev" {
		return fmt.Sprintf("dev-%s", GitCommit[:8])
	}
	return Version
}

// IsRelease проверяет, является ли это релизной версией
func IsRelease() bool {
	return Version != "dev" && !strings.Contains(Version, "-")
}
