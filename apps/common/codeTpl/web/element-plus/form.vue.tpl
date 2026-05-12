<template>
  <el-dialog
    v-model="visible"
    :title="isEdit ? '编辑' : '新增'"
    :close-on-click-modal="false"
    width="600px"
    @close="handleClose"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
      {{range .Columns}}<el-form-item label="{{.Comment}}" prop="{{.Name}}">
        <el-input v-model="form.{{.Name}}" placeholder="请输入{{.Comment}}" />
      </el-form-item>
      {{end}}    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, nextTick } from "vue";
import { get{{.ModelName}}Detail, add{{.ModelName}}, edit{{.ModelName}} } from "./api";
import { message } from "@/utils/message";
import type { FormInstance } from "element-plus";

const emit = defineEmits(["refresh"]);
const visible = ref(false);
const isEdit = ref(false);
const submitLoading = ref(false);
const formRef = ref<FormInstance>();
const editId = ref("");

const form = reactive({
  {{range .Columns}}  {{.Name}}: "" as string | number,
  {{end}}
});

const rules = {
  {{range .Columns}}{{if .Required}}  {{.Name}}: [{ required: true, message: "请输入{{.Comment}}", trigger: "blur" }],
  {{end}}{{end}}
};

function open(id?: string) {
  visible.value = true;
  nextTick(async () => {
    if (id) {
      isEdit.value = true;
      editId.value = id;
      const res = await get{{.ModelName}}Detail({ id });
      const data = res.data ?? {};
      Object.assign(form, data);
      {{range .Columns}}{{if eq .Type "bool"}}  form.{{.Name}} = data.{{.Name}} ? true : false;
      {{end}}{{end}}
    }
  });
}

function handleClose() {
  formRef.value?.resetFields();
  isEdit.value = false;
  editId.value = "";
  visible.value = false;
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;
  submitLoading.value = true;
  try {
    if (isEdit.value) {
      await edit{{.ModelName}}({ id: editId.value, ...form });
    } else {
      await add{{.ModelName}}(form);
    }
    message.success(isEdit.value ? "修改成功" : "新增成功");
    emit("refresh");
    handleClose();
  } finally {
    submitLoading.value = false;
  }
}

defineExpose({ open });
</script>
