package templates

// Loaders ...
var Loaders = `package gen

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader"
)

// GetLoaders ...
func GetLoaders(db *DB) map[string]*dataloader.Loader {
	loaders := map[string]*dataloader.Loader{}
	
	{{range $object := .Model.ObjectEntities}}
	{{$object.TableName}}BatchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result
		
		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}

		items := &[]{{$object.Name}}{}
		res := db.Query().Find(items, "id IN (?)", ids)
		if res.Error != nil && !res.RecordNotFound() {
			return []*dataloader.Result{
				&dataloader.Result{Error: res.Error},
			}
		}

		itemMap := make(map[string]{{$object.Name}}, len(keys))
		for _, item := range *items {
			itemMap[item.ID] = item
		}

		for _, key := range keys {
			id := key.String()
			item, ok := itemMap[id]
			if !ok {
				results = append(results, &dataloader.Result{
					Data: nil,
					Error: nil,
					// Error: fmt.Errorf("{{$object.Name}} with id '%s' not found", id),
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  &item,
					Error: nil,
				})
			}
		}
		return results
	}

	loaders["{{$object.Name}}"] = dataloader.NewBatchedLoader({{$object.TableName}}BatchFn, dataloader.WithClearCacheOnBatch())

	{{end}}

	return loaders
}
`
