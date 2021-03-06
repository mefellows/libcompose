package project

import (
	"fmt"
	"strings"
	"testing"

	yaml "github.com/cloudfoundry-incubator/candiedyaml"

	"github.com/stretchr/testify/assert"
)

type StructStringorslice struct {
	Foo Stringorslice
}

type TestConfig struct {
	SystemContainers map[string]*ServiceConfig
}

func newTestConfig() TestConfig {
	return TestConfig{
		SystemContainers: map[string]*ServiceConfig{
			"udev": {
				Image:      "udev",
				Restart:    "always",
				Net:        "host",
				Privileged: true,
				DNS:        Stringorslice{[]string{"8.8.8.8", "8.8.4.4"}},
				Environment: MaporEqualSlice{[]string{
					"DAEMON=true",
				}},
				Labels: SliceorMap{map[string]string{
					"io.rancher.os.detach": "true",
					"io.rancher.os.scope":  "system",
				}},
				VolumesFrom: []string{
					"system-volumes",
				},
				Ulimits: Ulimits{
					Elements: []Ulimit{
						{
							Name: "nproc",
							ulimitValues: ulimitValues{
								Soft: 65557,
								Hard: 65557,
							},
						},
					},
				},
			},
			"system-volumes": {
				Image:      "state",
				Net:        "none",
				ReadOnly:   true,
				Privileged: true,
				Labels: SliceorMap{map[string]string{
					"io.rancher.os.createonly": "true",
					"io.rancher.os.scope":      "system",
				}},
				Volumes: []string{
					"/dev:/host/dev",
					"/var/lib/rancher/conf:/var/lib/rancher/conf",
					"/etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt.rancher",
					"/lib/modules:/lib/modules",
					"/lib/firmware:/lib/firmware",
					"/var/run:/var/run",
					"/var/log:/var/log",
				},
				LogDriver: "json-file",
			},
		},
	}
}

func TestMarshalConfig(t *testing.T) {
	config := newTestConfig()
	bytes, err := yaml.Marshal(config)
	assert.Nil(t, err)

	config2 := TestConfig{}

	err = yaml.Unmarshal(bytes, &config2)
	assert.Nil(t, err)

	assert.Equal(t, config, config2)
}

func TestMarshalServiceConfig(t *testing.T) {
	configPtr := newTestConfig().SystemContainers["udev"]
	bytes, err := yaml.Marshal(configPtr)
	assert.Nil(t, err)

	configPtr2 := &ServiceConfig{}

	err = yaml.Unmarshal(bytes, configPtr2)
	assert.Nil(t, err)

	assert.Equal(t, configPtr, configPtr2)
}

func TestStringorsliceYaml(t *testing.T) {
	str := `{foo: [bar, baz]}`

	s := StructStringorslice{}
	yaml.Unmarshal([]byte(str), &s)

	assert.Equal(t, []string{"bar", "baz"}, s.Foo.parts)

	d, err := yaml.Marshal(&s)
	assert.Nil(t, err)

	s2 := StructStringorslice{}
	yaml.Unmarshal(d, &s2)

	assert.Equal(t, []string{"bar", "baz"}, s2.Foo.parts)
}

type StructSliceorMap struct {
	Foos SliceorMap `yaml:"foos,omitempty"`
	Bars []string   `yaml:"bars"`
}

type StructCommand struct {
	Entrypoint Command `yaml:"entrypoint,flow,omitempty"`
	Command    Command `yaml:"command,flow,omitempty"`
}

func TestSliceOrMapYaml(t *testing.T) {
	str := `{foos: [bar=baz, far=faz]}`

	s := StructSliceorMap{}
	yaml.Unmarshal([]byte(str), &s)

	assert.Equal(t, map[string]string{"bar": "baz", "far": "faz"}, s.Foos.parts)

	d, err := yaml.Marshal(&s)
	assert.Nil(t, err)

	s2 := StructSliceorMap{}
	yaml.Unmarshal(d, &s2)

	assert.Equal(t, map[string]string{"bar": "baz", "far": "faz"}, s2.Foos.parts)
}

var sampleStructSliceorMap = `
foos:
  io.rancher.os.bar: baz
  io.rancher.os.far: true
bars: []
`

func TestUnmarshalSliceOrMap(t *testing.T) {
	s := StructSliceorMap{}
	err := yaml.Unmarshal([]byte(sampleStructSliceorMap), &s)
	assert.Equal(t, fmt.Errorf("Cannot unmarshal 'true' of type bool into a string value"), err)
}

func TestStr2SliceOrMapPtrMap(t *testing.T) {
	s := map[string]*StructSliceorMap{"udav": {
		Foos: SliceorMap{map[string]string{"io.rancher.os.bar": "baz", "io.rancher.os.far": "true"}},
		Bars: []string{},
	}}
	d, err := yaml.Marshal(&s)
	assert.Nil(t, err)

	s2 := map[string]*StructSliceorMap{}
	yaml.Unmarshal(d, &s2)

	assert.Equal(t, s, s2)
}

