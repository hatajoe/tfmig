package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

var (
	tf string
)

func main() {
	tf = os.Getenv("TFMIG_TF_PATH")
	if tf == "" {
		tf = "terraform"
	}

	cmd := exec.Command(tf, "version")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	src := flag.String("s", "", "source location path of the terraform project")
	dst := flag.String("d", "", "destination location path of the terrafrom project")
	workspace := flag.String("w", "", "terraform workspace name")

	flag.Parse()

	if *src == "" {
		log.Fatal("source location path (-s) is must be specified to use tfmig")
	}

	if *dst == "" {
		log.Fatal("destination location path (-d) is must be specified to use tfmig")
	}

	tmpDir := "/tmp/tfmig"
	if err := os.Mkdir("/tmp/tfmig", 0755); err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if *workspace != "" {
		out, err := terraform(*src, "workspace", "select", *workspace)
		if err != nil {
			errorExit(out, err)
		}
		fmt.Fprintln(os.Stdout, out)
	}

	out, err := terraform(*src, "state", "list")
	if err != nil {
		errorExit(out, err)
	}
	srcStates := strings.Split(out, "\n")

	srcState, err := selectState(srcStates)
	if err != nil {
		log.Fatal(err)
	}

	fromState := "/tmp/tfmig/a.tfstate"
	out, err = terraform(*src, "state", "pull")
	if err != nil {
		errorExit(out, err)
	}
	if err := os.WriteFile(fromState, []byte(out), 0755); err != nil {
		log.Fatal(err)
	}

	toState := "/tmp/tfmig/b.tfstate"
	out, err = terraform(*dst, "state", "pull")
	if err != nil {
		errorExit(out, err)
	}
	if err := os.WriteFile(toState, []byte(out), 0755); err != nil {
		log.Fatal(err)
	}

	out, err = terraform(*dst, "state", "mv", fmt.Sprintf("-state=%s", fromState), fmt.Sprintf("-state-out=%s", toState), srcState, srcState)
	if err != nil {
		errorExit(out, err)
	}
	fmt.Fprintln(os.Stdout, out)

	out, err = terraform(*dst, "state", "push", toState)
	if err != nil {
		errorExit(out, err)
	}
	fmt.Fprintln(os.Stdout, out)

	out, err = terraform(*dst, "state", "list")
	if err != nil {
		errorExit(out, err)
	}
	fmt.Fprintln(os.Stdout, out)

	out, err = terraform(*src, "state", "push", fromState)
	if err != nil {
		errorExit(out, err)
	}
	fmt.Fprintln(os.Stdout, out)

	out, err = terraform(*src, "state", "list")
	if err != nil {
		errorExit(out, err)
	}
	fmt.Fprintln(os.Stdout, out)
}

func terraform(project string, args ...string) (string, error) {
	buf := []string{}

	log.Println(tf, fmt.Sprintf("-chdir=%s", project), args)
	cmd := exec.Command(tf, append([]string{fmt.Sprintf("-chdir=%s", project)}, args...)...)

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	stdoutScanner := bufio.NewScanner(cmdStdout)

	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	stderrScanner := bufio.NewScanner(cmdStderr)

	go func() {
		for _, s := range []*bufio.Scanner{stdoutScanner, stderrScanner} {
			go func(scanner *bufio.Scanner) {
				for scanner.Scan() {
					buf = append(buf, scanner.Text())
				}
			}(s)
		}
	}()

	if err := cmd.Start(); err != nil {
		return strings.Join(buf, "\n"), err
	}

	if err := cmd.Wait(); err != nil {
		return strings.Join(buf, "\n"), err
	}

	return strings.Join(buf, "\n"), err
}

func errorExit(out string, err error) {
	fmt.Fprintln(os.Stderr, out)
	log.Fatal(err)
}

func selectState(state []string) (string, error) {
	idx, err := fuzzyfinder.Find(
		state,
		func(i int) string {
			return state[i]
		},
	)
	if err != nil {
		return "", err
	}
	return state[idx], nil
}
