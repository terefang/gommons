package xlua

import rt "github.com/arnodel/golua/runtime"

func FromLuaValue(runtime rt.Runtime, value rt.Value, primitive bool) any {
	// convert only to a primitive value
	switch value.TypeName() {
	case "number":
		if _i, ok := value.TryInt(); ok {
			return _i
		}
		return value.AsFloat()
	case "boolean":
		return value.AsBool()
	case "string":
		return value.AsString()
	case "table":
		if !primitive {
			//_tab := value.AsTable()
			//_tab.Len()
			panic("not primitive")
		}
		fallthrough
	default:
		return value.AsString()
	}
}
