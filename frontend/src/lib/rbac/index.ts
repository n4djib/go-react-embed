export type UserWithRoles = {
  roles: string[];
} & {
  [key: string]: any;
};

type Role = {
  id: number;
  role: string;
};
type Permission = {
  id: number;
  permission: string;
  rule: string;
};
type RoleParent = {
  role_id: number;
  parent_id: number;
};
type PermissionParent = {
  permission_id: number;
  parent_id: number;
};
type RolePermission = {
  role_id: number;
  permission_id: number;
};

export class RBAC {
  roles: Role[];
  permissions: Permission[];
  role_parents: RoleParent[];
  permission_parents: PermissionParent[];
  role_permissions: RolePermission[];

  constructor() {
    this.roles = [];
    this.permissions = [];
    this.role_parents = [];
    this.permission_parents = [];
    this.role_permissions = [];
  }

  SetRoles(roles: Role[]) {
    this.roles = roles;
  }
  SetPermissions(permissions: Permission[]) {
    this.permissions = permissions;
  }
  SetRoleParents(role_parents: RoleParent[]) {
    this.role_parents = role_parents;
  }
  SetPermissionParents(permission_parents: PermissionParent[]) {
    this.permission_parents = permission_parents;
  }
  SetRolePermissions(role_permissions: RolePermission[]) {
    this.role_permissions = role_permissions;
  }

  getRole(id: number) {
    for (let i = 0; i < this.roles.length; i++) {
      const role = this.roles[i];
      if (role.id == id) {
        return role;
      }
    }
    return null;
  }
  getPermission(id: number) {
    for (let i = 0; i < this.permissions.length; i++) {
      const permission = this.permissions[i];
      if (permission.id == id) {
        return permission;
      }
    }
    return null;
  }
  getRoleParents(id: number) {
    const roles: Role[] = [];
    for (let i = 0; i < this.role_parents.length; i++) {
      const role_parent = this.role_parents[i];
      if (role_parent.role_id == id) {
        const role = this.getRole(role_parent.parent_id);
        roles.push(role!);
      }
    }
    return roles;
  }
  getPermissionParents(id: number) {
    const permissions: Permission[] = [];
    for (let i = 0; i < this.permission_parents.length; i++) {
      const permission_parent = this.permission_parents[i];
      if (permission_parent.permission_id == id) {
        const permission = this.getPermission(permission_parent.parent_id);
        if (permission!.rule.trim() == "") {
          permissions.push(permission!);
        }
      }
    }
    for (let i = 0; i < this.permission_parents.length; i++) {
      const permission_parent = this.permission_parents[i];
      if (permission_parent.permission_id == id) {
        const permission = this.getPermission(permission_parent.parent_id);
        if (permission!.rule.trim() != "") {
          permissions.push(permission!);
        }
      }
    }
    return permissions;
  }
  getPermissionRoles(id: number): Role[] {
    const roles: Role[] = [];
    for (let i = 0; i < this.role_permissions.length; i++) {
      const role_permission = this.role_permissions[i];
      if (role_permission.permission_id == id) {
        const role = this.getRole(role_permission.role_id);
        roles.push(role!);
      }
    }
    return roles;
  }
  getPermissionByName(name: string): Permission | null {
    for (let i = 0; i < this.permissions.length; i++) {
      const permission = this.permissions[i];
      if (permission.permission == name) {
        return permission;
      }
    }
    return null;
  }

  getParentRolesLoop(roles: Role[], child: Role) {
    if (roleExists(roles, child)) {
      return;
    }
    roles.push(child);
    const parents = this.getRoleParents(child.id);
    for (let i = 0; i < parents.length; i++) {
      this.getParentRolesLoop(roles, parents[i]);
    }
  }

  permissionsTraversal(user: UserWithRoles, resource: object, _rbac: RBAC) {
    const visitedPerissions: Permission[] = [];
    const foundRoles: Role[] = [];
    let breaked = false;

    const closureFunc = function (child: Permission): boolean {
      if (breaked) {
        console.log("- breaked at start.", child);
        return breaked;
      }

      if (permissionVisited(visitedPerissions, child)) {
        return false;
      }

      // check rule is true
      const result = RunRule(user, resource, child.rule);
      if (result === false) {
        return false;
      }

      // const perm = child.id + ", " + child.permission + ", " + child.rule;
      // console.log("+ next:[" + perm + "]", result);

      visitedPerissions.push(child);

      // get roles related to permissions
      // if user has appropriate role we break
      let userRoles: string[] = [];
      if (user.hasOwnProperty("roles")) {
        userRoles = user.roles;
      }
      const roles = _rbac.getPermissionRoles(child.id);
      for (let i = 0; i < roles.length; i++) {
        if (!roleExists(foundRoles, roles[i])) {
          foundRoles.push(roles[i]);
        }
      }
      const hasRole = checkUserHasRole(userRoles, roles);
      if (hasRole) {
        breaked = true;
        return true;
      }

      // we next go to parents
      const parents = _rbac.getPermissionParents(child.id);
      for (let i = 0; i < parents.length; i++) {
        if (breaked) {
          console.log("- breaked in loop.", parents[i]);
          return true;
        }
        closureFunc(parents[i]);
      }
      return breaked;
    };

    return { closureFunc, foundRoles };
  }

  IsAllowed(
    user: UserWithRoles,
    resource: object,
    permission: string
  ): boolean {
    // find permission
    const firstPermission = this.getPermissionByName(permission);
    if (!firstPermission) {
      console.log("can't find permission: ", permission);
      return false;
    }

    // check inputs (roles)
    if (!user.hasOwnProperty("roles")) {
      console.log("roles not found in user: ", permission);
      return false;
    }

    const { closureFunc: nextInChainFunc, foundRoles } =
      this.permissionsTraversal(user, resource, this);

    const allowed = nextInChainFunc(firstPermission);
    // console.log("foundRoles:", foundRoles);

    const roles: Role[] = [];
    for (let i = 0; i < foundRoles.length; i++) {
      const child = foundRoles[i];
      this.getParentRolesLoop(roles, child);
    }
    // console.log("roles:", roles);

    return allowed;
  }
}

//
//
//
//

function roleExists(roles: Role[], role: Role): boolean {
  for (let i = 0; i < roles.length; i++) {
    if (roles[i].id === role.id) {
      return true;
    }
  }
  return false;
}

function permissionVisited(
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

function checkUserHasRole(userRoles: string[], roles: Role[]): boolean {
  for (let i = 0; i < roles.length; i++) {
    for (let j = 0; j < userRoles.length; j++) {
      if (roles[i].role === userRoles[j]) {
        return true;
      }
    }
  }
  return false;
}

const otherCode = `
function listHasValue(obj, val) {
	var values = Object.values(obj);
	for(var i = 0; i < values.length; i++){
		if(values[i] === val) {
			return true;
		}
	}
	return false;
}
`;
export function RunRule(user: object, resource: object, rule: string): boolean {
  if (rule.trim() === "") return true;

  // evaluate rule
  const defaultRuleCode = `
    return ${rule};
  `;
  const ruleFunc = new Function(
    "user",
    "resource",
    otherCode + defaultRuleCode
  );
  const result = ruleFunc(user, resource);
  return result;

  // return false;
}
