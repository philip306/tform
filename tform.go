package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

func main() {
	insVer := CheckVersion()
	Versions(insVer).Sort()
	version := ParseMain()
	var v ConstraintStr = ConstraintStr(version)
	accepted, err := v.Parse()
	if err != nil {
		fmt.Println(err)
	}
	var selectedVer string = ""
	for _, ver := range insVer {
		if accepted.Allows(ver) {
			selectedVer = ver.String()
			break
		}
	}
	if selectedVer == "" {
		fmt.Println("Terraform version satisfying: " + version + " was not found on your system.")
		fmt.Println("Check your chocolately lib or try 'choco install terraform --version [VERSION] -my' and try again.")
		os.Exit(1)
	}

	chocoPath := "C:\\programdata\\chocolatey\\lib\\terraform." + selectedVer + "\\tools\\terraform.exe"

	var args []string
	args = append(args, chocoPath)
	flag.Parse()
	params := flag.Args()
	args = append(args, params...)

	output := fmt.Sprint("~~Executing with Terraform version: ", selectedVer)
	fmt.Println(output)
	cmd := exec.Command(chocoPath)
	cmd.Args = args
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func ParseMain() (version string) {
	module, _ := tfconfig.LoadModule(".")
	tfconstraint := module.RequiredCore[0]

	return tfconstraint
}

func Exists(name string, version string) {
	_, err := os.Stat(name)
	if err == nil {
		return
	}
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Terraform version: " + version + " was not found on your system.")
		fmt.Println("Check your chocolately lib or try 'choco install terraform --version " + version + " -my' and try again.")
		os.Exit(1)
	}
}

func CheckVersion() (ver []Version) {
	file, err := os.Open("C:\\programdata\\chocolatey\\lib\\")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// probably a better way to get the installed versions of tf without interating over all files in this dir.
	list, _ := file.Readdirnames(0)
	for _, name := range list {
		if strings.HasPrefix(name, "terraform") {
			name = strings.TrimPrefix(name, "terraform.")
			if name != "" && name != "terraform" {
				version2, err := VersionStr(name).Parse()
				if err != nil {
					fmt.Println(err)
				}
				ver = append(ver, version2)
			}
		}
	}
	return ver
}

// func Rightmost(insver []string, ver string) {
// 	re := regexp.MustCompile(`.*([0-9]+)\.([0-9]+)\.([0-9]+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+)?$`)
// 	fmt.Println(strings.TrimPrefix(re.FindStringSubmatch(ver)[0], "~> "))
// 	fmt.Println(insver)
// 	v1, err := version.NewVersion(insver[0])
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	v2, err := version.NewVersion(strings.TrimPrefix(re.FindStringSubmatch(ver)[0], "~> "))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(v1.LessThan(v2))
// 	fmt.Println()
// }

// https://github.com/hashicorp/go-version

// https://github.com/hashicorp/terraform/blob/main/internal/plugin/discovery/version.go
