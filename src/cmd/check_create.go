package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/opslevel-go/v2023"
	"github.com/spf13/cobra"
)

func createCheck(input CheckInputType, usePrompts bool) (*opslevel.Check, error) {
	var output *opslevel.Check
	var err error
	clientGQL := getClientGQL()
	opslevel.Cache.CacheCategories(clientGQL)
	opslevel.Cache.CacheLevels(clientGQL)
	opslevel.Cache.CacheTeams(clientGQL)
	opslevel.Cache.CacheFilters(clientGQL)
	err = input.resolveCategoryAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	err = input.resolveLevelAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	err = input.resolveTeamAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	err = input.resolveFilterAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	switch input.Kind {
	case opslevel.CheckTypeHasOwner:
		output, err = clientGQL.CreateCheckServiceOwnership(*input.AsServiceOwnershipCreateInput())

	case opslevel.CheckTypeHasRecentDeploy:
		output, err = clientGQL.CreateCheckHasRecentDeploy(*input.AsHasRecentDeployCreateInput())

	case opslevel.CheckTypeServiceProperty:
		output, err = clientGQL.CreateCheckServiceProperty(*input.AsServicePropertyCreateInput())

	case opslevel.CheckTypeHasServiceConfig:
		output, err = clientGQL.CreateCheckServiceConfiguration(*input.AsServiceConfigurationCreateInput())

	case opslevel.CheckTypeHasDocumentation:
		output, err = clientGQL.CreateCheckHasDocumentation(*input.AsHasDocumentationCreateInput())

	case opslevel.CheckTypeHasRepository:
		output, err = clientGQL.CreateCheckRepositoryIntegrated(*input.AsRepositoryIntegratedCreateInput())

	case opslevel.CheckTypeToolUsage:
		output, err = clientGQL.CreateCheckToolUsage(*input.AsToolUsageCreateInput())

	case opslevel.CheckTypeTagDefined:
		output, err = clientGQL.CreateCheckTagDefined(*input.AsTagDefinedCreateInput())

	case opslevel.CheckTypeRepoFile:
		output, err = clientGQL.CreateCheckRepositoryFile(*input.AsRepositoryFileCreateInput())

	case opslevel.CheckTypeRepoGrep:
		output, err = clientGQL.CreateCheckRepositoryGrep(*input.AsRepositoryGrepCreateInput())

	case opslevel.CheckTypeRepoSearch:
		output, err = clientGQL.CreateCheckRepositorySearch(*input.AsRepositorySearchCreateInput())

	case opslevel.CheckTypeManual:
		output, err = clientGQL.CreateCheckManual(*input.AsManualCreateInput())

	case opslevel.CheckTypeGeneric:
		opslevel.Cache.CacheIntegrations(clientGQL)
		err = input.resolveIntegrationAliases(clientGQL, usePrompts)
		cobra.CheckErr(err)
		output, err = clientGQL.CreateCheckCustomEvent(*input.AsCustomEventCreateInput())

	case opslevel.CheckTypeAlertSourceUsage:
		output, err = clientGQL.CreateCheckAlertSourceUsage(*input.AsAlertSourceUsageCreateInput())

	case opslevel.CheckTypeGitBranchProtection:
		output, err = clientGQL.CreateCheckGitBranchProtection(*input.AsGitBranchProtectionCreateInput())

	case opslevel.CheckTypeServiceDependency:
		output, err = clientGQL.CreateCheckServiceDependency(*input.AsServiceDependencyCreateInput())
	}
	cobra.CheckErr(err)
	if output == nil {
		return nil, fmt.Errorf("unknown error - no check data returned")
	}
	return output, err
}

func (self *CheckInputType) AsCustomEventCreateInput() *opslevel.CheckCustomEventCreateInput {
	self.Spec["resultMessage"] = self.Spec["message"]
	payload := &opslevel.CheckCustomEventCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsServiceOwnershipCreateInput() *opslevel.CheckServiceOwnershipCreateInput {
	payload := &opslevel.CheckServiceOwnershipCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsHasRecentDeployCreateInput() *opslevel.CheckHasRecentDeployCreateInput {
	payload := &opslevel.CheckHasRecentDeployCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsServicePropertyCreateInput() *opslevel.CheckServicePropertyCreateInput {
	payload := &opslevel.CheckServicePropertyCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsServiceConfigurationCreateInput() *opslevel.CheckServiceConfigurationCreateInput {
	payload := &opslevel.CheckServiceConfigurationCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsHasDocumentationCreateInput() *opslevel.CheckHasDocumentationCreateInput {
	payload := &opslevel.CheckHasDocumentationCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsRepositoryIntegratedCreateInput() *opslevel.CheckRepositoryIntegratedCreateInput {
	payload := &opslevel.CheckRepositoryIntegratedCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsToolUsageCreateInput() *opslevel.CheckToolUsageCreateInput {
	payload := &opslevel.CheckToolUsageCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsTagDefinedCreateInput() *opslevel.CheckTagDefinedCreateInput {
	payload := &opslevel.CheckTagDefinedCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsRepositoryFileCreateInput() *opslevel.CheckRepositoryFileCreateInput {
	payload := &opslevel.CheckRepositoryFileCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsRepositoryGrepCreateInput() *opslevel.CheckRepositoryGrepCreateInput {
	payload := &opslevel.CheckRepositoryGrepCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsRepositorySearchCreateInput() *opslevel.CheckRepositorySearchCreateInput {
	payload := &opslevel.CheckRepositorySearchCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsManualCreateInput() *opslevel.CheckManualCreateInput {
	payload := &opslevel.CheckManualCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsAlertSourceUsageCreateInput() *opslevel.CheckAlertSourceUsageCreateInput {
	payload := &opslevel.CheckAlertSourceUsageCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsGitBranchProtectionCreateInput() *opslevel.CheckGitBranchProtectionCreateInput {
	payload := &opslevel.CheckGitBranchProtectionCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsServiceDependencyCreateInput() *opslevel.CheckServiceDependencyCreateInput {
	payload := &opslevel.CheckServiceDependencyCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}
