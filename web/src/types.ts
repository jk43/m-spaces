export interface ErrorResponse {
  code: number;
  error: string;
  field: string;
  id: string;
  structField: string;
}

export type ErrorResponses = ErrorResponse | ErrorResponse[];

export type ResponseResult = 'success' | 'error' | 'verification_required';

export interface Response<T> {
  result: ResponseResult;
  data: T;
}

export interface LoginSuccessResponse {
  info: {
    email: string;
    firstName: string;
    lastName: string;
    profileImage: string;
  };
  metadata: { [key: string]: StringAnyType };
  accessToken: string;
  refreshToken: string;
  mfa: boolean | null;
  mfaToken: string | null;
}

export type LoginResponse =
  | Response<ErrorResponse>
  | Response<LoginSuccessResponse>;

export type FormError = {
  error: boolean;
  message: string;
};

// export type FormErrors = FormErrors[];

export interface FormErrors {
  inputErrors: {
    [name: string]: FormError;
  };
  formErrors: FormError[];
  //
}

export type OrganizationConfigElement = {
  name: string;
  value: any;
};

export type OrganizationItem = {
  label: string;
  icon: string;
  to: string;
  order: number;
  description: string;
  elementId: string;
  elementClass: string;
  when: any;
  options: StringAnyType;
  subMenu: OrganizationItem[];
};

export interface OrganizationItems {
  topMenu: OrganizationItem[];
  userDropdownMenu: OrganizationItem[];
  adminLeftNavMenu: OrganizationItem[];
  userLeftNavMenu: OrganizationItem[];
}
//
export interface FormInput {
  key: string;
  name: string;
  type: string;
  description: string;
  rules: any;
  options: [LabelValuePair];
  defaultValue: any;
  editable: boolean; // If the value is false, the frontend will just output value.
  shareable: boolean; // If the value is true, value will be share with front-end(vue store/JWT token)
  order: number;
  show: boolean;
}

export interface LabelValuePair {
  label: string;
  value: string;
}

// export interface FormInputShowOptionFor {

// }

export interface OrganizationSettingResponse {
  active: boolean;
  host: string;
  cdnAddr: string;
  name: string;
  settings: {
    forms: FormInput[];
  };
}

export interface KeyValueData {
  key: string;
  value: string | number | boolean | null;
}

export type StringAnyType = { [key: string]: NameValueData };

export interface KeyValueNameData extends KeyValueData {
  name: string;
}

export interface NameValueData {
  name: string;
  value: any;
}

export type UserSetting = KeyValueData;

export interface User {
  email: string;
  firstName: string;
  lastName: string;
  metadata: { [key: string]: StringAnyType };
}

export interface FormBuilderSetting {
  type: string;
  settings: any;
  data: StringAnyType;
  error: FormError;
}

// every payload to server must be wrapped in this
export interface Payload<T> {
  data: T;
}

export interface DeletePayload<T> {
  data: {
    data: T;
  };
}

export interface FormDataRequst {
  file: File;
  service: string;
  serviceCtx: string;
}

export type ApiMethod = 'get' | 'post' | 'put' | 'delete';

export interface VerificationInstruction {
  verificationName: string;
  URL: string;
  keyName: string;
  method: string;
  message: string;
  resend: {
    payload: any;
    method: ApiMethod;
    URL: string;
  };
}

export type DialogCloseType = '' | 'cancel' | 'ok';

// export interface LeftMenuItem {
//   icon: string;
//   label: string;
//   separator: boolean;
//   to: string;
// }

// Type definitions
export interface FileItem {
  id: string;
  name: string;
  size: number;
  contentType: string;
  s3Path: string;
  created_at: string;
  loading?: boolean;
  thumbnail?: string | null;
}
