package main

import (
	"flag"
	"os"

	"github.com/nice-pink/goutil/pkg/repo"
)

func main() {
	clone := flag.Bool("clone", false, "")
	push := flag.Bool("push", false, "")
	filename := flag.String("filename", "", "")
	flag.Parse()

	sshKeyPath := os.Getenv("SSH_KEY_PATH")
	repository := os.Getenv("REPO")
	repositoryUser := os.Getenv("REPO_USER")
	repositoryEmail := os.Getenv("REPO_EMAIL")
	g := repo.NewGitRepo(sshKeyPath, repositoryUser, repositoryEmail)

	if *clone {
		g.Clone(repository, "bin/repo", "main", true, false)
	} else if *push {
		g.ResetToRemoteHead("bin/repo")
		g.PullLocalRepo("bin/repo")
		filename := "bin/repo/" + *filename
		file, _ := os.Create(filename)
		file.Close()
		g.CommitPushLocalRepo("bin/repo", "commit push", true)
	}

}
