package temple_test

import (
	"bytes"
	"testing"

	"github.com/cheekybits/is"
	"github.com/tripism/temple"
)

func TestTemple(t *testing.T) {
	is := is.New(t)

	tpl, err := temple.Process("test")
	is.NoErr(err)
	is.OK(tpl)
	is.OK(tpl["site.welcome.about.nested"])

	data := map[string]interface{}{"Name": "Mat"}
	var buf bytes.Buffer
	is.NoErr(tpl["site.welcome.about.nested"].Execute(&buf, data))
	is.Equal(buf.String(), `<base>Hello Mat.</base>`)

	buf.Reset()
	is.NoErr(tpl["site.welcome"].Execute(&buf, data))
	is.Equal(buf.String(), `<base>Welcome</base>`)

}
