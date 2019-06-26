package generate

import (
	"encoding/json"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type Secret struct {
	APIVersion string            `json:"apiVersion,omitempty"`
	Kind       string            `json:"kind,omitempty"`
	Type       string            `json:"type,omitempty"`
	ObjectMeta Metadata          `json:"metadata,omitempty"`
	Data       map[string][]byte `json:"data,omitempty"`
}

type Metadata struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

func CopyKubernetesSecret(s *v1.Secret) (*Secret, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	secret := &Secret{}
	if err := json.Unmarshal(j, secret); err != nil {
		return nil, err
	}
	return secret, nil
}

func (s *Secret) AsYAML() ([]byte, error) {
	return yaml.Marshal(s)
}
