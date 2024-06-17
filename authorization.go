package main

import (
	"encoding/json"
	"fmt"
	"go-react-embed/models"
	"go-react-embed/rbac"
)


var (
	RBAC rbac.RBAC
)


func getRbacData() (
	[]models.GetRolesRow, 
	[]models.GetPermissionsRow,
	[]models.GetRoleParentsRow,
	[]models.GetPermissionParentsRow,
	[]models.RolePermission,
	error,
) {
	roles, err := models.QUERIES.GetRoles(models.CTX)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	permissions, err := models.QUERIES.GetPermissions(models.CTX)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	roleParents, err := models.QUERIES.GetRoleParents(models.CTX)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	permissionParents, err := models.QUERIES.GetPermissionParents(models.CTX)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	rolePermissions, err := models.QUERIES.GetRolePermissions(models.CTX)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return roles, permissions, roleParents, permissionParents, rolePermissions, nil
}

func setupRBAC() (rbac.RBAC, error) {
	rbacAuth := rbac.New()

	roles, permissions, roleParents, permissionParents, rolePermissions, err := getRbacData()
	if err != nil {
		return nil, err
	}

	recastedRoles := &[]rbac.Role{}
	recastedPermission := &[]rbac.Permission{}
	recastedRoleParent := &[]rbac.RoleParent{}
	recastedPermissionParent := &[]rbac.PermissionParent{}
	recastedRolePermission := &[]rbac.RolePermission{}

	recast(roles, recastedRoles)
	recast(permissions, recastedPermission)
	recast(roleParents, recastedRoleParent)
	recast(permissionParents, recastedPermissionParent)
	recast(rolePermissions, recastedRolePermission)

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
	//   function rule(user, ressource) {
	// 	console.log("set at main");
	// 	return %s;
	//   }`)

	fmt.Println("-rbacAuth:", rbacAuth)

	return rbacAuth, nil
}

func recast(a, b interface{}) error {
    js, err := json.Marshal(a)
    if err != nil {
        return err
    }
    return json.Unmarshal(js, b)
}
