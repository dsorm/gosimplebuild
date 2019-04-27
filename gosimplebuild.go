package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type gobuildvar struct {
	GOARCH string
	GOOS   string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var gobuildvars [8]gobuildvar
	gobuildvars[0] = gobuildvar{GOARCH: "386", GOOS: "windows"}
	gobuildvars[1] = gobuildvar{GOARCH: "amd64", GOOS: "windows"}
	gobuildvars[2] = gobuildvar{GOARCH: "386", GOOS: "linux"}
	gobuildvars[3] = gobuildvar{GOARCH: "amd64", GOOS: "linux"}
	gobuildvars[4] = gobuildvar{GOARCH: "arm64", GOOS: "linux"}
	gobuildvars[5] = gobuildvar{GOARCH: "arm", GOOS: "linux"}
	gobuildvars[6] = gobuildvar{GOARCH: "386", GOOS: "darwin"}
	gobuildvars[7] = gobuildvar{GOARCH: "amd64", GOOS: "darwin"}

	prefix := ""
	if len(os.Args) == 2 {
		for {
			fmt.Printf("Prefix for build outputs [%v]: ", strings.TrimSuffix(os.Args[1], ".go"))
			scanner.Scan()
			prefix = scanner.Text()
			if prefix == "" {
				prefix = strings.TrimSuffix(os.Args[1], ".go")
			} else {
				if strings.Contains(prefix, "/") {
					fmt.Println("Invalid characters, please try again.")
					continue
				}
			}
			break
		}

		for _, v := range gobuildvars {
			if len(v.GOARCH) <= 1 {
				break
			}
			fmt.Printf("Building %v_%v_%v...\n", prefix, v.GOOS, v.GOARCH)
			var executableName string
			if v.GOOS == "windows" {
				executableName = fmt.Sprintf("%v_%v_%v.exe", prefix, v.GOOS, v.GOARCH)
			} else {
				executableName = fmt.Sprintf("%v_%v_%v", prefix, v.GOOS, v.GOARCH)
			}
			cmd := exec.Command("go", "build", "-ldflags=-w", "-o="+executableName, os.Args[1])
			newEnv := append(os.Environ(), "GOARCH="+v.GOARCH, "GOOS="+v.GOOS)
			cmd.Env = newEnv
			if cmdOutput, err := cmd.CombinedOutput(); err != nil {
				fmt.Printf(string(cmdOutput))
			}
			/*
				fmt.Println("Stripping...")
				if output, err := exec.Command("strip", executableName).CombinedOutput(); err != nil {
					fmt.Printf(string(output))
				}
			*/

		}
		fmt.Println("Done!")
	} else {
		fmt.Println("Not enough arguments to run.")
	}
}
