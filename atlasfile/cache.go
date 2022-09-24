package atlasfile

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/brunoscheufler/atlas/helper"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type cachedFile struct {
	File     Atlasfile `json:"file"`
	CachedAt string    `json:"cachedAt"`
	Hash     string    `json:"hash"`
	Version  string    `json:"version"`
}

func cacheAtlasfile(ctx context.Context, logger logrus.FieldLogger, version, atlasDirPath string, atlasfile *Atlasfile) error {
	cacheFile := filepath.Join(atlasDirPath, "cache.json")

	hash, err := computeAtlasfileHash(atlasDirPath)
	if err != nil {
		return fmt.Errorf("could not compute hash: %w", err)
	}

	cachedFile := cachedFile{
		File:     *atlasfile,
		CachedAt: time.Now().Format(time.RFC3339),
		Hash:     hash,
		Version:  version,
	}

	fileBytes, err := json.Marshal(cachedFile)
	if err != nil {
		return fmt.Errorf("could not marshal atlasfile: %w", err)
	}

	err = os.WriteFile(cacheFile, fileBytes, 0644)
	if err != nil {
		return fmt.Errorf("could not write cache file: %w", err)
	}

	return nil
}

func getCachedAtlasfile(ctx context.Context, logger logrus.FieldLogger, version string, atlasDirPath string) (*Atlasfile, error) {
	cacheFile := filepath.Join(atlasDirPath, "cache.json")

	if !helper.FileExists(cacheFile) {
		return nil, nil
	}

	fileBytes, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, fmt.Errorf("could not read cache file: %w", err)
	}

	cachedFile := &cachedFile{}
	err = json.Unmarshal(fileBytes, cachedFile)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal cache file: %w", err)
	}

	shouldInvalidate, err := shouldInvalidateAtlasfile(ctx, logger, atlasDirPath, *cachedFile, version)
	if err != nil {
		return nil, fmt.Errorf("could not check if cache is valid: %w", err)
	}

	if shouldInvalidate {
		err := os.Remove(cacheFile)
		if err != nil {
			return nil, fmt.Errorf("could not remove cache file: %w", err)
		}

		return nil, nil
	}

	logger.WithFields(logrus.Fields{
		"atlasDirPath": atlasDirPath,
		"cachedAt":     cachedFile.CachedAt,
		"hash":         cachedFile.Hash,
	}).Debugf("using cached atlasfile")

	return &cachedFile.File, nil
}

func shouldInvalidateAtlasfile(ctx context.Context, logger logrus.FieldLogger, dirPath string, file cachedFile, version string) (bool, error) {
	// Check if cache was produced by older version
	if file.Version != version {
		logger.WithFields(logrus.Fields{
			"cachedVersion":  file.Version,
			"currentVersion": version,
		}).Debugln("mismatch in version, invalidating cache")
		return true, nil
	}

	// Parse time from ISO
	cachedAt, err := time.Parse(time.RFC3339, file.CachedAt)
	if err != nil {
		return false, fmt.Errorf("could not parse cachedAt time: %w", err)
	}

	// Check if cache is older than 1 day
	if time.Since(cachedAt) > time.Hour*24 {
		logger.WithField("cachedAt", file.CachedAt).Debugln("cache is older than 1 day, invalidating")
		return true, nil
	}

	currentHash, err := computeAtlasfileHash(dirPath)
	if err != nil {
		return false, fmt.Errorf("could not compute hash: %w", err)
	}

	if currentHash != file.Hash {
		logger.WithFields(logrus.Fields{
			"cachedHash":  file.Hash,
			"currentHash": currentHash,
		}).Debugln("mismatch in hash, invalidating cache")
		return true, nil
	}

	return false, nil
}

func computeAtlasfileHash(dir string) (string, error) {
	var hashBytes []byte

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		// Do not cache the cache file itself
		if d.Name() == "cache.json" {
			return nil
		}

		var shouldCache bool

		knownFiles := []string{
			"go.mod",
			"go.sum",
			"package.json",
			"package-lock.json",
			"yarn.lock",
			"pnpm-lock.yaml",
			"pnpm-workspace.yaml",
		}

		for _, knownFile := range knownFiles {
			if d.Name() == knownFile {
				shouldCache = true
			}
		}

		knownExtensions := []string{
			".go",
			".js",
			".ts",
			".tsx",
			".jsx",
			".json",
			".yaml",
			".yml",
			".toml",
		}

		for _, knownExtension := range knownExtensions {
			if filepath.Ext(d.Name()) == knownExtension {
				shouldCache = true
			}
		}

		if !shouldCache {
			return nil
		}

		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read file: %w", err)
		}

		hashBytes = append(hashBytes, fileBytes...)

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("could not walk directory: %w", err)
	}

	h := sha256.New()
	h.Write(hashBytes)

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
