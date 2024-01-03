// -------------------------------------------
// @file      : tpl4go.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/19 下午4:11
// -------------------------------------------

package main

var tpl4go = `
{{/**************************************************************************/}}

{{define "imports"}}
import(
    {{range $index, $element := .}} {{$index}} "{{$element.Name}}"
    {{end}}
)
{{end}}


{{/**************************************************************************/}}

{{define "script"}}


{{end}}

{{/**************************************************************************/}}

{{define "error"}}
{{$Enum := symbol .Name}}
// /////////////////////////////////////////////////////////////////////////////////////////////////////////

// {{$Enum}} is an autogenerated enum {{printComments .}} 
type {{$Enum}} int32

const (
{{range .SortedValues}}    {{$Enum}}{{symbol .Name}} {{$Enum}} = {{.Value}} {{printCommentsToLine .}}
{{end}} )

// String is an autogenerated method, implementing fmt.Stringer
func (val {{$Enum}}) String() string {
    switch val {  {{range .SortedValues}}
    case {{.Value}}:
        return "{{$Enum}}{{symbol .Name}}"{{end}}
    }
    return fmt.Sprintf("Unknown{{$Enum}}(%d)",val)
}

// Error is an autogenerated method, implementing error
func (val {{$Enum}}) Error() string {
    switch val {  {{range .SortedValues}}
    case {{.Value}}:
        return "{{$Enum}}{{symbol .Name}}"{{end}}
    }
    return fmt.Sprintf("Unknown{{$Enum}}(%d)",val)
}

// Unmarshal{{$Enum}} is an autogenerated function, reading the enum from a byte slice
func Unmarshal{{$Enum}}(data []byte) ({{$Enum}}, error) {
	if len(data) != 4 {
		return 0, gserrors.Newf("unmarshal {{$Enum}}, data length is not 4")
	}
	i := int32(data[0]) | int32(data[1])<<8 | int32(data[2])<<16 | int32(data[3])<<24
	return {{$Enum}}(i), nil
}



// Marshal{{$Enum}} is an autogenerated function, writing the enum to a byte slice
func Marshal{{$Enum}}(v {{$Enum}}) []byte {
	data := make([]byte, 4)
	data[0] = byte(v)
	data[1] = byte(v >> 8)
	data[2] = byte(v >> 16)
	data[3] = byte(v >> 24)
	return data
}


{{end}}

{{/**************************************************************************/}}


{{define "enum"}}
{{$Enum := symbol .Name}}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////

// {{$Enum}} is an autogenerated enum {{printComments .}} 
type {{$Enum}} int32

const (
{{range .SortedValues}}    {{$Enum}}{{symbol .Name}} {{$Enum}} = {{.Value}}
{{end}})

// String is an autogenerated method, implementing fmt.Stringer
func (val {{$Enum}}) String() string {
    switch val {  {{range .SortedValues}}
    case {{.Value}}:
        return "{{$Enum}}{{symbol .Name}}"{{end}}
    }
    return fmt.Sprintf("Unknown{{$Enum}}(%d)", val)
}

// Unmarshal{{$Enum}} is an autogenerated function, reading the enum from a byte slice
func Unmarshal{{$Enum}}(data []byte) ({{$Enum}}, error) {
	if len(data) != 4 {
		return 0, gserrors.Newf("unmarshal {{$Enum}}, data length is not 4")
	}
	i := int32(data[0]) | int32(data[1])<<8 | int32(data[2])<<16 | int32(data[3])<<24
	return {{$Enum}}(i), nil
}

// Marshal{{$Enum}} is an autogenerated function, writing the enum to a byte slice
func Marshal{{$Enum}}(v {{$Enum}}) []byte {
	data := make([]byte, 4)
	data[0] = byte(v)
	data[1] = byte(v >> 8)
	data[2] = byte(v >> 16)
	data[3] = byte(v >> 24)
	return data
}

{{end}}

{{/**************************************************************************/}}


{{define "arrayInit"}}func() {{typeName .}} {
    val := {{defaultVal .Element}}
    var buff {{typeName .}}
    for i := uint16(0); i < {{.Length}}; i ++ {
        buff[i] = val
    }
    return buff
}(){{end}}

{{/**************************************************************************/}}


{{define "struct"}}

{{$Struct := symbol .Name}}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////

{{$Receiver := lowerFirst $Struct}}


// {{$Struct}} is an autogenerated struct {{printComments .}} 
type {{$Struct}} struct { {{range .Fields}}
    {{symbol .Name}} {{typeName .Type}} {{printCommentsToLine .}} {{end}}
}

// New{{$Struct}} is an autogenerated constructor, creating a new {{$Struct}}
func New{{$Struct}}() *{{$Struct}} {
    return &{{$Struct}}{  {{range .Fields}}
        {{symbol .Name}}: {{defaultVal .Type}},     {{end}}
    }
}

// Size is an autogenerated function, returning the size of the struct
func (m *{{$Struct}})Size() int {
	if m == nil {
		return 0
	}
	n := 1
	var l int 
	_ = l
	{{range .Fields}}// {{.Name}} {{typeName .Type}}
	{{calTypeSize .}}
	{{end}}return n
}

// Marshal is an autogenerated function, marshalling the struct to a byte slice
func (m *{{$Struct}})Marshal() []byte {
	size := m.Size()
	data := make([]byte, size)
	m.MarshalToSizedBuffer(data[:size])
	return data
}

// MarshalTo is an autogenerated function, marshalling the struct to a byte slice
func (m *{{$Struct}})MarshalTo(data []byte) {
	size := m.Size()
	m.MarshalToSizedBuffer(data[:size])
	return
}

// MarshalToSizedBuffer is an autogenerated function, marshalling the struct to a byte slice
func (m *{{$Struct}})MarshalToSizedBuffer(data []byte) int {
	// flag
	data[0] = 0xFE
	i := 1
	{{range .Fields}}// {{.Name}} {{typeName .Type}}
	{{writeType .}}
	{{end}}
	return i
}

// Unmarshal is an autogenerated function, unmarshalling the struct from a byte slice
func (m *{{$Struct}})Unmarshal(data []byte) (err error) {
	defer func(){
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	// flag
	_ = data[0]
	l := len(data)
	i := 1
	for i < l {
		var fieldID uint16
	    i, fieldID = gsnet.ReadFieldID(data, i)
		switch fieldID {
		{{range .Fields}}case {{.ID}}:
			{{readType .}}
		{{end}}}
	}
	return
}

// CopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. m must be non-nil.
func (m *{{$Struct}})CopyInto(out *{{$Struct}}) {
	*out = *m
	{{range .Fields}}{{copyType .}}
	{{end}}return
}

// Copy is an autogenerated deepcopy function, copying the receiver, creating a new {{$Struct}}.
func (m *{{$Struct}})Copy() *{{$Struct}} {
	if m == nil {
		return nil
	}
	out := new({{$Struct}})
	m.CopyInto(out)
	return out
}

// Marshal{{$Struct}} is an autogenerated function, marshalling the struct to a byte slice
func Marshal{{$Struct}}(m *{{$Struct}}) []byte {
	return m.Marshal()
}


// Unmarshal{{$Struct}} is an autogenerated function, unmarshalling the struct from a byte slice
func Unmarshal{{$Struct}}(data []byte) (*{{$Struct}}, error) {
	m := New{{$Struct}}()
	err := m.Unmarshal(data)
	if err != nil {
		return nil, gserrors.Newf("unmarshal {{$Struct}} err:%s", err.Error())
	}
	return m, nil
}

// Read{{$Struct}} is an autogenerated function, read a {{$Struct}} from io.Reader
func Read{{$Struct}}(reader io.Reader) (*{{$Struct}}, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}
	size, err := gsnet.UnmarshalUint32(buf)
	if err != nil {
		return nil, err
	}
	data := make([]byte, size)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}
	return Unmarshal{{$Struct}}(data)	
}

// Write{{$Struct}} is an autogenerated function, write a {{$Struct}} to io.Writer
func Write{{$Struct}}(writer io.Writer, m *{{$Struct}}) error {
	size := m.Size()
	data := gsnet.MarshalUint32(uint32(size))
	_, err := writer.Write(data)
	if err != nil {
		return err
	}
	data = Marshal{{$Struct}}(m)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	return nil
}


{{end}}

{{/**************************************************************************/}}



{{/**************************************************************************/}}


{{define "table"}}
{{$Table := symbol .Name}}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////

// {{$Table}} is an autogenerated struct
type {{$Table}} struct { 
    {{range .Fields}} {{symbol .Name}} {{typeName .Type}}
    {{end}}
}


// New{{$Table}} is an autogenerated constructor, creating a new {{$Table}}
func New{{$Table}}() *{{$Table}} {
    return &{{$Table}}{
        {{range .Fields}} {{symbol .Name}}: {{defaultVal .Type}},
        {{end}} }
}

{{end}}

{{/**************************************************************************/}}


{{define "service"}}
{{$Service := symbol .Name}}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////

// {{$Service}}TypeName a unique name of the service
const {{$Service}}TypeName = "{{.Path}}"

// I{{$Service}} is an autogenerated interface
type I{{$Service}} interface {
{{range .Methods}}    {{symbol .Name}}{{params .Params}}{{returnParams .Return}}{{"\n"}}{{end}}}

//{{$Service}}Builder service builder used for building {{$Service}} service
type {{$Service}}Builder struct {
    localServiceBuilder func(service gsdock.IService) (I{{$Service}},error)
}

// New{{$Service}}Builder creating a new {{$Service}}Builder
func New{{$Service}}Builder(localServiceBuilder func(service gsdock.IService)(I{{$Service}}, error)) gsdock.ITypeBuilder {
    return &{{$Service}}Builder{
        localServiceBuilder:localServiceBuilder,
    }
}

// String implementing fmt.Stringer
func (builder *{{$Service}}Builder) String() string {
    return "{{.Path}}"
}

// NewService creating a new {{$Service}} Service
func (builder *{{$Service}}Builder) NewService(
	name string, id gsdock.ID, context interface{}) (gsdock.IService, error) {
    c := &{{$Service}}Service{
        id: id,
        name: name,
        typename: builder.String(),
        context: context,
        timeout: config.RPCTimeout(),
    }
    var err error
    c.I{{$Service}}, err = builder.localServiceBuilder(c)
    return c, err
}

// NewRemoteService creating a new {{$Service}} RemoteService
func (builder *{{$Service}}Builder) NewRemoteService(
	remote gsdock.IRemote, name string, lid gsdock.ID, 
	rid gsdock.ID, context interface{}) gsdock.IRemoteService {
    return &{{$Service}}RemoteService{
        name: name,
        remote: remote,
        context: context,
        lid: lid,
        rid: rid,
        typename: builder.String(),
        timeout: config.RPCTimeout(),
    }
}


// {{$Service}}Service a local service inherited {{$Service}}
type {{$Service}}Service struct {
    I{{$Service}}
    id gsdock.ID
    name string
    typename string
    timeout time.Duration
    context interface{}
}

// String implementing fmt.Stringer
func (service *{{$Service}}Service) String() string {
    return service.name
}

// Name service name
func (service *{{$Service}}Service) Name() string {
	return service.name
}

// ID service id 
func (service *{{$Service}}Service) ID() gsdock.ID {
    return service.id
}

// Type service type
func (service *{{$Service}}Service) Type() string {
    return service.typename
}

// Context service context
func (service *{{$Service}}Service) Context() interface{} {
    return service.context
}

// Call the specified method of the service
func (service *{{$Service}}Service) Call(call *gsnet.Call) (callReturn *gsnet.Return, err error) {
    defer func(){
        if e := recover(); e != nil {
            err = gserrors.New(e.(error).Error())
        }
		if err != nil {
			log.Error("{{$Service}}Service#Call err: %s", err.Error())
		}
    }()
    switch call.MethodID { {{range .Methods}}
   		case {{.ID}}:  {{$Name := symbol .Name}}
		// {{$Name}}
        if len(call.Params) != {{.InputParams}} {
            err = gserrors.NewfWith(gsdock.ErrRPC, "{{$Service}}::{{$Name}} expect {{.InputParams}} params but got :%d", len(call.Params))
            return
        }
		{{range .Params}} var param{{.ID}} {{typeName .Type}}
		param{{.ID}}, err = {{unmarshalType .Type}}(call.Params[{{.ID}}])
		if err != nil {
			return 
		}
		{{end}} {{range .Return}} var ret{{.ID}} {{typeName .Type}}
		{{end}} {{returnArgs .}} service.I{{$Service}}.{{$Name}}{{callParams .Params}}
        if err != nil {
            return
        }
        {{if .Return}} callReturn = &gsnet.Return{
            ID: call.ID,
            ServiceID: call.ServiceID,
        }
        {{range .Return}} data{{.ID}} := {{marshalType .Type}}(ret{{.ID}})
        callReturn.Params = append(callReturn.Params, data{{.ID}})
        {{end}}{{end}}return{{end}}
	}
    err = gserrors.NewfWith(gsdock.ErrRPC, "unknown {{$Service}}Service#%d method", call.MethodID)
    return
}

{{range .Methods}} {{$Name := symbol .Name}}
// {{$Name}} method of service {{$Service}}
func (service *{{$Service}}Service){{$Name}}{{params .Params}}{{returnParams .Return}}{
    call := &gsnet.Call{
        ServiceID: uint32(service.id),
        MethodID: {{.ID}},
    }
{{range .Params}} param{{.ID}} := {{marshalType .Type}}(arg{{.ID}})
    call.Params = append(call.Params, param{{.ID}})
    {{end}}
	{{if .Return}} future := make(chan *gsnet.Return,1)
    go func(){
        callReturn, err1 := service.Call(call)
        if err1 == nil {
            future <- callReturn
        }
    }()
    select {
        case callReturn := <- future:
            if len(callReturn.Params) != {{.ReturnParams}} {
                err = gserrors.NewfWith(gsdock.ErrRPC,"{{$Service}}Service#{{$Name}} expect {{.ReturnParams}} return params but got :%d",len(callReturn.Params))
                return
            }
            {{range .Return}} ret{{.ID}}, err = {{unmarshalType .Type}}(callReturn.Params[{{.ID}}])
            if err != nil {
                err = gserrors.NewWith(err, "unmarshal {{$Service}}Service#{{$Name}} return{{.ID}} {{typeName .Type}} err")
                return
            }
            {{end}}case <- time.After(service.timeout):
            err = gsdock.ErrTimeout
            return
    }
    {{else}}
    go func(){ 
		_, _ = service.Call(call) 
	}()
    {{end}}return
}
{{end}}



// {{$Service}}RemoteService a remote service inherited gsdock.IRemote
type {{$Service}}RemoteService struct {
    remote gsdock.IRemote
    rid gsdock.ID // remote id
    lid gsdock.ID // local id
    name string
    typename string
    context interface{}
    timeout time.Duration
}

// String implementing fmt.Stringer
func (service *{{$Service}}RemoteService) String() string {
    return service.name
}

// Name remote service name
func (service *{{$Service}}RemoteService) Name() string {
    return service.name
}

// ID remote service local id
func (service *{{$Service}}RemoteService) ID() gsdock.ID {
    return service.lid
}

// RemoteID remote service remote id
func (service *{{$Service}}RemoteService) RemoteID() gsdock.ID {
    return service.rid
}

// Remote remote service IRemote
func (service *{{$Service}}RemoteService) Remote() gsdock.IRemote {
    return service.remote
}


// Type remote service type
func (service *{{$Service}}RemoteService) Type() string {
    return service.typename
}

// Context remote service context
func (service *{{$Service}}RemoteService) Context() interface{} {
    return service.context
}

// Call remote service call
func (service *{{$Service}}RemoteService) Call(call *gsnet.Call) (callReturn *gsnet.Return, err error) {
    defer func(){
        if e := recover(); e != nil {
            err = gserrors.New(e.(error).Error())
        }
    }()
    switch call.MethodID { 
	{{range .Methods}} {{$Name := .Name}} case {{.ID}}:
		// {{$Name}}
        {{if .Return}} var future gsdock.Future
        future, err = service.remote.Wait(service, call, service.timeout)
        if err != nil {
            err = gserrors.NewWith(err, "call {{$Service}}RemoteService#{{$Name}} err")
            return
        }
        result := <-future
        if result.Timeout {
            err = gsdock.ErrTimeout
            return
        }
        callReturn = result.CallReturn
        if len(callReturn.Params) != {{.ReturnParams}} {
            err = gserrors.NewfWith(gsdock.ErrRPC, "{{$Service}}RemoteService#{{$Name}} expect {{.ReturnParams}} return params but got :%d", len(callReturn.Params))
            return
        }
        return
        {{else}} err = service.remote.Post(service,call)
        if err != nil {
            err = gserrors.NewWith(err, "post {{$Service}}RemoteService#{{$Name}} err")
            return
        }
        return 
		{{end}}{{end}} }
    err = gserrors.NewfWith(gsdock.ErrRPC, "unknown {{$Service}}RemoteService#%d method", call.MethodID)
    return
}

{{range .Methods}}
{{$Name := symbol .Name}}
// {{$Name}} methods of remote service
func (service *{{$Service}}RemoteService){{$Name}}{{params .Params}}{{returnParams .Return}}{
    call := &gsnet.Call{
        ServiceID: uint32(service.rid),
        MethodID: {{.ID}},
    }
    {{range .Params}} param{{.ID}} := {{marshalType .Type}}(arg{{.ID}})
    call.Params = append(call.Params, param{{.ID}})
    {{end}}
    {{if .Return}} var future gsdock.Future
    future,err = service.remote.Wait(service, call, service.timeout)
    if err != nil {
        err = gserrors.NewWith(err, "call {{$Service}}RemoteService#{{$Name}} err")
        return
    }
    result := <-future
    if result.Timeout {
        err = gsdock.ErrTimeout
        return
    }
    callReturn := result.CallReturn
    if len(callReturn.Params) != {{.ReturnParams}} {
        err = gserrors.NewfWith(gsdock.ErrRPC, "{{$Service}}RemoteService#{{$Name}} expect {{.ReturnParams}} return params but got :%d", len(callReturn.Params))
        return
    }
    {{range .Return}} ret{{.ID}}, err = {{unmarshalType .Type}}(callReturn.Params[{{.ID}}])
    if err != nil {
        err = gserrors.NewWith(err, "unmarshal {{$Service}}RemoteService#{{$Name}} return{{.ID}} {{typeName .Type}} err")
        return
    }
	{{end}} {{else}} err = service.remote.Post(service, call)
    if err != nil {
        err = gserrors.NewWith(err, "post {{$Service}}RemoteService#{{$Name}} error")
        return
    }
	{{end}} return
}
{{end}}

{{end}}
`
