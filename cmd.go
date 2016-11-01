package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

type exitStatusError struct {
	error
	exitCode int
}

func (err exitStatusError) exit() {
	if err.exitCode != 0 {
		os.Exit(err.exitCode)
	} else {
		log.Fatal(err)
	}
}

func runOrExit(cmd *exec.Cmd) *exitStatusError {
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v", cmd)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return &exitStatusError{exiterr, status.ExitStatus()}
			}
		} else {
			return &exitStatusError{exiterr, 1}
		}
	}

	return nil
}

func gitCmd(gitExe string, args []string) *exec.Cmd {
	if verbose {
		log.Println("Running:", gitExe, args)
	}
	cmd := exec.Command(gitExe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func findGit() (string, error) {
	return exec.LookPath("git")
}
