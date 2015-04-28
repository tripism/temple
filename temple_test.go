package temple_test

import (
	"bytes"
	"testing"

	"github.com/cheekybits/is"
	"github.com/tripism/temple"
)

func TestTemple(t *testing.T) {
	is := is.New(t)

	tpl, err := temple.New("test")
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
