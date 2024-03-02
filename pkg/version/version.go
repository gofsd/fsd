package version

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
)

type Version struct {
	Version                    string
	Prerelease                 string
	Build                      string
	Hash                       string
	major, minor, patch, build int
}

func FromString(v string) Version {
	var version Version
	if semver.IsValid(v) {
		version.Version = strings.Split(semver.Canonical(v), "-")[0]
		version.Prerelease = semver.Prerelease(v)
		version.Build = semver.Build(v)
		parts := strings.Split(strings.ReplaceAll(version.Version, "v", ""), ".")
		for i, vp := range parts {
			if i == 0 {
				version.major, _ = strconv.Atoi(vp)
			} else if i == 1 {
				version.minor, _ = strconv.Atoi(vp)
			} else {
				version.patch, _ = strconv.Atoi(vp)
			}
		}

		if b, e := strconv.Atoi(version.Build); e == nil {
			version.build = b
		}

	}
	return version
}

func (version Version) ToString() string {
	return fmt.Sprintf("%s%s%s", version.Version, version.Prerelease, version.Build)
}

func (version Version) GetMajor() string {
	return semver.Major(version.Version)
}

func (version Version) GetMajorMinor() string {
	return semver.MajorMinor(version.Version)
}

func (version *Version) UpBuild() {
	if build, err := strconv.Atoi(version.Build); err == nil || version.Build == "" {
		version.build = build + 1
		version.Build = fmt.Sprintf("+%d", version.build)
	}
}

func (version *Version) UpPatch() {
	version.patch = version.patch + 1
}

func (version *Version) UpMinor() {
	version.minor = version.minor + 1
}

func (version *Version) UpMajor() {
	version.major = version.major + 1
}

func (version *Version) Update() {
	var s string
	if semver.Compare(version.Version, fmt.Sprintf("v%d.%d.%d", version.major, version.minor, version.patch)) == 0 {
		s = fmt.Sprintf("v%d.%d.%d%s%s", version.major, version.minor, version.patch, version.Prerelease, version.Build)
	} else {
		s = fmt.Sprintf("v%d.%d.%d%s%s", version.major, version.minor, version.patch, version.Prerelease, version.Build)
	}
	*version = FromString(s)
}

func (version *Version) ToJson() (b []byte, e error) {
	b, e = json.Marshal(version)
	return b, e
}

func (version *Version) ToPretifiedJson() (b []byte, e error) {
	b, e = json.MarshalIndent(version, "", "  ")
	return b, e
}
