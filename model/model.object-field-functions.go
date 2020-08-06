package model

// ObjectFieldAggregation specifies which aggregation functions are supported for given field
type ObjectFieldAggregation struct {
	Name string
}

func (o *ObjectField) Aggregations() []ObjectFieldAggregation {
	res := []ObjectFieldAggregation{
		{Name: "Min"},
		{Name: "Max"},
		{Name: "Avg"},
	}
	return res
}
