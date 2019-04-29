package goradix

import (
	"fmt"
	"os"
	"testing"
)

func TestAlphaOffsetOfMismatch(t *testing.T) {
	cases := []struct {
		left, right string
		want        int
	}{
		// NOTE: While these two empty strings match, the function behaves as if
		// they shared no bytes. Which is true.
		{
			"",
			"",
			0,
		},
		{
			"alpha",
			"",
			0,
		},
		{
			"",
			"bravo",
			0,
		},
		{
			"alpha",
			"bravo",
			0,
		},
		{
			"felicity",
			"fred",
			1,
		},
		{
			"sam",
			"sally",
			2,
		},
		{
			"sam",
			"sam",
			3,
		},
		{
			"sam",
			"samuel",
			3,
		},
	}

	for _, item := range cases {
		got := offsetOfMismatch(item.left, item.right)
		if got != item.want {
			t.Errorf("%q %q; GOT: %v; WANT: %v", item.left, item.right, got, item.want)
		}
	}
}

func TestAlphaLoadEmptyTrie(t *testing.T) {
	root := new(Alpha)

	t.Run("EmptyString", func(t *testing.T) {
		_, ok := root.Load("")
		if got, want := ok, false; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("NonEmptyString", func(t *testing.T) {
		_, ok := root.Load("non-empty")
		if got, want := ok, false; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})
}

func TestAlphaLoadEmptyStringFromTrieWithEmptyString(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				value: 1,
			},
		},
	}

	value, ok := root.Load("")
	if got, want := ok, true; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := value, 1; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaLoadNonEmptyTrie(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				prefix: "robert",
				children: []*Alpha{
					// NOTE: Without element with "" prefix, "robert" is not a
					// valid key.
					&Alpha{
						prefix: "a",
						children: []*Alpha{
							&Alpha{
								value: 1,
							},
						},
					},
					&Alpha{
						prefix: "o",
						children: []*Alpha{
							&Alpha{
								value: 2,
							},
						},
					},
				},
			},
			&Alpha{
				prefix: "sam",
				children: []*Alpha{
					// NOTE: With element with "" prefix, "sam" is a valid key.
					&Alpha{
						prefix: "",
						value:  3,
					},
					&Alpha{
						prefix: "antha",
						children: []*Alpha{
							&Alpha{
								value: 4,
							},
						},
					},
					&Alpha{
						prefix: "uel",
						children: []*Alpha{
							&Alpha{
								value: 5,
							},
						},
					},
				},
			},
		},
	}

	t.Run("ExistingKeys", func(t *testing.T) {
		t.Run("roberta", func(t *testing.T) {
			value, ok := root.Load("roberta")
			if got, want := ok, true; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := value, 1; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("roberto", func(t *testing.T) {
			value, ok := root.Load("roberto")
			if got, want := ok, true; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := value, 2; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("sam", func(t *testing.T) {
			value, ok := root.Load("sam")
			if got, want := ok, true; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := value, 3; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("samantha", func(t *testing.T) {
			value, ok := root.Load("samantha")
			if got, want := ok, true; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := value, 4; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("samuel", func(t *testing.T) {
			value, ok := root.Load("samuel")
			if got, want := ok, true; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := value, 5; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
	})

	t.Run("MissingKey", func(t *testing.T) {
		t.Run("alpha", func(t *testing.T) {
			_, ok := root.Load("alpha")
			if got, want := ok, false; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("robert", func(t *testing.T) {
			_, ok := root.Load("robert")
			if got, want := ok, false; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("zulu", func(t *testing.T) {
			_, ok := root.Load("zulu")
			if got, want := ok, false; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
	})
}

func TestAlphaStoreEmptyStringInEmptyTrie(t *testing.T) {
	root := new(Alpha)

	root.Store("", 1)
	t.Log("\n" + string(root.Bytes()))

	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].value, 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaStoreEmptyStringInTrieWithEmptyString(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				value: 1,
			},
		},
	}

	root.Store("", 2)
	t.Log("\n" + string(root.Bytes()))

	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].value, 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaStoreEmptyStringInTrieWithEmptyStringAndFoo(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				value: 1,
			},
			&Alpha{
				prefix:   "foo",
				children: []*Alpha{&Alpha{value: 2}},
			},
		},
	}

	root.Store("", 3)
	t.Log("\n" + string(root.Bytes()))

	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].value, 3; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaStoreEmptyStringInTrieWithFoo(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				prefix:   "foo",
				children: []*Alpha{&Alpha{value: 2}},
			},
		},
	}

	root.Store("", 3)
	t.Log("\n" + string(root.Bytes()))

	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].value, 3; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaStoreFooInEmptyTrie(t *testing.T) {
	root := new(Alpha)

	root.Store("foo", 1)
	t.Log("\n" + string(root.Bytes()))

	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].prefix, "foo"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[0].value, 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaStoreFooInTrieWithBar(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				prefix:   "bar",
				children: []*Alpha{&Alpha{value: 1}},
			},
		},
	}

	root.Store("foo", 2)
	t.Log("\n" + string(root.Bytes()))

	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := root.children[0].prefix, "bar"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[0].value, 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := root.children[1].prefix, "foo"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[1].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[1].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[1].children[0].value, 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[1].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaStoreFooInTrieWithEmptyStringAndBar(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				prefix: "",
				value:  1,
			},
			&Alpha{
				prefix:   "bar",
				children: []*Alpha{&Alpha{value: 2}},
			},
		},
	}

	root.Store("foo", 3)
	t.Log("\n" + string(root.Bytes()))

	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 3; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// ""
	if got, want := root.children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].value, 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// "bar"
	if got, want := root.children[1].prefix, "bar"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[1].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[1].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[1].children[0].value, 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[1].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// "foo"
	if got, want := root.children[2].prefix, "foo"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[2].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[2].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[2].children[0].value, 3; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[2].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaStoreFooInTrieWithEmptyStringAndFooAndBar(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				prefix: "",
				value:  1,
			},
			&Alpha{
				prefix:   "bar",
				children: []*Alpha{&Alpha{value: 2}},
			},
			&Alpha{
				prefix:   "foo",
				children: []*Alpha{&Alpha{value: 3}},
			},
		},
	}

	root.Store("foo", 4)
	t.Log("\n" + string(root.Bytes()))

	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 3; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// ""
	if got, want := root.children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].value, 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// "bar"
	if got, want := root.children[1].prefix, "bar"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[1].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[1].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[1].children[0].value, 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[1].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// "foo"
	if got, want := root.children[2].prefix, "foo"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[2].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[2].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[2].children[0].value, 4; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[2].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaStoreFooInTrieWithEmptyStringAndFoo(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				prefix: "",
				value:  1,
			},
			&Alpha{
				prefix:   "foo",
				children: []*Alpha{&Alpha{value: 2}},
			},
		},
	}

	root.Store("foo", 3)
	t.Log("\n" + string(root.Bytes()))

	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// ""
	if got, want := root.children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].value, 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// "foo"
	if got, want := root.children[1].prefix, "foo"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[1].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[1].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[1].children[0].value, 3; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[1].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaStoreFooInTrieWithEmptyString(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				prefix: "",
				value:  1,
			},
		},
	}

	root.Store("foo", 2)
	t.Log("\n" + string(root.Bytes()))

	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// ""
	if got, want := root.children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].value, 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// "foo"
	if got, want := root.children[1].prefix, "foo"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[1].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[1].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[1].children[0].value, 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[1].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaStoreFooInTrieWithFoo(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				prefix:   "foo",
				children: []*Alpha{&Alpha{value: 1}},
			},
		},
	}

	root.Store("foo", 2)
	t.Log("\n" + string(root.Bytes()))

	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// "foo"
	if got, want := root.children[0].prefix, "foo"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[0].value, 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaStoreRobertInTrieWithRob(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				prefix:   "rob",
				children: []*Alpha{&Alpha{value: 1}},
			},
		},
	}

	root.Store("robert", 2)
	t.Log("\n" + string(root.Bytes()))

	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// "rob"
	if got, want := root.children[0].prefix, "rob"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[0].value, 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// "robert"
	if got, want := root.children[0].children[1].prefix, "ert"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children[1].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[1].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[1].children[0].value, 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children[1].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaStoreSallyInTrieWithSam(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				prefix:   "sam",
				children: []*Alpha{&Alpha{value: 1}},
			},
		},
	}

	root.Store("sally", 2)
	t.Log("\n" + string(root.Bytes()))

	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// "sa"
	if got, want := root.children[0].prefix, "sa"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// "sally"
	if got, want := root.children[0].children[0].prefix, "lly"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children[0].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[0].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[0].children[0].value, 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children[0].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// "sam"
	if got, want := root.children[0].children[1].prefix, "m"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children[1].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[1].children[0].prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := root.children[0].children[1].children[0].value, 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children[1].children[0].children), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestAlphaStoreOrder(t *testing.T) {
	root := new(Alpha)
	root.Store("/Cabinda/Earle/Dabih", 1)
	root.Store("/Baalath/Dabih/Cabinda", 2)
	root.Store("/Aaron/Dabih/Cabinda", 3)
	// root.Display()

	// root
	if got, want := root.prefix, ""; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// "/"
	if got, want := root.children[0].prefix, "/"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children), 3; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// Aaron
	if got, want := root.children[0].children[0].prefix, "Aaron/Dabih/Cabinda"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children[0].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// Baalath
	if got, want := root.children[0].children[1].prefix, "Baalath/Dabih/Cabinda"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children[1].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	// Cabinda
	if got, want := root.children[0].children[2].prefix, "Cabinda/Earle/Dabih"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := len(root.children[0].children[2].children), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}

	t.Log(string(root.Bytes()))
}

