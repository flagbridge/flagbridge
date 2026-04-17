package provider

import "fmt"

// Registry holds the Providers available to the AI router. Registration is
// typically done at server boot from config. Registry is NOT safe for
// concurrent registration — callers are expected to register all providers
// before serving requests. Get is safe for concurrent use.
type Registry struct {
	providers map[string]Provider
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{providers: make(map[string]Provider)}
}

// Register adds a provider under its Name(). Later registrations overwrite
// earlier ones for the same name — handy for swapping in a FakeProvider in tests.
func (r *Registry) Register(p Provider) {
	r.providers[p.Name()] = p
}

// Get returns the provider registered under name, or an error if absent.
func (r *Registry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("ai provider %q not registered", name)
	}
	return p, nil
}

// Names returns the registered provider names in insertion-order-indeterminate
// form. Intended for debugging and /health output.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.providers))
	for n := range r.providers {
		names = append(names, n)
	}
	return names
}
