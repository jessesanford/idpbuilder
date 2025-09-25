package k8s

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/yaml"
)

// Deserializer handles deserialization of Kubernetes objects from YAML/JSON
type Deserializer struct {
	scheme     *runtime.Scheme
	codecs     serializer.CodecFactory
	universal  runtime.Decoder
}

// NewDeserializer creates a new Deserializer with the provided scheme
func NewDeserializer(scheme *runtime.Scheme) *Deserializer {
	if scheme == nil {
		scheme = NewSchemeBuilder().Build()
	}

	codecs := serializer.NewCodecFactory(scheme)
	universal := codecs.UniversalDeserializer()

	return &Deserializer{
		scheme:    scheme,
		codecs:    codecs,
		universal: universal,
	}
}

// Decode deserializes YAML or JSON data into a runtime.Object
func (d *Deserializer) Decode(data []byte) (runtime.Object, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	// Try to determine if this is JSON or YAML
	isJSON := d.isJSON(data)

	var objData []byte
	var err error

	if isJSON {
		objData = data
	} else {
		// Convert YAML to JSON
		objData, err = yaml.YAMLToJSON(data)
		if err != nil {
			return nil, fmt.Errorf("failed to convert YAML to JSON: %w", err)
		}
	}

	// Decode using the universal deserializer
	obj, _, err := d.universal.Decode(objData, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decode object: %w", err)
	}

	return obj, nil
}

// DecodeInto deserializes data into a provided runtime.Object
func (d *Deserializer) DecodeInto(data []byte, obj runtime.Object) error {
	if len(data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	if obj == nil {
		return fmt.Errorf("target object cannot be nil")
	}

	// Try to determine if this is JSON or YAML
	isJSON := d.isJSON(data)

	var objData []byte
	var err error

	if isJSON {
		objData = data
	} else {
		// Convert YAML to JSON
		objData, err = yaml.YAMLToJSON(data)
		if err != nil {
			return fmt.Errorf("failed to convert YAML to JSON: %w", err)
		}
	}

	// Decode into the provided object
	_, _, err = d.universal.Decode(objData, nil, obj)
	if err != nil {
		return fmt.Errorf("failed to decode into object: %w", err)
	}

	return nil
}

// isJSON checks if data appears to be JSON format
func (d *Deserializer) isJSON(data []byte) bool {
	// Simple heuristic: JSON typically starts with { or [
	// This is not foolproof but good enough for our use case
	if len(data) == 0 {
		return false
	}

	// Skip whitespace
	for _, b := range data {
		if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
			continue
		}
		return b == '{' || b == '['
	}

	return false
}

// GetScheme returns the scheme used by this deserializer
func (d *Deserializer) GetScheme() *runtime.Scheme {
	return d.scheme
}

// ValidateObject performs basic validation on a decoded object
func (d *Deserializer) ValidateObject(obj runtime.Object) error {
	if obj == nil {
		return fmt.Errorf("object cannot be nil")
	}

	// Check if the object is known to our scheme
	gvks, _, err := d.scheme.ObjectKinds(obj)
	if err != nil {
		return fmt.Errorf("object type not recognized by scheme: %w", err)
	}

	if len(gvks) == 0 {
		return fmt.Errorf("object has no associated GroupVersionKind")
	}

	return nil
}

// DecodeMultiple deserializes multiple YAML documents separated by ---
func (d *Deserializer) DecodeMultiple(data []byte) ([]runtime.Object, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	// Split by YAML document separator
	docs := splitYAMLDocuments(string(data))
	objects := make([]runtime.Object, 0, len(docs))

	for docIndex, doc := range docs {
		if len(doc) == 0 {
			continue
		}

		obj, err := d.Decode([]byte(doc))
		if err != nil {
			return nil, fmt.Errorf("failed to decode document %d: %w", docIndex, err)
		}

		objects = append(objects, obj)
	}

	return objects, nil
}

// splitYAMLDocuments splits a YAML stream into individual documents
func splitYAMLDocuments(data string) []string {
	// Split by YAML document separator (---)
	docs := strings.Split(data, "---")
	result := make([]string, 0, len(docs))

	for _, doc := range docs {
		trimmed := strings.TrimSpace(doc)
		if len(trimmed) > 0 {
			result = append(result, trimmed)
		}
	}

	return result
}