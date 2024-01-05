package git

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

var UserName string = ""
var UserEmail string = ""
var SshKeyPath string = ""

func Setup(user string, email string, keypath string) {
	UserName = user
	UserEmail = email
	SshKeyPath = keypath
}

// Pull repo
func PullLocalRepo(path string) error {
	auth, err := ssh.NewPublicKeysFromFile("git", SshKeyPath, "")
	if err != nil {
		panic(err)
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	workDir, err := repo.Worktree()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = workDir.Pull(&git.PullOptions{
		RemoteName:   "origin",
		SingleBranch: true,
		Depth:        1,
		Auth:         auth,
		Force:        true,
	})
	if err == git.NoErrAlreadyUpToDate {
		// do nothing
	} else if err != nil {
		fmt.Println(err)
	}
	return err
}

// Pull, commit and push repo.
func CommitPushLocalRepo(path string, message string, pull bool) error {
	// Open key file for auth
	auth, err := ssh.NewPublicKeysFromFile("git", SshKeyPath, "")
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

	// Pull remote
	if pull {
		err = workDir.Pull(&git.PullOptions{
			// RemoteName:   "origin",
			SingleBranch: true,
			Depth:        1,
			Auth:         auth,
			Force:        true,
		})
		if err == git.NoErrAlreadyUpToDate {
			// do nothing
		} else if err != nil {
			fmt.Println("could not pull.")
			fmt.Println(err)
		}
	}

	// Get status
	status, err := workDir.Status()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Add all files, with changes.
	for path, _ := range status {
		fmt.Println("Added: " + path)
		workDir.Add(path)
	}

	// Commit changes
	commit, err := workDir.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  UserName,
			Email: UserEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Print commit object
	obj, err := repo.CommitObject(commit)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(obj)

	// Push
	err = repo.Push(&git.PushOptions{
		// RemoteName: "origin",
		Auth: auth,
	})
	if err != nil {
		fmt.Println(err)
		fmt.Println("could not push.")
		return err
	}

	// Success!
	fmt.Println("Success!")
	return nil
}

// Reset local repo to origin/main HEAD and clean unstaged files.
func ResetToRemoteHead(path string) error {
	// Open key file for auth
	auth, err := ssh.NewPublicKeysFromFile("git", SshKeyPath, "")
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
		RemoteName: "origin",
		Auth:       auth,
	})

	// Get remote head
	remoteRef, err := repo.Reference(plumbing.Main, true)
	if err != nil {
		fmt.Print("Could not get remote head.")
		return err
	}

	// Reset
	err = workDir.Reset(&git.ResetOptions{
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
