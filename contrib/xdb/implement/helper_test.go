package implement

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/zhiyunliu/golibs/xreflect"
	"github.com/zhiyunliu/golibs/xtypes"
)

var (
	ptrInt *int = new(int)
)

func init() {
	*ptrInt = 3
}

type dbparam struct {
	A string
}

func (p dbparam) ToDbParam() map[string]any {
	return map[string]any{"a": p.A}
}

type jsonStructParam struct {
	A    string             `json:"a"`
	B    int                `json:"b"`
	C    []int              `json:"c"`
	D    *[]int             `json:"d"`
	E    []string           `json:"e"`
	Dec  xtypes.Decimal     `json:"dec"`
	Dec2 *xtypes.Decimal    `json:"dec2"`
	Map  map[string]string  `json:"map"`
	Map2 *map[string]string `json:"map2"`
	Map3 StrMap             `json:"map3"`
	Map4 *StrMap            `json:"map4"`
	Obj  dbstructParam      `json:"obj"`
	Obj2 *dbstructParam     `json:"obj2"`
	Obj3 dbstructParamStr   `json:"obj3"`
	Obj4 *dbstructParamStr  `json:"obj4"`
}

type dbstructJson struct {
	A string `json:"a"`
}

type dbstructParam struct {
	A string `db:"a"`
}

type dbstructParamStr struct {
	A string `db:"a"`
}

func (p dbstructParamStr) String() string {
	return p.A
}

type StrMap map[string]string

func (m StrMap) String() string {
	mapBytes, _ := json.Marshal(m)
	return string(mapBytes)
}

type NotImpl struct{ NotImpl any }

type Impl struct{ impl any }

func (i *Impl) Scan(v any) error {
	i.impl = v
	return nil
}

type Binary []uint8

func Test_analyzeParamFields(t *testing.T) {
	decVal := xtypes.NewDecimalFromInt(10)
	mapVal := map[string]string{"m": "m"}
	var smap StrMap = mapVal
	tests := []struct {
		name       string
		input      any
		wantParams xtypes.XMap
		wantErr    bool
	}{
		{name: "1.", input: jsonStructParam{
			A:    "1",
			B:    2,
			C:    []int{1, 2},
			D:    &[]int{1, 2},
			E:    []string{"a", "b"},
			Dec:  decVal,
			Dec2: &decVal,
			Map:  mapVal,
			Map2: &mapVal,
			Map3: smap,
			Map4: &smap,
			Obj:  dbstructParam{A: "obj"},
			Obj2: &dbstructParam{A: "obj2"},
			Obj3: dbstructParamStr{A: "obj3"},
			Obj4: &dbstructParamStr{A: "obj4"},
		}, wantParams: map[string]any{
			"a":    "1",
			"b":    int64(2),
			"c":    []int{1, 2},
			"d":    []int{1, 2},
			"e":    []string{"a", "b"},
			"dec":  "10",
			"dec2": "10",
			"map":  nil,
			"map2": nil,
			"map3": `{"m":"m"}`,
			"map4": `{"m":"m"}`,
			"obj":  nil,
			"obj2": nil,
			"obj3": "obj3",
			"obj4": "obj4",
		}, wantErr: false},

		{name: "2.", input: &jsonStructParam{
			A:    "1",
			B:    2,
			C:    []int{1, 2},
			D:    &[]int{1, 2},
			E:    []string{"a", "b"},
			Dec:  decVal,
			Dec2: &decVal,
			Map:  mapVal,
			Map2: &mapVal,
			Map3: smap,
			Map4: &smap,
			Obj:  dbstructParam{A: "obj"},
			Obj2: &dbstructParam{A: "obj2"},
			Obj3: dbstructParamStr{A: "obj3"},
			Obj4: &dbstructParamStr{A: "obj4"},
		}, wantParams: map[string]any{
			"a":    "1",
			"b":    int64(2),
			"c":    []int{1, 2},
			"d":    []int{1, 2},
			"e":    []string{"a", "b"},
			"dec":  "10",
			"dec2": "10",
			"map":  nil,
			"map2": nil,
			"map3": `{"m":"m"}`,
			"map4": `{"m":"m"}`,
			"obj":  nil,
			"obj2": nil,
			"obj3": "obj3",
			"obj4": "obj4",
		}, wantErr: false},
		{name: "3.", input: &anonymousB{
			IAPtr:          ptrInt,
			anonymousInner: &anonymousInner{Str: "str", Int: 1},
		},
			wantParams: map[string]any{
				"iaptr": int64(3),
				"str":   "str",
				"int":   int64(0),
			},
			wantErr: false,
		},
		{name: "4.", input: &anonymousB{
			IAPtr: ptrInt,
			Int:   2,
		},
			wantParams: map[string]any{
				"iaptr": int64(3),
				"int":   int64(2),
			},
			wantErr: false,
		},
		{name: "5.", input: &anonymousB{
			IAPtr:          ptrInt,
			Int:            2,
			anonymousInner: &anonymousInner{Str: "str", Int: 1},
		},
			wantParams: map[string]any{
				"iaptr": int64(3),
				"int":   int64(2),
				"str":   "str",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParams, err := analyzeParamFields(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("analyzeParamFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotParams, tt.wantParams) {
				t.Errorf("analyzeParamFields() = %v, want %v", gotParams, tt.wantParams)
			}

			for k, v := range gotParams {
				if !reflect.DeepEqual(v, tt.wantParams[k]) {
					t.Errorf("analyzeParamFields() %s = %v, want %v", k, v, tt.wantParams[k])
				}
			}
		})
	}
}

func TestResolveParams(t *testing.T) {

	tests := []struct {
		name       string
		input      any
		wantParams xtypes.XMap
		wantErr    bool
	}{
		{name: "1.", input: map[string]any{"a": 1, "b": 2}, wantParams: map[string]any{"a": 1, "b": 2}, wantErr: false},
		{name: "2.", input: xtypes.XMap{"a": 1, "b": 2}, wantParams: map[string]any{"a": 1, "b": 2}, wantErr: false},
		{name: "3.", input: xtypes.SMap{"a": "1", "b": "2"}, wantParams: map[string]any{"a": "1", "b": "2"}, wantErr: false},
		{name: "4.", input: map[string]string{"a": "1", "b": "2"}, wantParams: map[string]any{"a": "1", "b": "2"}, wantErr: false},
		{name: "5.", input: dbparam{A: "1"}, wantParams: map[string]any{"a": "1"}, wantErr: false},
		{name: "6.", input: &dbparam{A: "1"}, wantParams: map[string]any{"a": "1"}, wantErr: false},
		{name: "7.", input: dbstructParam{A: "1"}, wantParams: map[string]any{"a": "1"}, wantErr: false},
		{name: "8.", input: &dbstructParam{A: "1"}, wantParams: map[string]any{"a": "1"}, wantErr: false},
		{name: "9.", input: dbstructJson{A: "1"}, wantParams: map[string]any{"a": "1"}, wantErr: false},
		{name: "10.", input: &dbstructJson{A: "1"}, wantParams: map[string]any{"a": "1"}, wantErr: false}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParams, err := ResolveParams(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotParams, tt.wantParams) {
				t.Errorf("ResolveParams() = %v, want %v", gotParams, tt.wantParams)
			}
		})
	}
}

