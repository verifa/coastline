package server

var loginOutputModule = `
package auth

result = {
	"admin": admin,
	"allow": allow,
	"deny": deny,
}

default admin = false
default allow = false
default deny = false
`
var loginPolicies = []Policy{
	{
		Name: "Admins",
		Type: LoginPolicyType,
		Module: `
package auth

admins := ["test"]

admin { input.user.name == admins[_] }
allow = true
`,
	},
}
