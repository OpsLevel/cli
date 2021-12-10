package common

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gosimple/slug"
	"github.com/opslevel/opslevel-go"
	"github.com/rs/zerolog/log"
)

type AliasCacher struct {
	mutex        sync.Mutex
	Tiers        map[string]opslevel.Tier
	Lifecycles   map[string]opslevel.Lifecycle
	Teams        map[string]opslevel.Team
	Categories   map[string]opslevel.Category
	Levels       map[string]opslevel.Level
	Filters      map[string]opslevel.Filter
	Integrations map[string]opslevel.Integration
}

func (c *AliasCacher) TryGetTier(alias string) (*opslevel.Tier, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Tiers[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *AliasCacher) TryGetLifecycle(alias string) (*opslevel.Lifecycle, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Lifecycles[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *AliasCacher) TryGetTeam(alias string) (*opslevel.Team, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Teams[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *AliasCacher) TryGetCategory(alias string) (*opslevel.Category, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Categories[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *AliasCacher) TryGetLevel(alias string) (*opslevel.Level, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Levels[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *AliasCacher) TryGetFilter(alias string) (*opslevel.Filter, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Filters[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *AliasCacher) TryGetIntegration(alias string) (*opslevel.Integration, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Integrations[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *AliasCacher) doCacheTiers(client *opslevel.Client) {
	log.Info().Msg("Caching 'Tier' lookup table from OpsLevel API ...")

	data, dataErr := client.ListTiers()
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Tier' from OpsLevel API - REASON: %s", dataErr.Error())
	}
	for _, item := range data {
		c.Tiers[string(item.Alias)] = item
	}
}

func (c *AliasCacher) doCacheLifecycles(client *opslevel.Client) {
	log.Info().Msg("Caching 'Lifecycle' lookup table from OpsLevel API ...")

	data, dataErr := client.ListLifecycles()
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Lifecycle' from OpsLevel API - REASON: %s", dataErr.Error())
	}
	for _, item := range data {
		c.Lifecycles[string(item.Alias)] = item
	}
}

func (c *AliasCacher) doCacheTeams(client *opslevel.Client) {
	log.Info().Msg("Caching 'Team' lookup table from OpsLevel API ...")

	data, dataErr := client.ListTeams()
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Team' from OpsLevel API - REASON: %s", dataErr.Error())
	}

	for _, item := range data {
		c.Teams[string(item.Alias)] = item
	}
}

func (c *AliasCacher) doCacheCategories(client *opslevel.Client) {
	log.Info().Msg("Caching 'Category' lookup table from OpsLevel API ...")

	data, dataErr := client.ListCategories()
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Category' from OpsLevel API - REASON: %s", dataErr.Error())
	}

	for _, item := range data {
		// TODO: Categories need an alias
		key := strings.ToLower(string(item.Name))
		c.Categories[key] = item
	}
}

func (c *AliasCacher) doCacheLevels(client *opslevel.Client) {
	log.Info().Msg("Caching 'Level' lookup table from OpsLevel API ...")

	data, dataErr := client.ListLevels()
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Level' from OpsLevel API - REASON: %s", dataErr.Error())
	}

	for _, item := range data {
		c.Levels[string(item.Alias)] = item
	}
}

func (c *AliasCacher) doCacheFilters(client *opslevel.Client) {
	log.Info().Msg("Caching 'Filter' lookup table from OpsLevel API ...")

	data, dataErr := client.ListFilters()
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Filter' from OpsLevel API - REASON: %s", dataErr.Error())
	}

	for _, item := range data {
		c.Filters[slug.Make(item.Name)] = item
	}
}

func (c *AliasCacher) doCacheIntegrations(client *opslevel.Client) {
	log.Info().Msg("Caching 'Integration' lookup table from OpsLevel API ...")

	data, dataErr := client.ListIntegrations()
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Integration' from OpsLevel API - REASON: %s", dataErr.Error())
	}

	for _, item := range data {
		c.Integrations[fmt.Sprintf("%s-%s", slug.Make(item.Type), slug.Make(item.Name))] = item
	}
}

func (c *AliasCacher) CacheTiers(client *opslevel.Client) {
	c.mutex.Lock()
	c.doCacheTiers(client)
	c.mutex.Unlock()
}

func (c *AliasCacher) CacheLifecycles(client *opslevel.Client) {
	c.mutex.Lock()
	c.doCacheLifecycles(client)
	c.mutex.Unlock()
}

func (c *AliasCacher) CacheTeams(client *opslevel.Client) {
	c.mutex.Lock()
	c.doCacheTeams(client)
	c.mutex.Unlock()
}

func (c *AliasCacher) CacheCategories(client *opslevel.Client) {
	c.mutex.Lock()
	c.doCacheCategories(client)
	c.mutex.Unlock()
}

func (c *AliasCacher) CacheLevels(client *opslevel.Client) {
	c.mutex.Lock()
	c.doCacheLevels(client)
	c.mutex.Unlock()
}

func (c *AliasCacher) CacheFilters(client *opslevel.Client) {
	c.mutex.Lock()
	c.doCacheFilters(client)
	c.mutex.Unlock()
}

func (c *AliasCacher) CacheIntegrations(client *opslevel.Client) {
	c.mutex.Lock()
	c.doCacheIntegrations(client)
	c.mutex.Unlock()
}

func (c *AliasCacher) CacheAll(client *opslevel.Client) {
	c.mutex.Lock()
	c.doCacheTiers(client)
	c.doCacheLifecycles(client)
	c.doCacheTeams(client)
	c.doCacheCategories(client)
	c.doCacheLevels(client)
	c.doCacheFilters(client)
	c.doCacheIntegrations(client)
	c.mutex.Unlock()
}

func createAliasCache() *AliasCacher {
	return &AliasCacher{
		mutex:        sync.Mutex{},
		Tiers:        make(map[string]opslevel.Tier),
		Lifecycles:   make(map[string]opslevel.Lifecycle),
		Teams:        make(map[string]opslevel.Team),
		Categories:   make(map[string]opslevel.Category),
		Levels:       make(map[string]opslevel.Level),
		Filters:      make(map[string]opslevel.Filter),
		Integrations: make(map[string]opslevel.Integration),
	}
}

var AliasCache = createAliasCache()
