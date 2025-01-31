package page

import "testing"

func TestPageString(t *testing.T) {
	page := Page{
		Title:     "Test",
		Content:   []byte("Test"),
		Template:  []byte("Test"),
		Params:    map[string]interface{}{"testKey": "testParam"},
		AssetPath: "test/asset",
		BuildPath: "test/build",
	}

	expected := `Title: Test
Content: Test
Template: Test
Params: map[testKey:testParam]
AssetPath: test/asset
BuildPath: test/build`
	if page.String() != expected {
		t.Errorf("Expected %s, got %s", expected, page.String())
	}
}
