package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/luno/jettison/errors"
)

func CloneRepo(url, dest string) error {
	cloneOptions := &git.CloneOptions{
		URL: url,
	}

	_, err := git.PlainClone(dest, false, cloneOptions)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
