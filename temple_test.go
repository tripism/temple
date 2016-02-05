package temple

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cheekybits/is"
)

func TestTemple(t *testing.T) {
	is := is.New(t)

	tpl, err := New("test/v1")
	is.NoErr(err)
	is.OK(tpl)
	_, ok := tpl.GetOK("site.welcome.about.nested")
	is.True(ok)

	data := map[string]interface{}{"Name": "Mat"}
	var buf bytes.Buffer
	is.NoErr(tpl.Get("site.welcome.about.nested").Execute(&buf, data))
	is.Equal(buf.String(), `<base>Hello Mat.</base>`)

	buf.Reset()
	is.NoErr(tpl.Get("site.welcome").Execute(&buf, data))
	is.Equal(buf.String(), `<base>Welcome</base>`)

}

func TestTempleFuncs(t *testing.T) {
	is := is.New(t)

	tpl, err := NewFuncs("test/funcs", map[string]interface{}{"title": strings.Title})
	is.NoErr(err)
	is.OK(tpl)
	theTemplate, ok := tpl.GetOK("site.welcome.about.funcs")
	is.True(ok)

	data := map[string]interface{}{"Name": "mat"}
	var buf bytes.Buffer
	is.NoErr(theTemplate.Execute(&buf, data))
	is.Equal(buf.String(), `<base>Hello Mat.</base>`)

	buf.Reset()
	is.NoErr(tpl.Get("site.welcome").Execute(&buf, data))
	is.Equal(buf.String(), `<base>Welcome</base>`)

}
