package model

// HasFederatedTypes ...
func (m *Model) HasFederatedTypes() bool {
	for _, o := range m.Objects() {
		if o.IsFederatedType() {
			return true
		}
	}

	return false
}

// IsFederatedType ...
func (o *Object) IsFederatedType() bool {
	return o.HasDirective("key")
}
