package method

type MethodList map[string][]*Method

func (ml MethodList) Add(typeName string, m *Method) {
	if _, ok := ml[typeName]; !ok {
		ml[typeName] = make([]*Method, 0)
	}
	ml[typeName] = append(ml[typeName], m)
}

type Method struct {
	Name string
	Args *Struct
	Result *Struct
}

func NewMethod(n string, a *Struct, r *Struct) *Method {
	return &Method{n, a, r}
}

type Struct struct{
	Pack   string
	Name   string
	Prefix string
}

func NewStruct(p string, n string) *Struct{
	return &Struct{p, n, ""}
}

func (s *Struct) SetPrefix(p string) {
	s.Prefix = p
}