type StructMaporslice struct {
	Foo MaporEqualSlice
}

func contains(list []string, item string) bool {
	for _, test := range list {
		if test == item {
			return true
		}
	}
	return false
}

func TestMaporsliceYaml(t *testing.T) {
	str := `{foo: {bar: baz, far: faz}}`

	s := StructMaporslice{}
	yaml.Unmarshal([]byte(str), &s)

	assert.Equal(t, 2, len(s.Foo.parts))
	assert.True(t, contains(s.Foo.parts, "bar=baz"))
	assert.True(t, contains(s.Foo.parts, "far=faz"))

	d, err := yaml.Marshal(&s)
	assert.Nil(t, err)

	s2 := StructMaporslice{}
	yaml.Unmarshal(d, &s2)

	assert.Equal(t, 2, len(s2.Foo.parts))
	assert.True(t, contains(s2.Foo.parts, "bar=baz"))
	assert.True(t, contains(s2.Foo.parts, "far=faz"))
}

var sampleStructCommand = `command: bash`

func TestUnmarshalCommand(t *testing.T) {
	s := &StructCommand{}
	err := yaml.Unmarshal([]byte(sampleStructCommand), s)

	assert.Nil(t, err)
	assert.Equal(t, []string{"bash"}, s.Command.Slice())
	assert.Nil(t, s.Entrypoint.Slice())

	bytes, err := yaml.Marshal(s)
	assert.Nil(t, err)

	s2 := &StructCommand{}
	err = yaml.Unmarshal(bytes, s2)

	assert.Nil(t, err)
	assert.Equal(t, []string{"bash"}, s2.Command.Slice())
	assert.Nil(t, s2.Entrypoint.Slice())
}

var sampleEmptyCommand = `{}`

func TestUnmarshalEmptyCommand(t *testing.T) {
	s := &StructCommand{}
	err := yaml.Unmarshal([]byte(sampleEmptyCommand), s)

	assert.Nil(t, err)
	assert.Nil(t, s.Command.Slice())

	bytes, err := yaml.Marshal(s)
	assert.Nil(t, err)
	assert.Equal(t, "entrypoint: []\ncommand: []", strings.TrimSpace(string(bytes)))

	s2 := &StructCommand{}
	err = yaml.Unmarshal(bytes, s2)

	assert.Nil(t, err)
	assert.Nil(t, s2.Command.Slice())
}

func TestMarshalUlimit(t *testing.T) {
	ulimits := []struct {
		ulimits  *Ulimits
		expected string
	}{
		{
			ulimits: &Ulimits{
				Elements: []Ulimit{
					{
						ulimitValues: ulimitValues{
							Soft: 65535,
							Hard: 65535,
						},
						Name: "nproc",
					},
				},
			},
			expected: `nproc: 65535
`,
		},
		{
			ulimits: &Ulimits{
				Elements: []Ulimit{
					{
						Name: "nofile",
						ulimitValues: ulimitValues{
							Soft: 20000,
							Hard: 40000,
						},
					},
				},
			},
			expected: `nofile:
  soft: 20000
  hard: 40000
`,
		},
	}

	for _, ulimit := range ulimits {

		bytes, err := yaml.Marshal(ulimit.ulimits)

		assert.Nil(t, err)
		assert.Equal(t, ulimit.expected, string(bytes), "should be equal")
	}
}

func TestUnmarshalUlimits(t *testing.T) {
	ulimits := []struct {
		yaml     string
		expected *Ulimits
	}{
		{
			yaml: "nproc: 65535",
			expected: &Ulimits{
				Elements: []Ulimit{
					{
						Name: "nproc",
						ulimitValues: ulimitValues{
							Soft: 65535,
							Hard: 65535,
						},
					},
				},
			},
		},
		{
			yaml: `nofile:
  soft: 20000
  hard: 40000`,
			expected: &Ulimits{
				Elements: []Ulimit{
					{
						Name: "nofile",
						ulimitValues: ulimitValues{
							Soft: 20000,
							Hard: 40000,
						},
					},
				},
			},
		},
		{
			yaml: `nproc: 65535
nofile:
  soft: 20000
  hard: 40000`,
			expected: &Ulimits{
				Elements: []Ulimit{
					{
						Name: "nofile",
						ulimitValues: ulimitValues{
							Soft: 20000,
							Hard: 40000,
						},
					},
					{
						Name: "nproc",
						ulimitValues: ulimitValues{
							Soft: 65535,
							Hard: 65535,
						},
					},
				},
			},
		},
	}

	for _, ulimit := range ulimits {
		actual := &Ulimits{}
		err := yaml.Unmarshal([]byte(ulimit.yaml), actual)

		assert.Nil(t, err)
		assert.Equal(t, ulimit.expected, actual, "should be equal")
	}
}
