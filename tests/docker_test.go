package main

import (
	"ecapture/user/config"
	"testing"
)

// readelf -sW ./testdata/docker_aarch64/docker | grep $FUNCTION
// 6649: 0000000000635cb0   928 FUNC    LOCAL  DEFAULT    6 crypto/tls.(*Conn).Read

// in ida:
// .text:0000000000635CB0
// .text:0000000000635CB0
// .text:0000000000635CB0  ; retval_635CB0 __golang crypto_tls__ptr_Conn_Read(crypto_tls_Conn_0 *c, _slice_uint8 b)
// .text:0000000000635CB0    crypto_tls._ptr_Conn.Read               ; CODE XREF: crypto_tls._ptr_Conn.Read+390↓j
// .text:0000000000635CB0

const readFuncAddress = 0x0000000000635cb0

func TestDockerSymbol_ByElfSymbol(t *testing.T) {
	path := "./docker_binary/arm64"
	gotls := config.NewGoTLSConfig()
	gotls.Path = path
	if err := gotls.Init(); err != nil {
		t.Fatal("init error:", err)
	}
	symbol, err := gotls.GetSymbol(config.GoTlsReadFunc)
	if err != nil {
		t.Fatalf("get symbol %s error: %v", config.GoTlsReadFunc, err)
	}
	if symbol.Value != readFuncAddress {
		t.Fatalf("error symbol address: got[0x%08x] vs want[0x%08x]", symbol.Value, readFuncAddress)
	}
}

func TestDockerSymbol_BySymbolTable(t *testing.T) {
	path := "./docker_binary/arm64"
	gotls := config.NewGoTLSConfig()
	gotls.Path = path
	if err := gotls.Init(); err != nil {
		t.Fatal("init error:", err)
	}
	// i think this error  will in this function.
	goSymTab, err := gotls.ReadTable()
	if err != nil {
		t.Fatalf("build symbTable error: %v", err)
	}
	// when this symTab builded, this function address is error,
	// this functino index = 6626,
	// you can set condition breakpoint in `/usr/lib/go/src/debug/gosym/pclntab.go:308` this code build this functino.
	f := goSymTab.LookupFunc(config.GoTlsReadFunc)
	if f == nil {
		t.Fatalf("con't looup symbol %s", config.GoTlsReadFunc)
	}
	if f.Value != readFuncAddress {
		t.Fatalf("error symbol address: got[0x%08x] vs want[0x%08x], dec: 0x%08x", f.Value, readFuncAddress, readFuncAddress-f.Value)
	}
}
