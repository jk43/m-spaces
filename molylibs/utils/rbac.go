package utils

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/moly-space/molylibs/pb"
	csbsa "github.com/qiangmzsx/string-adapter/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserRole string

const (
	RoleGuest      UserRole = "guest"
	RoleMember     UserRole = "member"
	RoleAdmin      UserRole = "admin"
	RoleSuperAdmin UserRole = "super"
)

type Ctx string
type Service string

const (
	CtxOrgMenu         Ctx = "menu"
	CtxOrgForm         Ctx = "form"
	CtxOrgFormEdit     Ctx = "formEdit"
	CtxOrgFormEditView Ctx = "formEditView"
	CtxRBACFilter      Ctx = "filter"
	CtxRBACRole        Ctx = "role"
)

const (
	ServiceOrg     Service = "organization"
	ServiceRBAC    Service = "rbac"
	ServiceMessage Service = "message"
	ServiceBroker  Service = "broker"
)

type Casbin struct {
	Host      string
	Service   Service
	PolicyCtx Ctx
	ModelCtx  Ctx
	Enforcer  *casbin.Enforcer
}

func NewCasbin(host string, service Service, policyCtx Ctx, modelCtx Ctx) *Casbin {
	return &Casbin{Host: host, Service: service, PolicyCtx: policyCtx, ModelCtx: modelCtx}
}

