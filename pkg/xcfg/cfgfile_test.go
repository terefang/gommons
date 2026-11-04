package xcfg

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseCfgString_DirectivesAndComments(t *testing.T) {
	input := `
	// Double slash comment
	-- Double dash comment
	# Hash comment
	; Semicolon comment
	: Colon comment
	% Percent comment
	! Exclamation comment
	$ Dollar comment
	* Asterisk comment

	server 127.0.0.1 8080
	timeout 30s
	`

	cfg, err := ParseCfgString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Configs) != 2 {
		t.Fatalf("expected 2 root configs, got %d", len(cfg.Configs))
	}

	// Verify 'server' command
	if cfg.Configs[0].Cmd != "server" {
		t.Errorf("expected cmd 'server', got %q", cfg.Configs[0].Cmd)
	}
	expectedArgs0 := []string{"127.0.0.1", "8080"}
	if !reflect.DeepEqual(cfg.Configs[0].Args, expectedArgs0) {
		t.Errorf("expected args %v, got %v", expectedArgs0, cfg.Configs[0].Args)
	}

	// Verify 'timeout' command
	if cfg.Configs[1].Cmd != "timeout" {
		t.Errorf("expected cmd 'timeout', got %q", cfg.Configs[1].Cmd)
	}
	expectedArgs1 := []string{"30s"}
	if !reflect.DeepEqual(cfg.Configs[1].Args, expectedArgs1) {
		t.Errorf("expected args %v, got %v", expectedArgs1, cfg.Configs[1].Args)
	}
}

func TestParseCfgString_Quoting(t *testing.T) {
	input := `
	double_quote "hello world" "foo bar"
	single_quote 'hello world' 'foo bar'
	backtick_quote ` + "`hello world` `foo bar`" + `
	mixed_quotes "double" 'single' ` + "`backtick`" + `
	`

	cfg, err := ParseCfgString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		cmd  string
		want []string
	}{
		{"double_quote", []string{"hello world", "foo bar"}},
		{"single_quote", []string{"hello world", "foo bar"}},
		{"backtick_quote", []string{"hello world", "foo bar"}},
		{"mixed_quotes", []string{"double", "single", "backtick"}},
	}

	for i, tt := range tests {
		got := cfg.Configs[i]
		if got.Cmd != tt.cmd {
			t.Errorf("[%d] expected cmd %q, got %q", i, tt.cmd, got.Cmd)
		}
		if !reflect.DeepEqual(got.Args, tt.want) {
			t.Errorf("[%d] expected args %v, got %v", i, tt.want, got.Args)
		}
	}
}

func TestParseCfgString_Blocks(t *testing.T) {
	input := `
	route /api {
		backend localhost:8080
		methods GET POST
	}
	`

	cfg, err := ParseCfgString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Configs) != 1 {
		t.Fatalf("expected 1 root config, got %d", len(cfg.Configs))
	}

	parent := cfg.Configs[0]
	if parent.Cmd != "route" || !reflect.DeepEqual(parent.Args, []string{"/api"}) {
		t.Errorf("unexpected parent: cmd=%q, args=%v", parent.Cmd, parent.Args)
	}

	if len(parent.SubConfigs) != 2 {
		t.Fatalf("expected 2 sub-configs, got %d", len(parent.SubConfigs))
	}

	child1 := parent.SubConfigs[0]
	if child1.Cmd != "backend" || !reflect.DeepEqual(child1.Args, []string{"localhost:8080"}) {
		t.Errorf("unexpected child1: cmd=%q, args=%v", child1.Cmd, child1.Args)
	}

	child2 := parent.SubConfigs[1]
	if child2.Cmd != "methods" || !reflect.DeepEqual(child2.Args, []string{"GET", "POST"}) {
		t.Errorf("unexpected child2: cmd=%q, args=%v", child2.Cmd, child2.Args)
	}
}

func TestParseCfgString_LineContinuation(t *testing.T) {
	input := `
	listen 127.0.0.1 \
		8080 \
		tls
	`

	cfg, err := ParseCfgString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Configs) != 1 {
		t.Fatalf("expected 1 root config, got %d", len(cfg.Configs))
	}

	got := cfg.Configs[0]
	if got.Cmd != "listen" {
		t.Errorf("expected cmd 'listen', got %q", got.Cmd)
	}

	expectedArgs := []string{"127.0.0.1", "8080", "tls"}
	if !reflect.DeepEqual(got.Args, expectedArgs) {
		t.Errorf("expected args %v, got %v", expectedArgs, got.Args)
	}
}

