package model

func (m *Model) HasFederatedTypes() bool {
	for _, o := range m.Objects() {
		if o.HasDirective("key") {
			return true
		}
	}

	return false
}
