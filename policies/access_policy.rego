package authz

# example minimalistic input:
#{
#    "project": {
#        "id": "11111111-0000000000-22222222",
#        "name": "dummy-project",
#        "owner": "bob"
#    },
#    "user": {
#        "name": "bob"
#    }
#}

# could be set to true by default as well
default allow := false

# sets
allowed_users := {"dev"}

result = {
	"allow": allow,
	"input_user": input,
}

# check if allowed_users includes the current user
allow {
	allowed_users[input.user.name]
}

inputs := input.user