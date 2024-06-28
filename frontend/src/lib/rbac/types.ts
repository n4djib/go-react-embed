export type UserWithRoles = {
  roles: string[];
} & {
  [key: string]: any;
};

export type Role = {
  id: number;
  role: string;
};

export type Permission = {
  id: number;
  permission: string;
  rule: string;
};

export type RoleParent = {
  role_id: number;
  parent_id: number;
};

export type PermissionParent = {
  permission_id: number;
  parent_id: number;
};

export type RolePermission = {
  role_id: number;
  permission_id: number;
};
