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

var defaultProjPath string = "."
var defaultLibPath string = "C:\\programdata\\chocolatey\\lib\\"

func main() {

	err := ExecuteTform(defaultProjPath, defaultLibPath, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ExecuteTform(projPath string, libPath string, outDisable bool) (e error) {
	selectedVer, err := ChooseVer(projPath, libPath)
	if err != nil {
		return err
	}

	chocoPath := "C:\\programdata\\chocolatey\\lib\\terraform." + selectedVer + "\\tools\\terraform.exe"

	// setup command with args to be executed
	var args []string
	args = append(args, chocoPath)
	flag.Parse()
	params := flag.Args()
	args = append(args, params...)

	output := fmt.Sprint("~~Executing with Terraform version: ", selectedVer)
	fmt.Println(output)
	cmd := exec.Command(chocoPath)
	cmd.Args = args
	if !outDisable {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.Stdin = os.Stdin
	cmd.Run()

	return nil
}

func ChooseVer(projpath string, chocolib string) (selVer string, e error) {

	// parse local tf project to get version contrain
	if !Exists(projpath) {
		return selVer, errors.New("File path: " + projpath + " not found on system")
	}
	version := ParseMain(projpath)
	var v ConstraintStr = ConstraintStr(version)
	accepted, err := v.Parse()
	if err != nil {
		return selVer, err
	}

	// get locally installed versions of tf
	insVer, err := CheckInstVersion(chocolib)
	if err != nil {
		return selVer, err
	}
	Versions(insVer).Sort()

	// compare version contrains to locally isntalled versions to get latest version that satisfies contraint
	selVer = ""
	for _, ver := range insVer {
		if accepted.Allows(ver) {
			selVer = ver.String()
			break
		}
	}

	if selVer == "" {
		// fmt.Println("Terraform version satisfying: " + version + " was not found on your system.")
		// fmt.Println("Check your chocolately lib or try 'choco install terraform --version [VERSION] -my' and try again.")
		e := errors.New("No version satisfying constraint: " + version + " found on system.  Check your chocolately lib or try 'choco install terraform --version [VERSION] -my' and try again.")
		return selVer, e
	}

	return selVer, e
}

// returns the version constraint from the tf project that tform is being executed from
func ParseMain(path string) (version string) {
	module, _ := tfconfig.LoadModule(path)
	tfconstraint := module.RequiredCore[0]

	return tfconstraint
}

// returns list of installed versions of tf
func CheckInstVersion(chocolib string) (ver []Version, e error) {
	file, err := os.Open(chocolib)
	if err != nil {
		return ver, err
	}
	defer file.Close()

	// probably a better way to get the installed versions of tf without iterating over all files in this dir.
	list, _ := file.Readdirnames(0)
	for _, name := range list {
		if strings.HasPrefix(name, "terraform") {
			name = strings.TrimPrefix(name, "terraform.")
			if name != "" && name != "terraform" {
				version2, err := VersionStr(name).Parse()
				if err != nil {
					return ver, err
				}
				ver = append(ver, version2)
			}
		}
	}
	return ver, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else {
		return false
	}

}
