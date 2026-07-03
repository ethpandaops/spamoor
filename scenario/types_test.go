package scenario

import (
	"testing"
)

func TestPluginDescriptorGetAllScenarios(t *testing.T) {
	d1 := &Descriptor{Name: "one"}
	d2 := &Descriptor{Name: "two"}
	d3 := &Descriptor{Name: "three"}
	d4 := &Descriptor{Name: "four"}

	plugin := &PluginDescriptor{
		Name: "test-plugin",
		Categories: []*Category{
			{
				Name:        "cat1",
				Descriptors: []*Descriptor{d1},
				Children: []*Category{
					{Name: "sub1", Descriptors: []*Descriptor{d2, d3}},
				},
			},
			{Name: "cat2", Descriptors: []*Descriptor{d4}},
		},
	}

	all := plugin.GetAllScenarios()
	if len(all) != 4 {
		t.Fatalf("expected 4 scenarios across nested categories, got %d", len(all))
	}

	expected := []string{"one", "two", "three", "four"}
	for i, desc := range all {
		if desc.Name != expected[i] {
			t.Fatalf("expected scenario %q at position %d, got %q", expected[i], i, desc.Name)
		}
	}

	if got := (&PluginDescriptor{}).GetAllScenarios(); len(got) != 0 {
		t.Fatalf("expected no scenarios for empty plugin descriptor, got %d", len(got))
	}
}
