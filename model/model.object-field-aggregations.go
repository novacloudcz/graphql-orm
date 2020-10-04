package model

// ObjectFieldAggregation specifies which aggregation functions are supported for given field
type ObjectFieldAggregation struct {
	Name string
}

func (o *ObjectField) Aggregations() []ObjectFieldAggregation {
	res := []ObjectFieldAggregation{
		{Name: "Min"},
		{Name: "Max"},
	}
	if o.IsNumeric() {
		res = append(res, ObjectFieldAggregation{Name: "Avg"})
		res = append(res, ObjectFieldAggregation{Name: "Sum"})
	}
	return res
}
