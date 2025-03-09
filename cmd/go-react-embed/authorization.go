package main

import (
	"context"
	"encoding/json"
	"go-react-embed/models"
	"go-react-embed/rbac"
)

func getRbacData(ctx context.Context, queries *models.Queries) (
	[]models.GetRolesRow, 
	[]models.GetPermissionsRow,
	[]models.RoleParent,
	[]models.PermissionParent,
	[]models.RolePermission,
	error,
) {
	roles, err := queries.GetRoles(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	permissions, err := queries.GetPermissions(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	roleParents, err := queries.GetRoleParents(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	permissionParents, err := queries.GetPermissionParents(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	rolePermissions, err := queries.GetRolePermissions(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return roles, permissions, roleParents, permissionParents, rolePermissions, nil
}

// func setupRBAC(ctx context.Context, queries *models.Queries) (rbac.RBAC, error) {
func setupRBAC(ctx context.Context, queries *models.Queries) (rbac.RBAC, error) {
	rbacAuth := rbac.New()

	roles, permissions, roleParents, permissionParents, rolePermissions, err := getRbacData(ctx, queries)
	if err != nil {
		return nil, err
	}

	recastedRoles := &[]rbac.Role{}
	recastedPermission := &[]rbac.Permission{}
	recastedRoleParent := &[]rbac.RoleParent{}
	recastedPermissionParent := &[]rbac.PermissionParent{}
	recastedRolePermission := &[]rbac.RolePermission{}

	// recasting type from 
	recast(roles, recastedRoles)
	recast(permissions, recastedPermission)
	recast(roleParents, recastedRoleParent)
	recast(permissionParents, recastedPermissionParent)
	recast(rolePermissions, recastedRolePermission)

	// setting RBAC authorization
	rbacAuth.SetRoles(*recastedRoles)
	rbacAuth.SetPermissions(*recastedPermission)
	rbacAuth.SetRoleParents(*recastedRoleParent)
	rbacAuth.SetPermissionParents(*recastedPermissionParent)
	rbacAuth.SetRolePermissions(*recastedRolePermission)
	
	// rbacAuth.SetRuleEvalCode(
	//   `function listHasValue(obj, val) {
	// 	var values = Object.values(obj);
	// 	for(var i = 0; i < values.length; i++){
	// 	  if(values[i] === val) {
	// 		return true;
	// 	  }
	// 	}
	// 	return false;
	//   }
	//   function rule(user, resource) {
	// 	console.log("set at main");
	// 	return %s;
	//   }`)

	return rbacAuth, nil
}

func recast(a, b interface{}) error {
    js, err := json.Marshal(a)
    if err != nil {
        return err
    }
    return json.Unmarshal(js, b)
}
