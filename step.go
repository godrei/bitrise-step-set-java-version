package main

import (
	"fmt"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-steputils/tools"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/env"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/bitrise-step-set-java-version/javasetter"
)

// Input is the Set java version step configuration
type Input struct {
	JavaVersion string `env:"set_java_version,opt[11,8]"`
}

// Config ...
type Config struct {
	javaVersion string
}

// JavaSelector ...
type JavaSelector struct {
	inputParser   stepconf.InputParser
	envRepository env.Repository
	logger        log.Logger
	cmdFactory    command.Factory
}

// NewJavaSelector ...
func NewJavaSelector(inputParser stepconf.InputParser, envRepository env.Repository, logger log.Logger, cmdFactory command.Factory) *JavaSelector {
	return &JavaSelector{inputParser: inputParser, envRepository: envRepository, logger: logger, cmdFactory: cmdFactory}
}

// ProcessConfig ...
func (j JavaSelector) ProcessConfig() (Config, error) {
	var input Input
	err := j.inputParser.Parse(&input)
	if err != nil {
		return Config{}, err
	}
	stepconf.Print(input)
	return Config{
		javaVersion: input.JavaVersion,
	}, nil
}

func (j JavaSelector) printJavaVersion() error {
	//
	// java -version
	cmd := j.cmdFactory.Create(
		"java",
		[]string{
			"-version",
		},
		nil,
	)

	j.logger.Printf("$ %s", cmd.PrintableCommandArgs())
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	j.logger.Printf(out)
	return err
}

func (j JavaSelector) printJavaCVersion() error {
	//
	// javac -version
	cmd := j.cmdFactory.Create(
		"javac",
		[]string{
			"-version",
		},
		nil,
	)

	j.logger.Println()
	j.logger.Printf("$ %s", cmd.PrintableCommandArgs())
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	j.logger.Printf(out)
	return err
}

// Run ...
func (j JavaSelector) Run(cfg Config) (javasetter.Result, error) {
	versionToSet := javasetter.JavaVersion(cfg.javaVersion)
	setter := javasetter.New(j.logger, j.cmdFactory)
	result, err := setter.SetJava(versionToSet)

	j.logger.Println()
	j.logger.Infof("Global java & javac versions the after the command run")
	if err := j.printJavaVersion(); err != nil {
		return javasetter.Result{}, err
	}
	if err := j.printJavaCVersion(); err != nil {
		return javasetter.Result{}, err
	}

	return result, err
}

// Export ...
func (j JavaSelector) Export(result javasetter.Result) error {
	if string(result.JavaHome) == "" {
		return nil
	}

	j.logger.Println()
	j.logger.Infof("Export step outputs")
	j.logger.Printf("- Exporting JAVA_HOME=%s", result.JavaHome)
	if err := tools.ExportEnvironmentWithEnvman("JAVA_HOME", result.JavaHome); err != nil {
		return fmt.Errorf("failed to export environment variable: %s", "JAVA_HOME")
	}
	return nil
}
