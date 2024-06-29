import { DefaultEngine } from "./engine";
import {
  Permission,
  PermissionParent,
  Role,
  RoleParent,
  RolePermission,
  UserWithRoles,
  PermissionsMap,
  RolesMap,
} from "./types";
import { checkUserHasRole, roleExist } from "./utils";

export interface EvalEngine {
  RunRule: (user: UserWithRoles, resource: object, rule: string) => boolean;
  SetOtherCode: (code: string) => void;
  SetRuleCode: (code: string) => void;
}

export class RBAC {
  roles: Role[];
  permissions: Permission[];
  role_parents: RoleParent[];
  permission_parents: PermissionParent[];
  role_permissions: RolePermission[];
  evalEngine: EvalEngine;

  constructor(ee: EvalEngine | null = null) {
    this.roles = [];
    this.permissions = [];
    this.role_parents = [];
    this.permission_parents = [];
    this.role_permissions = [];

    if (ee === null) {
      this.evalEngine = new DefaultEngine();
    } else {
      this.evalEngine = ee;
    }
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

  GetEvalEngine(): EvalEngine {
    return this.evalEngine;
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

  collectRoles(foundRoles: RolesMap): Role[] {
    const roles: Role[] = [];
    const _rbac = this;

    function dfs(child: Role) {
      if (roleExist(roles, child)) {
        return;
      }
      roles.push(child);
      const parents = _rbac.getRoleParents(child.id);
      for (let i = 0; i < parents.length; i++) {
        dfs(parents[i]);
      }
    }
    for (const key in foundRoles) {
      const child = foundRoles[key];
      dfs(child);
    }
    return roles;
  }

  hasPermission(
    user: UserWithRoles,
    resource: object,
    firstPermission: Permission
  ) {
    const visitedPerissions: PermissionsMap = {};
    const foundRoles: RolesMap = {};
    let breaked = false;
    const _rbac = this;

    function dfs(child: Permission): boolean {
      if (breaked) {
        return breaked;
      }
      if (child.id in visitedPerissions) {
        return false;
      }
      // check rule is true
      const result = _rbac.evalEngine.RunRule(user, resource, child.rule);
      if (result === false) {
        return false;
      }
      // const perm = child.id + ", " + child.permission + ", " + child.rule;
      // console.log("+ next:[" + perm + "]", result);

      visitedPerissions[child.id] = child;

      // get roles related to permissions
      // if user has appropriate role we break
      const roles = _rbac.getPermissionRoles(child.id);
      for (let i = 0; i < roles.length; i++) {
        // if (!foundRoles.hasOwnProperty(roles[i].id)) {
        const role = roles[i];
        foundRoles[role.id] = role;
        // }
      }

      const hasRole = checkUserHasRole(user.roles, roles);
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
        dfs(parents[i]);
      }
      return breaked;
    }

    const allowed = dfs(firstPermission);
    return { allowed, foundRoles };
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

    const { allowed, foundRoles } = this.hasPermission(
      user,
      resource,
      firstPermission
    );

    const roles = this.collectRoles(foundRoles);

    if (allowed === false) {
      // check again if user has role if breaked allowed is false
      const hasRole = checkUserHasRole(user.roles, roles);
      return hasRole;
    } else {
      return allowed;
    }
  }
}
