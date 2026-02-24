package resolve

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Workfiles struct {
	TempDirectory   string
	Domains         string
	MassdnsPublic   string
	MassdnsTrusted  string
	Temporary       string
	PublicResolvers string
	TrustedResolvers string
	WildcardRoots   string
}

func (w *Workfiles) Close() {
	if w.TempDirectory != "" {
		os.RemoveAll(w.TempDirectory)
	}
}

type DefaultWorkfileCreator struct {
	osMkdirTemp func(dir string, pattern string) (string, error)
	osCreate    func(name string) (*os.File, error)
}

func NewDefaultWorkfileCreator() *DefaultWorkfileCreator {
	return &DefaultWorkfileCreator{
		osMkdirTemp: ioutil.TempDir,
		osCreate:    os.Create,
	}
}

func (w *DefaultWorkfileCreator) Create() (*Workfiles, error) {
	files := &Workfiles{}
	dir, err := w.osMkdirTemp("", "resodns.")
	if err != nil {
		return nil, fmt.Errorf("unable to create temporary work directory: %w", err)
	}
	files.TempDirectory = dir
	sep := string(filepath.Separator)
	if files.Domains, err = w.createFile(dir + sep + "domains.txt"); err != nil {
		return nil, err
	}
	if files.MassdnsPublic, err = w.createFile(dir + sep + "massdns_public.txt"); err != nil {
		return nil, err
	}
	if files.MassdnsTrusted, err = w.createFile(dir + sep + "massdns_trusted.txt"); err != nil {
		return nil, err
	}
	if files.Temporary, err = w.createFile(dir + sep + "temporary.txt"); err != nil {
		return nil, err
	}
	if files.PublicResolvers, err = w.createFile(dir + sep + "resolvers.txt"); err != nil {
		return nil, err
	}
	if files.TrustedResolvers, err = w.createFile(dir + sep + "trusted.txt"); err != nil {
		return nil, err
	}
	if files.WildcardRoots, err = w.createFile(dir + sep + "wildcards.txt"); err != nil {
		return nil, err
	}
	return files, nil
}

func (w *DefaultWorkfileCreator) createFile(path string) (string, error) {
	file, err := w.osCreate(path)
	if err != nil {
		return "", fmt.Errorf("unable to create temporary file %s: %w", path, err)
	}
	defer file.Close()
	return file.Name(), nil
}
