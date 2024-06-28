import { Permission, Role } from "./types";

export function roleExists(roles: Role[], role: Role): boolean {
  for (let i = 0; i < roles.length; i++) {
    if (roles[i].id === role.id) {
      return true;
    }
  }
  return false;
}

// FIXME try to match the naming
export function permissionVisited(
  permissions: Permission[],
  permission: Permission
): boolean {
  for (let i = 0; i < permissions.length; i++) {
    const current = permissions[i];
    if (current.id == permission.id) {
      return true;
    }
  }
  return false;
}

export function checkUserHasRole(userRoles: string[], roles: Role[]): boolean {
  for (let i = 0; i < roles.length; i++) {
    for (let j = 0; j < userRoles.length; j++) {
      if (roles[i].role === userRoles[j]) {
        return true;
      }
    }
  }
  return false;
}
