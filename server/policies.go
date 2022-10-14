package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
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

func NewPolicyEngine() *PolicyEngine {
	policies := make(map[PolicyType][]Policy)
	policies[LoginPolicyType] = loginPolicies

	return &PolicyEngine{
		Policies: policies,
	}
}

func (pe PolicyEngine) EvaluateLoginRequest(userInfo UserInfo) (bool, error) {

	policies := pe.Policies[LoginPolicyType]

	var options = []func(r *rego.Rego){
		rego.Query("data.auth.result"),
		rego.Module("defaults", loginOutputModule),
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
