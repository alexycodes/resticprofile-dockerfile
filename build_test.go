package main

import (
	"errors"
	"testing"
	"time"
)

func TestGetDockerArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		tag        string
		dockerfile string
		platform   string
		latest     string
		wantArgs   []string
	}{
		{
			name:       "non-latest",
			tag:        "0.32.0",
			dockerfile: "tags/0.32.0/Dockerfile",
			platform:   "linux/amd64",
			latest:     "0.33.0",
			wantArgs: []string{
				"buildx", "build",
				"--push",
				"--no-cache",
				"--file", "tags/0.32.0/Dockerfile",
				"--platform", "linux/amd64",
				"--tag", image + ":0.32.0",
				".",
			},
		},
		{
			name:       "latest",
			tag:        "0.33.0",
			dockerfile: "tags/0.33.0/Dockerfile",
			platform:   "linux/amd64",
			latest:     "0.33.0",
			wantArgs: []string{
				"buildx", "build",
				"--push",
				"--no-cache",
				"--file", "tags/0.33.0/Dockerfile",
				"--platform", "linux/amd64",
				"--tag", image + ":0.33.0",
				"--tag", image + ":latest",
				".",
			},
		},
		{
			name:       "multi-platform",
			tag:        "0.33.0",
			dockerfile: "tags/0.33.0/Dockerfile",
			platform:   "linux/amd64,linux/arm64",
			latest:     "0.33.0",
			wantArgs: []string{
				"buildx", "build",
				"--push",
				"--no-cache",
				"--file", "tags/0.33.0/Dockerfile",
				"--platform", "linux/amd64,linux/arm64",
				"--tag", image + ":0.33.0",
				"--tag", image + ":latest",
				".",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := getDockerArgs(tt.tag, tt.dockerfile, tt.platform, tt.latest)

			if len(got) != len(tt.wantArgs) {
				t.Errorf("Expected %d args, got %d: %v", len(tt.wantArgs), len(got), got)
			}

			for i := range min(len(got), len(tt.wantArgs)) {
				if got[i] != tt.wantArgs[i] {
					t.Errorf("At index %d: expected %q, got %q", i, tt.wantArgs[i], got[i])
				}
			}
		})
	}
}

func TestSortTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   []string
		want    []string
		wantErr bool
	}{
		{
			name:    "realistic list",
			input:   []string{"0.32.0", "0.33.0", "0.10.0", "0.9.0", "0.6.0", "0.6.1"},
			want:    []string{"0.6.0", "0.6.1", "0.9.0", "0.10.0", "0.32.0", "0.33.0"},
			wantErr: false,
		},
		{
			name:    "missing parts",
			input:   []string{"0.33", "0", "0.33.1", "0.33.0"},
			want:    []string{"0.33.0", "0.33.1", "0.33", "0"},
			wantErr: false,
		},
		{
			name:    "single",
			input:   []string{"0.33.0"},
			want:    []string{"0.33.0"},
			wantErr: false,
		},
		{
			name:    "empty",
			input:   []string{},
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "invalid",
			input:   []string{"0.6.z", "0.6.1"},
			want:    []string{"0.6.z", "0.6.1"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := sortTags(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected an error, got nil")
				}

				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %s", err.Error())
			}

			if len(tt.input) != len(tt.want) {
				t.Errorf("Expected length %d, got %d", len(tt.want), len(tt.input))
			}

			for i := range len(tt.input) {
				if tt.input[i] != tt.want[i] {
					t.Errorf("At index %d: expected %s, got %s", i, tt.want[i], tt.input[i])
				}
			}
		})
	}
}

func TestRetry(t *testing.T) {
	t.Parallel()

	t.Run("succeeds on first attempt", func(t *testing.T) {
		t.Parallel()

		calls := 0

		fn := func() error {
			calls++
			return nil
		}

		retryErr := retry(fn, 3, 0)

		if retryErr != nil {
			t.Errorf("Expected no error, got %v", retryErr)
		}

		if calls != 1 {
			t.Errorf("Expected 1 call, got %d", calls)
		}
	})

	t.Run("succeeds on second attempt", func(t *testing.T) {
		t.Parallel()

		calls := 0

		fn := func() error {
			calls++
			if calls < 2 {
				return errors.New("temporary failure")
			}
			return nil
		}

		retryErr := retry(fn, 3, 0)

		if retryErr != nil {
			t.Errorf("Expected no error, got %v", retryErr)
		}

		if calls != 2 {
			t.Errorf("Expected 2 calls, got %d", calls)
		}
	})

	t.Run("fails after all attempts exhausted", func(t *testing.T) {
		t.Parallel()

		calls := 0
		err := errors.New("permanent failure")

		fn := func() error {
			calls++
			return err
		}

		retryErr := retry(fn, 3, 0)

		if retryErr == nil {
			t.Error("Expected an error, got nil")
		}

		if !errors.Is(retryErr, err) {
			t.Errorf("Expected error to wrap %v, got %v", err, retryErr)
		}

		if calls != 3 {
			t.Errorf("Expected 3 calls, got %d", calls)
		}
	})

	t.Run("sleeps between retries", func(t *testing.T) {
		t.Parallel()

		calls := 0
		delay := 50 * time.Millisecond
		start := time.Now()

		fn := func() error {
			calls++
			return errors.New("fail")
		}

		_ = retry(fn, 3, delay)

		elapsed := time.Since(start)

		// 3 attempts = 2 sleeps
		minExpected := 2 * delay

		if elapsed < minExpected {
			t.Errorf("Expected at least %v elapsed, got %v", minExpected, elapsed)
		}
	})
}
