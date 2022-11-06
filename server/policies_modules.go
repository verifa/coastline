package server

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
