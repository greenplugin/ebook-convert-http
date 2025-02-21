package contract

// String returns the string representation of the Port type.
func (p *Port) String() string {
	if p == nil {
		return ""
	}
	return string(*p)
}

// Set sets the value of the port with possible validation (validation can be added if needed).
func (p *Port) Set(value string) error {
	*p = Port(value)
	return nil
}

// Type returns the type of the flag.
func (p *Port) Type() string {
	return "port"
}
