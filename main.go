package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	bflag "github.com/terinjokes/bakelite/internal/flag"
	"golang.org/x/sync/semaphore"
)

type kvs map[string]string

func (o kvs) Strings() []string {
	str := []string{}
	for k, v := range o {
		str = append(str, k+"="+v)
	}

	return str
}

var cgo bool
var ldflags string
var platformFields []string

func main() {
	flags := flag.NewFlagSet("bakelite", flag.ExitOnError)
	flags.BoolVar(&cgo, "cgo", false, "enables cgo (may require your own toolchain).")
	flags.StringVar(&ldflags, "ldflags", "", "arguments to pass on each go tool compile invocation.")
	flags.Var((*bflag.StringsValue)(&platformFields), "platforms", "modify the list of platforms built")
	flags.Usage = func() {
		fmt.Println("usage: bakelite [build flags] [packages]")
		fmt.Println(`
Bakelite compiles the packages named by the import paths for multiple GOOS and
GOARCH combinations. It does not install their results.

When compiling a package, Bakelite writes the result to output files named
after the source directory in the form "$package_$goos_$goarch". The '.exe'
suffix is added when writing a Windows executable.

Multiple packages may be given to Bakelite, the result of each are saved as
described in the preceding paragraph.

The build flags recognized by Bakelite:

	-ldflags 'flag list'
		arguments to pass on each go tool compile invocation.

The Bakelite specific flags:

	-cgo
		passes CGO_ENABLED=1 to the build environment.
		May require a build toolchain for each GOOS and GOARCH combination.
	-platforms 'platform list'
		modify the built platforms.
		Platforms are prefixed with "-" to remove from the set and "+" to add
		to the set. They can be specified sparsely as just the OS, or as a
		complete GOOS/GOARCH declaration. If the special platform "-" is
		provided as the first platform, the default set is disabled. See below
		for the default list of platforms.

By default Bakelite builds for the following platforms:

	darwin/386
	darwin/amd64
	dragonfly/amd64
	freebsd/386
	freebsd/amd64
	linux/386
	linux/amd64
	linux/ppc64
	linux/ppc64le
	linux/mips
	linux/mipsle
	linux/mips64
	linux/mips64le
	netbsd/386
	netbsd/amd64
	openbsd/386
	openbsd/amd64
	plan9/386
	plan9/amd64
	solaris/amd64
	windows/386
	windows/amd64

All the flags that take a list of arguments accept a space-separated
list of strings.

For more about specifying packages, see 'go help packages'.
For more about calling between Go and C/C++, run 'go help c'.

See also: go build, go install, go clean.
			`)
	}
	if err := flags.Parse(os.Args[1:]); err != nil {
		flags.Usage()
		os.Exit(-1)
	}

	packages := flags.Args()
	if len(packages) == 0 {
		fmt.Println("fatal: Expected at least one package.")
		os.Exit(-1)
	}

	plBuilder, _ := NewPlatformBuilder()
	if len(platformFields) > 0 && platformFields[0] == "-" {
		platformFields = platformFields[1:]
	} else {
		plBuilder = plBuilder.WithDefaults()
	}
	if len(platformFields) != 0 {
		plBuilder = parsePlatforms(plBuilder, platformFields)
	}

	platforms := plBuilder.Build()

	environ := parseEnvironment(os.Environ())

	var (
		parallelJobs = runtime.NumCPU()
		sem          = semaphore.NewWeighted(int64(parallelJobs))
		ctx          = context.Background()
	)

	fmt.Printf("info: running bakelite with %d jobs\n", parallelJobs)

	var errored bool
	for _, platform := range platforms {
		for _, pkg := range packages {
			if err := sem.Acquire(ctx, 1); err != nil {
				log.Printf("failed to acquire semaphore: %s", err)
				errored = true
				break
			}

			go func(platform Platform, pkg string) {
				defer sem.Release(1)
				err := build(ctx, environ, platform, pkg)

				if err != nil {
					errored = true
				}
			}(platform, pkg)
		}
	}

	if err := sem.Acquire(ctx, int64(parallelJobs)); err != nil {
		log.Printf("failed to acquire semaphore: %s", err)
		errored = true
	}

	if errored {
		os.Exit(1)
	}
}

func build(ctx context.Context, environ kvs, platform Platform, pkg string) error {
	name := fmt.Sprintf("%s-%s-%s", filepath.Base(pkg), platform.OS, platform.Arch)

	if platform.OS == OS_WINDOWS {
		name += ".exe"
	}

	env := kvs{}
	for key, val := range environ {
		env[key] = val
	}

	env["GOOS"] = string(platform.OS)
	env["GOARCH"] = string(platform.Arch)

	if cgo {
		env["CGO_ENABLED"] = "1"
	} else {
		env["CGO_ENABLED"] = "0"
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	args := []string{
		"build",
		"-o",
		name,
	}

	if ldflags != "" {
		args = append(args, "-ldflags", ldflags)
	}

	args = append(args, pkg)

	cmd := exec.CommandContext(context.Background(), "go", args...)
	cmd.Env = env.Strings()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	fmt.Printf("info: Running build for %s @ %s/%s…\n", pkg, platform.OS, platform.Arch)
	err := cmd.Run()

	if err != nil {
		log.Printf("fatal: There was an error! goos='%s' goarch='%s' err='%s' stdout='%s' stderr='%s'", platform.OS, platform.Arch, err, stdout.String(), stderr.String())
	}

	return err
}

func parsePlatforms(plBuilder *PlatformBuilder, fields []string) *PlatformBuilder {
	for _, f := range fields {
		switch f[0] {
		case '-':
			if strings.ContainsRune(f, '/') {
				sp := strings.Split(f[1:], "/")
				p := Platform{
					OS:   OS(sp[0]),
					Arch: Arch(sp[1]),
				}

				plBuilder = plBuilder.WithoutPlatform(p)
			} else {
				plBuilder = plBuilder.WithoutOS(OS(f[1:]))
			}
		case '+':
			if strings.ContainsRune(f, '/') {
				sp := strings.Split(f[1:], "/")
				p := Platform{
					OS:   OS(sp[0]),
					Arch: Arch(sp[1]),
				}

				plBuilder = plBuilder.WithPlatform(p)
			} else {
				plBuilder = plBuilder.WithOS(OS(f[1:]))
			}
		}
	}

	return plBuilder
}

func parseEnvironment(environ []string) kvs {
	env := kvs{}
	for _, s := range environ {
		split := strings.SplitN(s, "=", 2)
		env[split[0]] = split[1]
	}

	return env
}
