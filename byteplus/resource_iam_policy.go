package byteplus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/byteplus-sdk/byteplus-go-sdk-v2/byteplus"
	"github.com/byteplus-sdk/byteplus-go-sdk-v2/byteplus/bytepluserr"
	"github.com/byteplus-sdk/byteplus-go-sdk-v2/service/iam"
	"github.com/cenkalti/backoff"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const maxLength = 6144

var (
	_ resource.Resource              = &iamPolicyResource{}
	_ resource.ResourceWithConfigure = &iamPolicyResource{}
	//_ resource.ResourceWithImportState = &iamPolicyResource{}
)

func NewIamPolicyResource() resource.Resource {
	return &iamPolicyResource{}
}

type iamPolicyResource struct {
	client *iam.IAM
}

type iamPolicyResourceModel struct {
	AttachedPolicies       types.List      `tfsdk:"attached_policies"`
	CombinedPolices        []*policyDetail `tfsdk:"combined_policies"`
	AttachedPoliciesDetail []*policyDetail `tfsdk:"attached_policies_detail"`
	UserName               types.String    `tfsdk:"user_name"`
}

type policyDetail struct {
	PolicyName     types.String `tfsdk:"policy_name"`
	PolicyDocument types.String `tfsdk:"policy_document"`
}

func (r *iamPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_policy"
}

func (r *iamPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a IAM Policy resource that manages policy content " +
			"exceeding character limits by splitting it into smaller segments. " +
			"These segments are combined to form a complete policy attached to " +
			"the user. However, the policy that exceed the maximum length of a " +
			"policy, they will be attached directly to the user.",
		Attributes: map[string]schema.Attribute{
			"attached_policies": schema.ListAttribute{
				Description: "The IAM policies to attach to the user.",
				Required:    true,
				ElementType: types.StringType,
			},
			"combined_policies": schema.ListNestedAttribute{
				Description: "A list of combined policies that are attached to users.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"policy_name": schema.StringAttribute{
							Description: "The policy name.",
							Computed:    true,
						},
						"policy_document": schema.StringAttribute{
							Description: "The policy document of the RAM policy.",
							Computed:    true,
						},
					},
				},
			},
			"attached_policies_detail": schema.ListNestedAttribute{
				Description: "A list of policies. Used to compare whether policy has been changed outside of Terraform",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"policy_name": schema.StringAttribute{
							Description: "The policy name.",
							Computed:    true,
						},
						"policy_document": schema.StringAttribute{
							Description: "The policy document of the RAM policy.",
							Computed:    true,
						},
					},
				},
			},
			"user_name": schema.StringAttribute{
				Description: "The name of the IAM user that attached to the policy.",
				Required:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *iamPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(byteplusClients).iamClient
}

