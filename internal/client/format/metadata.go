package format

import (
	"fmt"

	"github.com/gophkeeper/gophkeeper/proto"
)

// MetadataToDisplayLines возвращает строки для отображения метаданных ("  key: value").
func MetadataToDisplayLines(metadata []*proto.Metadata) []string {
	if len(metadata) == 0 {
		return nil
	}
	lines := make([]string, 0, len(metadata))
	for _, md := range metadata {
		if md == nil {
			continue
		}
		lines = append(lines, fmt.Sprintf("  %s: %s", md.Key, md.Value))
	}
	return lines
}
