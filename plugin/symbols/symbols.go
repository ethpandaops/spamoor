// Package loader provides dynamic scenario loading using Yaegi.
package symbols

import "reflect"

// Symbols holds all extracted symbols for the Yaegi interpreter.
// Each package's symbols are added via init() functions in the generated files.
var Symbols = make(map[string]map[string]reflect.Value)
