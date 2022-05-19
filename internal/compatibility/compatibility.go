package compatibility

const (
	deprecatedPrefix = "DEPRECATED: "
	removedPrefix    = "REMOVED: "
	changedPrefix    = "CHANGED: "
)

// Compatibility holds compatibility details for configuration
type Compatibility struct {
	deprecated []string
	removed    []string
	changed    []string
	writer     Consoler
}

// NewCompatibility creates a new compatibility checker
func NewCompatibility(w Consoler) *Compatibility {
	return &Compatibility{
		deprecated: []string{},
		removed:    []string{},
		changed:    []string{},
		writer:     w,
	}
}

// AddDeprecated adds a deprecated field to the compatibility list
func (c *Compatibility) AddDeprecated(deprecated ...string) {
	c.deprecated = append(c.deprecated, deprecated...)
}

// AddRemoved adds a removed field to the compatibility list
func (c *Compatibility) AddRemoved(removed ...string) {
	c.removed = append(c.removed, removed...)
}

// AddChanged adds a changed field to the compatibility list
func (c *Compatibility) AddChanged(changed ...string) {
	c.changed = append(c.changed, changed...)
}

// Report return the compatibility issues
func (c *Compatibility) Report() {
	for _, field := range c.removed {
		c.writer.Warn(removedPrefix, field)
	}

	for _, field := range c.deprecated {
		c.writer.Warn(deprecatedPrefix, field)
	}

	for _, field := range c.changed {
		c.writer.Warn(changedPrefix, field)
	}
}
