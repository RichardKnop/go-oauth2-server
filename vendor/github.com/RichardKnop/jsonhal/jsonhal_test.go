package jsonhal_test

import (
	"bytes"
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/RichardKnop/jsonhal"

	"github.com/stretchr/testify/assert"
)

// HelloWorld is a simple test struct
type HelloWorld struct {
	jsonhal.Hal
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Foobar is a simple test struct
type Foobar struct {
	jsonhal.Hal
	ID   uint      `json:"id"`
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

// Qux is a simple test struct
type Qux struct {
	jsonhal.Hal
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

var expectedJSON = []byte(`{
	"id": 1,
	"name": "Hello World"
}`)

var expectedJSON2 = []byte(`{
	"_links": {
		"self": {
			"href": "/v1/hello/world/1"
		}
	},
	"id": 1,
	"name": "Hello World"
}`)

var expectedJSON3 = []byte(`{
	"_links": {
		"next": {
			"href": "/v1/hello/world?offset=4\u0026limit=2"
		},
		"previous": {
			"href": "/v1/hello/world?offset=0\u0026limit=2"
		},
		"self": {
			"href": "/v1/hello/world?offset=2\u0026limit=2"
		}
	},
	"_embedded": {
		"foobar": {
			"_links": {
				"self": {
					"href": "/v1/foo/bar/1"
				}
			},
			"id": 1,
			"name": "Foo bar 1",
            "date":"2017-09-12T08:45:20Z"
		}
	},
	"id": 1,
	"name": "Hello World"
}`)

var expectedJSON4 = []byte(`{
	"_links": {
		"self": {
			"href": "/v1/hello/world/1"
		}
	},
	"_embedded": {
		"foobars": [
			{
				"_links": {
					"self": {
						"href": "/v1/foo/bar/1"
					}
				},
				"id": 1,
				"name": "Foo bar 1",
                "date":"2017-09-12T08:45:20Z"
			},
			{
				"_links": {
					"self": {
						"href": "/v1/foo/bar/2"
					}
				},
				"id": 2,
				"name": "Foo bar 2",
            	"date":"2017-09-12T08:45:20Z"
			}
		]
	},
	"id": 1,
	"name": "Hello World"
}`)

var expectedJSON5 = []byte(`{
	"_links": {
		"self": {
			"href": "/v1/hello/world/1"
		}
	},
	"_embedded": {
		"foobars": [
			{
				"_links": {
					"self": {
						"href": "/v1/foo/bar/1"
					}
				},
				"id": 1,
				"name": "Foo bar 1",
            	"date":"2017-09-12T08:45:20Z"
			},
			{
				"_links": {
					"self": {
						"href": "/v1/foo/bar/2"
					}
				},
				"id": 2,
				"name": "Foo bar 2",
            	"date":"2017-09-12T08:45:20Z"
			}
		],
		"quxes": [
			{
				"_links": {
					"self": {
						"href": "/v1/qux/1"
					}
				},
				"id": 1,
				"name": "Qux 1"
			},
			{
				"_links": {
					"self": {
						"href": "/v1/qux/2"
					}
				},
				"id": 2,
				"name": "Qux 2"
			}
		]
	},
	"id": 1,
	"name": "Hello World"
}`)

func TestHal(t *testing.T) {
	t.Parallel()

	var (
		helloWorld *HelloWorld
		expected   *bytes.Buffer
		actual     []byte
		err        error
		foobar     *Foobar
		foobars    []*Foobar
		quxes      []*Qux
	)

	// Let's test the simplest scenario without links
	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}

	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))

	// Let's add a self link
	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)

	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON2)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))

	// Let's add more links and a single embedded resource
	date, _ := time.Parse(time.RFC3339, "2017-09-12T08:45:20Z")
	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}
	helloWorld.SetLink(
		"self", // name
		"/v1/hello/world?offset=2&limit=2", // href
		"", // title
	)
	helloWorld.SetLink(
		"next", // name
		"/v1/hello/world?offset=4&limit=2", // href
		"", // title
	)
	helloWorld.SetLink(
		"previous",                         // name
		"/v1/hello/world?offset=0&limit=2", // href
		"", // title
	)
	foobar = &Foobar{ID: 1, Name: "Foo bar 1", Date: date}
	foobar.SetLink("self", "/v1/foo/bar/1", "")
	helloWorld.SetEmbedded("foobar", jsonhal.Embedded(foobar))

	// Assert JSON after marshalling is as expected
	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON3)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))

	// Let's test embedded resources
	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}
	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)

	// Add embedded foobars
	foobars = []*Foobar{
		{
			Hal: jsonhal.Hal{
				Links: map[string]*jsonhal.Link{
					"self": {Href: "/v1/foo/bar/1"},
				},
			},
			ID:   1,
			Name: "Foo bar 1",
			Date: date,
		},
		{
			Hal: jsonhal.Hal{
				Links: map[string]*jsonhal.Link{
					"self": {Href: "/v1/foo/bar/2"},
				},
			},
			ID:   2,
			Name: "Foo bar 2",
			Date: date,
		},
	}
	helloWorld.SetEmbedded("foobars", jsonhal.Embedded(foobars))

	// Assert JSON after marshalling is as expected
	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON4)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))

	// Let's test multiple embedded resources
	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}
	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)

	// Add embedded foobars
	foobars = []*Foobar{
		{
			Hal: jsonhal.Hal{
				Links: map[string]*jsonhal.Link{
					"self": {Href: "/v1/foo/bar/1"},
				},
			},
			ID:   1,
			Name: "Foo bar 1",
			Date: date,
		},
		{
			Hal: jsonhal.Hal{
				Links: map[string]*jsonhal.Link{
					"self": {Href: "/v1/foo/bar/2"},
				},
			},
			ID:   2,
			Name: "Foo bar 2",
			Date: date,
		},
	}
	helloWorld.SetEmbedded("foobars", jsonhal.Embedded(foobars))

	// Add embedded quxes
	quxes = []*Qux{
		{
			Hal: jsonhal.Hal{
				Links: map[string]*jsonhal.Link{
					"self": {Href: "/v1/qux/1"},
				},
			},
			ID:   1,
			Name: "Qux 1",
		},
		{
			Hal: jsonhal.Hal{
				Links: map[string]*jsonhal.Link{
					"self": {Href: "/v1/qux/2"},
				},
			},
			ID:   2,
			Name: "Qux 2",
		},
	}
	helloWorld.SetEmbedded("quxes", jsonhal.Embedded(quxes))

	// Assert JSON after marshalling is as expected
	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON5)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))
}