type val struct {
	BB         bool              `json:"bool"`
	IA         int               `json:"ia"`
	IB         int32             `json:"ib"`
	IC         int64             `json:"ic"`
	IU         uint64            `json:"iu"`
	FA         float32           `json:"fa"`
	FB         float32           `json:"fb"`
	Str        string            `json:"str"`
	Bytes      []byte            `json:"bytes"`
	IntArray   []int             `json:"ints"`
	FloatArray []float32         `json:"floats"`
	Impl       Impl              `json:"impl"`
	NotImpl    NotImpl           `json:"notimpl"`
	MapStr     map[string]string `json:"mapstr"`
	MapAny     map[string]any    `json:"mapany"`
	XMap       xtypes.XMap       `json:"xmap"`
	SMap       xtypes.SMap       `json:"smap"`
	XMaps      xtypes.XMaps      `json:"xmaps"`
	Binary     Binary            `json:"binary"`
	Any        any               `json:"any"`
	Dec        xtypes.Decimal    `json:"dec"`
	DecPtr     *xtypes.Decimal   `json:"decptr"`
	BBPtr      *bool             `json:"boolptr"`
	IAPtr      *int              `json:"iaptr"`
	IBPtr      *int32            `json:"ibptr"`
	ICPtr      *int64            `json:"icptr"`
	IUPtr      *uint64           `json:"iuptr"`
	FAPtr      *float32          `json:"faptr"`
	FBPtr      *float32          `json:"fbptr"`
	StrPtr     *string           `json:"strptr"`
	ImplPtr    *Impl             `json:"implptr"`
	Time       time.Time         `json:"time"`
	TimePtr    *time.Time        `json:"timeptr"`
	MapAnyPtr  *map[string]any   `json:"mapanyptr"`
	XMapPtr    *xtypes.XMap      `json:"xmapptr"`
}

