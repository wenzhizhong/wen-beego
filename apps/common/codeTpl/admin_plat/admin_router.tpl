// ***** 路由需要手动粘贴 *****

// 1
import admin{{.MenuModule | snakeToPascal}} "WenBeego/apps/{{.AppModule}}/controllers/{{.MenuModule}}"

// 2
{{$appModelLower := replaceAll .AppModule "admin_" "" -}}
{{$appModelLower = $appModelLower | snakeToCamel -}}
{{$menuModule := .MenuModule | snakeToPascal -}}
{{$routerFuncName := printf "%s%sSlices()" $appModelLower $menuModule}}
func {{$routerFuncName}} []beego.LinkNamespace {
    return []beego.LinkNamespace{
        // 测试路由
        {{$old := printf "/%s" .AppModule -}}
        {{$new := "" -}}
        {{$pathRead := printf "%s/get" .ApiUrlPrefix -}}
        {{$pathAdd := printf "%s/add" .ApiUrlPrefix -}}
        {{$pathEdit := printf "%s/edit" .ApiUrlPrefix -}}
        {{$pathDel := printf "%s/del" .ApiUrlPrefix -}}
        {{$pathDetail := printf "%s/detail" .ApiUrlPrefix -}}

        {{$pathRead = replaceAll $pathRead $old $new  -}}
        {{$pathAdd = replaceAll $pathAdd $old $new  -}}
        {{$pathEdit = replaceAll $pathEdit $old $new  -}}
        {{$pathDel = replaceAll $pathDel $old $new  -}}
        {{$pathDetail = replaceAll $pathDetail $old $new  -}}

        // {{$menuModule}} begin
        beego.NSCtrlGet("{{$pathRead}}", (*admin{{$menuModule}}.{{.ModelName}}Controller).Get),
        beego.NSCtrlPost("{{$pathAdd}}", (*admin{{$menuModule}}.{{.ModelName}}Controller).Add),
        beego.NSCtrlPost("{{$pathEdit}}", (*admin{{$menuModule}}.{{.ModelName}}Controller).Edit),
        beego.NSCtrlPost("{{$pathDel}}", (*admin{{$menuModule}}.{{.ModelName}}Controller).Del),
        beego.NSCtrlGet("{{$pathDetail}}", (*admin{{$menuModule}}.{{.ModelName}}Controller).Detail),
        // {{$menuModule}} end
    }
}
// 3
func init() {
    // ...
    allNamespaces = append(allNamespaces, {{$routerFuncName}}...)

}