func TestGetLink(t *testing.T) {
	t.Parallel()

	helloWorld := new(HelloWorld)

	var (
		link *jsonhal.Link
		err  error
	)

	// Test when object has no links
	link, err = helloWorld.GetLink("self")
	assert.Nil(t, link)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Link \"self\" not found", err.Error())
	}

	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)

	// Test getting a bogus link
	link, err = helloWorld.GetLink("bogus")
	assert.Nil(t, link)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Link \"bogus\" not found", err.Error())
	}

	// Test getting a valid link
	link, err = helloWorld.GetLink("self")
	assert.Nil(t, err)
	if assert.NotNil(t, link) {
		assert.Equal(t, "/v1/hello/world/1", link.Href)
		assert.Equal(t, "", link.Title)
	}
}

func TestDeleteLink(t *testing.T) {
	t.Parallel()

	helloWorld := new(HelloWorld)
	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)
	link, err := helloWorld.GetLink("self")
	assert.NotNil(t, link)
	assert.NoError(t, err)

	helloWorld.DeleteLink("self")
	link, err = helloWorld.GetLink("self")
	assert.Nil(t, link)
	assert.EqualError(t, err, "Link \"self\" not found")
}

func TestGetEmbedded(t *testing.T) {
	t.Parallel()

	helloWorld := new(HelloWorld)

	var (
		embedded jsonhal.Embedded
		err      error
		foobars  []*Foobar
	)

	// Test when object has no embedded resources
	embedded, err = helloWorld.GetEmbedded("foobars")
	assert.Nil(t, embedded)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Embedded \"foobars\" not found", err.Error())
	}

	// Add embedded foobars
	foobars = []*Foobar{
		{ID: 1, Name: "Foo bar 1"},
		{ID: 2, Name: "Foo bar 2"},
	}
	helloWorld.SetEmbedded("foobars", jsonhal.Embedded(foobars))

	// Test geting bogus embedded resources
	embedded, err = helloWorld.GetEmbedded("bogus")
	assert.Nil(t, embedded)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Embedded \"bogus\" not found", err.Error())
	}

	// Test geting valid embedded resources
	embedded, err = helloWorld.GetEmbedded("foobars")
	assert.Nil(t, err)
	if assert.NotNil(t, embedded) {
		reflectedValue := reflect.ValueOf(embedded)
		expectedType := reflect.SliceOf(reflect.TypeOf(new(Foobar)))
		if assert.Equal(t, expectedType, reflectedValue.Type()) {
			assert.Equal(t, 2, reflectedValue.Len())
		}
	}
}