type scannerVal struct {
	IA int   `json:"ia"`
	IB int32 `json:"ib"`
	IC int64 `json:"ic"`
}

func (v *scannerVal) StructScan(vals ...any) error {
	v.IA, _ = vals[0].(int)
	v.IB, _ = vals[1].(int32)
	v.IC, _ = vals[2].(int64)
	return nil
}

func Test_fillRowToStruct(t *testing.T) {

	var (
		testVal1 *val = &val{}
	)

	tests := []struct {
		name       string
		fields     *xreflect.StructFields
		reflectVal reflect.Value
		result     any
		vals       map[string]any
		wantErr    bool
	}{
		{name: "1.", result: testVal1, vals: map[string]any{
			"bool":      true,
			"ia":        1,
			"ib":        2,
			"ic":        3,
			"iu":        4,
			"fa":        1.1,
			"fb":        2.2,
			"str":       "str",
			"bytes":     []byte{65, 66},
			"ints":      []int{1, 2},
			"floats":    []float32{3.1, 4.2},
			"impl":      "impl",
			"notimpl":   "notimpl",
			"mapstr":    "mapstr",
			"mapany":    "mapany",
			"xmap":      `{"xmap":"1"}`,
			"smap":      `{"smap":"2"}`,
			"xmaps":     `[{"xmap":"1"},{"smap":"2"}]`,
			"binary":    []uint8{1, 2, 3},
			"any":       "any",
			"dec":       []byte("10.2"),
			"decptr":    []byte("10.3"),
			"boolptr":   true,
			"iaptr":     11,
			"ibptr":     12,
			"icptr":     13,
			"iuptr":     14,
			"faptr":     2.1,
			"fbptr":     2.2,
			"strptr":    "abc",
			"implptr":   "implptr",
			"time":      time.Now(),
			"timeptr":   time.Now(),
			"mapanyptr": "mapanyptr",
			"xmapptr":   `{"xmapptr":"1"}`,
		}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.reflectVal = reflect.ValueOf(tt.result)
			tt.fields = xreflect.CachedTypeFields(tt.reflectVal.Type())

			cols := make([]string, len(tt.fields.List))
			vals := make([]any, len(tt.fields.List))
			for i, k := range tt.fields.List {
				cols[i] = k.Name
				vals[i] = tt.vals[k.Name]
			}

			if err := scanInToStruct(tt.fields, tt.reflectVal, cols, vals); (err != nil) != tt.wantErr {
				t.Errorf("fillRowToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fillRowToStruct_Scanner(t *testing.T) {

	var (
		testVal1 *scannerVal = &scannerVal{}
	)

	tests := []struct {
		name       string
		fields     *xreflect.StructFields
		reflectVal reflect.Value
		result     any
		vals       []any
		wantErr    bool
	}{
		{name: "1.", result: testVal1, vals: []any{
			1, 2, 3,
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.reflectVal = reflect.ValueOf(tt.result)
			tt.fields = xreflect.CachedTypeFields(tt.reflectVal.Type())

			cols := make([]string, len(tt.fields.List))
			for i, k := range tt.fields.List {
				cols[i] = k.Name
			}

			if err := scanInToStruct(tt.fields, tt.reflectVal, cols, tt.vals); (err != nil) != tt.wantErr {
				t.Errorf("fillRowToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type anonymousA struct {
	anonymousInner
	AStr string `json:"astr"`
	AInt int    `json:"aint"`
}

type anonymousB struct {
	IAPtr *int `json:"iaptr"`
	Int   int  `json:"int"`
	*anonymousInner
}

type anonymousInner struct {
	Str string `json:"str"`
	Int int    `json:"int"`
}

func Test_fillRowToStruct_anonymousA(t *testing.T) {

	b := anonymousB{anonymousInner: &anonymousInner{}}

	json.Unmarshal([]byte(`{"iaptr":1,"int":2,"str":"str"}`), &b)

	var (
		testVal1 *anonymousA = &anonymousA{}
		ptrInt   *int        = new(int)
	)
	*ptrInt = 3
	tests := []struct {
		name       string
		fields     *xreflect.StructFields
		reflectVal reflect.Value
		result     any
		vals       map[string]any
		wantErr    bool
		wantVal    *anonymousA
	}{
		{name: "1.", result: testVal1,
			vals: map[string]any{
				"str":  "strval",
				"int":  2,
				"astr": "astrval",
				"aint": 12,
			},
			wantVal: &anonymousA{AStr: "astrval", AInt: 12, anonymousInner: anonymousInner{Str: "strval", Int: 2}},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.reflectVal = reflect.ValueOf(tt.result)
			tt.fields = xreflect.CachedTypeFields(tt.reflectVal.Type())

			vals := make([]any, len(tt.fields.List))
			cols := make([]string, len(tt.fields.List))
			for i, k := range tt.fields.List {
				cols[i] = k.Name
				vals[i] = tt.vals[cols[i]]
			}

			if err := scanInToStruct(tt.fields, tt.reflectVal, cols, vals); (err != nil) != tt.wantErr {
				t.Errorf("fillRowToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.result, tt.wantVal) {
				t.Errorf("fillRowToStruct() DeepEqual Got= %v, wantVal %v", tt.result, tt.wantVal)
			}
		})
	}
}

func Test_fillRowToStruct_anonymousB(t *testing.T) {
	var (
		testVal1 *anonymousB = &anonymousB{anonymousInner: &anonymousInner{}}
		ptrInt   *int        = new(int)
	)
	*ptrInt = 3

	//json.Unmarshal([]byte(`{"iaptr":3,"str":"aaa","int":1}`), testVal1)

	tests := []struct {
		name       string
		fields     *xreflect.StructFields
		reflectVal reflect.Value
		result     any
		vals       map[string]any
		wantErr    bool
		wantVal    *anonymousB
	}{
		{name: "1.", result: testVal1,
			vals: map[string]any{
				"iaptr": 3,
				"str":   "strval",
				"int":   2,
				"astr":  "astrval",
				"aint":  12,
			},
			wantVal: &anonymousB{IAPtr: ptrInt, Int: 2, anonymousInner: &anonymousInner{Str: "strval", Int: 2}},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.reflectVal = reflect.ValueOf(tt.result)
			tt.fields = xreflect.CachedTypeFields(tt.reflectVal.Type())

			vals := make([]any, len(tt.fields.List))
			cols := make([]string, len(tt.fields.List))
			i := 0
			for _, k := range tt.fields.ExactName {
				cols[i] = k.Name
				vals[i] = tt.vals[k.Name]
				i++
			}

			if err := scanInToStruct(tt.fields, tt.reflectVal, cols, vals); (err != nil) != tt.wantErr {
				t.Errorf("fillRowToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.result, tt.wantVal) {
				t.Errorf("fillRowToStruct() DeepEqual Got= %v, wantVal %v", tt.result, tt.wantVal)
			}
		})
	}
}

func Benchmark_fillRowToStruct(b *testing.B) {
	var (
		testVal1 *val = &val{}
	)

	tt := struct {
		name       string
		fields     *xreflect.StructFields
		reflectVal reflect.Value
		result     any
		vals       map[string]any
		wantErr    bool
	}{name: "1.", result: testVal1, vals: map[string]any{
		"bool":      true,
		"ia":        1,
		"ib":        2,
		"ic":        3,
		"iu":        4,
		"fa":        1.1,
		"fb":        2.2,
		"str":       "str",
		"bytes":     []byte{65, 66},
		"ints":      []int{1, 2},
		"floats":    []float32{1.0, 2.0},
		"impl":      "impl",
		"notimpl":   "notimpl",
		"mapstr":    "mapstr",
		"mapany":    "mapany",
		"xmap":      `{"xmap":"1"}`,
		"smap":      `{"smap":"2"}`,
		"xmaps":     `[{"xmap":"1"},{"smap":"2"}]`,
		"binary":    []uint8{1, 2, 3},
		"any":       "any",
		"dec":       []byte("10.2"),
		"decptr":    []byte("10.3"),
		"boolptr":   true,
		"iaptr":     11,
		"ibptr":     12,
		"icptr":     13,
		"iuptr":     14,
		"faptr":     2.1,
		"fbptr":     2.2,
		"strptr":    "abc",
		"implptr":   "implptr",
		"time":      time.Now(),
		"timeptr":   time.Now(),
		"mapanyptr": "mapanyptr",
		"xmapptr":   `{"xmapptr":"1"}`,
	}, wantErr: false}

	for i := 0; i < b.N; i++ {

		tt.reflectVal = reflect.ValueOf(tt.result)
		tt.fields = xreflect.CachedTypeFields(tt.reflectVal.Type())

		cols := make([]string, len(tt.fields.List))
		vals := make([]any, len(cols))
		for i, k := range tt.fields.List {
			cols[i] = k.Name
			vals[i] = tt.vals[k.Name]
		}

		if err := scanInToStruct(tt.fields, tt.reflectVal, cols, vals); (err != nil) != tt.wantErr {
			b.Errorf("fillRowToStruct() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
