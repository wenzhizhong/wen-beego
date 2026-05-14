import { reactive } from "vue";
import type { FormRules } from "element-plus";

/** 自定义表单规则校验 */
export const formRules = reactive(<FormRules>{
{{range .Columns}}{{if and .Required (not .SkipForm)}}  {{.Name}}: [{ required: true, message: "{{if .Comment}}{{.Comment}}{{else}}{{.Name}}{{end}}为必填项", trigger: "blur" }],
{{end}}{{end}}
});