func TestUnmarshalingAndDecodeEmbedded(t *testing.T) {
	t.Parallel()

	// Simplest case
	hw := new(HelloWorld)
	if err := json.Unmarshal(expectedJSON, hw); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, uint(1), hw.ID)
	assert.Equal(t, "Hello World", hw.Name)

	// Single embedded object

	hw = new(HelloWorld)
	if err := json.Unmarshal(expectedJSON3, hw); err != nil {
		log.Fatal(err)
	}

	next, err := hw.GetLink("next")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "", next.Title)
	assert.Equal(t, "/v1/hello/world?offset=4\u0026limit=2", next.Href)

	previous, err := hw.GetLink("previous")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "", previous.Title)
	assert.Equal(t, "/v1/hello/world?offset=0\u0026limit=2", previous.Href)

	self, err := hw.GetLink("self")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "", self.Title)
	assert.Equal(t, "/v1/hello/world?offset=2\u0026limit=2", self.Href)

	f := new(Foobar)
	assert.NoError(t, hw.DecodeEmbedded("foobar", f))
	assert.Equal(t, uint(1), f.ID)
	assert.Equal(t, "Foo bar 1", f.Name)
	assert.Equal(t, "2017-09-12T08:45:20Z", f.Date.Format(time.RFC3339))

	// Slice of embedded objects

	hw = new(HelloWorld)
	if err := json.Unmarshal(expectedJSON4, hw); err != nil {
		log.Fatal(err)
	}

	foobars := make([]*Foobar, 0)
	assert.NoError(t, hw.DecodeEmbedded("foobars", &foobars))
	assert.Equal(t, uint(1), foobars[0].ID)
	assert.Equal(t, "Foo bar 1", foobars[0].Name)
	assert.Equal(t, uint(2), foobars[1].ID)
	assert.Equal(t, "Foo bar 2", foobars[1].Name)
}

func TestCountEmbedded(t *testing.T) {
	t.Parallel()

	hw := new(HelloWorld)
	hw.SetEmbedded("foobar", jsonhal.Embedded(new(Foobar)))
	c, err := hw.CountEmbedded("foobar")
	assert.Equal(t, 0, c)
	assert.Equal(t, "Embedded object is not a slice or a map", err.Error())

	hw = new(HelloWorld)
	hw.SetEmbedded("foobars", jsonhal.Embedded(make([]*Foobar, 0)))
	c, err = hw.CountEmbedded("foobars")
	assert.Equal(t, 0, c)
	assert.NoError(t, err)

	foobars := []*Foobar{
		{ID: 1, Name: "Foo bar 1"},
		{ID: 2, Name: "Foo bar 2"},
	}
	hw.SetEmbedded("foobars", jsonhal.Embedded(foobars))
	c, err = hw.CountEmbedded("foobars")
	assert.Equal(t, 2, c)
	assert.NoError(t, err)

	bars := make(map[string]string, 2)
	bars["a"] = "b"
	bars["c"] = "d"
	hw.SetEmbedded("bars", jsonhal.Embedded(bars))
	c, err = hw.CountEmbedded("bars")
	assert.Equal(t, 2, c)
	assert.NoError(t, err)
}

func TestDeleteEmbedded(t *testing.T) {
	t.Parallel()

	helloWorld := new(HelloWorld)
	var (
		embedded jsonhal.Embedded
		err      error
		foobars  []*Foobar
	)

	// Add embedded foobars
	foobars = []*Foobar{
		{ID: 1, Name: "Foo bar 1"},
		{ID: 2, Name: "Foo bar 2"},
	}
	helloWorld.SetEmbedded("foobars", jsonhal.Embedded(foobars))

	// Test geting valid embedded resources
	embedded, err = helloWorld.GetEmbedded("foobars")
	assert.NoError(t, err)
	assert.NotNil(t, embedded)

	helloWorld.DeleteEmbedded("foobars")
	embedded, err = helloWorld.GetEmbedded("bogus")
	assert.Nil(t, embedded)
	assert.EqualError(t, err, "Embedded \"bogus\" not found")
}
