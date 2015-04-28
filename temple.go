package temple

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const defaultRootTemplateName = "_entry"

// Temple represents a map of Templates with their dot notation
// names.
type Temple map[string]*Template

// Process processes the root temple folder.
func Process(root string) (Temple, error) {
	temple := make(Temple)
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil // skip-non directories
		}
		if root == p {
			return nil // skip root
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		name := strings.Replace(rel, "/", ".", -1)
		// process the template
		tpl, err := parse(root, p)
		if err != nil {
			return err
		}
		temple[name] = tpl
		return nil
	})
	if err != nil {
		return nil, err
	}
	return temple, nil
}

// Template represents a single temple Template.
type Template struct {
	*template.Template
	foundlist map[string]struct{}
	// RootTemplateName is the name of the template that will be
	// rendered when Execute is called.
	RootTemplateName string
}

func parse(root, path string) (*Template, error) {
	tpl := &Template{}
	if err := tpl.parse(root, path); err != nil {
		return nil, err
	}
	return tpl, nil
}

// Execute applies a parsed template to the specified data object, writing the output to wr.
// It calls ExecuteTemplate on the underlying Template with the appropriate
// RootTemplateName.
func (c *Template) Execute(wr io.Writer, data interface{}) error {
	return c.Template.ExecuteTemplate(wr, c.RootTemplateName, data)
}

func (c *Template) parse(root, path string) error {
	if c.Template == nil {
		c.Template = template.New(defaultRootTemplateName)
	}
	if c.foundlist == nil {
		c.foundlist = make(map[string]struct{})
	}
	items, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		if isTempleFile(item.Name()) {
			name := templateName(item.Name())
			// skip it if we already have it
			if _, present := c.foundlist[name]; present {
				continue
			}
			// base always becomes the entry template
			if len(c.RootTemplateName) == 0 || name == "base" {
				c.RootTemplateName = name
			}
			filepath := filepath.Join(path, item.Name())
			b, err := ioutil.ReadFile(filepath)
			if err != nil {
				return err
			}
			s := string(b)
			s = "{{define \"" + name + "\"}}" + s + "{{end}}"
			_, err = c.Parse(s)
			if err != nil {
				return err
			}
			c.foundlist[name] = struct{}{}
		}
	}
	relpath, err := filepath.Rel(root, path)
	if err != nil {
		return err
	}
	if relpath != "." {
		// we haven't reached root - keep climbing
		up := filepath.Dir(path)
		if err := c.parse(root, up); err != nil {
			return err
		}
	}
	return nil
}

func isTempleFile(p string) bool {
	return strings.Contains(p, ".temple")
}

func templateName(p string) string {
	p = filepath.Base(p)
	i := strings.Index(p, ".temple")
	p = strings.ToLower(p[0:i])
	return p
}