// Create implements resource.Resource.
func (r *iamPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *iamPolicyResourceModel
	getPlanDiags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(getPlanDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, currentPoliciesList, err := r.createPolicy(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to Create the Policy.",
			err.Error(),
		)
		return
	}

	state := &iamPolicyResourceModel{}
	state.AttachedPolicies = plan.AttachedPolicies
	state.CombinedPolices = policy
	state.AttachedPoliciesDetail = currentPoliciesList
	state.UserName = plan.UserName

	if err := r.attachPolicyToUser(state); err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to Attach Policy to User.",
			err.Error(),
		)
		return
	}

	_, errReadPolicyDiags := r.readCombinedPolicy(state)
	resp.Diagnostics.Append(errReadPolicyDiags)
	if resp.Diagnostics.HasError() {
		return
	}

	setStateDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *iamPolicyResourceModel
	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This state will be using to compare with the current state.
	var oriState *iamPolicyResourceModel
	getOriStateDiags := req.State.Get(ctx, &oriState)
	resp.Diagnostics.Append(getOriStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	warnReadPolicyDiags, errReadPolicyDiags := r.readCombinedPolicy(state)
	resp.Diagnostics.Append(warnReadPolicyDiags, errReadPolicyDiags)
	if resp.Diagnostics.HasError() {
		return
	}

	setStateDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	warnReadPolicyDiags, errReadPolicyDiags = r.readAttachedPolicy(state, true)

	resp.Diagnostics.Append(warnReadPolicyDiags, errReadPolicyDiags)
	if resp.Diagnostics.HasError() {
		return
	}

	setStateDiags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if warnReadPolicyDiags == nil {
		compareEachPolicyDiags := r.compareEachPolicy(state, oriState)
		resp.Diagnostics.Append(compareEachPolicyDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	setStateDiags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *iamPolicyResourceModel
	getPlanDiags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(getPlanDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	warnReadPolicyDiags, errReadPolicyDiags := r.readAttachedPolicy(plan, false) //to prevent removal of combined policies, if user inputs non-existing attached policies
	resp.Diagnostics.Append(warnReadPolicyDiags, errReadPolicyDiags)
	if resp.Diagnostics.HasError() {
		return
	}

	removePolicyDiags := r.removePolicy(state)
	resp.Diagnostics.Append(removePolicyDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, currentPoliciesList, err := r.createPolicy(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to Update the Policy.",
			err.Error(),
		)
		return
	}

	state.AttachedPolicies = plan.AttachedPolicies
	state.CombinedPolices = policy
	state.AttachedPoliciesDetail = currentPoliciesList
	state.UserName = plan.UserName

	if err := r.attachPolicyToUser(state); err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to Attach Policy to User.",
			err.Error(),
		)
		return
	}

	_, errReadPolicyDiags = r.readCombinedPolicy(state)
	resp.Diagnostics.Append(errReadPolicyDiags)
	if resp.Diagnostics.HasError() {
		return
	}

	setStateDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements resource.Resource.
func (r *iamPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *iamPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	removePolicyDiags := r.removePolicy(state)
	resp.Diagnostics.Append(removePolicyDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamPolicyResource) createPolicy(plan *iamPolicyResourceModel) (policiesList []*policyDetail, currentPoliciesList []*policyDetail, err error) {
	combinedPolicyStatements, notCombinedPolicies, currentPoliciesStatements, err := r.combinePolicyDocument(plan)
	if err != nil {
		return nil, nil, err
	}

	createPolicy := func() error {
		for i, policy := range combinedPolicyStatements {
			policyName := plan.UserName.ValueString() + "-" + strconv.Itoa(i+1)

			createPolicyRequest := &iam.CreatePolicyInput{
				PolicyName:     byteplus.String(policyName),
				PolicyDocument: byteplus.String(policy),
			}

			if _, err := r.client.CreatePolicy(createPolicyRequest); err != nil {
				return handleAPIError(err)
			}
		}

		return nil
	}

	reconnectBackoff := backoff.NewExponentialBackOff()
	reconnectBackoff.MaxElapsedTime = 30 * time.Second
	err = backoff.Retry(createPolicy, reconnectBackoff)

	if err != nil {
		return nil, nil, err
	}

	for i, policies := range combinedPolicyStatements {
		policyName := plan.UserName.ValueString() + "-" + strconv.Itoa(i+1)

		policyObj := &policyDetail{
			PolicyName:     types.StringValue(policyName),
			PolicyDocument: types.StringValue(policies),
		}

		policiesList = append(policiesList, policyObj)
	}

	// These policies will be attached directly to the user since splitting the
	// policy "statement" will be hitting the limitation of "maximum number of
	// attached policies" easily.
	for _, policy := range notCombinedPolicies {
		policyObj := &policyDetail{
			PolicyName:     types.StringValue(policy.policyName),
			PolicyDocument: types.StringValue(policy.policyDocument),
		}

		policiesList = append(policiesList, policyObj)
	}

	// These policies are used for comparing whether there is a differerence
	// between current policies in state file and in the console
	for _, policy := range currentPoliciesStatements {
		policyObj := &policyDetail{
			PolicyName:     types.StringValue(strings.Trim(policy.policyName, "\"")),
			PolicyDocument: types.StringValue(policy.policyDocument),
		}

		currentPoliciesList = append(currentPoliciesList, policyObj)
	}

	return policiesList, currentPoliciesList, nil
}

func (r *iamPolicyResource) readCombinedPolicy(state *iamPolicyResourceModel) (warnDiagnostics, errDiagnostics diag.Diagnostic) {
	policyDetailsState := []*policyDetail{}

	var warning, err error
	getPolicy := func() error {
		for _, combinedPolicy := range state.CombinedPolices {
			getPolicyRequest := iam.GetPolicyInput{
				PolicyName: byteplus.String(combinedPolicy.PolicyName.ValueString()),
				PolicyType: byteplus.String("Custom"),
			}

			getPolicyResponse, errGetCombinedPolicy := r.client.GetPolicy(&getPolicyRequest)
			if errGetCombinedPolicy != nil {
				if errGetCombinedPolicy.(bytepluserr.Error).Code() == "PolicyNotExist" {
					// To detect if policy has been deleted after being attached.
					warning = errors.Join(warning, handleAPIError(errGetCombinedPolicy))
				} else {
					err = errors.Join(err, handleAPIError(errGetCombinedPolicy))
				}
				continue
			}

			// Sometimes combined policies may be removed accidentally by human mistake or API error.
			if getPolicyResponse != nil && getPolicyResponse.Policy != nil {
				if getPolicyResponse.Policy.PolicyName != nil && *getPolicyResponse.Policy.PolicyDocument != "" {
					policyDetail := policyDetail{
						PolicyName:     types.StringValue(*getPolicyResponse.Policy.PolicyName),
						PolicyDocument: types.StringValue(*getPolicyResponse.Policy.PolicyDocument),
					}
					policyDetailsState = append(policyDetailsState, &policyDetail)
				}
			}
		}
		return nil
	}

	reconnectBackoff := backoff.NewExponentialBackOff()
	reconnectBackoff.MaxElapsedTime = 30 * time.Second
	err = backoff.Retry(getPolicy, reconnectBackoff)
	if err != nil {
		errDiagnostics = diag.NewErrorDiagnostic(
			"[API ERROR] Failed to Read Combined Policy",
			err.Error(),
		)
		return nil, errDiagnostics
	}

	if warning != nil {
		warnDiagnostics = diag.NewWarningDiagnostic(
			"Combined Policies could not be found.",
			"The combined policies attached to the user may be deleted due to human mistake or API error. This resource will be re-created.\n\n"+
				warning.Error(),
		)

		state.AttachedPolicies = types.ListNull(types.StringType) //This is to ensure Update() is called
	}

	state.CombinedPolices = policyDetailsState

	return warnDiagnostics, errDiagnostics
}

func (r *iamPolicyResource) readAttachedPolicy(state *iamPolicyResourceModel, inRead bool) (warnDiagnostics, errDiagnostics diag.Diagnostic) {
	attachedPolicies := state.AttachedPolicies.Elements()
	policyDetailsState, warning, err := r.fetchPolicies(attachedPolicies, inRead)

	if err != nil {
		errDiagnostics = diag.NewErrorDiagnostic(
			"[API ERROR] Failed to Read Attached Policy",
			err.Error(),
		)
	}

	if warning != nil {
		warnDiagnostics = diag.NewWarningDiagnostic(
			"One (or more) of the Attached Policy could not be found.",
			"The policy used for Combined Policies may be deleted due to human mistake or API error.\n\n"+
				warning.Error(),
		)
	}

	state.AttachedPoliciesDetail = policyDetailsState
	if warnDiagnostics != nil {
		state.AttachedPolicies = types.ListNull(types.StringType) // Ensure Update() is called
	}

	return warnDiagnostics, errDiagnostics
}

func (r *iamPolicyResource) fetchPolicies(attachedPolicies []attr.Value, inRead bool) (policyDetailsState []*policyDetail, errNotExist, errOther error) {
	getPolicyResponse := &iam.GetPolicyOutput{}

	var errGetEachPolicy error

	getPolicy := func() error {
	OuterLoop:
		for _, policy := range attachedPolicies {
			getPolicyRequest := &iam.GetPolicyInput{
				PolicyName: byteplus.String(strings.Trim(policy.String(), "\"")),
				PolicyType: byteplus.String("Custom"),
			}

			for {
				getPolicyResponse, errGetEachPolicy = r.client.GetPolicy(getPolicyRequest)
				if errGetEachPolicy != nil {
					switch errGetEachPolicy.(bytepluserr.Error).Code() {
					case "PolicyNotExist":
						if *getPolicyRequest.PolicyType == "Custom" {
							*getPolicyRequest.PolicyType = "System" // Switch to System and Retry
							continue
						} else if inRead {
							errNotExist = errors.Join(errNotExist, handleAPIError(errGetEachPolicy))
						} else {
							errOther = errors.Join(errOther, handleAPIError(errGetEachPolicy))
						}
					default:
						errOther = errors.Join(errOther, handleAPIError(errGetEachPolicy))
					}
					continue OuterLoop
				}
				break
			}

			if getPolicyResponse != nil && getPolicyResponse.Policy != nil {
				if getPolicyResponse.Policy.PolicyName != nil && getPolicyResponse.Policy.PolicyDocument != nil {
					policyDetail := policyDetail{
						PolicyName:     types.StringValue(*getPolicyResponse.Policy.PolicyName),
						PolicyDocument: types.StringValue(*getPolicyResponse.Policy.PolicyDocument),
					}
					policyDetailsState = append(policyDetailsState, &policyDetail)
				}
			}
		}
		return nil
	}

	reconnectBackoff := backoff.NewExponentialBackOff()
	reconnectBackoff.MaxElapsedTime = 30 * time.Second
	backoff.Retry(getPolicy, reconnectBackoff)

	return policyDetailsState, errNotExist, errOther
}

func (r *iamPolicyResource) compareEachPolicy(newState, oriState *iamPolicyResourceModel) diag.Diagnostics {
	var driftedPolicies []string

	for _, oldPolicyDetailState := range oriState.AttachedPoliciesDetail {
		for _, currPolicyDetailState := range newState.AttachedPoliciesDetail {
			if oldPolicyDetailState.PolicyName.String() == currPolicyDetailState.PolicyName.String() {
				if oldPolicyDetailState.PolicyDocument.String() != currPolicyDetailState.PolicyDocument.String() {
					driftedPolicies = append(driftedPolicies, oldPolicyDetailState.PolicyName.String())
				}
			}
		}
	}

	if len(driftedPolicies) > 0 {
		driftedPoliciesMessage := fmt.Sprintf(
			"The following policies have drifted: %s. It may be caused by modifying the .json file outside of Terraform.",
			strings.Join(driftedPolicies, ", "),
		)

		newState.AttachedPolicies = types.ListNull(types.StringType) // Set the state to trigger an update.
		return diag.Diagnostics{
			diag.NewWarningDiagnostic(
				"Policy Drift Detected.",
				driftedPoliciesMessage,
			),
		}
	}

	return nil
}

func (r *iamPolicyResource) removePolicy(state *iamPolicyResourceModel) diag.Diagnostics {
	removePolicy := func() error {
		for _, combinedPolicy := range state.CombinedPolices {
			detachPolicyFromUserRequest := &iam.DetachUserPolicyInput{
				PolicyType: byteplus.String("Custom"),
				PolicyName: byteplus.String(combinedPolicy.PolicyName.ValueString()),
				UserName:   byteplus.String(state.UserName.ValueString()),
			}

			deletePolicyRequest := &iam.DeletePolicyInput{
				PolicyName: byteplus.String(combinedPolicy.PolicyName.ValueString()),
			}

			if _, err := r.client.DetachUserPolicy(detachPolicyFromUserRequest); err != nil {
				return handleAPIError(err)
			}

			if _, err := r.client.DeletePolicy(deletePolicyRequest); err != nil {
				return handleAPIError(err)
			}
		}

		return nil
	}

	reconnectBackoff := backoff.NewExponentialBackOff()
	reconnectBackoff.MaxElapsedTime = 30 * time.Second
	err := backoff.Retry(removePolicy, reconnectBackoff)
	if err != nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"[API ERROR] Failed to Delete Policy",
				err.Error(),
			),
		}
	}

	return nil
}

