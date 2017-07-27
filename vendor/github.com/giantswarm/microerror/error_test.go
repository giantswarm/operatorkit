package microerror

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/juju/errgo"
)

func TestMask_Nil(t *testing.T) {
	err := Mask(nil)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestMaskf_Nil(t *testing.T) {
	err := Maskf(nil, "test annotation")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestStack(t *testing.T) {
	tests := []struct {
		desc     string
		depth    int
		newError func() error
	}{
		{
			desc:  "Mask depth=1 constructor=New",
			depth: 1,
			newError: func() error {
				err := New("test")
				return err
			},
		},
		{
			desc:  "Mask depth=2 constructor=New",
			depth: 2,
			newError: func() error {
				err := New("test")
				err = Mask(err)
				return err
			},
		},
		{
			desc:  "Mask/Maskf depth=3 constructor=Newf",
			depth: 3,
			newError: func() error {
				err := Newf("%s", "test")
				err = Mask(err)
				err = Maskf(err, "3")
				return err

			},
		},
		{
			desc:  "Mask depth=1 constructor=fmt.Errorf",
			depth: 1,
			newError: func() error {
				err := fmt.Errorf("test")
				err = Mask(err)
				return err
			},
		},
		{
			desc:  "Mask depth=3 constructor=fmt.Errorf",
			depth: 3,
			newError: func() error {
				err := fmt.Errorf("test")
				err = Mask(err)
				err = Mask(err)
				err = Mask(err)
				return err
			},
		},
		{
			desc:  "Maskf depth=3 constructor=fmt.Errorf",
			depth: 3,
			newError: func() error {
				err := fmt.Errorf("test")
				err = Maskf(err, "1")
				err = Maskf(err, "2")
				err = Maskf(err, "3")
				return err
			},
		},
	}

	for i, tc := range tests {
		err := tc.newError()

		var depth int
		for {
			// Check err location.
			if err, ok := err.(errgo.Locationer); ok {
				file := filepath.Base(err.Location().File)
				wfile := "error_test.go"
				if file != wfile {
					t.Errorf("#%d %s: expected  %s, got %s", i, tc.desc, wfile, file)
				}
			}

			if cerr, ok := err.(errgo.Wrapper); ok {
				depth++
				err = cerr.Underlying()
			} else {
				break
			}
		}

		if tc.depth != depth {
			t.Fatalf("#%d %s: expected depth = %d, got %d", i, tc.desc, tc.depth, depth)
		}

	}
}
