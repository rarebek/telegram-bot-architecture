package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MessageService handles loading and retrieving message templates
type MessageService struct {
	commandTemplates map[string]interface{}
	keywordTemplates map[string]interface{}
}

// NewMessageService creates a new MessageService instance
func NewMessageService() (*MessageService, error) {
	service := &MessageService{
		commandTemplates: make(map[string]interface{}),
		keywordTemplates: make(map[string]interface{}),
	}

	// Load command templates
	if err := service.loadJSONFile("templates/messages/commands.json", &service.commandTemplates); err != nil {
		return nil, fmt.Errorf("failed to load command templates: %w", err)
	}

	// Load keyword templates
	if err := service.loadJSONFile("templates/messages/keywords.json", &service.keywordTemplates); err != nil {
		return nil, fmt.Errorf("failed to load keyword templates: %w", err)
	}

	return service, nil
}

// loadJSONFile loads a JSON file into the provided destination map
func (s *MessageService) loadJSONFile(path string, dest interface{}) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// GetCommandTemplate retrieves a command template by name
func (s *MessageService) GetCommandTemplate(command string) (map[string]interface{}, error) {
	if template, ok := s.commandTemplates[command].(map[string]interface{}); ok {
		return template, nil
	}
	return nil, fmt.Errorf("command template not found: %s", command)
}

// GetCommandText gets the text for a specific command
func (s *MessageService) GetCommandText(command string) (string, error) {
	template, err := s.GetCommandTemplate(command)
	if err != nil {
		return "", err
	}

	if text, ok := template["text"].(string); ok {
		return text, nil
	}

	return "", fmt.Errorf("text not found for command: %s", command)
}

// GetVersionText gets the text for a specific Go version
func (s *MessageService) GetVersionText(version string) (string, error) {
	template, err := s.GetCommandTemplate("version")
	if err != nil {
		return "", err
	}

	versions, ok := template["versions"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("versions not found in template")
	}

	// Check if the requested version exists
	if versionData, ok := versions[version].(map[string]interface{}); ok {
		if text, ok := versionData["text"].(string); ok {
			return text, nil
		}
	}

	// Use default template if version not found
	if defaultData, ok := versions["default"].(map[string]interface{}); ok {
		if text, ok := defaultData["text"].(string); ok {
			// Replace version placeholder
			text = strings.Replace(text, "${version}", version, -1)
			return text, nil
		}
	}

	return "", fmt.Errorf("text not found for version: %s", version)
}

// GetKeywordResponse gets the appropriate response for a message based on keywords
func (s *MessageService) GetKeywordResponse(text string) string {
	lowerText := strings.ToLower(text)

	// Check each keyword category
	for category, data := range s.keywordTemplates {
		if category == "default" {
			continue
		}

		categoryData, ok := data.(map[string]interface{})
		if !ok {
			continue
		}

		keywordsData, ok := categoryData["keywords"].([]interface{})
		if !ok {
			continue
		}

		// Convert keywords to strings and check
		for _, kw := range keywordsData {
			keyword, ok := kw.(string)
			if !ok {
				continue
			}

			if strings.Contains(lowerText, strings.ToLower(keyword)) {
				if responseText, ok := categoryData["text"].(string); ok {
					return responseText
				}
			}
		}
	}

	// Return default response if no keywords matched
	if defaultData, ok := s.keywordTemplates["default"].(map[string]interface{}); ok {
		if text, ok := defaultData["text"].(string); ok {
			return text
		}
	}

	return "Savolingiz uchun rahmat!"
}

// GetMessageForAdmin appends admin-specific content to a message if needed
func (s *MessageService) GetMessageForAdmin(command string, isAdmin bool) (string, error) {
	// Get the base message text
	text, err := s.GetCommandText(command)
	if err != nil {
		return "", err
	}

	// If this is an admin and there is admin-specific content, append it
	if isAdmin {
		template, err := s.GetCommandTemplate(command)
		if err != nil {
			return text, nil
		}

		if adminText, ok := template["admin_text"].(string); ok {
			text += adminText
		}
	}

	return text, nil
}
