/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	c "github.com/gofsd/fsd/pkg/cmd"
	"github.com/gofsd/fsd/pkg/version"
	v "github.com/gofsd/fsd/pkg/version"

	"github.com/spf13/cobra"
)

var Output, SemVer = "string", v.FromString("v0.0.0-lc+1")

const VersionValidationError = "version validation error"

// versionCmd represents the version command
func Version(rootCmd *cobra.Command) *cobra.Command {
	var prelease, hash, build string
	cmd := &cobra.Command{}
	rootCmd.AddCommand(cmd)
	cmd.Use = "version"
	cmd.Short = "Show and/or update version info"
	cmd.Long = `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`
	cmd.RunE = func(cmd *cobra.Command, args []string) (e error) {
		SemVer = version.FromString(SemVer.Version)
		if build != "" {
			SemVer.Build = build
		}
		if prelease != "" {
			SemVer.Prerelease = prelease
		}
		if hash != "" {
			SemVer.Hash = hash
		}
		SemVer.Update()
		c.Set(cmd).
			Equal(SemVer.Version, "").
			Throw(VersionValidationError).
			Equal(Output, "byte").
			HandleB(SemVer.ToJson).
			Equal(Output, "json").
			HandleB(SemVer.ToPretifiedJson).
			Equal(Output, "string").
			JustStr(SemVer.ToString).Error()
		return e
	}

	cmd.Flags().StringVarP(&SemVer.Version, "version", "v", SemVer.Version, "Version")
	cmd.Flags().StringVar(&prelease, "pre-release", prelease, "Prelease name")
	cmd.Flags().StringVar(&hash, "hash", hash, "Git commit hash")
	cmd.Flags().StringVar(&build, "build", build, "Build number")
	cmd.Flags().StringVarP(&Output, "output", "o", Output, "Output types: json, string(default)")
	return cmd
}

func VersionUpdate(rootCmd *cobra.Command) *cobra.Command {
	var (
		version                    string
		major, minor, patch, build bool
	)
	cmd := &cobra.Command{}
	rootCmd.AddCommand(cmd)
	cmd.Use = "up"
	cmd.Short = "Update version"
	cmd.Long = `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`
	cmd.RunE = func(cmd *cobra.Command, args []string) (e error) {
		semVer := v.FromString(version)
		c.Set(cmd).
			Equal(semVer.Version, "").
			Throw(VersionValidationError).
			Equal(major, true).
			JustFn(semVer.UpMajor).
			Equal(minor, true).
			JustFn(semVer.UpMinor).
			Equal(patch, true).
			JustFn(semVer.UpPatch).
			Equal(build, true).
			JustFn(semVer.UpBuild).
			Equal(true, true).
			JustFn(semVer.Update).
			Equal(Output, "byte").
			HandleB(semVer.ToJson).
			Equal(Output, "json").
			HandleB(semVer.ToPretifiedJson).
			Equal(Output, "string").
			JustStr(semVer.ToString).
			Error()
		return e
	}

	cmd.Flags().StringVarP(&version, "version", "v", version, "Version")
	cmd.Flags().BoolVarP(&major, "major", "j", major, "Major version")
	cmd.Flags().BoolVarP(&minor, "minor", "n", minor, "Minor version")
	cmd.Flags().BoolVarP(&patch, "patch", "t", patch, "Patch version")
	cmd.Flags().BoolVarP(&build, "build", "i", build, "Build number")
	cmd.Flags().StringVarP(&Output, "output", "o", Output, "Output types: json, string(default)")
	return cmd
}

func init() {
	VersionUpdate(Version(MainCmd))
}
