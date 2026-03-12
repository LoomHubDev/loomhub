package sync

import (
	"slices"

	"github.com/LoomHubDev/loomhub/internal/models"
)

type DiffEntry struct {
	Path      string `json:"path"`
	Status    string `json:"status"` // "added", "modified", "deleted"
	SourceRef string `json:"source_ref,omitempty"`
	TargetRef string `json:"target_ref,omitempty"`
	SourceContent string `json:"source_content,omitempty"`
	TargetContent string `json:"target_content,omitempty"`
}

// CompareStreams returns a list of entity diffs between source and target streams.
func (h *Handler) CompareStreams(loom *models.Loom, sourceStream, targetStream string) ([]DiffEntry, error) {
	db, err := h.loomDBs.Open(loom.DiskPath)
	if err != nil {
		return []DiffEntry{}, nil
	}

	sourceStreamID := resolveStreamID(db, sourceStream)
	targetStreamID := resolveStreamID(db, targetStream)

	// Get entities for source stream
	sourceEntities := make(map[string]Entity)
	rows, err := db.Query(`
		SELECT DISTINCT e.id, e.space_id, e.path, e.object_ref, e.size, e.mod_time
		FROM entities e
		WHERE e.space_id IN (
			SELECT DISTINCT space_id FROM operations WHERE stream_id = ?
		)
		ORDER BY e.path ASC`, sourceStreamID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var e Entity
			if err := rows.Scan(&e.ID, &e.SpaceID, &e.Path, &e.ObjectRef, &e.Size, &e.ModTime); err == nil {
				sourceEntities[e.Path] = e
			}
		}
		rows.Close()
	}

	// Get entities for target stream
	targetEntities := make(map[string]Entity)
	rows, err = db.Query(`
		SELECT DISTINCT e.id, e.space_id, e.path, e.object_ref, e.size, e.mod_time
		FROM entities e
		WHERE e.space_id IN (
			SELECT DISTINCT space_id FROM operations WHERE stream_id = ?
		)
		ORDER BY e.path ASC`, targetStreamID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var e Entity
			if err := rows.Scan(&e.ID, &e.SpaceID, &e.Path, &e.ObjectRef, &e.Size, &e.ModTime); err == nil {
				targetEntities[e.Path] = e
			}
		}
		rows.Close()
	}

	var diffs []DiffEntry

	// Files in source but not in target → added
	// Files in both but different object_ref → modified
	for path, src := range sourceEntities {
		tgt, exists := targetEntities[path]
		if !exists {
			entry := DiffEntry{Path: path, Status: "added", SourceRef: src.ObjectRef}
			content, err := h.objects.Read(src.ObjectRef)
			if err == nil && isText(content) {
				entry.SourceContent = string(content)
			}
			diffs = append(diffs, entry)
		} else if src.ObjectRef != tgt.ObjectRef {
			entry := DiffEntry{
				Path:      path,
				Status:    "modified",
				SourceRef: src.ObjectRef,
				TargetRef: tgt.ObjectRef,
			}
			srcContent, err := h.objects.Read(src.ObjectRef)
			if err == nil && isText(srcContent) {
				entry.SourceContent = string(srcContent)
			}
			tgtContent, err := h.objects.Read(tgt.ObjectRef)
			if err == nil && isText(tgtContent) {
				entry.TargetContent = string(tgtContent)
			}
			diffs = append(diffs, entry)
		}
	}

	// Files in target but not in source → deleted
	for path, tgt := range targetEntities {
		if _, exists := sourceEntities[path]; !exists {
			entry := DiffEntry{Path: path, Status: "deleted", TargetRef: tgt.ObjectRef}
			content, err := h.objects.Read(tgt.ObjectRef)
			if err == nil && isText(content) {
				entry.TargetContent = string(content)
			}
			diffs = append(diffs, entry)
		}
	}

	if diffs == nil {
		diffs = []DiffEntry{}
	}
	return diffs, nil
}

func isText(data []byte) bool {
	if len(data) > 512*1024 { // skip files > 512KB
		return false
	}
	return !slices.Contains(data, 0)
}
