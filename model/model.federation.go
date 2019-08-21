package model

func (m *Model) HasFederatedTypes() bool {
	for _, o := range m.Objects() {
		if o.IsFederatedType() {
			return true
		}
	}

	return false
}

func (o *Object) IsFederatedType() bool {
	return o.HasDirective("key")
}
