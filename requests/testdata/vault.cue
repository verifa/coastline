package testdata

request: #VaultAdminAccess: {
	kind: "VaultAdminAccess"
	serviceSelector: {
		matchLabels: {
			tool: "vault"
		}
	}

	spec: {
		// Path to Vault mount.
		path: =~"^[A-Za-z0-9-_\\/]+$"
		// Vault role to assume.
		role: =~"^[A-Za-z0-9-_]+$"
	}
}

workflow: vaultAdminAccess: {
	input: request.#VaultAdminAccess

	step: login: {
		token: input.spec.path + "/" + input.spec.role
	}

	output: {
		token: step.login.token
	}
}
