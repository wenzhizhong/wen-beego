import { http } from "@/utils/http";
import type { Result, ResultTable } from "@/utils/http/types";

const BASE_URL = "{{.ApiUrlPrefix}}";

export const get{{.ModelName}}List = (params?: object) => {
  return http.request<ResultTable>("get", BASE_URL + "/get", { params });
};

export const get{{.ModelName}}Detail = (data?: object) => {
  return http.request<Result>("post", BASE_URL + "/detail", { data });
};

export const add{{.ModelName}} = (data?: object) => {
  return http.request<Result>("post", BASE_URL + "/add", { data });
};

export const edit{{.ModelName}} = (data?: object) => {
  return http.request<Result>("post", BASE_URL + "/edit", { data });
};

export const del{{.ModelName}} = (data?: object) => {
  return http.request<Result>("post", BASE_URL + "/del", { data });
};
