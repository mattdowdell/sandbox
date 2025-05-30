package otelx

import (
	"path"
	"runtime"
	"runtime/debug"
	"strings"
)

const (
	unknownVersion = "(unknown)"
)

// packageName returns the package name of the caller, using skip to ignore 0 or more stack frames
// to identify the correct caller.
func packageName(skip int) string {
	//nolint:dogsled // only need the program counter
	pc, _, _, _ := runtime.Caller(1 + skip)
	fn := runtime.FuncForPC(pc).Name()

	i := strings.LastIndexByte(fn, '/')
	if i < 0 {
		i = 0
	}

	j := strings.IndexByte(fn[i:], '.')
	if j < 0 {
		j = 0
	}

	return fn[:i+j]
}

// packageVersion returns the version of the given package, or "(unknown)" if no version could be
// identified.
func packageVersion(pkg string) string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return unknownVersion
	}

	return extractVersion(pkg, info)
}

// extractVersion returns the version of the given package using the given build info, or
// "(unknown)" if no version could be identified.
//
// Build info is passed via an argument to support unit testing with runtime/debug.ParseBuildInfo.
func extractVersion(pkg string, info *debug.BuildInfo) string {
	mods := map[string]string{}

	if isParent(info.Main.Path, pkg) {
		mods[info.Main.Path] = info.Main.Version
	}

	for _, mod := range info.Deps {
		if isParent(mod.Path, pkg) {
			mods[mod.Path] = mod.Version
		}
	}

	if len(mods) == 1 {
		for _, v := range mods {
			return v
		}
	}

	for ; pkg != "."; pkg = path.Dir(pkg) {
		if version, ok := mods[pkg]; ok {
			return version
		}
	}

	return unknownVersion
}

func isParent(mod, pkg string) bool {
	return strings.HasPrefix(pkg, mod)
}