type simplePolicy struct {
	policyName     string
	policyDocument string
}

func (r *iamPolicyResource) combinePolicyDocument(plan *iamPolicyResourceModel) (finalPolicyDocument []string, excludedPolicy []simplePolicy, currentPolicyList []simplePolicy, err error) {
	attachedPolicies := plan.AttachedPolicies.Elements()
	policyDetailsState, _, err := r.fetchPolicies(attachedPolicies, false)

	const policyKeywordLen = 30

	if err != nil {
		return nil, nil, nil, err
	}

	currentLength := 0
	currentPolicyDocument := ""
	appendedPolicyDocument := make([]string, 0)

	for _, detail := range policyDetailsState {
		tempPolicyDocument := detail.PolicyDocument.ValueString()

		currentPolicyList = append(currentPolicyList, simplePolicy{
			policyName:     detail.PolicyName.ValueString(),
			policyDocument: tempPolicyDocument,
		})

		// If the policy itself have more than 6144 characters, then skip the combine
		// policy part since splitting the policy "statement" will be hitting the
		// limitation of "maximum number of attached policies" easily.
		if len(tempPolicyDocument) > maxLength {
			excludedPolicy = append(excludedPolicy, simplePolicy{
				policyName:     detail.PolicyName.ValueString(),
				policyDocument: tempPolicyDocument,
			})
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(tempPolicyDocument), &data); err != nil {
			return nil, nil, nil, err
		}

		statementArr := data["Statement"].([]interface{})
		statementBytes, err := json.Marshal(statementArr)
		if err != nil {
			return nil, nil, nil, err
		}

		finalStatement := strings.Trim(string(statementBytes), "[]")
		currentLength += len(finalStatement)

		// Before further proceeding the current policy, we need to add a number of 30 to simulate the total length of completed policy to check whether it is already execeeded the max character length of 6144.
		// Number of 30 indicates the character length of neccessary policy keyword such as "Version" and "Statement" and some JSON symbols ({}, [])
		if (currentLength + policyKeywordLen) > maxLength {
			currentPolicyDocument = strings.TrimSuffix(currentPolicyDocument, ",")
			appendedPolicyDocument = append(appendedPolicyDocument, currentPolicyDocument)
			currentPolicyDocument = finalStatement + ","
			currentLength = len(finalStatement)
		} else {
			currentPolicyDocument += finalStatement + ","
		}
	}

	if len(currentPolicyDocument) > 0 {
		currentPolicyDocument = strings.TrimSuffix(currentPolicyDocument, ",")
		appendedPolicyDocument = append(appendedPolicyDocument, currentPolicyDocument)
	}

	for _, policy := range appendedPolicyDocument {
		finalPolicyDocument = append(finalPolicyDocument, fmt.Sprintf(`{"Version":"1","Statement":[%v]}`, policy))
	}

	return finalPolicyDocument, excludedPolicy, currentPolicyList, nil
}

func (r *iamPolicyResource) attachPolicyToUser(state *iamPolicyResourceModel) (err error) {
	attachPolicyToUser := func() error {
		for _, combinedPolicy := range state.CombinedPolices {
			attachPolicyToUserRequest := &iam.AttachUserPolicyInput{
				PolicyType: byteplus.String("Custom"),
				PolicyName: byteplus.String(combinedPolicy.PolicyName.ValueString()),
				UserName:   byteplus.String(state.UserName.ValueString()),
			}

			if _, err := r.client.AttachUserPolicy(attachPolicyToUserRequest); err != nil {
				return handleAPIError(err)
			}
		}
		return nil
	}

	reconnectBackoff := backoff.NewExponentialBackOff()
	reconnectBackoff.MaxElapsedTime = 30 * time.Second
	return backoff.Retry(attachPolicyToUser, reconnectBackoff)
}

func handleAPIError(err error) error {
	if _t, ok := err.(bytepluserr.Error); ok {
		if isAbleToRetry(_t.Code()) {
			return err
		} else {
			return backoff.Permanent(err)
		}
	} else {
		return err
	}
}
