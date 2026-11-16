package spamoor

import (
	"strings"
	"testing"
)

func TestRedactRPCURL(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no userinfo", "http://localhost:8545", "http://localhost:8545"},
		{"user and password", "https://user:secret@rpc.example.com/path?x=1", "https://user:xxxxx@rpc.example.com/path?x=1"},
		{"username only", "https://apikey@rpc.example.com", "https://xxxxx@rpc.example.com"},
		{"websocket", "wss://user:secret@rpc.example.com", "wss://user:xxxxx@rpc.example.com"},
		{"unparseable", "https://user:sec ret%zz@rpc.example.com/path", "https://rpc.example.com/path"},
		{"at sign in password", "group(a)https://user:p@ss@rpc.example.com/a@b", "group(a)https://rpc.example.com/a@b"},
		{"headers prefix", "group(a)headers(Authorization:Bearer token)http://localhost:8545", "group(a)headers(xxxxx)http://localhost:8545"},
		{"no scheme", "user:secret@rpc.example.com:8545", "rpc.example.com:8545"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RedactRPCURL(tt.in); got != tt.want {
				t.Fatalf("RedactRPCURL(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestNewClientRedactsURLParseError(t *testing.T) {
	_, err := NewClient(&ClientOptions{RpcHost: "https://user:secret@rpc.example.com:port/"})
	if err == nil {
		t.Fatal("expected an error for an invalid URL")
	}

	if strings.Contains(err.Error(), "secret") {
		t.Fatalf("error contains the password: %v", err)
	}
}
