package main

// Arch represents a Go Arch.
type Arch string

const (
	ARCH_AMD64    Arch = "amd64"
	ARCH_386      Arch = "386"
	ARCH_ARM      Arch = "arm"
	ARCH_ARM64    Arch = "arm64"
	ARCH_PPC64    Arch = "ppc64"
	ARCH_PPC64LE  Arch = "ppc64le"
	ARCH_MIPS     Arch = "mips"
	ARCH_MIPSLE   Arch = "mipsle"
	ARCH_MIPS64   Arch = "mips64"
	ARCH_MIPS64LE Arch = "mips64le"
)

// OS represents a Go OS.
type OS string

const (
	OS_ANDROID   OS = "android"
	OS_DARWIN    OS = "darwin"
	OS_DRAGONFLY OS = "dragonfly"
	OS_FREEBSD   OS = "freebsd"
	OS_LINUX     OS = "linux"
	OS_NETBSD    OS = "netbsd"
	OS_OPENBSD   OS = "openbsd"
	OS_PLAN9     OS = "plan9"
	OS_SOLARIS   OS = "solaris"
	OS_WINDOWS   OS = "windows"
)

// Platform represents an OS/Arch combination.
type Platform struct {
	OS   OS
	Arch Arch
}

// PlatformBuilder provides a fluent way to build up (or tear down) a list of
// Go platforms. The built list of platforms can be retrieved from the Build method
// of a returned builder.
type PlatformBuilder struct {
	platforms map[Platform]bool
}

// NewPlatformBuilder returns a PlatformBuilder with no platforms configured.
func NewPlatformBuilder() (*PlatformBuilder, error) {
	return &PlatformBuilder{
		platforms: make(map[Platform]bool),
	}, nil
}

// WithoutPlatform returns a new PlatformBuilder with platform removed.
func (p *PlatformBuilder) WithoutPlatform(platform Platform) *PlatformBuilder {
	pm := make(map[Platform]bool)

	for k, v := range p.platforms {
		if k != platform {
			pm[k] = v
		}
	}

	return &PlatformBuilder{
		platforms: pm,
	}
}

// WithoutOS returns a new PlatformBuilder with all platforms of the OS's
// platforms removed.
func (p *PlatformBuilder) WithoutOS(os OS) *PlatformBuilder {
	pm := make(map[Platform]bool)

	for k, v := range p.platforms {
		if k.OS != os {
			pm[k] = v
		}
	}

	return &PlatformBuilder{
		platforms: pm,
	}
}

// WithPlatform returns a new PlatformBuilder with the platform added.
func (p *PlatformBuilder) WithPlatform(platform Platform) *PlatformBuilder {
	pm := make(map[Platform]bool)

	for k, v := range p.platforms {
		pm[k] = v
	}

	pm[platform] = true

	return &PlatformBuilder{
		platforms: pm,
	}
}

// WithOS returns a new PlatformBuilder with the OS's default platforms added.
func (p *PlatformBuilder) WithOS(os OS) *PlatformBuilder {
	pm := make(map[Platform]bool)

	for k, v := range p.platforms {
		pm[k] = v
	}

	var pl []Platform
	switch os {
	case OS_DARWIN:
		pl = defaultDarwin()
	case OS_DRAGONFLY:
		pl = defaultDragonfly()
	case OS_FREEBSD:
		pl = defaultFreeBSD()
	case OS_LINUX:
		pl = defaultLinux()
	case OS_NETBSD:
		pl = defaultNetBSD()
	case OS_OPENBSD:
		pl = defaultOpenBSD()
	case OS_PLAN9:
		pl = defaultPlan9()
	case OS_SOLARIS:
		pl = defaultSolaris()
	case OS_WINDOWS:
		pl = defaultWindows()
	default:
		pl = []Platform{}
	}

	for _, k := range pl {
		pm[k] = true
	}

	return &PlatformBuilder{
		platforms: pm,
	}
}

// WithDefaults returns a new PlatformBuilder with just Bakelite's default
// platforms.
func (p *PlatformBuilder) WithDefaults() *PlatformBuilder {
	pm := make(map[Platform]bool)

	pls := []Platform{
	//{OS_ANDROID, ARCH_ARM},
	}

	pls = append(pls, defaultDarwin()...)
	pls = append(pls, defaultDragonfly()...)
	pls = append(pls, defaultFreeBSD()...)
	pls = append(pls, defaultLinux()...)
	pls = append(pls, defaultNetBSD()...)
	pls = append(pls, defaultOpenBSD()...)
	pls = append(pls, defaultPlan9()...)
	pls = append(pls, defaultSolaris()...)
	pls = append(pls, defaultWindows()...)

	for _, k := range pls {
		pm[k] = true
	}

	return &PlatformBuilder{
		platforms: pm,
	}
}

// Build returns a new list of Platforms.
func (p *PlatformBuilder) Build() []Platform {
	pl := []Platform{}

	for k, _ := range p.platforms {
		pl = append(pl, k)
	}

	return pl
}
