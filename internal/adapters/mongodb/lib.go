package mongodb

import (
	"reflect"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func getUpdatedFields(i any) bson.D {
	v := reflect.ValueOf(i).Elem()
	t := reflect.TypeOf(i).Elem()

	fields := bson.D{}
	for i := range v.NumField() {
		fieldVal := v.Field(i)
		fieldTyp := t.Field(i)
		bsonTag := fieldTyp.Tag.Get("bson")

		if fieldVal.IsNil() {
			continue
		}

		fields = append(fields, bson.E{Key: bsonTag, Value: fieldVal.Elem().Interface()})
	}

	return fields
}