func TestAlphaKeys(t *testing.T) {
	root := new(Alpha)
	root.Store("/Cabinda/Earle/Dabih/Aaron/Aaron/Baalath", 0)
	root.Store("/Baalath/Dabih/Cabinda/Dabih/Aaron/Dabih", 1)
	root.Store("/Aaron/Dabih/Cabinda/Earle/Dabih/Baalath", 2)
	root.Store("/Aaron/Baalath/Dabih/Cabinda/Baalath/Baalath", 3)
	root.Store("/Dabih/Dabih/Earle/Dabih/Cabinda/Baalath", 4)
	root.Store("/Baalath/Cabinda/Aaron/Aaron/Cabinda/Dabih", 5)
	root.Store("/Cabinda/Earle/Aaron/Baalath/Baalath/Aaron", 6)
	root.Store("/Cabinda/Cabinda/Earle/Baalath/Baalath/Baalath", 7)
	root.Store("/Baalath/Earle/Dabih/Dabih/Dabih/Aaron", 8)
	root.Store("/Dabih/Earle/Dabih/Earle/Aaron/Cabinda", 9)
	root.Store("/Cabinda/Earle/Dabih/Aaron/Baalath/Aaron", 10)
	root.Display(os.Stderr)

	t.Run("empty key", func(t *testing.T) {
		t.Skip("hide output")
		keys := root.Keys("", -1)
		if got, want := len(keys), 11; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("entire key", func(t *testing.T) {
		t.Skip("hide output")
		keys := root.Keys("/Baalath/Cabinda/Aaron/Aaron/Cabinda/Dabih", -1)
		if got, want := len(keys), 1; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := keys[0], "/Baalath/Cabinda/Aaron/Aaron/Cabinda/Dabih"; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("prefix of several keys", func(t *testing.T) {
		t.Skip("not working")
		keys := root.Keys("/Cabinda/Earle", -1)
		t.Log(keys)
		if got, want := len(keys), 3; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := keys[0], "/Cabinda/Earle/Aaron/Baalath/Baalath/Aaron"; got != want {
			t.Fatalf("\n\tGOT:\n\t\t%v\n\tWANT:\n\t\t%v", got, want)
		}
		if got, want := keys[1], "/Cabinda/Earle/Dabih/Aaron/Aaron/Baalath"; got != want {
			t.Fatalf("\n\tGOT:\n\t\t%v\n\tWANT:\n\t\t%v", got, want)
		}
		if got, want := keys[2], "/Cabinda/Earle/Dabih/Aaron/Baalath/Aaron"; got != want {
			t.Fatalf("\n\tGOT:\n\t\t%v\n\tWANT:\n\t\t%v", got, want)
		}
	})
}

func TestAlphaStoreOrder2(t *testing.T) {
	root := new(Alpha)
	root.Store("/Cabinda/Earle/Dabih/Aaron/Aaron/Baalath/Earle/Earle/Aaron/Baalath", 0)
	root.Store("/Baalath/Dabih/Cabinda/Dabih/Aaron/Dabih/Baalath/Baalath/Baalath/Baalath", 1)
	root.Store("/Aaron/Dabih/Cabinda/Earle/Dabih/Baalath/Earle/Aaron/Dabih/Aaron", 2)
	root.Store("/Aaron/Baalath/Dabih/Cabinda/Baalath/Baalath/Dabih/Earle/Dabih/Aaron", 3)
	root.Store("/Dabih/Dabih/Earle/Dabih/Cabinda/Baalath/Cabinda/Dabih/Dabih/Dabih", 4)
	root.Store("/Baalath/Cabinda/Aaron/Aaron/Cabinda/Dabih/Cabinda/Cabinda/Cabinda/Earle", 5)
	root.Store("/Cabinda/Earle/Aaron/Baalath/Baalath/Aaron/Aaron/Earle/Dabih/Dabih", 6)
	root.Store("/Cabinda/Cabinda/Earle/Baalath/Baalath/Baalath/Earle/Dabih/Aaron/Dabih", 7)
	root.Store("/Baalath/Earle/Dabih/Dabih/Dabih/Aaron/Aaron/Cabinda/Aaron/Earle", 8)
	root.Store("/Dabih/Earle/Dabih/Earle/Aaron/Cabinda/Aaron/Baalath/Baalath/Baalath", 9)
	root.Display(os.Stderr)
	fmt.Println(string(root.Bytes()))
}

// delete empty string from trie with empty string and no others
// delete empty string from trie with empty string and others
// delete key from trie without key
// delete key from trie that does not cause roll up
// delete key from trie that does cause roll up

func TestAlphaDelete(t *testing.T) {
	root := &Alpha{
		children: []*Alpha{
			&Alpha{
				prefix: "robert",
				children: []*Alpha{
					// NOTE: Without element with "" prefix, "robert" is not a
					// valid key.
					&Alpha{
						prefix: "a",
						children: []*Alpha{
							&Alpha{
								value: 1,
							},
						},
					},
					&Alpha{
						prefix: "o",
						children: []*Alpha{
							&Alpha{
								value: 2,
							},
						},
					},
				},
			},
			&Alpha{
				prefix: "sam",
				children: []*Alpha{
					// NOTE: With element with "" prefix, "sam" is a valid key.
					&Alpha{
						prefix: "",
						value:  3,
					},
					&Alpha{
						prefix: "antha",
						children: []*Alpha{
							&Alpha{
								value: 4,
							},
						},
					},
					&Alpha{
						prefix: "uel",
						children: []*Alpha{
							&Alpha{
								value: 5,
							},
						},
					},
				},
			},
		},
	}

	// ought to do nothing
	t.Run("zulu", func(t *testing.T) {
		root.Delete("zulu")
		t.Log("\n" + string(root.Bytes()))

		_, ok := root.Load("zulu")
		if got, want := ok, false; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	// ought to not cause roll up, because we still have sam and samuel
	t.Run("samantha", func(t *testing.T) {
		root.Delete("samantha")
		t.Log("\n" + string(root.Bytes()))

		_, ok := root.Load("samantha")
		if got, want := ok, false; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}

		if got, want := root.prefix, ""; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := len(root.children), 2; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}

		// "robert"
		if got, want := root.children[0].prefix, "robert"; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := len(root.children[0].children), 2; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		// "roberta"
		if got, want := root.children[0].children[0].prefix, "a"; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		// "roberto"
		if got, want := root.children[0].children[1].prefix, "o"; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		// "sam"
		if got, want := root.children[1].prefix, "sam"; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := len(root.children[1].children), 2; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := root.children[1].children[0].prefix, ""; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := root.children[1].children[0].value, 3; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		// "samuel"
		if got, want := root.children[1].children[1].prefix, "uel"; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
	})

	// ought to cause roll up
	t.Run("roberta", func(t *testing.T) {
		root.Delete("roberta")
		t.Log("\n" + string(root.Bytes()))

		_, ok := root.Load("roberta")
		if got, want := ok, false; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}

		if got, want := root.prefix, ""; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := len(root.children), 2; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}

		// "roberto"
		if got, want := root.children[0].prefix, "roberto"; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := len(root.children[0].children), 1; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := root.children[0].children[0].prefix, ""; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := root.children[0].children[0].value, 2; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := len(root.children[0].children[0].children), 0; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}

		// "sam"
		if got, want := root.children[1].prefix, "sam"; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := len(root.children[1].children), 2; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := root.children[1].children[0].prefix, ""; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := root.children[1].children[0].value, 3; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		// "samuel"
		if got, want := root.children[1].children[1].prefix, "uel"; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
	})

	// ought to cause roll up
	t.Run("sam", func(t *testing.T) {
		root.Delete("sam")

		_, ok := root.Load("sam")
		if got, want := ok, false; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}

		if got, want := root.prefix, ""; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := len(root.children), 2; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}

		// "roberto"
		if got, want := root.children[0].prefix, "roberto"; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := len(root.children[0].children), 1; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := root.children[0].children[0].prefix, ""; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := root.children[0].children[0].value, 2; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := len(root.children[0].children[0].children), 0; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}

		// "samuel"
		if got, want := root.children[1].prefix, "samuel"; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := len(root.children[1].children), 1; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := root.children[1].children[0].prefix, ""; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
		if got, want := root.children[1].children[0].value, 5; got != want {
			t.Fatalf("GOT: %v; WANT: %v", got, want)
		}
	})
}
