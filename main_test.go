package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParsePlatforms(t *testing.T) {
	fields := []string{"+windows", "+plan9", "-windows/386"}
	plb, _ := NewPlatformBuilder()

	plb = parsePlatforms(plb, fields)

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

func TestParsePlatformsRemoveOS(t *testing.T) {
	fields := []string{"+windows", "+plan9", "-windows"}
	plb, _ := NewPlatformBuilder()

	plb = parsePlatforms(plb, fields)

	pl := plb.Build()
	expected := []Platform{
		{OS("plan9"), Arch("386")},
		{OS("plan9"), Arch("amd64")},
	}

	if diff := cmp.Diff(expected, pl, cmpopts.SortSlices(lessPlatforms)); diff != "" {
		t.Errorf("manifest differs. (-got +want):\n%s", diff)
	}
}
func TestParsePlatformRemoveNothing(t *testing.T) {
	fields := []string{"+windows", "-"}
	plb, _ := NewPlatformBuilder()

	plb = parsePlatforms(plb, fields)

	pl := plb.Build()
	expected := []Platform{
		{OS("windows"), Arch("amd64")},
		{OS("windows"), Arch("386")},
	}

	if diff := cmp.Diff(expected, pl, cmpopts.SortSlices(lessPlatforms)); diff != "" {
		t.Errorf("manifest differs. (-got +want):\n%s", diff)
	}
}
