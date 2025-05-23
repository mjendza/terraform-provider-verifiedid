// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-msgraph/internal/clients"
	"github.com/microsoft/terraform-provider-msgraph/internal/docstrings"
	"github.com/microsoft/terraform-provider-msgraph/internal/dynamic"
	"github.com/microsoft/terraform-provider-msgraph/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                     = &MSGraphResource{}
	_ resource.ResourceWithImportState      = &MSGraphResource{}
	_ resource.ResourceWithConfigValidators = &MSGraphResource{}
	_ resource.ResourceWithModifyPlan       = &MSGraphResource{}
)

func NewMSGraphResource() resource.Resource {
	return &MSGraphResource{}
}

// MSGraphResource defines the resource implementation.
type MSGraphResource struct {
	client *clients.MSGraphClient
}

func (r *MSGraphResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

// MSGraphResourceModel describes the resource data model.
type MSGraphResourceModel struct {
	Id                    types.String      `tfsdk:"id"`
	ApiVersion            types.String      `tfsdk:"api_version"`
	Url                   types.String      `tfsdk:"url"`
	Body                  types.Dynamic     `tfsdk:"body"`
	CreateQueryParameters types.Map         `tfsdk:"create_query_parameters"`
	UpdateQueryParameters types.Map         `tfsdk:"update_query_parameters"`
	ReadQueryParameters   types.Map         `tfsdk:"read_query_parameters"`
	DeleteQueryParameters types.Map         `tfsdk:"delete_query_parameters"`
	ResponseExportValues  map[string]string `tfsdk:"response_export_values"`
	Output                types.Dynamic     `tfsdk:"output"`
}

func (r *MSGraphResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

func (r *MSGraphResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource can manage any Microsoft Graph API resource.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: docstrings.ResourceID(),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"url": schema.StringAttribute{
				MarkdownDescription: docstrings.Url("resource"),
				Required:            true,
			},

			"api_version": schema.StringAttribute{
				MarkdownDescription: docstrings.ApiVersion(),
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("v1.0", "beta"),
				},
				Default: stringdefault.StaticString("v1.0"),
			},

			"body": schema.DynamicAttribute{
				MarkdownDescription: docstrings.Body(),
				Optional:            true,
			},

			"create_query_parameters": schema.MapAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Optional:            true,
				MarkdownDescription: "A mapping of query parameters to be sent with the create request.",
			},

			"update_query_parameters": schema.MapAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Optional:            true,
				MarkdownDescription: "A mapping of query parameters to be sent with the update request.",
			},

			"read_query_parameters": schema.MapAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Optional:            true,
				MarkdownDescription: "A mapping of query parameters to be sent with the read request.",
			},

			"delete_query_parameters": schema.MapAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Optional:            true,
				MarkdownDescription: "A mapping of query parameters to be sent with the delete request.",
			},

			"response_export_values": schema.MapAttribute{
				MarkdownDescription: docstrings.ResponseExportValues(),
				Optional:            true,
				ElementType:         types.StringType,
			},

			"output": schema.DynamicAttribute{
				MarkdownDescription: docstrings.Output(),
				Computed:            true,
			},
		},
	}
}

func (r *MSGraphResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if v, ok := req.ProviderData.(*clients.Client); ok {
		r.client = v.MSGraphClient
	}
}

func (r *MSGraphResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	var plan *MSGraphResourceModel
	if response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...); response.Diagnostics.HasError() {
		return
	}

	var state *MSGraphResourceModel
	if response.Diagnostics.Append(request.State.Get(ctx, &state)...); response.Diagnostics.HasError() {
		return
	}

	if plan == nil || state == nil {
		return
	}

	if strings.Contains(plan.Url.ValueString(), "/$ref") {
		if !dynamic.SemanticallyEqual(plan.Body, state.Body) {
			response.RequiresReplace.Append(path.Root("body"))
		}
		if !reflect.DeepEqual(plan.ResponseExportValues, state.ResponseExportValues) {
			response.RequiresReplace.Append(path.Root("response_export_values"))
		}
		if !reflect.DeepEqual(plan.ApiVersion, state.ApiVersion) {
			response.RequiresReplace.Append(path.Root("api_version"))
		}
	}
}

