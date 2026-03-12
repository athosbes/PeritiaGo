package artifacts

import (
	"log"

	"github.com/athosbes/PeritiaGo/internal/models"
	"golang.org/x/sys/windows/registry"
)

// rot13 decodes UserAssist key names.
func rot13(input string) string {
	b := []byte(input)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] = 'A' + (b[i]-'A'+13)%26
		} else if b[i] >= 'a' && b[i] <= 'z' {
			b[i] = 'a' + (b[i]-'a'+13)%26
		}
	}
	return string(b)
}

// ParseUserAssist pulls recently executed GUI applications for the current user.
func ParseUserAssist() []models.Artifact {
	var artifacts []models.Artifact
	basePath := `Software\Microsoft\Windows\CurrentVersion\Explorer\UserAssist`

	k, err := registry.OpenKey(registry.CURRENT_USER, basePath, registry.READ|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		log.Printf("Failed to open UserAssist: %v\n", err)
		return artifacts
	}
	defer k.Close()

	subkeys, err := k.ReadSubKeyNames(-1)
	if err != nil {
		return artifacts
	}

	for _, sub := range subkeys {
		countPath := basePath + `\` + sub + `\Count`
		cKey, err := registry.OpenKey(registry.CURRENT_USER, countPath, registry.QUERY_VALUE)
		if err != nil {
			continue
		}

		valNames, err := cKey.ReadValueNames(-1)
		if err != nil {
			cKey.Close()
			continue
		}

		for _, vName := range valNames {
			decoded := rot13(vName)
			artifacts = append(artifacts, models.Artifact{
				Name:        vName,
				Type:        "UserAssist",
				Path:        `HKCU\` + countPath,
				Description: "Executed GUI artifact",
				Value:       decoded,
			})
		}
		cKey.Close()
	}

	log.Printf("Extracted %d UserAssist entries\n", len(artifacts))
	return artifacts
}
