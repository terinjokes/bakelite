package main

import (
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func ExampleNewPlatformBuilder() {
	plb, _ := NewPlatformBuilder()

	plb = plb.
		WithOS(OS("windows")).
		WithoutPlatform(Platform{OS("windows"), Arch("amd64")}).
		WithPlatform(Platform{OS("plan9"), Arch("amd64")})

	pl := plb.Build()
	sort.SliceStable(pl, func(i, j int) bool {
		return lessPlatforms(pl[i], pl[j])
	})

	fmt.Println(pl)

	// Output:
	// [{plan9 amd64} {windows 386}]
}

func TestWithOS(t *testing.T) {
	plb, _ := NewPlatformBuilder()

	plb = plb.WithOS(OS("windows"))

	pl := plb.Build()
	expected := defaultWindows()

	if diff := cmp.Diff(expected, pl, cmpopts.SortSlices(lessPlatforms)); diff != "" {
		t.Errorf("manifest differs. (-got +want):\n%s", diff)
	}
}

func TestWithoutOS(t *testing.T) {
	plb, _ := NewPlatformBuilder()

	plb = plb.
		WithPlatform(
			Platform{OS("windows"), Arch("amd64")},
		).
		WithPlatform(
			Platform{OS("plan9"), Arch("amd64")},
		).
		WithoutOS(OS("windows"))

	pl := plb.Build()
	expected := []Platform{{OS("plan9"), Arch("amd64")}}

	if diff := cmp.Diff(expected, pl, cmpopts.SortSlices(lessPlatforms)); diff != "" {
		t.Errorf("manifest differs. (-got +want):\n%s", diff)
	}
}

func TestWithPlatform(t *testing.T) {
	plb, _ := NewPlatformBuilder()

	plb = plb.WithPlatform(Platform{OS("windows"), Arch("arm64")})

	pl := plb.Build()
	expected := []Platform{{OS("windows"), Arch("arm64")}}

	if diff := cmp.Diff(expected, pl, cmpopts.SortSlices(lessPlatforms)); diff != "" {
		t.Errorf("manifest differs. (-got +want):\n%s", diff)
	}
}

func TestWithoutPlaform(t *testing.T) {
	plb, _ := NewPlatformBuilder()

	plb = plb.
		WithOS(OS("windows")).
		WithOS(OS("plan9")).
		WithoutPlatform(Platform{OS("windows"), Arch("386")})

	pl := plb.Build()
	expected := []Platform{
		{OS("windows"), Arch("amd64")},
		{OS("plan9"), Arch("386")},
		{OS("plan9"), Arch("amd64")},
	}

	if diff := cmp.Diff(expected, pl, cmpopts.SortSlices(lessPlatforms)); diff != "" {
		t.Errorf("manifest differs. (-got +want):\n%s", diff)
	}
}

func TestWithDefaults(t *testing.T) {
	plb, _ := NewPlatformBuilder()

	plb = plb.WithDefaults()

	pl := plb.Build()

	if len(pl) == 0 {
		t.Errorf("expected more defaults!")
	}
}

func lessPlatforms(x, y Platform) bool {
	if x.OS < y.OS {
		return true
	}
	if x.OS > y.OS {
		return false
	}
	return x.Arch < y.Arch
}