func (r *MSGraphResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model *MSGraphResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	data, err := dynamic.ToJSON(model.Body)
	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal body", err.Error())
		return
	}
	var requestBody interface{}
	if err = json.Unmarshal(data, &requestBody); err != nil {
		resp.Diagnostics.AddError("Failed to unmarshal body", err.Error())
		return
	}

	options := clients.NewRequestOptions(nil, AsMapOfLists(model.CreateQueryParameters))
	responseBody, err := r.client.Create(ctx, model.Url.ValueString(), model.ApiVersion.ValueString(), requestBody, options)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create resource", err.Error())
		return
	}

	if strings.HasSuffix(model.Url.ValueString(), "/$ref") { // extract the id from the response body
		if requestMap, ok := requestBody.(map[string]interface{}); ok {
			if idValue, ok := requestMap["@odata.id"]; ok {
				if idString, ok := idValue.(string); ok {
					uuidValue := idString[strings.LastIndex(idString, "/")+1:]
					model.Id = types.StringValue(uuidValue)
				}
			}
		}
	} else {
		responseId := ""
		if responseBody != nil {
			if responseMap, ok := responseBody.(map[string]interface{}); ok {
				if idValue, ok := responseMap["id"]; ok && idValue != nil {
					if idString, ok := idValue.(string); ok {
						responseId = idString
					}
				}
			}
		}

		model.Id = types.StringValue(responseId)
		options = clients.NewRequestOptions(nil, AsMapOfLists(model.ReadQueryParameters))
		responseBody, err = r.client.Read(ctx, fmt.Sprintf("%s/%s", model.Url.ValueString(), model.Id.ValueString()), model.ApiVersion.ValueString(), options)
		if err != nil {
			resp.Diagnostics.AddError("Failed to read data source", err.Error())
			return
		}
	}

	model.Output = types.DynamicValue(buildOutputFromBody(responseBody, model.ResponseExportValues))

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *MSGraphResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model *MSGraphResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	data, err := dynamic.ToJSON(model.Body)
	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal body", err.Error())
		return
	}
	var requestBody interface{}
	if err = json.Unmarshal(data, &requestBody); err != nil {
		resp.Diagnostics.AddError("Failed to unmarshal body", err.Error())
		return
	}

	options := clients.NewRequestOptions(nil, AsMapOfLists(model.UpdateQueryParameters))
	_, err = r.client.Update(ctx, fmt.Sprintf("%s/%s", model.Url.ValueString(), model.Id.ValueString()), model.ApiVersion.ValueString(), requestBody, options)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create resource", err.Error())
		return
	}

	options = clients.NewRequestOptions(nil, AsMapOfLists(model.ReadQueryParameters))
	responseBody, err := r.client.Read(ctx, fmt.Sprintf("%s/%s", model.Url.ValueString(), model.Id.ValueString()), model.ApiVersion.ValueString(), options)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read data source", err.Error())
		return
	}
	model.Output = types.DynamicValue(buildOutputFromBody(responseBody, model.ResponseExportValues))
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *MSGraphResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model *MSGraphResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	if model.ApiVersion.ValueString() == "" {
		model.ApiVersion = types.StringValue("v1.0")
	}

	if !strings.HasSuffix(model.Url.ValueString(), "/$ref") {
		options := clients.NewRequestOptions(nil, AsMapOfLists(model.ReadQueryParameters))
		responseBody, err := r.client.Read(ctx, fmt.Sprintf("%s/%s", model.Url.ValueString(), model.Id.ValueString()), model.ApiVersion.ValueString(), options)
		if err != nil {
			if utils.ResponseErrorWasNotFound(err) {
				tflog.Info(ctx, fmt.Sprintf("Error reading %q - removing from state", model.Id.ValueString()))
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError("Failed to read data source", err.Error())
			return
		}
		model.Output = types.DynamicValue(buildOutputFromBody(responseBody, model.ResponseExportValues))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *MSGraphResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model *MSGraphResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	var itemUrl string
	if strings.HasSuffix(model.Url.ValueString(), "/$ref") {
		itemUrl = strings.ReplaceAll(model.Url.ValueString(), "/$ref", fmt.Sprintf("/%s/$ref", model.Id.ValueString()))
	} else {
		itemUrl = fmt.Sprintf("%s/%s", model.Url.ValueString(), model.Id.ValueString())
	}

	options := clients.NewRequestOptions(nil, AsMapOfLists(model.DeleteQueryParameters))
	err := r.client.Delete(ctx, itemUrl, model.ApiVersion.ValueString(), options)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete resource", err.Error())
		return
	}
}

func (r *MSGraphResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var id, urlValue string
	parsedUrl, err := url.Parse(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse URL", err.Error())
		return
	}

	apiVersion := "v1.0"
	if parsedUrl.Query().Get("api-version") != "" {
		apiVersion = parsedUrl.Query().Get("api-version")
	}

	if strings.HasSuffix(parsedUrl.Path, "/$ref") {
		reqIdWithoutRef := strings.TrimSuffix(parsedUrl.Path, "/$ref")
		lastIndex := strings.LastIndex(reqIdWithoutRef, "/")
		id = reqIdWithoutRef[lastIndex+1:]
		urlValue = reqIdWithoutRef[0:lastIndex]
		urlValue = strings.TrimPrefix(urlValue, "/")
		urlValue = fmt.Sprintf("%s/$ref", urlValue)
	} else {
		lastIndex := strings.LastIndex(parsedUrl.Path, "/")
		id = parsedUrl.Path[lastIndex+1:]
		urlValue = strings.TrimPrefix(parsedUrl.Path[0:lastIndex], "/")
	}

	model := &MSGraphResourceModel{
		Id:                    types.StringValue(id),
		Url:                   types.StringValue(urlValue),
		ApiVersion:            types.StringValue(apiVersion),
		CreateQueryParameters: types.MapNull(types.ListType{ElemType: types.StringType}),
		UpdateQueryParameters: types.MapNull(types.ListType{ElemType: types.StringType}),
		ReadQueryParameters:   types.MapNull(types.ListType{ElemType: types.StringType}),
		DeleteQueryParameters: types.MapNull(types.ListType{ElemType: types.StringType}),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func buildOutputFromBody(body interface{}, paths map[string]string) attr.Value {
	var output interface{}
	output = make(map[string]interface{})
	for pathKey, path := range paths {
		part := utils.ExtractObjectJMES(body, pathKey, path)
		if part == nil {
			continue
		}
		output = utils.MergeObject(output, part)
	}
	data, err := json.Marshal(output)
	if err != nil {
		return nil
	}
	out, err := dynamic.FromJSONImplied(data)
	if err != nil {
		return nil
	}
	return out
}
