/* eslint-disable @typescript-eslint/no-explicit-any */
import request from "@/lib/api";

export async function getInternalDict(type: string) {
  if (type === 'file_check_task_status') {
    return Promise.resolve({
      code: 0,
      data: [
        { value: 0, label: '待处理' },
        { value: 1, label: '处理中' },
        { value: 2, label: '处理完成' },
        { value: 3, label: '处理失败' },
        { value: 4, label: '文件不存在' },
        { value: 5, label: '文件已存在' },
      ],
      msg: '成功',
    });
  }
  const params = {
    dictType: type,
  };
  return request({
    url: "/dict/internal/type",
    method: "get",
    params: params,
  });
}