func (csb *Casbin) GetCasbinPolicy() (*csbsa.Adapter, error) {
	in := pb.RBACRequest{
		Host:      csb.Host,
		Service:   string(csb.Service),
		ModelCtx:  string(csb.ModelCtx),
		PolicyCtx: string(csb.PolicyCtx),
	}
	conn, err := grpc.Dial(os.Getenv("RBAC_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewRBACServiceClient(conn)
	res, err := c.GetPolicy(context.Background(), &in)
	if err != nil {
		return nil, err
	}
	pol := strings.Join(res.Policy, "\n")
	sa := csbsa.NewAdapter(pol)
	return sa, nil
}

func (csb *Casbin) GetCasbinModel() (model.Model, error) {
	in := pb.RBACRequest{
		Host:     csb.Host,
		Service:  string(csb.Service),
		ModelCtx: string(csb.ModelCtx),
	}
	conn, err := grpc.Dial(os.Getenv("RBAC_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewRBACServiceClient(conn)
	res, err := c.GetModel(context.Background(), &in)
	model, err := model.NewModelFromString(res.Model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (csb *Casbin) GetCasbinModeAndPolicy() (model.Model, *csbsa.Adapter, error) {
	in := pb.RBACRequest{
		Host:      csb.Host,
		Service:   string(csb.Service),
		ModelCtx:  string(csb.ModelCtx),
		PolicyCtx: string(csb.PolicyCtx),
	}
	conn, err := grpc.Dial(os.Getenv("RBAC_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	c := pb.NewRBACServiceClient(conn)
	res, err := c.GetModelAndPolicy(context.Background(), &in)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()
	m, err := model.NewModelFromString(res.Model)
	if err != nil {
		return nil, nil, err
	}
	pol := strings.Join(res.Policy, "\n")
	sa := csbsa.NewAdapter(pol)
	return m, sa, nil
}

func (csb *Casbin) GetCasbinEnforcer() (*casbin.Enforcer, error) {
	if csb.Enforcer != nil {
		return csb.Enforcer, nil
	}
	// p1, err := csb.GetCasbinPolicy()
	// if err != nil {
	// 	return nil, err
	// }
	// m1, err := csb.GetCasbinModel()
	// if err != nil {
	// 	return nil, err
	// }
	m, p, err := csb.GetCasbinModeAndPolicy()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(m, p)
	if err != nil {
		return nil, err
	}
	csb.Enforcer = e
	return e, nil
}

func (csb *Casbin) GetPermissions(role UserRole) ([][]string, error) {
	e, err := csb.GetCasbinEnforcer()
	if err != nil {
		return nil, err
	}
	return e.GetImplicitPermissionsForUser(string(role))
}

func (csb *Casbin) ParsePermission(role UserRole, f func([][]string) []string) ([]string, error) {
	perm, err := csb.GetPermissions(role)
	if err != nil {
		return nil, err
	}
	return f(perm), nil
}

func (csb *Casbin) HasLink(userRole UserRole, role UserRole) (bool, error) {
	e, err := csb.GetCasbinEnforcer()
	if err != nil {
		return false, err
	}
	if ok, _ := e.GetRoleManager().HasLink(string(userRole), string(role)); ok {
		return true, nil
	}
	return false, nil
}

type CasbinManager struct {
	Service   Service
	PolicyCtx Ctx
	ModelCtx  Ctx
	Casbins   map[string]*Casbin
}

func (cs *CasbinManager) GetCasbin(host string) *Casbin {
	if cs.Casbins == nil {
		cs.Casbins = make(map[string]*Casbin)
	}
	// e, ok := cs.Casbins[host]
	// if !ok {
	// 	TermDebugging(`ok`, nil)
	// 	e = &Casbin{
	// 		Host:      host,
	// 		Service:   cs.Service,
	// 		PolicyCtx: cs.PolicyCtx,
	// 		ModelCtx:  cs.ModelCtx,
	// 	}
	// 	cs.Casbins[host] = e
	// }
	if e, ok := cs.Casbins[host]; ok {
		return e
	}

	e := &Casbin{
		Host:      host,
		Service:   cs.Service,
		PolicyCtx: cs.PolicyCtx,
		ModelCtx:  cs.ModelCtx,
	}
	cs.Casbins[host] = e
	return e
}

func (cs *CasbinManager) GetEnforcer(host string) (*casbin.Enforcer, error) {
	csb := cs.GetCasbin(host)
	return csb.GetCasbinEnforcer()
}

// type Permission struct {
// 	Sub string
// 	Obj string
// 	Act string
// }
// type Permissions = []Permission

// // to delete
// func GetPermissions(service Service, host string, unit Ctx, role UserRole, permissionParser func([][]string) Permissions) (Permissions, error) {
// 	csb := &Casbin{
// 		Host:      host,
// 		Service:   service,
// 		PolicyCtx: unit}
// 	perm, err := csb.GetPermissions(role)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return permissionParser(perm), nil
// }

// type FormPermissionPackage struct {
// 	AccessibleForms Permissions
// }

// // to delete
// func NewFormPermissionPackage(host string, role UserRole) (*FormPermissionPackage, error) {
// 	return (&FormPermissionPackage{}).GetFormPerm(host, role)
// }

// // to delete
// func (p *FormPermissionPackage) GetFormPerm(host string, role UserRole) (*FormPermissionPackage, error) {
// 	access, err := GetPermissions(ServiceOrg, host, CtxOrgForm, role, GeneralPermissionParser)
// 	if err != nil {
// 		return nil, err
// 	}
// 	//to delete
// 	//edit, err := GetPermissions(ServiceOrg, host, CtxOrgFormEdit, role, GeneralPermissionParser)
// 	pkg := &FormPermissionPackage{
// 		AccessibleForms: access,
// 	}
// 	return pkg, nil
// }

// // to delete
// func (p *FormPermissionPackage) IsAccessible(form string) bool {
// 	for _, f := range p.AccessibleForms {
// 		if f.Act == form {
// 			return true
// 		}
// 	}
// 	return false
// }

// func GeneralPermissionParser(perm [][]string) Permissions {
// 	var permissions Permissions
// 	for _, p := range perm {
// 		permissions = append(permissions, Permission{p[0], p[1], p[2]})
// 	}
// 	return permissions
// }

func GeneralPermissionParser(perm [][]string) []string {
	var perms []string
	for _, p := range perm {
		perms = append(perms, p[2])
	}
	return perms
}

func PermissionParserAsJoin(perm [][]string) []string {
	output := make([]string, 0)
	for _, p := range perm {
		output = append(output, fmt.Sprintf("%s:%s:%s", p[0], p[1], p[2]))
	}
	return output
}
