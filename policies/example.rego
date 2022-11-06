package example

import future.keywords.if
import future.keywords.in

# example input:
#{
#    "project": {
#        "id": "11111111-0000000000-22222222",
#        "labels": [
#            "dev",
#            "ops"
#        ],
#        "name": "dummy-project",
#        "owner": "bob"
#    },
#    "request": {
#        "epoch_time": 1666549548,
#        "remote_ip": "1.2.3.4"
#    },
#    "service": {
#        "id": "00000000-11111111-22222222",
#        "labels": [
#            "buzz",
#            "word"
#        ],
#        "name": "artifactory"
#    },
#    "session": {
#        "admin": false,
#        "email": "bob@localhost",
#        "groups": [
#            "builders",
#            "superdevops",
#            "devsecops"
#        ],
#        "name": "bob"
#    }
#}

# read can be set to true by default as well
default read := false

#write implies read
default write := false

# explicit deny, useful when restricting using IP or time (or other non-user attribute)
default deny_write := false

# sets
writer_groups := {"superdevops"}
reader_groups := {"superdevops"}
writer_users := {"bob"}
reader_users := {"bob"}

# global admin (or should admins skip policy checks?)
write if input.session.admin == true
# admin cannot be blocked explicitly either
deny_write := false if input.session.admin == true

# owner always write
write if input.session.name == input.project.owner
# owner cannot be blocked explicitly (by default)
# how to implement this best? cannot just set it to false, since any other rule might come out as true
#deny_write := false if input.session.name == input.project.owner

# check if allowed_users includes the current user
write if writer_users[input.session.name]
read if reader_users[input.session.name]

#check if allowed groups includes group of the user
write if {
	some group in input.session.groups
		group in writer_groups
}

read if {
	some group in input.session.groups
		group in reader_groups
}

#deny if service is artifactory
deny_write if input.service.name == "artifactory"
