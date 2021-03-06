package main

import (
	"fmt"
	"os"
  "flag"

	"gopkg.in/src-d/go-git.v4"
  . "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// let's start with some g

func main() {
  remotePtr := flag.String("remote", "", "a remote git repository to inspect")
  pathPtr := flag.String("path", ".", "a local git repository to inspect")

  flag.Parse()
  var r *git.Repository
  var err error
  if *remotePtr != "" {
    Info("git clone " + *remotePtr)
    r, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: *remotePtr,
	})
  } else {
    Info("using git against" + *pathPtr)
    r, err = git.PlainOpen(*pathPtr)
  }
	CheckIfError(err)

	Info("git log")

	ref, err := r.Head()
	CheckIfError(err)

	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	Info("git ls-tree -r HEAD")

	// ... retrive the tree from the commit
	tree, err := commit.Tree()
	CheckIfError(err)

	// ... get the files iterator and print the file
	tree.Files().ForEach(func(f *object.File) error {
		fmt.Printf("%s\n", f.Name)
		cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), FileName: &f.Name})
		CheckIfError(err)

		err = cIter.ForEach(func(c *object.Commit) error {
			fmt.Printf("\t%s\t%s\n", c.Hash, c.Author.When)
			return nil
		})
		return nil
	})
}

func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
