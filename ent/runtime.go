// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/google/uuid"
	"github.com/verifa/coastline/ent/approval"
	"github.com/verifa/coastline/ent/project"
	"github.com/verifa/coastline/ent/request"
	"github.com/verifa/coastline/ent/schema"
	"github.com/verifa/coastline/ent/service"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	approvalFields := schema.Approval{}.Fields()
	_ = approvalFields
	// approvalDescIsAutomated is the schema descriptor for is_automated field.
	approvalDescIsAutomated := approvalFields[1].Descriptor()
	// approval.DefaultIsAutomated holds the default value on creation for the is_automated field.
	approval.DefaultIsAutomated = approvalDescIsAutomated.Default.(bool)
	// approvalDescApprover is the schema descriptor for approver field.
	approvalDescApprover := approvalFields[2].Descriptor()
	// approval.ApproverValidator is a validator for the "approver" field. It is called by the builders before save.
	approval.ApproverValidator = approvalDescApprover.Validators[0].(func(string) error)
	// approvalDescID is the schema descriptor for id field.
	approvalDescID := approvalFields[0].Descriptor()
	// approval.DefaultID holds the default value on creation for the id field.
	approval.DefaultID = approvalDescID.Default.(func() uuid.UUID)
	projectFields := schema.Project{}.Fields()
	_ = projectFields
	// projectDescName is the schema descriptor for name field.
	projectDescName := projectFields[1].Descriptor()
	// project.NameValidator is a validator for the "name" field. It is called by the builders before save.
	project.NameValidator = projectDescName.Validators[0].(func(string) error)
	// projectDescID is the schema descriptor for id field.
	projectDescID := projectFields[0].Descriptor()
	// project.DefaultID holds the default value on creation for the id field.
	project.DefaultID = projectDescID.Default.(func() uuid.UUID)
	requestFields := schema.Request{}.Fields()
	_ = requestFields
	// requestDescType is the schema descriptor for type field.
	requestDescType := requestFields[1].Descriptor()
	// request.TypeValidator is a validator for the "type" field. It is called by the builders before save.
	request.TypeValidator = requestDescType.Validators[0].(func(string) error)
	// requestDescRequestedBy is the schema descriptor for requested_by field.
	requestDescRequestedBy := requestFields[2].Descriptor()
	// request.RequestedByValidator is a validator for the "requested_by" field. It is called by the builders before save.
	request.RequestedByValidator = requestDescRequestedBy.Validators[0].(func(string) error)
	// requestDescID is the schema descriptor for id field.
	requestDescID := requestFields[0].Descriptor()
	// request.DefaultID holds the default value on creation for the id field.
	request.DefaultID = requestDescID.Default.(func() uuid.UUID)
	serviceFields := schema.Service{}.Fields()
	_ = serviceFields
	// serviceDescName is the schema descriptor for name field.
	serviceDescName := serviceFields[1].Descriptor()
	// service.NameValidator is a validator for the "name" field. It is called by the builders before save.
	service.NameValidator = serviceDescName.Validators[0].(func(string) error)
	// serviceDescID is the schema descriptor for id field.
	serviceDescID := serviceFields[0].Descriptor()
	// service.DefaultID holds the default value on creation for the id field.
	service.DefaultID = serviceDescID.Default.(func() uuid.UUID)
}
