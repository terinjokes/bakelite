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

	"golang.org/x/sync/semaphore"
)

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

type Platform struct {
	OS   OS
	Arch Arch
}

type kvs map[string]string

func (o kvs) Strings() []string {
	str := []string{}
	for k, v := range o {
		str = append(str, k+"="+v)
	}

	return str
}

var cgo bool

func main() {
	// TODO: enable ARM after supporting GOARM
	// TODO: probably should make this configurable…
	platforms := []Platform{
		//{OS_ANDROID, ARCH_ARM},
		{OS_DARWIN, ARCH_386},
		{OS_DARWIN, ARCH_AMD64},
		//{OS_DARWIN, ARCH_ARM},
		//{OS_DARWIN, ARCH_ARM64},
		{OS_DRAGONFLY, ARCH_AMD64},
		{OS_FREEBSD, ARCH_386},
		{OS_FREEBSD, ARCH_AMD64},
		//{OS_FREEBSD, ARCH_ARM},
		{OS_LINUX, ARCH_386},
		{OS_LINUX, ARCH_AMD64},
		//{OS_LINUX, ARCH_ARM},
		//{OS_LINUX, ARCH_ARM64},
		{OS_LINUX, ARCH_PPC64},
		{OS_LINUX, ARCH_PPC64LE},
		{OS_LINUX, ARCH_MIPS},
		{OS_LINUX, ARCH_MIPSLE},
		{OS_LINUX, ARCH_MIPS64},
		{OS_LINUX, ARCH_MIPS64LE},
		{OS_NETBSD, ARCH_386},
		{OS_NETBSD, ARCH_AMD64},
		//{OS_NETBSD, ARCH_ARM},
		{OS_OPENBSD, ARCH_386},
		{OS_OPENBSD, ARCH_AMD64},
		//{OS_OPENBSD, ARCH_ARM},
		{OS_PLAN9, ARCH_386},
		{OS_PLAN9, ARCH_AMD64},
		{OS_SOLARIS, ARCH_AMD64},
		{OS_WINDOWS, ARCH_386},
		{OS_WINDOWS, ARCH_AMD64},
	}

	flags := flag.NewFlagSet("bakelite", flag.ExitOnError)
	flags.BoolVar(&cgo, "cgo", false, "Enables cgo (BYOTC)")
	flags.Usage = func() {
		fmt.Println("bakelite - dev")
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

	var (
		parallelJobs = runtime.NumCPU()
		sem          = semaphore.NewWeighted(int64(parallelJobs))
		out          = make([]error, len(platforms)*len(packages))
		ctx          = context.Background()
	)

	fmt.Printf("info: running bakelite with %d jobs\n", parallelJobs)

	for i, platform := range platforms {
		for j, pkg := range packages {
			if err := sem.Acquire(ctx, 1); err != nil {
				log.Printf("failed to acquire semaphore: %s", err)
				break
			}

			go func(i int, j int, platform Platform, pkg string) {
				defer sem.Release(1)
				out[i*j] = build(ctx, platform, pkg)
			}(i, j, platform, pkg)
		}
	}

	if err := sem.Acquire(ctx, int64(parallelJobs)); err != nil {
		log.Printf("failed to acquire semaphore: %s", err)
	}
}

func build(ctx context.Context, platform Platform, pkg string) error {
	name := fmt.Sprintf("%s-%s-%s", filepath.Base(pkg), platform.OS, platform.Arch)

	env := kvs{
		"GOOS":   string(platform.OS),
		"GOARCH": string(platform.Arch),
		"GOROOT": os.Getenv("GOROOT"),
		"GOPATH": os.Getenv("GOPATH"),
	}

	if cgo {
		env["CGO_ENABLED"] = "1"
	} else {
		env["CGO_ENABLED"] = "0"
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.CommandContext(context.Background(), "go", "build", "-o", name, pkg)
	cmd.Env = env.Strings()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	fmt.Printf("info: Running build for %s @ %s/%s…\n", pkg, platform.OS, platform.Arch)
	err := cmd.Run()

	if err != nil {
		log.Printf("fatal: There was an error! err='%s' stdout='%s' stderr='%s'", err, stdout.String(), stderr.String())
	}

	return err
}
