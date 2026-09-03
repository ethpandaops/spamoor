package spamoor

import (
	"slices"
	"strings"
	"testing"
)

// Client groups are additive on top of the implicit "default" group, and a
// "-name" token removes one. Removing "default" is the only way to keep a
// client out of the selections that name no group.
func TestNewClientGroupParsing(t *testing.T) {
	tests := []struct {
		name    string
		rpchost string
		want    []string
		wantErr string
	}{
		{
			name:    "no prefix stays in default",
			rpchost: "http://localhost:8545",
			want:    []string{"default"},
		},
		{
			name:    "a named group is added to default",
			rpchost: "group(builder)http://localhost:8545",
			want:    []string{"default", "builder"},
		},
		{
			name:    "several groups, comma separated",
			rpchost: "group(a,b)http://localhost:8545",
			want:    []string{"default", "a", "b"},
		},
		{
			name:    "duplicates are ignored",
			rpchost: "group(a)group(a)http://localhost:8545",
			want:    []string{"default", "a"},
		},
		{
			name:    "default can be removed",
			rpchost: "group(builder,-default)http://localhost:8545",
			want:    []string{"builder"},
		},
		{
			name:    "removal is order independent",
			rpchost: "group(-default,builder)http://localhost:8545",
			want:    []string{"builder"},
		},
		{
			name:    "removal works across prefixes",
			rpchost: "group(builder)group(-default)http://localhost:8545",
			want:    []string{"builder"},
		},
		{
			name:    "removing an absent group is a no-op",
			rpchost: "group(-nope)http://localhost:8545",
			want:    []string{"default"},
		},
		{
			name:    "removing every group is refused",
			rpchost: "group(-default)http://localhost:8545",
			wantErr: "no client groups left",
		},
		{
			name:    "a bare dash is refused",
			rpchost: "group(-)http://localhost:8545",
			wantErr: "expected a group name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(&ClientOptions{RpcHost: tt.rpchost})

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected an error containing %q, got groups %v", tt.wantErr, client.GetClientGroups())
				}

				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected an error containing %q, got %q", tt.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := client.GetClientGroups(); !slices.Equal(got, tt.want) {
				t.Fatalf("groups = %v, want %v", got, tt.want)
			}

			for _, group := range tt.want {
				if !client.HasGroup(group) {
					t.Errorf("HasGroup(%q) = false, want true", group)
				}
			}
		})
	}
}

// A client removed from "default" must be invisible to the group-less
// selections used by wallet funding, deployments and scenarios that set no
// client group.
func TestClientRemovedFromDefaultIsNotSelectedByDefault(t *testing.T) {
	reserved, err := NewClient(&ClientOptions{RpcHost: "group(private,-default)http://localhost:8545"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reserved.HasGroup("") {
		t.Error(`HasGroup("") = true, want false (an empty group means default)`)
	}

	if reserved.HasGroup("default") {
		t.Error(`HasGroup("default") = true, want false`)
	}

	if !reserved.HasGroup("private") {
		t.Error(`HasGroup("private") = false, want true`)
	}

	shared, err := NewClient(&ClientOptions{RpcHost: "group(private)http://localhost:8546"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !shared.HasGroup("default") {
		t.Error("without the removal the client must stay in default")
	}
}
