package main

import ( // "path/filepath"
	// "delta"
	// dasjson "github.com/GeoNet/delta/internal/das-json"

	"reflect"
	"time"

	"github.com/GeoNet/delta/meta"
)

func collectMatching[T any](Input []T, Match_on string, ID_to_match string, Span_to_match meta.Span) (matching []T) {
	// fmt.Printf("trying to match on %s \n", Match_on)
	for _, object := range Input {

		v := reflect.ValueOf(object)
		t := reflect.TypeOf(object)
		fieldNum := v.NumField()

		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		name := v.FieldByName(Match_on)

		var s meta.Span
		ignoreSpan := true
		spans_match := false

		for i := 0; i < fieldNum; i++ {
			field := v.Field(i)
			fieldName := t.Field(i).Name
			isSet := field.IsValid() && !field.IsZero()
			if fieldName == "Span" && isSet {
				ignoreSpan = false
			}
			// fmt.Printf("Struct %s | field %s | is set? %t \n\n", t.Name(), fieldName, isSet)
		}
		if ignoreSpan {
			spans_match = true
		} else {
			if v.FieldByName("Start").Type() == reflect.TypeOf(time.Time{}) && v.FieldByName("End").Type() == reflect.TypeOf(time.Time{}) {
				s = meta.Span{Start: v.FieldByName("Start").Interface().(time.Time),
					End: v.FieldByName("End").Interface().(time.Time)}
				spans_match = Span_to_match.Overlaps(s)
			}
		}
		// fmt.Print("object ", name.String(), " ", ID_to_match, " ", spans_match, "\n")
		if name.String() == ID_to_match && spans_match {
			// fmt.Print("Matched\n")
			matching = append(matching, object)
		} // else {
		// 	fmt.Print("Failed\n")
		// }
		//#else {
		// fmt.Printf("Spans to compare: %s | %s | returns %t  \n\n", s, Span_to_match, Span_to_match.Overlaps(s))
		// fmt.Printf("Spans to compare: %s | %s | returns %t  \n\n", name.String(), ID_to_match, name.String() == ID_to_match)
		//}

	}
	return matching
}
