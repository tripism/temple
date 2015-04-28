package temple

import (
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const defaultRootTemplateName = "temple"

// Temple represents a map of Templates keyed by their dot notation
// names.
type Temple struct {
	root      string
	lock      sync.RWMutex
	templates map[string]*Template
	files     []string
}

// Get gets a Template by name.
func (t *Temple) Get(name string) *Template {
	tpl, _ := t.GetOK(name)
	return tpl
}

// GetOK gets a Template by name and returns whether one
// was found or not.
func (t *Temple) GetOK(name string) (*Template, bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	tpl, ok := t.templates[name]
	return tpl, ok
}

// Root gets the root directory for the templates.
func (t *Temple) Root() string {
	return t.root
}

// Files gets the template files that were loaded.
func (t *Temple) Files() []string {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.files
}

// Reload reloads all templates.
func (t *Temple) Reload() error {
	root := t.Root()
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
		tpl := &Template{}
		if err := tpl.parse(root, p); err != nil {
			return err
		}
		t.lock.Lock()
		t.templates[name] = tpl
		t.files = append(t.files, tpl.Files...)
		t.lock.Unlock()
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// New walks directories starting at root and generates a Temple
// object containing all compiled templates.
func New(root string) (*Temple, error) {
	temple := &Temple{
		root:      root,
		templates: make(map[string]*Template),
	}
	err := temple.Reload()
	return temple, err
}

// Template represents a single temple Template.
type Template struct {
	*template.Template
	foundlist map[string]struct{}
	// RootTemplateName is the name of the template that will be
	// rendered when Execute is called.
	RootTemplateName string
	// Files represents the files that make up this template.
	Files []string
}

// Execute applies a parsed template to the specified data object, writing the output to wr.
// It calls ExecuteTemplate on the underlying Template with the appropriate
// RootTemplateName.
func (t *Template) Execute(wr io.Writer, data interface{}) error {
	return t.Template.ExecuteTemplate(wr, t.RootTemplateName, data)
}

func (t *Template) parse(root, path string) error {
	if t.Template == nil {
		t.Template = template.New(defaultRootTemplateName)
	}
	if t.foundlist == nil {
		t.foundlist = make(map[string]struct{})
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
			if _, present := t.foundlist[name]; present {
				continue
			}
			// base always becomes the entry template
			if len(t.RootTemplateName) == 0 || name == "base" {
				t.RootTemplateName = name
			}
			filepath := filepath.Join(path, item.Name())
			t.Files = append(t.Files, filepath)
			b, err := ioutil.ReadFile(filepath)
			if err != nil {
				return err
			}
			s := string(b)
			s = "{{define \"" + name + "\"}}" + s + "{{end}}"
			_, err = t.Parse(s)
			if err != nil {
				return err
			}
			t.foundlist[name] = struct{}{}
		}
	}
	relpath, err := filepath.Rel(root, path)
	if err != nil {
		return err
	}
	if relpath != "." {
		// we haven't reached root - keep climbing
		up := filepath.Dir(path)
		if err := t.parse(root, up); err != nil {
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
