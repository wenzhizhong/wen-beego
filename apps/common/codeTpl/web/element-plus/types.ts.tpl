interface FormItemProps {
{{range .Columns}}{{if not .SkipForm}}  {{.Name}}: {{.TsType}};
{{end}}{{end}}}
interface FormProps {
  formInline: FormItemProps;
}

export type { FormItemProps, FormProps };
