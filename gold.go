package gold

import (
	"bytes"
	"os"
	"testing"

	"github.com/sebdah/goldie/v2"
	"gopkg.in/yaml.v3"
)

type Gold struct {
	*goldie.Goldie
}

func New(t *testing.T, opts ...goldie.Option) *Gold {
	defaults := []goldie.Option{
		goldie.WithDiffEngine(goldie.ColoredDiff),
		goldie.WithTestNameForDir(true),
	}

	return &Gold{
		Goldie: goldie.New(t, append(defaults, opts...)...),
	}
}

// AssertYaml compares the actual yaml data received with expected data in the
// golden files. If the update flag is set, it will also update the golden
// file.
//
// `name` refers to the name of the test and it should typically be unique
// within the package. Also it should be a valid file name (so keeping to
// `a-z0-9\-\_` is a good idea).
func (g *Gold) AssertYaml(t *testing.T, name string, actual interface{}, indent int) {
	sb := &bytes.Buffer{}
	e := yaml.NewEncoder(sb)
	e.SetIndent(indent)
	err := e.Encode(actual)
	if err != nil {
		t.Fatalf("unable to marshal actual '%v' for test '%s'; err=%v\n", actual, t.Name(), err)
	}

	actualYaml := sb.String()
	_ = actualYaml

	ys := sb.Bytes()
	g.Assert(t, name, normalizeLF(ys))
}

// normalizeLF normalizes line feed character set across os (es)
// \r\n (windows) & \r (mac) into \n (unix)
func normalizeLF(d []byte) []byte {
	// if empty / nil return as is
	if len(d) == 0 {
		return d
	}
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)
	return d
}

func (g *Gold) Get(t *testing.T, name string) (string, []byte) {
	path := g.GoldenFileName(t, name)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unable to read golden master '%s' for test '%s' at path '%s'; err=%v\n", name, t.Name(), path, err)
	}
	return string(b), b
}
