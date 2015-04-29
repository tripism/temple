package temple

import (
	"bytes"
	"log"
	"testing"

	"github.com/cheekybits/is"
	//"github.com/tripism/temple"
)

func TestTemple(t *testing.T) {
	is := is.New(t)

	tpl, err := New("test")
	is.NoErr(err)
	is.OK(tpl)
	log.Println(tpl.templates)
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
