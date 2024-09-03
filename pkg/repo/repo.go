package repo

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/nice-pink/goutil/pkg/log"
)

type GitRepo struct {
	userName   string
	userEmail  string
	sshKeyPath string
}

func NewGitRepo(sshKeyPath string, userName string, userEmail string) *GitRepo {
	return &GitRepo{
		userName:   userName,
		userEmail:  userEmail,
		sshKeyPath: sshKeyPath,
	}
}

// NOTE:
// For save usage do:
// - Clone() (shallow)
// - PullLocalRepo()
// - Do your changes to the repo
// - CommitPushLocalRepo()

// Clone
func (g *GitRepo) Clone(url string, dest string, branch string, shallow bool, repoSubfolder bool) error {
	// set clone options
	cloneOpt := &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	// setup ssh auth
	if g.sshKeyPath != "" {
		auth, err := ssh.NewPublicKeysFromFile("git", g.sshKeyPath, "")
		if err != nil {
			log.Err(err, "key")
			return err
		}
		cloneOpt.Auth = auth
	}

	if branch != "" {
		cloneOpt.ReferenceName = plumbing.NewBranchReferenceName(branch)
	}

	if shallow {
		cloneOpt.Depth = 1
		cloneOpt.SingleBranch = true
		cloneOpt.ShallowSubmodules = true
	}

	// clone repo
	if repoSubfolder {
		path := strings.Split(url, "/")
		dest = filepath.Join(dest, path[len(path)-1])
	}
	r, err := git.PlainClone(dest, false, cloneOpt)
	if err != nil {
		log.Err(err, "clone")
		return err
	}

	// ... retrieving the branch being pointed by HEAD
	_, err = r.Head()
	if err != nil {
		log.Err(err)
	}
	return err
}

// Pull repo
func (g *GitRepo) PullLocalRepo(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Err(err, "open")
		return err
	}

	workDir, err := repo.Worktree()
	if err != nil {
		log.Err(err, "worktree")
		return err
	}

	// fetch opt
	fetchOpt := &git.FetchOptions{
		Depth: 1,
		Force: true,
	}

	// pull opt
	pullOpt := &git.PullOptions{
		SingleBranch: true,
		Depth:        1,
		Force:        true,
	}

	// setup ssh auth
	if g.sshKeyPath != "" {
		auth, err := ssh.NewPublicKeysFromFile("git", g.sshKeyPath, "")
		if err != nil {
			panic(err)
		}
		pullOpt.Auth = auth
		fetchOpt.Auth = auth
	}

	// Fetch remote
	repo.Fetch(fetchOpt)

	err = workDir.Pull(pullOpt)
	if err == git.NoErrAlreadyUpToDate {
		// do nothing
	} else if err != nil {
		log.Err(err, "pull")
		// err = g.ResetToRemoteHead(path)
		// if err != nil {
		// 	log.Err(err)
		// } else {
		// 	err = workDir.Pull(pullOpt)
		// 	if err != nil {
		// 		log.Err(err)
		// 	}
		// }
	} else {
		log.Info("Pulled repo.")
	}

	return err
}

// Pull, commit and push repo.
func (g *GitRepo) CommitPushLocalRepo(path string, message string, verbose bool) error {
	// Open folder as git repo
	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Err(err, "open")
		return err
	}

	// Get worktree
	workDir, err := repo.Worktree()
	if err != nil {
		log.Err(err, "worktree")
		return err
	}
	// if verbose {
	// 	log.Info("workDir:")
	//  log.Info(workDir)
	// }

	// Pull remote
	// Note: The pull overwrites the current repo and undoes all changes!
	// if pull {
	// 	err = workDir.Pull(&git.PullOptions{
	// 		// RemoteName:   "origin",
	// 		SingleBranch: true,
	// 		Depth:        1,
	// 		Auth:         auth,
	// 		Force:        true,
	// 	})
	// 	if err == git.NoErrAlreadyUpToDate {
	// 		// do nothing
	// 	} else if err != nil {
	// 		log.Err(err, "pull")
	// 	}
	// 	if verbose {
	// 		log.Info("Pulled repo.")
	// 	}
	// }

	// Get status
	status, err := workDir.Status()
	if err != nil {
		log.Err(err, "status")
		return err
	}
	if verbose {
		log.Info("status:")
		log.Info(status.String())
	}

	// Add all files, with changes.
	for path := range status {
		fmt.Println("Added: " + path)
		workDir.Add(path)
	}

	// Commit changes
	commit, err := workDir.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  g.userName,
			Email: g.userEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Err(err, "commit")
		return err
	}

	// Print commit object
	obj, err := repo.CommitObject(commit)
	if err != nil {
		log.Err(err, "commit object")
		return err
	}
	fmt.Println(obj)

	pushOpt := &git.PushOptions{}

	// setup ssh auth
	if g.sshKeyPath != "" {
		// Open key file for auth
		auth, err := ssh.NewPublicKeysFromFile("git", g.sshKeyPath, "")
		if err != nil {
			panic(err)
		}
		pushOpt.Auth = auth
	}

	// Push
	err = repo.Push(pushOpt)
	if err != nil {
		log.Err(err, "push")
		return err
	}

	// Success!
	fmt.Println("Success!")
	return nil
}

// Reset local repo to origin/main HEAD and clean unstaged files.
func (g *GitRepo) ResetToRemoteHead(path string) error {
	// Open key file for auth
	auth, err := ssh.NewPublicKeysFromFile("git", g.sshKeyPath, "")
	if err != nil {
		panic(err)
	}

	// Open folder as git repo
	repo, err := git.PlainOpen(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Get worktree
	workDir, err := repo.Worktree()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Fetch remote
	err = repo.Fetch(&git.FetchOptions{
		Auth: auth,
	})
	if err != nil {
		log.Err(err, "fetch")
	}

	// Get remote head
	remoteRef, err := repo.Reference(plumbing.Main, true)
	if err != nil {
		fmt.Print("Could not get remote head.")
		return err
	}

	// Reset
	workDir.Reset(&git.ResetOptions{
		Mode:   git.HardReset,
		Commit: remoteRef.Hash(),
	})

	// Clean git repo
	err = workDir.Clean(&git.CleanOptions{})
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

// func CheckoutCommit(url string, commit string, destPath string, recursive bool) {
// 	// Open key file for auth
// 	auth, err := ssh.NewPublicKeysFromFile("git", SshKeyPath, "")
// 	if err != nil {
// 		panic(err)
// 	}

// 	repo, err := git.PlainClone(destPath, false, &git.CloneOptions{
// 		URL:      url,
// 		Progress: os.Stdout,
// 		Auth:     auth,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	// get workdir
// 	workDir, err := repo.Worktree()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// checkout commit
// 	err = workDir.Checkout(&git.CheckoutOptions{
// 		Hash:  plumbing.NewHash(commit),
// 		Depth: 1,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }
