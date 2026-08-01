package xlua

import (
    rt "github.com/arnodel/golua/runtime"
)

func RegisterTableFunctions(r *rt.Runtime) error {

    _tab := r.GlobalEnv().Get(rt.StringValue("table")).AsTable()

    // luazdf/tab/isarray/isarray.lua
    r.SetEnvGoFunc(_tab, "is_array", tableIsArray, 1, false)
    // luazdf/tab/isdict/isdict.lua
    r.SetEnvGoFunc(_tab, "is_map", tableIsMap, 1, false)
    r.SetEnvGoFunc(_tab, "is_mixed", tableIsMixed, 1, false)

    r.SetEnvGoFunc(_tab, "is_array_only", tableIsArrayOnly, 1, false)
    r.SetEnvGoFunc(_tab, "is_map_only", tableIsMapOnly, 1, false)

    // luazdf/arr/collectk/collectk.lua
    // we take table instead of iterator
    r.SetEnvGoFunc(_tab, "ikeys", tableArrayKeys, 1, false)
    r.SetEnvGoFunc(_tab, "akeys", tableMapKeys, 1, false)
    r.SetEnvGoFunc(_tab, "keys", tableKeys, 1, false)

    // luazdf/arr/collectv/collectv.lua
    // we take table instead of iterator
    r.SetEnvGoFunc(_tab, "ivalues", tableArrayValues, 1, false)
    r.SetEnvGoFunc(_tab, "avalues", tableMapValues, 1, false)
    r.SetEnvGoFunc(_tab, "values", tableValues, 1, false)

    // luazdf/tab/isempty/isempty.lua
    r.SetEnvGoFunc(_tab, "is_empty", tableIsEmpty, 1, false)
    // luazdf/tab/isfilled/isfilled.lua
    r.SetEnvGoFunc(_tab, "is_filled", tableIsFilled, 1, false)
    // luazdf/arr/append/append.lua
    r.SetEnvGoFunc(_tab, "append", tableAppend, 1, true)
    // luazdf/arr/appendall/appendall.lua
    r.SetEnvGoFunc(_tab, "append_all", tableAppendAll, 1, true)

    return nil
}

// function append( arr, v [, ...] ) --> arr
func tableAppend(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }

    TableAppend(_tab, c.Etc()...)

    return c.PushingNext(t.Runtime, rt.BoolValue(true)), nil
}

func tableAppendAll(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }

    TableAppendAll(_tab, c.Etc()...)

    return c.PushingNext(t.Runtime, rt.BoolValue(true)), nil
}

func tableArrayKeys(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }
    _ret := rt.NewTable()
    for _, v := range TableArrayKeys(_tab) {
        TableAppend(_ret, v)
    }
    return c.PushingNext(t.Runtime, rt.AsValue(_ret)), nil
}

func tableMapKeys(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }
    _ret := rt.NewTable()
    for _, v := range TableMapKeys(_tab) {
        TableAppend(_ret, v)
    }
    return c.PushingNext(t.Runtime, rt.AsValue(_ret)), nil
}

func tableKeys(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }
    _ret := rt.NewTable()
    for _, v := range TableKeys(_tab) {
        TableAppend(_ret, v)
    }
    return c.PushingNext(t.Runtime, rt.AsValue(_ret)), nil
}

func tableArrayValues(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }
    _ret := rt.NewTable()
    for _, v := range TableArrayValues(_tab) {
        TableAppend(_ret, v)
    }
    return c.PushingNext(t.Runtime, rt.AsValue(_ret)), nil
}

func tableMapValues(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }
    _ret := rt.NewTable()
    for _, v := range TableMapValues(_tab) {
        TableAppend(_ret, v)
    }
    return c.PushingNext(t.Runtime, rt.AsValue(_ret)), nil
}

func tableValues(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }
    _ret := rt.NewTable()
    for _, v := range TableValues(_tab) {
        TableAppend(_ret, v)
    }
    return c.PushingNext(t.Runtime, rt.AsValue(_ret)), nil
}

func tableIsEmpty(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }
    return c.PushingNext(t.Runtime, rt.BoolValue(TableIsEmpty(_tab))), nil
}

func tableIsFilled(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }
    return c.PushingNext(t.Runtime, rt.BoolValue(!TableIsEmpty(_tab))), nil
}

func tableIsArray(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }
    return c.PushingNext(t.Runtime, rt.BoolValue(TableIsArray(_tab))), nil
}

func tableIsMap(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }
    return c.PushingNext(t.Runtime, rt.BoolValue(TableIsMap(_tab))), nil
}

func tableIsArrayOnly(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }
    return c.PushingNext(t.Runtime, rt.BoolValue(TableIsArrayOnly(_tab))), nil
}

func tableIsMapOnly(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }
    return c.PushingNext(t.Runtime, rt.BoolValue(TableIsMapOnly(_tab))), nil
}

func tableIsMixed(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    _tab, _err := c.TableArg(0)
    if _err != nil {
        return c.PushingNext(t.Runtime, rt.BoolValue(false)), nil
    }
    return c.PushingNext(t.Runtime, rt.BoolValue(TableIsMixed(_tab))), nil
}

// TableIsEmpty reports whether t contains any key-value pairs.
//
// - TableIsEmpty returns true if and only if t contains no values.
// - A table containing at least one key-value pair always returns false.
func TableIsEmpty(t *rt.Table) bool {
    var _k rt.Value = rt.NilValue
    var _ok bool = true

    for _ok {
        _k, _, _ok = t.Next(_k)
        if !_ok || _k.IsNil() {
            break
        }
        return false
    }
    return true
}

// TableIsArray reports whether t contains values stored using array-style keys.
//
// - TableIsArray returns true if and only if t contains at least one key with an integer value.
// - Keys stored as non-integer values do not affect the result.
// - An empty table always returns false.
// - The result is based on the table's visible keys rather than its internal storage layout.
func TableIsArray(t *rt.Table) bool {
    // this does not work because of the way
    // the allocation is optimized
    // ----------------------------------
    // return t.mixedTable.array.itemCount() > 0
    // ----------------------------------
    // we resort to a full/partial table scan
    var _k rt.Value = rt.NilValue
    var _ok bool = true

    for _ok {
        _k, _, _ok = t.Next(_k)
        if !_ok || _k.IsNil() {
            break
        }
        // only if keys are integer numbers
        if _k.TypeName() == "number" {
            switch _k.Interface().(type) {
            case int64:
                return true
            }
        }
    }
    return false
}

// TableIsMap reports whether t contains values stored using map-style keys.
//
// - TableIsMap returns true if and only if t contains at least one key that is not an integer value.
// - Keys with integer values are treated as array-style keys and do not affect the result.
// - An empty table always returns false.
// - The result is based on the table's visible keys rather than its internal storage layout.
func TableIsMap(t *rt.Table) bool {
    var _k rt.Value = rt.NilValue
    var _ok bool = true

    for _ok {
        _k, _, _ok = t.Next(_k)
        if !_ok || _k.IsNil() {
            break
        }
        switch _k.Interface().(type) {
        case int64:
        default:
            return true
        }
    }
    return false
}

// TableIsMixed reports whether t contains both array-style and map-style keys.
//
// - returns true if and only if TableIsArray and TableIsMap both return true.
// - A table with only array-style keys or only map-style keys returns false.
// - An empty table always returns false.
func TableIsMixed(t *rt.Table) bool {
    return TableIsArray(t) && TableIsMap(t)
}

// TableIsArrayOnly reports whether t contains only array-style keys.
//
// - returns true if and only if TableIsArray returns true and TableIsMap returns false.
func TableIsArrayOnly(t *rt.Table) bool {
    return TableIsArray(t) && !TableIsMap(t)
}

// TableIsMapOnly reports whether t contains only map-style keys.
//
// - returns true if and only if TableIsMap returns true and TableIsArray returns false.
func TableIsMapOnly(t *rt.Table) bool {
    return !TableIsArray(t) && TableIsMap(t)
}

// TableMapKeys returns the keys stored in the map part of t.
//
// - Only keys that are not integer Values are included.
// - If the map part of t is empty, TableMapKeys returns an empty slice.
// - The order of the returned keys is the iteration order and should not be relied upon.
func TableMapKeys(t *rt.Table) []rt.Value {
    // we resort to a full/partial table scan
    keys := make([]rt.Value, 0)
    var _k rt.Value = rt.NilValue
    var _ok bool = true

    for _ok {
        _k, _, _ok = t.Next(_k)
        if !_ok || _k.IsNil() {
            break
        }
        switch _k.Interface().(type) {
        case int64:
        default:
            keys = append(keys, _k)
        }
    }
    return keys
}

func TableMapGoKeys(t *rt.Table) []any {
    vkeys := TableMapKeys(t)
    keys := make([]any, 0)

    for _, _k := range vkeys {
        switch _k.Interface().(type) {
        case int64:
        default:
            keys = append(keys, _k.Interface())
        }
    }

    return keys
}

// TableArrayKeys returns the keys stored in the array part of t.
//
// - Only keys with integer Values are included.
// - If the array part of t is empty, TableArrayKeys returns an empty slice.
// - The order of the returned keys is the iteration order and should not be relied upon.
func TableArrayKeys(t *rt.Table) []rt.Value {
    // we resort to a full/partial table scan
    keys := make([]rt.Value, 0)
    var _k rt.Value = rt.NilValue
    var _ok bool = true

    for _ok {
        _k, _, _ok = t.Next(_k)
        if !_ok || _k.IsNil() {
            break
        }
        switch _k.Interface().(type) {
        case int64:
            keys = append(keys, _k)
        }
    }
    return keys
}

// TableKeys returns the keys stored in t.
//
// - If the array part of t is empty, TableKeys returns an empty slice.
// - The order of the returned keys is the iteration order and should not be relied upon.
func TableKeys(t *rt.Table) []rt.Value {
    // we resort to a full/partial table scan
    keys := make([]rt.Value, 0)
    var _k rt.Value = rt.NilValue
    var _ok bool = true

    for _ok {
        _k, _, _ok = t.Next(_k)
        if !_ok || _k.IsNil() {
            break
        }
        keys = append(keys, _k)
    }
    return keys
}

func TableArrayGoKeys(t *rt.Table) []int64 {
    vkeys := TableArrayKeys(t)
    keys := make([]int64, 0)

    for _, _k := range vkeys {
        switch _k.Interface().(type) {
        case int64:
            keys = append(keys, _k.AsInt())
        }
    }

    return keys
}

// TableMapValues returns the values stored in the map part of t.
//
// - Only values whose keys are not integer Values are included.
// - If the map part of t is empty, TableMapValues returns an empty slice.
// - The order of the returned values is the iteration order and should not be relied upon.
func TableMapValues(t *rt.Table) []rt.Value {
    // we resort to a full/partial table scan
    keys := make([]rt.Value, 0)
    var _k rt.Value = rt.NilValue
    var _v rt.Value = rt.NilValue
    var _ok bool = true

    for _ok {
        _k, _v, _ok = t.Next(_k)
        if !_ok || _k.IsNil() {
            break
        }
        switch _k.Interface().(type) {
        case int64:
        default:
            keys = append(keys, _v)
        }
    }
    return keys
}

// TableArrayValues returns the values stored in the array part of t.
//
// - Only values whose keys are integer Values are included.
// - The order of the returned values is the iteration order and should not be relied upon.
func TableArrayValues(t *rt.Table) []rt.Value {
    // we resort to a full/partial table scan
    keys := make([]rt.Value, 0)
    var _k rt.Value = rt.NilValue
    var _v rt.Value = rt.NilValue
    var _ok bool = true

    for _ok {
        _k, _v, _ok = t.Next(_k)
        if !_ok || _k.IsNil() {
            break
        }
        switch _k.Interface().(type) {
        case int64:
            keys = append(keys, _v)
        }
    }
    return keys
}

// TableValues returns all values stored in t.
//
// - Both array and map values are included.
// - Each value in t is returned exactly once.
// - If t is empty, TableValues returns an empty slice.
// - The order of the returned values is the iteration order and should not be relied upon.
func TableValues(t *rt.Table) []rt.Value {
    // we resort to a full/partial table scan
    keys := make([]rt.Value, 0)
    var _k rt.Value = rt.NilValue
    var _v rt.Value = rt.NilValue
    var _ok bool = true

    for _ok {
        _k, _v, _ok = t.Next(_k)
        if !_ok || _k.IsNil() {
            break
        }
        keys = append(keys, _v)
    }
    return keys
}

// TableAppend appends one or more values to the array part of arr.
//
// - Values are appended in the order they are provided.
func TableAppend(t *rt.Table, v ...rt.Value) {
    for _, val := range v {
        i := t.Len()
        t.Set(rt.IntValue(i+1), val)
    }
}

// TableAppendAll appends one or more values to the array part of t.
//
// - Values are appended in the order they are provided.
// - If an argument is a table, each of its values is appended instead of the table itself.
// - Values from table arguments are appended in the order returned by TableValues.
// - Empty table arguments contribute no values.
func TableAppendAll(t *rt.Table, v ...rt.Value) {
    for _, val := range v {
        _t, _ok := val.TryTable()
        if _ok {
            for _, v2 := range TableValues(_t) {
                TableAppend(t, v2)
            }
            continue
        }
        TableAppend(t, val)
    }
}

// TableChunk splits the array values of t into chunks of the specified size.
//
// - Only values from the array part of t are included.
// - Values are grouped in the order returned by TableArrayValues.
// - Each chunk contains at most size values.
func TableChunk(t *rt.Table, size int) []rt.Table {
    if t.Len() == 0 {
        return make([]rt.Table, 0)
    }
    chunks := make([]rt.Table, 0)
    var chunk *rt.Table = rt.NewTable()
    for i, v := range TableArrayValues(t) {
        if i%size == 0 {
            chunks = append(chunks, *chunk)
            chunk = rt.NewTable()
        }
        TableAppend(chunk, v)
    }
    chunks = append(chunks, *chunk)
    return chunks
}
