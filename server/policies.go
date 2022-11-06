package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
	"github.com/verifa/coastline/policies"
	"github.com/verifa/coastline/server/oapi"
)

type PolicyType string

const (
	LoginPolicyType  PolicyType = "login"
	AccessPolicyType PolicyType = "access"
)

type Policy struct {
	Name   string
	Type   PolicyType
	Module string
}

type PolicyEngine struct {
	Policies map[PolicyType][]Policy
}

func NewPolicyEngine() (*PolicyEngine, error) {
	policies := make(map[PolicyType][]Policy)
	policies[LoginPolicyType] = loginPolicies

	return &PolicyEngine{
		Policies: policies,
	}, nil
}

func (pe PolicyEngine) EvaluateLoginRequest(userInfo oapi.User) (bool, error) {
	loginPolicy, err := policies.Policies.ReadFile("login_policy.rego")
	if err != nil {
		return false, fmt.Errorf("reading login policy file: %w", err)
	}
	policies := pe.Policies[LoginPolicyType]

	var options = []func(r *rego.Rego){
		rego.Query("data.auth.result"),
		rego.Module("defaults", string(loginPolicy)),
	}
	for _, policy := range policies {
		options = append(options, rego.Module(policy.Name, policy.Module))
	}

	ctx := context.TODO()
	query, err := rego.New(options...).PrepareForEval(ctx)
	if err != nil {
		return false, fmt.Errorf("preparing query: %w", err)
	}

	input := map[string]interface{}{
		"user": userInfo,
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, fmt.Errorf("evaluating query: %w", err)
	} else if len(results) == 0 {
		return false, errors.New("undefined result possibly due to bad policies")
	} else {
		exprs := results[0].Expressions
		for _, expr := range exprs {
			fmt.Println("expr: ", expr.Value.(map[string]interface{}))
		}
	}

	return false, nil
}

func (pe PolicyEngine) EvaluateAccessPolicy(userInfo oapi.User, project oapi.Project) (bool, error) {
	// should actually read something like <project>_access_policy.rego or fallback to default access_policy.rego
	// also should load these outside of the handler, but no idea how we will store these atm
	accessPolicy, err := policies.Policies.ReadFile("access_policy.rego")
	if err != nil {
		return false, fmt.Errorf("reading access policy file: %w", err)
	}

	var options = []func(r *rego.Rego){
		rego.Query("data.authz.result"),
		rego.Module("defaults", string(accessPolicy)),
	}

	ctx := context.TODO()
	query, err := rego.New(options...).PrepareForEval(ctx)
	if err != nil {
		return false, fmt.Errorf("preparing query: %w", err)
	}

	input := map[string]interface{}{
		"user":    userInfo,
		"project": project,
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, fmt.Errorf("evaluating query: %w", err)
	} else if len(results) == 0 {
		fmt.Println("err", err)
		return false, errors.New("undefined result possibly due to bad policies")
	} else {
		exprs := results[0].Expressions
		for _, expr := range exprs {
			fmt.Println("expr: ", expr.Value.(map[string]interface{}))
		}
	}
	resultMap := results[0].Expressions[0].Value.(map[string]interface{})
	allow := resultMap["allow"].(bool)
	fmt.Println("allow:", allow)

	return allow, nil
}
