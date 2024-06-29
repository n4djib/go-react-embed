import { Role } from "./types";

export function roleExist(roles: Role[], role: Role): boolean {
  for (let i = 0; i < roles.length; i++) {
    if (roles[i].id === role.id) {
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
