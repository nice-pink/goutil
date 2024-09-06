package repo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/nice-pink/goutil/pkg/log"
)

type RepoHandle struct {
	userName   string
	userEmail  string
	sshKeyPath string
	repo       *git.Repository
}

func NewRepoHandle(sshKeyPath string, userName string, userEmail string) *RepoHandle {
	return &RepoHandle{
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

// open

func (g *RepoHandle) Open(path string) error {
	var err error
	g.repo, err = git.PlainOpen(path)
	if err != nil {
		log.Err(err, "open")
		return err
	}
	return nil
}

// Clone
func (g *RepoHandle) Clone(url string, dest string, branch string, shallow bool, repoSubfolder bool) error {
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

	var err error
	g.repo, err = git.PlainClone(dest, false, cloneOpt)
	if err != nil {
		log.Err(err, "clone")
		return err
	}

	// ... retrieving the branch being pointed by HEAD
	_, err = g.repo.Head()
	if err != nil {
		log.Err(err)
	}
	return err
}

// Pull repo
func (g *RepoHandle) PullLocalRepo(path string) error {
	var err error
	g.repo, err = git.PlainOpen(path)
	if err != nil {
		log.Err(err, "open")
		return err
	}

	workDir, err := g.repo.Worktree()
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
	g.repo.Fetch(fetchOpt)

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
func (g *RepoHandle) CommitPushLocalRepo(path string, message string, verbose bool) error {
	// Open folder as git repo
	var err error
	g.repo, err = git.PlainOpen(path)
	if err != nil {
		log.Err(err, "open")
		return err
	}

	// Get worktree
	workDir, err := g.repo.Worktree()
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
	obj, err := g.repo.CommitObject(commit)
	if err != nil {
		log.Err(err, "commit object")
		return err
	}
	fmt.Println(obj)

	err = g.Push()
	if err != nil {
		return err
	}

	// Success!
	fmt.Println("Success!")
	return nil
}

func (g *RepoHandle) Push() error {
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
	err := g.repo.Push(pushOpt)
	if err != nil {
		log.Err(err, "push")
	}
	return err
}

// Reset local repo to origin/main HEAD and clean unstaged files.
func (g *RepoHandle) ResetToRemoteHead(path string) error {
	// Open key file for auth
	auth, err := ssh.NewPublicKeysFromFile("git", g.sshKeyPath, "")
	if err != nil {
		panic(err)
	}

	// Open folder as git repo
	g.repo, err = git.PlainOpen(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Get worktree
	workDir, err := g.repo.Worktree()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Fetch remote
	err = g.repo.Fetch(&git.FetchOptions{
		Auth: auth,
	})
	if err != nil {
		log.Err(err, "fetch")
	}

	// Get remote head
	remoteRef, err := g.repo.Reference(plumbing.Main, true)
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

// tag

func (g *RepoHandle) TagRepo(tag, msg string, push bool) error {
	if g.repo == nil {
		fmt.Println("No repo")
		return errors.New("no repo")
	}

	created, err := g.setTag(tag, msg)
	if err != nil {
		fmt.Println("create tag error")
		fmt.Println(err)
		return err
	}

	if !push || !created {
		return nil
	}

	err = g.PushTags()
	if err != nil {
		fmt.Println("push tag error")
		fmt.Println(err)
		return err
	}

	return nil
}

func (g *RepoHandle) TagExists(tag string) bool {
	tagFoundErr := "tag was found"
	// Info("git show-ref --tag")
	tags, err := g.repo.TagObjects()
	if err != nil {
		fmt.Println("get tags error")
		fmt.Println(err)
		return false
	}
	res := false
	err = tags.ForEach(func(t *object.Tag) error {
		if t.Name == tag {
			res = true
			return fmt.Errorf(tagFoundErr)
		}
		return nil
	})
	if err != nil && err.Error() != tagFoundErr {
		fmt.Println("iterate tags error")
		fmt.Println(err)
		return false
	}
	return res
}

func (g *RepoHandle) setTag(tag, msg string) (bool, error) {
	if g.TagExists(tag) {
		fmt.Println("tag", tag, "already exists")
		return false, nil
	}
	fmt.Println("Set tag", tag)
	h, err := g.repo.Head()
	if err != nil {
		fmt.Println("get HEAD error")
		fmt.Println(err)
		return false, err
	}
	// Info("git tag -a %s %s -m \"%s\"", tag, h.Hash(), tag)
	_, err = g.repo.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Message: msg,
	})

	if err != nil {
		fmt.Println("create tag error")
		fmt.Println(err)
		return false, err
	}

	return true, nil
}

func (g *RepoHandle) PushTags() error {

	po := &git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
	}

	// setup ssh auth
	if g.sshKeyPath != "" {
		auth, err := ssh.NewPublicKeysFromFile("git", g.sshKeyPath, "")
		if err != nil {
			log.Err(err, "key")
			return err
		}
		po.Auth = auth
	}

	// Info("git push --tags")
	err := g.repo.Push(po)

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Println("origin remote was up to date, no push done")
			return nil
		}
		fmt.Println("push to remote origin error")
		fmt.Println(err)
		return err
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
