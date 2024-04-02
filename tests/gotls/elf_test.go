package main

import (
	"ecapture/user/config"
	"testing"
)

func Test_GOTLSConfig(t *testing.T) {
	check := func(t *testing.T, path string, want []int, isPie bool) {
		gotls := config.NewGoTLSConfig()
		gotls.Path = path
		gotls.Check() // testcase contain arm64 and amd64, decode instruction will return error

		if gotls.IsPieBuildMode != isPie {
			t.Fatalf("this testcase build error! gotPIE(%v) vs wantPIE(%v)", gotls.IsPieBuildMode, isPie)
		}
		got := gotls.ReadTlsAddrs // function: "crypto/tls.(*Conn).Read"

		if len(got) != len(want) {
			t.Fatalf("leng error: want(%d) vs got(%d)", len(want), len(got))
		}
		for i, g := range got {
			if want[i] != g {
				t.Fatalf("item[%d]: want(%d) vs got(%d)", i, want[i], g)
			}
		}
	}

	t.Run("amd64 default", func(t *testing.T) {
		// go build -o ./amd64_default/main ./main.go
		path := "./testdata/amd64_default/main"
		check(t, path, []int{
			// this data get form cli: `./testdata/get_ret.sh ./testdata/normal_x86/main`
			0x1a47b3, 0x1a47da, 0x1a4856, 0x1a49ab, 0x1a49da, 0x1a4a49, 0x1a4a63,
		}, false)
	})

	t.Run("amd64 pie", func(t *testing.T) {
		// go build -buildmode=pie -o ./amd64_default/main ./main.go
		path := "./testdata/amd64_pie/main"
		check(t, path, []int{0x1a4e13, 0x1a4e3a, 0x1a4eb6, 0x1a500b, 0x1a503a, 0x1a50a9, 0x1a50c3}, true)
	})

	t.Run("ubuntu 22.04 default", func(t *testing.T) {
		// in docker ubuntu 22.04
		check(t, "./testdata/ubuntu_default/main", []int{
			0x19fbb3, 0x19fbda, 0x19fc56, 0x19fdab, 0x19fdda, 0x19fe49, 0x19fe63,
		}, false)
	})

	t.Run("ubuntu 22.04 pie", func(t *testing.T) {
		// in docker ubuntu 22.04
		check(t, "./testdata/ubuntu_pie/main", []int{
			0x1a0213, 0x1a023a, 0x1a02b6, 0x1a040b, 0x1a043a, 0x1a04a9, 0x1a04c3,
		}, true)
	})

	t.Run("arm64 default", func(t *testing.T) {
		// GOOS=linux GOARCH=arm64 go build -o ./arm64_pie/main ./main.go
		check(t, "./testdata/arm64_default/main", []int{
			0x171510, 0x171538, 0x17159c, 0x1716b8, 0x1716e8, 0x171748, 0x171764,
		}, false)
	})

	t.Run("arm64 pie", func(t *testing.T) {
		// GOOS=linux GOARCH=arm64 go build -buildmode=pie -o ./arm64_pie/main ./main.go
		check(t, "./testdata/arm64_pie/main", []int{
			0x171510, 0x171538, 0x17159c, 0x1716b8, 0x1716e8, 0x171748, 0x171764,
		}, true)
	})

}
