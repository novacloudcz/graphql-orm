package model

// HasAggregableColumn ...
func (o *Object) HasAggregableColumn() bool {
	for _, column := range o.Columns() {
		if column.IsAggregable() {
			return true
		}
	}
	return false
}

// AggregationsByField ...
func (o *Object) AggregationsByField() (res map[string]*ObjectFieldAggregation) {
	res = map[string]*ObjectFieldAggregation{}
	for _, column := range o.Columns() {
		if column.IsAggregable() {
			for _, agg := range column.Aggregations() {
				val := agg
				res[agg.FieldName()] = &val
			}
		}
	}
	return
}
