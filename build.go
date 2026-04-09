package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	tagsDir      = "tags"
	latestFile   = "latest"
	platformFile = "platform"
	image        = "alexycodes/resticprofile"
	maxAttempts  = 10
	retryDelay   = 10 * time.Second
)

func main() {
	runErr := run()
	if runErr != nil {
		slog.Error("Failed", "error", runErr)
		os.Exit(1)
	}
}

func run() error {
	latest, latestErr := readLatest()
	if latestErr != nil {
		return latestErr
	}

	tags, tagsErr := readTags()
	if tagsErr != nil {
		return tagsErr
	}

	for _, tag := range tags {
		tagErr := buildTag(tag, latest)
		if tagErr != nil {
			return tagErr
		}
	}

	return nil
}

func buildTag(tag, latest string) error {
	dir := filepath.Join(tagsDir, tag)
	dockerfile := filepath.Join(dir, "Dockerfile")

	if _, dockerfileStatErr := os.Stat(dockerfile); dockerfileStatErr != nil {
		return fmt.Errorf("no Dockerfile in %s", dir)
	}

	platform, platformErr := readPlatform(dir)
	if platformErr != nil {
		return platformErr
	}

	args := getDockerArgs(tag, dockerfile, platform, latest)

	slog.Info(
		"Building",
		"image", image,
		"tag", tag,
		"platform(s)", platform,
		"latest", tag == latest,
	)

	fn := func() error {
		cmd := exec.Command("docker", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	cmdErr := retry(fn, maxAttempts, retryDelay)

	return cmdErr
}

func getDockerArgs(tag, dockerfile, platform, latest string) []string {
	args := []string{
		"buildx", "build",
		"--push",
		"--no-cache",
		"--file", dockerfile,
		"--platform", platform,
		"--tag", image + ":" + tag,
	}

	if tag == latest {
		args = append(args, "--tag", image+":latest")
	}

	args = append(args, ".")

	return args
}

func readTags() ([]string, error) {
	entries, entriesErr := os.ReadDir(tagsDir)
	if entriesErr != nil {
		return []string{}, fmt.Errorf("failed reading %s: %w", tagsDir, entriesErr)
	}

	tags := make([]string, 0, len(entries))

	for _, e := range entries {
		if e.IsDir() {
			tags = append(tags, e.Name())
		}
	}

	sortErr := sortTags(tags)
	if sortErr != nil {
		return []string{}, sortErr
	}

	return tags, nil
}

func sortTags(tags []string) error {
	var sortErr error

	sort.Slice(tags, func(i, j int) bool {
		less, lessErr := versionLess(tags[i], tags[j])
		if lessErr != nil {
			sortErr = lessErr
			return false
		}
		return less
	})

	return sortErr
}

func readLatest() (string, error) {
	bytes, bytesErr := os.ReadFile(latestFile)
	if bytesErr != nil {
		return "", fmt.Errorf("failed reading %s: %w", latestFile, bytesErr)
	}

	return strings.TrimSpace(string(bytes)), nil
}

func readPlatform(dir string) (string, error) {
	file := filepath.Join(dir, platformFile)

	bytes, bytesErr := os.ReadFile(file)
	if bytesErr != nil {
		return "", fmt.Errorf("failed reading platform file: %w", bytesErr)
	}

	platform := strings.TrimSpace(string(bytes))

	return platform, nil
}

// versionLess compares semantic versions and returns true if a is less than b.
func versionLess(a, b string) (bool, error) {
	separator := "."
	partsA := strings.Split(a, separator)
	partsB := strings.Split(b, separator)
	partsCount := max(len(partsA), len(partsB))

	for i := range partsCount {
		if i >= len(partsA) {
			return false, nil
		}

		if i >= len(partsB) {
			return true, nil
		}

		valA, valAErr := strconv.Atoi(partsA[i])
		if valAErr != nil {
			return false, valAErr
		}

		valB, valBErr := strconv.Atoi(partsB[i])
		if valBErr != nil {
			return false, valBErr
		}

		if valA != valB {
			return valA < valB, nil
		}
	}

	return false, nil
}

func retry(fn func() error, attempts int, delay time.Duration) error {
	var fnErr error

	for i := 1; i <= attempts; i++ {
		fnErr = fn()
		if fnErr == nil {
			return nil
		}

		if i < attempts {
			slog.Warn(
				"Attempt failed, retrying",
				"attempt", i,
				"max", attempts,
				"delay", delay,
			)
			time.Sleep(delay)
		}
	}

	return fmt.Errorf("all %d attempts failed: %w", attempts, fnErr)
}
