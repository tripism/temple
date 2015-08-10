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
	root       string
	lock       sync.RWMutex
	templates  map[string]*Template
	files      []string
	onTemplate func(template *template.Template) (*template.Template, error)
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
		if isPartialsDir(p) {
			return nil // skip partial directories
		}
		if root == p {
			return nil // skip root
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		name := nameFromPath(rel)
		// process the template
		tpl := &Template{
			onTemplate: t.onTemplate,
		}
		if err := tpl.parse(root, p, true); err != nil {
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
	return NewOnTemplate(root, nil)
}

// NewOnTemplate walks directories starting at root and generates a Temple
// object containing all compiled templates, calling the OnTemplate callback for
// each template.
func NewOnTemplate(root string, onTemplate func(template *template.Template) (*template.Template, error)) (*Temple, error) {
	temple := &Temple{
		root:       root,
		templates:  make(map[string]*Template),
		onTemplate: onTemplate,
	}
	err := temple.Reload()
	return temple, err
}

// Template represents a single temple Template.
type Template struct {
	*template.Template
	onTemplate func(template *template.Template) (*template.Template, error)
	foundlist  map[string]struct{}
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

func (t *Template) parse(root, path string, climbup bool) error {
	if t.Template == nil {
		t.Template = template.New(defaultRootTemplateName)
		if t.onTemplate != nil {
			var err error
			if t.Template, err = t.onTemplate(t.Template); err != nil {
				return err
			}
		}
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

			// if these are components - load them regardless
			if isPartialsDir(item.Name()) {
				down := filepath.Join(path, item.Name())

				// process partials
				if err := filepath.Walk(down, func(p string, info os.FileInfo, err error) error {
					if info.IsDir() {
						return nil
					}
					if isTempleFile(p) {
						relpath, err := filepath.Rel(root, p)
						if err != nil {
							return err
						}
						name := nameFromPath(relpath)
						name = templateName(name)
						if err := t.parseFile(name, p); err != nil {
							return err
						}
					}
					return nil
				}); err != nil {
					return err
				}

			}

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
			if err := t.parseFile(name, filepath); err != nil {
				return err
			}
		}
	}
	if climbup {
		relpath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if relpath != "." {
			// we haven't reached root - keep climbing
			up := filepath.Dir(path)
			if err := t.parse(root, up, climbup); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Template) parseFile(name, path string) error {
	t.Files = append(t.Files, path)
	b, err := ioutil.ReadFile(path)
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

func nameFromPath(p string) string {
	if strings.Contains(p, "_") {
		segs := strings.Split(p, string(filepath.Separator))
		for i, seg := range segs {
			if strings.HasPrefix(seg, "_") {
				segs[i] = strings.TrimPrefix(seg, "_")
			}
		}
		return strings.Join(segs, ".")
	}
	return strings.Replace(p, string(filepath.Separator), ".", -1)
}

func isPartialsDir(p string) bool {
	segs := strings.Split(p, string(filepath.Separator))
	for _, seg := range segs {
		if strings.HasPrefix(seg, "_") {
			return true
		}
	}
	return false
}
