package releases

import (
	"platform.prodigy9.co/gitcmd"
	"platform.prodigy9.co/project"
	"platform.prodigy9.co/releases/dateref"
	"strings"
)

type Date struct{}

var _ Strategy = Date{}

func (d Date) List(cfg *project.Project) ([]*Release, error) {
	lines, err := gitcmd.ListTags(cfg.ConfigDir)
	if err != nil {
		return nil, err
	}

	var result []*Release
	for _, line := range strings.Split(lines, "\n") {
		if dateref.IsValid(line) {
			result = append(result, &Release{Name: line})
		}
	}
	return result, nil
}

func (d Date) Recover(cfg *project.Project, opts *Options) (*Release, error) {
	// get annotated tag and name
	if opts.Name == "" {
		tagname, err := gitcmd.Describe(cfg.ConfigDir)
		if err != nil {
			return nil, err
		} else if !dateref.IsValid(tagname) {
			return nil, ErrBadSemver
		}

		opts.Name = tagname
	}

	tagmsg, err := gitcmd.TagMessage(cfg.ConfigDir, opts.Name)
	if err != nil {
		return nil, err
	}

	return &Release{Name: opts.Name, Message: tagmsg}, nil
}

func (d Date) Generate(cfg *project.Project, opts *Options) (*Release, error) {
	if opts.Name == "" {
		opts.Name = dateref.Now()
	}

	tagmsg, err := gitcmd.TagMessage(cfg.ConfigDir, opts.Name)
	if err != nil {
		return nil, err
	}
	return &Release{Name: opts.Name, Message: tagmsg}, nil
}

func (d Date) Create(cfg *project.Project, rel *Release) error {
	if _, err := gitcmd.Tag(cfg.ConfigDir, rel.Name, rel.Message); err != nil {
		return err
	} else if branch, err := gitcmd.CurrentBranch(cfg.ConfigDir); err != nil {
		return err
	} else if remote, err := gitcmd.TrackingRemote(cfg.ConfigDir, branch); err != nil {
		return err
	} else if _, err := gitcmd.PushTag(cfg.ConfigDir, remote, rel.Name); err != nil {
		return err
	} else {
		return nil
	}
}
