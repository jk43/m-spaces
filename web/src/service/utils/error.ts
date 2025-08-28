import { FormErrors, ErrorResponse, ErrorResponses } from '../../../src/types';
import axios from 'axios';

export function resetServerErrors(obj: FormErrors) {
  for (const key in obj.inputErrors) {
    obj.inputErrors[key].error = false;
    obj.inputErrors[key].message = '';
  }
  obj.formErrors = [];
}

export function handleServerErrors(err: any, obj: FormErrors) {
  if (!axios.isAxiosError(err)) {
    return;
  }
  const data = err.response!.data.data;
  //single error
  if ('code' in data) {
    const axioErr = data as ErrorResponse;
    obj.formErrors.push({
      error: true,
      message: axioErr.error,
    });
    return;
  }
  const axioErr = data as ErrorResponses;
  for (const e of axioErr) {
    if (obj.inputErrors[e.field]) {
      obj.inputErrors[e.field].error = true;
      obj.inputErrors[e.field].message = e.error;
    } else {
      obj.formErrors.push({
        error: true,
        message: e.error,
      });
    }
  }
}
