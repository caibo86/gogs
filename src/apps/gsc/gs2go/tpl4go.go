// -------------------------------------------
// @file      : tpl4go.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/19 下午4:11
// -------------------------------------------

package gs2go

var tpl4go = `
{{/**************************************************************************/}}
{{define "imports"}}
import(
    {{range $index, $element := .}} {{$index}} "{{$element.Name}}"
    {{end}}
)
{{end}}
{{/**************************************************************************/}}
{{define "error"}}{{$Enum := symbol .Name}}
// {{$Enum}} 类型定义 gsc自动生成
type {{$Enum}} {{enumType .}}
// 枚举 {{$Enum}} 常量 gsc自动生成
const (
{{range .Values}}    {{$Enum}}{{symbol .Name}} {{$Enum}} = {{.Value}}
    {{end}}
)

// Write{{$Enum}} 将枚举写到输出流
func Write{{$Enum}}(writer yfnet.Writer, val {{$Enum}}) error{
    return {{enumWrite .}}
}

// WriteTag{{$Enum}} 将枚举写到输出流带标签
func WriteTag{{$Enum}}(writer yfnet.Writer, val {{$Enum}}) error{
    yfnet.WriteTag(writer,yfnet.Enum)
    return {{enumWrite .}}
}

// Read{{$Enum}} 从输入流读取枚举
func Read{{$Enum}}(reader yfnet.Reader)({{$Enum}}, error){
    val, err := {{enumRead .}}
    return {{$Enum}}(val),err
}

// String 实现 Stringer 接口
func (val {{$Enum}}) String() string {
    switch val {
        {{range .Values}}
        case {{.Value}}:
            return "{{$Enum}}{{symbol .Name}}"
        {{end}}
    }
    return fmt.Sprintf("Unknown(%d)",val)
}

// Error 实现 Error 接口
func (val {{$Enum}}) Error() string {
    switch val {
        {{range .Values}}
        case {{.Value}}:
            return "{{$Enum}}{{symbol .Name}}"
        {{end}}
    }
    return fmt.Sprintf("Unknown(%d)",val)
}

{{end}}


{{/**************************************************************************/}}
{{define "enum"}}{{$Enum := symbol .Name}}
// {{$Enum}} 类型 gsc自动生成
type {{$Enum}} {{enumType .}}
// 枚举 {{$Enum}} 常量 gsc自动生成
const (
{{range .Values}}    {{$Enum}}{{symbol .Name}} {{$Enum}} = {{.Value}}
{{end}})

// Write{{$Enum}} 将枚举写到输出流
func Write{{$Enum}}(writer yfnet.Writer, val {{$Enum}}) error{
    return {{enumWrite .}}
}

// WriteTag{{$Enum}} 将枚举写到输出流带标签
func WriteTag{{$Enum}}(writer yfnet.Writer, val {{$Enum}}) error{
    yfnet.WriteTag(writer,yfnet.Enum)
    return {{enumWrite .}}
}

// Read{{$Enum}} 从输入流读取枚举
func Read{{$Enum}}(reader yfnet.Reader)({{$Enum}}, error){
    val, err := {{enumRead .}}
    return {{$Enum}}(val),err
}

// String 实现 Stringer 接口
func (val {{$Enum}}) String() string {
    switch val {  {{range .Values}}
    case {{.Value}}:
        return "{{$Enum}}{{symbol .Name}}"{{end}}
    }
    return fmt.Sprintf("Unknown(%d)",val)
}
{{end}}
{{/**************************************************************************/}}
{{define "readMap"}}func(reader yfnet.Reader)({{typeName .}},error) {
    length, err1 := yfnet.ReadUint16(reader)
    if err1 != nil {
        return nil, err1
    }
    buff := make({{typeName .}})
    for i := uint16(0); i < length; i++ {
        key, err1 := {{readType .Key}}(reader)
        if err1 != nil {
            return buff, err1
        }
        value, err1 := {{readType .Value}}(reader)
        if err1 != nil {
            return buff, err1
        }
        buff[key] = value
    }
    return buff,nil
}{{end}}

{{define "readList"}}func(reader yfnet.Reader)({{typeName .}},error) {
    length, err1 := yfnet.ReadUint16(reader)
    if err1 != nil {
        return nil,err1
    }
    buff := make({{typeName .}},length)
    for i := uint16(0); i < length; i ++ {
        buff[i] ,err1 = {{readType .Element}}(reader)
        if err1 != nil {
            return buff,err1
        }
    }
    return buff,nil
}{{end}}

{{define "readByteList"}}func(reader yfnet.Reader)({{typeName .}},error) {
    length, err1 := yfnet.ReadUint16(reader)
    if err1 != nil {
        return nil,err1
    }
    buff := make({{typeName .}},length)
    err1 = yfnet.ReadBytes(reader,buff)
    return buff,err1
}{{end}}

{{define "readArray"}}func(reader yfnet.Reader)({{typeName .}},error) {
    var buff {{typeName .}}
    if err != nil {
        return buff,err
    }
    for i := uint16(0); i < {{.Length}}; i ++ {
        buff[i] ,err = {{readType .Element}}(reader)
        if err != nil {
            return buff,err
        }
    }
    return buff,nil
}{{end}}

{{define "readByteArray"}}func(reader yfnet.Reader)({{typeName .}},error) {
    var buff {{typeName .}}
    if err != nil {
        return buff,err
    }
    err = yfnet.ReadBytes(reader,buff[:])
    return buff,err
}{{end}}

{{define "writeMap"}}func(writer yfnet.Writer,val {{typeName .}})(error) {
    yfnet.WriteUint16(writer,uint16(len(val)))
    for k, v := range val {
        err1 := {{writeType .Key}}(writer, k)
        if err1 != nil {
            return err1
        }
        err1 = {{writeType .Value}}(writer, v)
        if err1 != nil {
            return err1
        }
    }
    return nil
}{{end}}

{{define "writeList"}}func(writer yfnet.Writer,val {{typeName .}})(error) {
    yfnet.WriteUint16(writer,uint16(len(val)))
    for _, c := range val {
        err1 := {{writeType .Element}}(writer,c)
        if err1 != nil {
            return err1
        }
    }
    return nil
}{{end}}

{{define "writeByteList"}}func(writer yfnet.Writer,val {{typeName .}})(error) {
    yfnet.WriteUint16(writer,uint16(len(val)))
    err1 := yfnet.WriteBytes(writer,val)
    return err1
}{{end}}


{{define "writeArray"}}func(writer yfnet.Writer,val {{typeName .}})(error) {
    for _,c:= range val {
        err := {{writeType .Element}}(writer,c)
        if err != nil {
            return err
        }
    }
    return nil
}{{end}}

{{define "writeByteArray"}}func(writer yfnet.Writer,val {{typeName .}})(error) {
    return yfnet.WriteBytes(writer,val[:])
}{{end}}

{{define "writeTagList"}}func(writer yfnet.Writer,val {{typeName .}})(error) {
    if len(val) == 0 {
        yfnet.WriteTag(writer,yfnet.None)
        return nil
    }
    yfnet.WriteTag(writer,yfnet.List)
    yfnet.WriteUint16(writer,uint16(len(val)))
    for _,c:= range val {
        err := {{writeType .Element}}(writer,c)
        if err != nil {
            return err
        }
    }
    return nil
}{{end}}

{{define "writeTagByteList"}}func(writer yfnet.Writer,val {{typeName .}})(error) {
    if len(val) == 0 {
        yfnet.WriteTag(writer,yfnet.None)
        return nil
    }
    err := yfnet.WriteTag(writer,yfnet.List)
    if err != nil {
        return err
    }
    err = yfnet.WriteUint16(writer,uint16(len(val)))
    if err != nil {
        return err
    }
    return yfnet.WriteBytes(writer,val)
}{{end}}


{{define "writeTagArray"}}func(writer yfnet.Writer,val {{typeName .}})(error) {
    yfnet.WriteTag(writer,yfnet.Array)
    for _,c:= range val {
        err := {{writeType .Element}}(writer,c)
        if err != nil {
            return err
        }
    }
    return nil
}{{end}}

{{define "writeTagByteArray"}}func(writer yfnet.Writer,val {{typeName .}})(error) {
    yfnet.WriteTag(writer,yfnet.Array)
    return yfnet.WriteBytes(writer,val[:])
}{{end}}

{{define "arrayInit"}}func() {{typeName .}} {
    val := {{defaultVal .Element}}
    var buff {{typeName .}}
    for i := uint16(0); i < {{.Length}}; i ++ {
        buff[i] = val
    }
    return buff
}(){{end}}

{{/**************************************************************************/}}
{{define "struct"}} {{$Struct := symbol .Name}}
// {{$Struct}} gsc自动生成
type {{$Struct}} struct { {{range .Fields}}
    {{symbol .Name}} {{typeName .Type}} {{end}}
}

// New{{$Struct}} 用默认值生成一个结构 gsc自动生成
func New{{$Struct}}() *{{$Struct}} {
    return &{{$Struct}}{  {{range .Fields}}
        {{symbol .Name}}: {{defaultVal .Type}},     {{end}}
    }
}

// Read{{$Struct}} 从输入流读取一个 {{$Struct}}  gsc自动生成
func Read{{$Struct}}(reader yfnet.Reader) (target *{{$Struct}},err error) {
    target = New{{$Struct}}()   {{range .Fields}}
    target.{{symbol .Name}}, err = {{readType .Type}}(reader)
    if err != nil {
        return
    }   {{end}}
    return
}

// Write{{$Struct}} 将 {{$Struct}} 写入到输出流 gsc自动生成
func Write{{$Struct}}(writer yfnet.Writer,val *{{$Struct}}) (err error) { {{range .Fields}}
    {{writeType .Type}}(writer,val.{{symbol .Name}})
    if err != nil {
        return
    }{{end}}
    return nil
}

// WriteTag{{$Struct}} 将 {{$Struct}} 写入到输出流带标签 gsc自动生成
func WriteTag{{$Struct}}(writer yfnet.Writer,val *{{$Struct}}) (error) {
    if val == nil {
        yfnet.WriteTag(writer,yfnet.None)
        return nil
    }
    yfnet.WriteTag(writer,yfnet.Struct)
    return Write{{$Struct}}(writer,val)
}
{{end}}
{{/**************************************************************************/}}
{{define "table"}}
{{$Table := symbol .Name}}

// {{$Table}} gsc自动生成
type {{$Table}} struct {
    {{range .Fields}}
    {{symbol .Name}} {{typeName .Type}}
    {{end}}
}


// New{{$Table}} 用默认值生成一个结构 gsc自动生成
func New{{$Table}}() *{{$Table}} {
    return &{{$Table}}{
        {{range .Fields}}
        {{symbol .Name}}: {{defaultVal .Type}},
        {{end}}
    }
}

// Read{{$Table}} 从输入流读取一个 {{$Table}} gsc自动生成
func Read{{$Table}}(reader yfnet.Reader) (target *{{$Table}},err error) {

    target = New{{$Table}}()

    {{range .Fields}}
    var tag{{.ID}} yfnet.Tag
    tag{{.ID}}, err = yfnet.ReadTag(reader)
    if err != nil {
        return
    }

    if tag{{.ID}} != yfnet.None {
        if tag{{.ID}} != {{tag .Type}} {
            return target,yferrors.Newf(yfnet.ErrDecode,"unmatch tag(%d,%d) :{{pos .}}",tag{{.ID}},{{tag .Type}})
        }
        target.{{symbol .Name}},err = {{readType .Type}}(reader)
        if err != nil {
            return
        }
    }
    {{end}}

    return
}

// WriteTag{{$Table}} 将 {{$Table}} 写入到输出流带标签 gsc自动生成
func WriteTag{{$Table}}(writer yfnet.Writer,val *{{$Table}}) (error) {

    if val == nil {
        yfnet.WriteTag(writer,yfnet.None)
        return nil
    }

    yfnet.WriteTag(writer,yfnet.Table)

    return Write{{$Table}}(writer,val)
}

// Write{{$Table}} 将 {{$Table}} 写入到输出流 gsc自动生成
func Write{{$Table}}(writer yfnet.Writer,val *{{$Table}}) (err error) {

    {{range .Fields}}

    {{writeTagType .Type}}(writer,val.{{symbol .Name}})

    if err != nil {

        return err
    }
    {{end}}

    return nil
}

{{end}}

{{/**************************************************************************/}}
{{define "contract"}}   {{$Contract := symbol .Name}}
// 服务名
var (
    {{$Contract}}TypeName = "{{.Path}}"
)

// I{{$Contract}} gsc自动生成
type I{{$Contract}} interface {
{{range .Methods}}    {{symbol .Name}}{{params .Params}}{{returnParams .Return}}{{"\n"}}{{end}}}

//{{$Contract}}Builder gsc自动生成
type {{$Contract}}Builder struct {
    lsbuilder func(service yfdocker.Service) (I{{$Contract}},error)
}

// New{{$Contract}}Builder gsc自动生成
func New{{$Contract}}Builder(lsbuilder func(service yfdocker.Service)(I{{$Contract}},error)) yfdocker.TypeBuilder {
    return &{{$Contract}}Builder{
        lsbuilder:lsbuilder,
    }
}

// String gsc自动生成
func (builder *{{$Contract}}Builder) String() string {
    return "{{.Path}}"
}

// NewService gsc自动生成
func (builder *{{$Contract}}Builder) NewService(name string, id yfdocker.ID, context interface{}) (yfdocker.Service,error) {
    c := &{{$Contract}}Service{
        id:id,
        name:name,
        typename:builder.String(),
        context:context,
        timeout:yfconfig.Seconds(fmt.Sprintf("yfdocker.rpc_timeout.%s",name),5),
    }
    var err error
    c.I{{$Contract}},err = builder.lsbuilder(c)
    return c,err
}

// NewRemoteService gsc自动生成
func (builder *{{$Contract}}Builder) NewRemoteService(remote yfdocker.Remote, name string, lid yfdocker.ID, rid yfdocker.ID, context interface{}) yfdocker.RemoteService {
    return &{{$Contract}}RemoteService{
        name:name,
        remote:remote,
        context: context,
        lid :lid,
        rid :rid,
        typename :builder.String(),
        timeout:yfconfig.Seconds(fmt.Sprintf("yfdocker.rpc_timeout.%s",name),5),
    }
}
{{/**************************************************************************/}}
// {{$Contract}}Service gsc自动生成
type {{$Contract}}Service struct {
    I{{$Contract}}
    id yfdocker.ID
    name string
    typename string
    timeout time.Duration
    context interface{}
}

// String gsc自动生成
func (service *{{$Contract}}Service) String() string {
    return service.name
}

// ID gsc自动生成
func (service *{{$Contract}}Service) ID() yfdocker.ID {
    return service.id
}

// Type gsc自动生成
func (service *{{$Contract}}Service) Type() string {
    return service.typename
}

// Context gsc自动生成
func (service *{{$Contract}}Service) Context() interface{} {
    return service.context
}

// Call gsc自动生成
func (service *{{$Contract}}Service) Call(call *yfnet.Call) (callReturn *yfnet.Return, err error) {
    defer func(){
        if e := recover(); e != nil {
            err = yferrors.New(e.(error))
        }
    }()
    switch call.Method {  {{range .Methods}}
    case {{.ID}}:{{$Name := symbol .Name}}
        if len(call.Params) != {{.InputParams}} {
            err = yferrors.Newf(yfdocker.ErrRPC,"{{$Contract}}::{{$Name}} expect {{.InputParams}} params but got :%d",len(call.Params))
            return
        }

{{range .Params}}        var arg{{.ID}} {{typeName .Type}}
        arg{{.ID}}, err = {{readType .Type}}(bytes.NewBuffer(call.Params[{{.ID}}].Content))
        if err != nil {
            err = yferrors.Newf(err,"read {{$Contract}}::{{$Name}} arg{{.ID}} err")
            return
        }{{"\n"}}{{end}}
{{range .Return}}        var ret{{.ID}} {{typeName .Type}}{{"\n"}}{{end}}
        {{returnargs .}} service.I{{$Contract}}.{{$Name}}{{callargs .Params}}
        if err != nil {
            return
        }
        {{if .Return}}
        callReturn = &yfnet.Return{
            ID : call.ID,
            Service:call.Service,
        }
        {{range .Return}}
        var buff{{.ID}} bytes.Buffer
        err = {{writeType .Type}}(&buff{{.ID}},ret{{.ID}})
        if err != nil {
            return
        }
        callReturn.Params = append(callReturn.Params,&yfnet.Param{Content:buff{{.ID}}.Bytes()})
        {{end}}{{end}}
        return{{end}}
    }
    err = yferrors.Newf(yfdocker.ErrRPC,"unknown {{$Contract}}#%d method",call.Method)
    return
}
{{range .Methods}} {{$Name := symbol .Name}}
// {{$Name}} gsc自动生成
func (service *{{$Contract}}Service){{$Name}}{{params .Params}}{{returnParams .Return}}{
    call := &yfnet.Call{
        Service:uint16(service.id),
        Method:{{.ID}},
    }
{{range .Params}}    var param{{.ID}} bytes.Buffer
    err = {{writeType .Type}}(&param{{.ID}},arg{{.ID}})
    if err != nil {
        return
    }
    call.Params = append(call.Params,&yfnet.Param{Content:param{{.ID}}.Bytes()})
    {{end}}
    {{if .Return}}
    future := make(chan *yfnet.Return,1)
    go func(){
        var callReturn *yfnet.Return
        callReturn,err = service.Call(call)
        if err == nil {
            future <- callReturn
        }
    }()
    select {
        case callreturn := <- future:
            if len(callreturn.Params) != {{.ReturnParams}} {
                err = yferrors.Newf(yfdocker.ErrRPC,"{{$Contract}}#{{$Name}} expect {{.ReturnParams}} return params but got :%d",len(callreturn.Params))
                return
            }
            {{range .Return}}
            ret{{.ID}},err = {{readType .Type}}(bytes.NewBuffer(callreturn.Params[{{.ID}}].Content))
            if err != nil {
                err = yferrors.Newf(err,"read {{$Contract}}#{{$Name}} return{{.ID}}")
                return
            }
            {{end}}
        case <- time.After(service.timeout):
            err = yfdocker.ErrTimeout
            return
    }
    {{else}}
    go func(){ service.Call(call) }()
    {{end}}
    return
}
{{end}}

{{/**************************************************************************/}}
// {{$Contract}}RemoteService gsc自动生成
type {{$Contract}}RemoteService struct {
    remote yfdocker.Remote
    rid yfdocker.ID
    lid yfdocker.ID
    name string
    typename string
    context interface{}
    timeout time.Duration
}


// String gsc自动生成
func (service *{{$Contract}}RemoteService) String() string {
    return service.name
}

// ID gsc自动生成
func (service *{{$Contract}}RemoteService) ID() yfdocker.ID {
    return service.lid
}

// RemoteID gsc自动生成
func (service *{{$Contract}}RemoteService) RemoteID() yfdocker.ID {
    return service.rid
}

// Remote gsc自动生成
func (service *{{$Contract}}RemoteService) Remote() yfdocker.Remote {
    return service.remote
}


// Type gsc自动生成
func (service *{{$Contract}}RemoteService) Type() string {
    return service.typename
}

// Context gsc自动生成
func (service *{{$Contract}}RemoteService) Context() interface{} {
    return service.context
}

// Call gsc自动生成
func (service *{{$Contract}}RemoteService) Call(call *yfnet.Call) (callReturn *yfnet.Return, err error) {

    defer func(){
        if e := recover(); e != nil {
            err = yferrors.New(e.(error))
        }
    }()

    switch call.Method {
    {{range .Methods}}
    {{$Name := .Name}}
    case {{.ID}}:
        {{if .Return}}
        var future yfdocker.Future
        future,err = service.remote.Wait(service,call,service.timeout)
        if err != nil {
            err = yferrors.Newf(err,"call {{$Contract}}#{{$Name}} error")
            return
        }
        result := <-future
        if result.Timeout {
            err = yfdocker.ErrTimeout
            return
        }

        callReturn =result.CallReturn

        if len(callReturn.Params) != {{.ReturnParams}} {
            err = yferrors.Newf(yfdocker.ErrRPC,"{{$Contract}}#{{$Name}} expect {{.ReturnParams}} return params but got :%d",len(callReturn.Params))
            return
        }

        return
        {{else}}
        err = service.remote.Post(service,call)
        if err != nil {
            err = yferrors.Newf(err,"post {{$Contract}}#{{$Name}} error")
            return
        }
        return
        {{end}}
    {{end}}
    }

    err = yferrors.Newf(yfdocker.ErrRPC,"unknown {{$Contract}}#%d method",call.Method)

    return
}

{{range .Methods}}

{{$Name := symbol .Name}}

// {{$Name}} gsc自动生成
func (service *{{$Contract}}RemoteService){{$Name}}{{params .Params}}{{returnParams .Return}}{
    call := &yfnet.Call{
        Service:uint16(service.rid),
        Method:{{.ID}},
    }

    {{range .Params}}
    var param{{.ID}} bytes.Buffer
    err = {{writeType .Type}}(&param{{.ID}},arg{{.ID}})
    if err != nil {
        return
    }
    call.Params = append(call.Params,&yfnet.Param{Content:param{{.ID}}.Bytes()})
    {{end}}

    {{if .Return}}
    var future yfdocker.Future
    future,err = service.remote.Wait(service,call,service.timeout)
    if err != nil {
        err = yferrors.Newf(err,"call {{$Contract}}#{{$Name}} error")
        return
    }
    result := <-future
    if result.Timeout {
        err = yfdocker.ErrTimeout
        return
    }

    callreturn :=result.CallReturn

    if len(callreturn.Params) != {{.ReturnParams}} {
        err = yferrors.Newf(yfdocker.ErrRPC,"{{$Contract}}#{{$Name}} expect {{.ReturnParams}} return params but got :%d",len(callreturn.Params))
        return
    }

    {{range .Return}}


    ret{{.ID}},err = {{readType .Type}}(bytes.NewBuffer(callreturn.Params[{{.ID}}].Content))
    if err != nil {
        err = yferrors.Newf(err,"read {{$Contract}}#{{$Name}} return{{.ID}}")
        return
    }
    {{end}}
    {{else}}
    err = service.remote.Post(service,call)
    if err != nil {
        err = yferrors.Newf(err,"post {{$Contract}}#{{$Name}} error")
        return
    }
    {{end}}

    return

}
{{end}}

{{end}}
`
