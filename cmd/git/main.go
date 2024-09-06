package main

import (
	"flag"
	"os"

	"github.com/nice-pink/goutil/pkg/repo"
)

func main() {
	repository := flag.String("repo", "git@github.com:nice-pink/git-test.git", "")
	clone := flag.Bool("clone", false, "")
	tag := flag.String("tag", "", "")
	push := flag.Bool("push", false, "")
	filename := flag.String("filename", "", "")
	flag.Parse()

	sshKeyPath := os.Getenv("SSH_KEY_PATH")
	repositoryUser := os.Getenv("REPO_USER")
	repositoryEmail := os.Getenv("REPO_EMAIL")
	g := repo.NewRepoHandle(sshKeyPath, repositoryUser, repositoryEmail)

	if *clone {
		g.Clone(*repository, "bin/repo", "main", true, false)
	}

	if *tag != "" {
		g.TagRepo(*tag, *tag, *push)
		return
	}

	if *push {
		g.ResetToRemoteHead("bin/repo")
		g.PullLocalRepo("bin/repo")
		filename := "bin/repo/" + *filename
		file, _ := os.Create(filename)
		file.Close()
		g.CommitPushLocalRepo("bin/repo", "commit push", true)
	}

}
