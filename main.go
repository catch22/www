package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func usage() {
	fmt.Printf("Usage: %s [DIR]\n", os.Args[0])
	// flag.PrintDefaults()
}

func getWorkdir() string {
	switch len(flag.Args()) {
	case 0:
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		return dir
	case 1:
		return flag.Args()[0]
	default:
		usage()
		os.Exit(2)
		panic("unreachable")
	}
}

func gitGetCurrentBranch(workdir string) string {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = workdir
	cmdOutput, err := cmd.Output()
	if err != nil {
		log.Fatal("Unable to determine current branch (not a git repository?)")
	}
	return strings.TrimSpace(string(cmdOutput))
}

func gitGetBranchRemote(workdir string, branch string) string {
	cmd := exec.Command("git", "config", "--get", fmt.Sprintf("branch.%s.remote", branch))
	cmd.Dir = workdir
	cmdOutput, err := cmd.Output()
	if err != nil {
		log.Fatal("Unable to determine remote for current branch")
	}
	return strings.TrimSpace(string(cmdOutput))
}

func gitGetRemoteURL(workdir string, remote string) string {
	cmd := exec.Command("git", "config", "--get", fmt.Sprintf("remote.%s.url", remote))
	cmd.Dir = workdir
	cmdOutput, err := cmd.Output()
	if err != nil {
		log.Fatal("Unable to determine URL for remote")
	}
	return strings.TrimSpace(string(cmdOutput))
}

func main() {
	log.SetFlags(0)

	// determine working directory
	flag.Usage = usage
	flag.Parse()
	workdir := getWorkdir()

	// get name of current branch
	branch := gitGetCurrentBranch(workdir)
	remote := gitGetBranchRemote(workdir, branch)
	url := gitGetRemoteURL(workdir, remote)

	// map repository URL to web URL
	const GithubPrefix = "git@github.com:"
	const OverleafPrefix = "https://git.overleaf.com/"
	const OverleafPrefixNew = "https://git@git.overleaf.com/"
	if strings.HasPrefix(url, GithubPrefix) {
		path := strings.TrimSuffix(strings.TrimPrefix(url, GithubPrefix), ".git")
		url = "https://github.com/" + path
	} else if strings.HasPrefix(url, OverleafPrefix) {
		id := strings.TrimPrefix(url, OverleafPrefix)
		url = "https://www.overleaf.com/project/" + id
	} else if strings.HasPrefix(url, OverleafPrefixNew) {
		id := strings.TrimPrefix(url, OverleafPrefixNew)
		url = "https://www.overleaf.com/project/" + id
	} else {
		log.Fatal("Unsupported remote URL: ", url)
	}

	// open URL	in browser
	log.Printf("%s => %s", workdir, url)
	err := exec.Command("open", url).Run()
	if err != nil {
		log.Fatal("Unable to open URL")
	}
}
