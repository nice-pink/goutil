package main

import (
	"flag"
	"os"

	"github.com/nice-pink/goutil/pkg/git"
)

func main() {
	clone := flag.Bool("clone", false, "")
	push := flag.Bool("push", false, "")
	filename := flag.String("filename", "", "")
	flag.Parse()

	sshKeyPath := os.Getenv("SSH_KEY_PATH")
	repo := os.Getenv("REPO")
	repoUser := os.Getenv("REPO_USER")
	repoEmail := os.Getenv("REPO_EMAIL")
	g := git.NewGit(sshKeyPath, repoUser, repoEmail)

	if *clone {
		g.Clone(repo, "bin/repo", "main", true, false)
	} else if *push {
		g.ResetToRemoteHead("bin/repo")
		g.PullLocalRepo("bin/repo")
		filename := "bin/repo/" + *filename
		file, _ := os.Create(filename)
		file.Close()
		g.CommitPushLocalRepo("bin/repo", "commit push", true)
	}

}
