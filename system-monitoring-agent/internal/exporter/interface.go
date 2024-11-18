package exporter

// Exporter defines the interface for all metric exporters
type Exporter interface {
	Export(data map[string]interface{}) error
	Close() error
}
