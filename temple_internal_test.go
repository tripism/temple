package temple

import (
	"bytes"
	"testing"

	"github.com/cheekybits/is"
)

func TestParseTemplate(t *testing.T) {
	is := is.New(t)

	c := &Template{}
	err := c.parse("test/site", "test/site/welcome/about/nested")
	is.NoErr(err)
	is.OK(c.Template)
	is.Equal(c.RootTemplateName, "base")

	data := map[string]interface{}{"Name": "Mat"}
	var buf bytes.Buffer
	is.NoErr(c.Execute(&buf, data))
	is.Equal(buf.String(), `<base>Hello Mat.</base>`)

}
