package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/opslevel-go/v2023"
	"github.com/spf13/cobra"
)

func updateCheck(input CheckInputType, usePrompts bool) (*opslevel.Check, error) {
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
		output, err = clientGQL.UpdateCheckServiceOwnership(*input.AsServiceOwnershipUpdateInput())

	case opslevel.CheckTypeHasRecentDeploy:
		output, err = clientGQL.UpdateCheckHasRecentDeploy(*input.AsHasRecentDeployUpdateInput())

	case opslevel.CheckTypeServiceProperty:
		output, err = clientGQL.UpdateCheckServiceProperty(*input.AsServicePropertyUpdateInput())

	case opslevel.CheckTypeHasServiceConfig:
		output, err = clientGQL.UpdateCheckServiceConfiguration(*input.AsServiceConfigurationUpdateInput())

	case opslevel.CheckTypeHasDocumentation:
		output, err = clientGQL.UpdateCheckHasDocumentation(*input.AsHasDocumentationUpdateInput())

	case opslevel.CheckTypeHasRepository:
		output, err = clientGQL.UpdateCheckRepositoryIntegrated(*input.AsRepositoryIntegratedUpdateInput())

	case opslevel.CheckTypeToolUsage:
		output, err = clientGQL.UpdateCheckToolUsage(*input.AsToolUsageUpdateInput())

	case opslevel.CheckTypeTagDefined:
		output, err = clientGQL.UpdateCheckTagDefined(*input.AsTagDefinedUpdateInput())

	case opslevel.CheckTypeRepoFile:
		output, err = clientGQL.UpdateCheckRepositoryFile(*input.AsRepositoryFileUpdateInput())

	case opslevel.CheckTypeRepoGrep:
		output, err = clientGQL.UpdateCheckRepositoryGrep(*input.AsRepositoryGrepUpdateInput())

	case opslevel.CheckTypeRepoSearch:
		output, err = clientGQL.UpdateCheckRepositorySearch(*input.AsRepositorySearchUpdateInput())

	case opslevel.CheckTypeManual:
		output, err = clientGQL.UpdateCheckManual(*input.AsManualUpdateInput())

	case opslevel.CheckTypeGeneric:
		opslevel.Cache.CacheIntegrations(clientGQL)
		err = input.resolveIntegrationAliases(clientGQL, usePrompts)
		cobra.CheckErr(err)
		output, err = clientGQL.UpdateCheckCustomEvent(*input.AsCustomEventUpdateInput())

	case opslevel.CheckTypeAlertSourceUsage:
		output, err = clientGQL.UpdateCheckAlertSourceUsage(*input.AsAlertSourceUsageUpdateInput())

	case opslevel.CheckTypeGitBranchProtection:
		output, err = clientGQL.UpdateCheckGitBranchProtection(*input.AsGitBranchProtectionUpdateInput())

	case opslevel.CheckTypeServiceDependency:
		output, err = clientGQL.UpdateCheckServiceDependency(*input.AsServiceDependencyUpdateInput())
	}
	cobra.CheckErr(err)
	if output == nil {
		return nil, fmt.Errorf("unknown error - no check data returned")
	}
	return output, err
}

func (self *CheckInputType) AsCustomEventUpdateInput() *opslevel.CheckCustomEventUpdateInput {
	self.Spec["resultMessage"] = self.Spec["message"]
	payload := &opslevel.CheckCustomEventUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsServiceOwnershipUpdateInput() *opslevel.CheckServiceOwnershipUpdateInput {
	payload := &opslevel.CheckServiceOwnershipUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsHasRecentDeployUpdateInput() *opslevel.CheckHasRecentDeployUpdateInput {
	payload := &opslevel.CheckHasRecentDeployUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsServicePropertyUpdateInput() *opslevel.CheckServicePropertyUpdateInput {
	payload := &opslevel.CheckServicePropertyUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsServiceConfigurationUpdateInput() *opslevel.CheckServiceConfigurationUpdateInput {
	payload := &opslevel.CheckServiceConfigurationUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsHasDocumentationUpdateInput() *opslevel.CheckHasDocumentationUpdateInput {
	payload := &opslevel.CheckHasDocumentationUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsRepositoryIntegratedUpdateInput() *opslevel.CheckRepositoryIntegratedUpdateInput {
	payload := &opslevel.CheckRepositoryIntegratedUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsToolUsageUpdateInput() *opslevel.CheckToolUsageUpdateInput {
	payload := &opslevel.CheckToolUsageUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsTagDefinedUpdateInput() *opslevel.CheckTagDefinedUpdateInput {
	payload := &opslevel.CheckTagDefinedUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsRepositoryFileUpdateInput() *opslevel.CheckRepositoryFileUpdateInput {
	payload := &opslevel.CheckRepositoryFileUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsRepositoryGrepUpdateInput() *opslevel.CheckRepositoryGrepUpdateInput {
	payload := &opslevel.CheckRepositoryGrepUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsRepositorySearchUpdateInput() *opslevel.CheckRepositorySearchUpdateInput {
	payload := &opslevel.CheckRepositorySearchUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsManualUpdateInput() *opslevel.CheckManualUpdateInput {
	payload := &opslevel.CheckManualUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsAlertSourceUsageUpdateInput() *opslevel.CheckAlertSourceUsageUpdateInput {
	payload := &opslevel.CheckAlertSourceUsageUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsGitBranchProtectionUpdateInput() *opslevel.CheckGitBranchProtectionUpdateInput {
	payload := &opslevel.CheckGitBranchProtectionUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckInputType) AsServiceDependencyUpdateInput() *opslevel.CheckServiceDependencyUpdateInput {
	payload := &opslevel.CheckServiceDependencyUpdateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}