func TestParseCfgString_BlockWithLineContinuation(t *testing.T) {
	input := `
server localhost {
  listen 127.0.0.1 \
    8080 \
    tls
}
`

	cfg, err := ParseCfgString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Configs) != 1 {
		t.Fatalf("expected 1 root config block, got %d", len(cfg.Configs))
	}

	parent := cfg.Configs[0]
	if parent.Cmd != "server" {
		t.Errorf("expected parent cmd 'server', got %q", parent.Cmd)
	}
	expectedParentArgs := []string{"localhost"}
	if !reflect.DeepEqual(parent.Args, expectedParentArgs) {
		t.Errorf("expected parent args %v, got %v", expectedParentArgs, parent.Args)
	}

	if len(parent.SubConfigs) != 1 {
		t.Fatalf("expected 1 sub-config inside block, got %d", len(parent.SubConfigs))
	}

	child := parent.SubConfigs[0]
	if child.Cmd != "listen" {
		t.Errorf("expected child cmd 'listen', got %q", child.Cmd)
	}

	expectedChildArgs := []string{"127.0.0.1", "8080", "tls"}
	if !reflect.DeepEqual(child.Args, expectedChildArgs) {
		t.Errorf("expected child args %v, got %v", expectedChildArgs, child.Args)
	}
}

func TestParseCfgString_ExtendedMultilineString(t *testing.T) {
	t.Run("Root extended string", func(t *testing.T) {
		input := "\nscript ```\n#!/bin/bash\necho 'hello'\nexit 0\n        ```\n"

		cfg, err := ParseCfgString(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(cfg.Configs) != 1 {
			t.Fatalf("expected 1 config, got %d", len(cfg.Configs))
		}

		block := cfg.Configs[0]
		if block.Cmd != "script" {
			t.Errorf("expected cmd 'script', got %q", block.Cmd)
		}

		expectedExt := "#!/bin/bash\necho 'hello'\nexit 0\n"
		if block.Ext != expectedExt {
			t.Errorf("expected Ext:\n%q\ngot:\n%q", expectedExt, block.Ext)
		}
	})

	t.Run("Nested extended string inside block", func(t *testing.T) {
		input := "\n		task setup {\n			run ```\necho 'task running'\n        ```\n		}\n"

		cfg, err := ParseCfgString(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(cfg.Configs) != 1 || len(cfg.Configs[0].SubConfigs) != 1 {
			t.Fatalf("expected nested structure, got: %+v", cfg.Configs)
		}

		child := cfg.Configs[0].SubConfigs[0]
		if child.Cmd != "run" {
			t.Errorf("expected child cmd 'run', got %q", child.Cmd)
		}

		expectedExt := "echo 'task running'\n"
		if child.Ext != expectedExt {
			t.Errorf("expected Ext:\n%q\ngot:\n%q", expectedExt, child.Ext)
		}
	})
}

func TestParseCfgString_Errors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		errContains string
	}{
		{
			name: "Nested blocks (invalid depth)",
			input: `
			parent {
				child {
					key value
				}
			}
			`,
			errContains: "invalid depth",
		},
		{
			name: "Malformed block close suffix",
			input: `
			server localhost}
			`,
			errContains: "invalid block-end",
		},
		{
			name: "Unclosed block",
			input: `
			route /api {
				backend localhost:8080
			`,
			errContains: "unclosed block",
		},
		{
			name:        "Unclosed multiline string",
			input:       "\nscript ```\necho 'test'\n",
			errContains: "unclosed multi-line string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCfgString(tt.input)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.errContains)
			}
			if !reflect.ValueOf(err).IsValid() || !reflect.ValueOf(err.Error()).IsValid() {
				t.Fatalf("invalid error returned")
			}
		})
	}
}

func TestParseCfgFile(t *testing.T) {
	t.Run("Valid file read", func(t *testing.T) {
		content := `
            app "My App"
            port 9000
        `

		tmpDir := t.TempDir()
		tmpPath := filepath.Join(tmpDir, "config.cdy")

		if err := os.WriteFile(tmpPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}

		cfg, err := ParseCfgFile(tmpPath)
		if err != nil {
			t.Fatalf("ParseCfgFile failed: %v", err)
		}

		if len(cfg.Configs) != 2 {
			t.Fatalf("expected 2 configs, got %d", len(cfg.Configs))
		}

		if cfg.Configs[0].Cmd != "app" || cfg.Configs[0].Args[0] != "My App" {
			t.Errorf("unexpected first block: %+v", cfg.Configs[0])
		}
	})

	t.Run("Non-existent file", func(t *testing.T) {
		_, err := ParseCfgFile("/path/does/not/exist/config.cdy")
		if err == nil {
			t.Errorf("expected file read error, got nil")
		}
	})
}
