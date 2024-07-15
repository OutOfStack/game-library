package facade

import (
	"context"
	"fmt"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
)

// GetCompanies returns companies
func (p *Provider) GetCompanies(ctx context.Context) ([]model.Company, error) {
	list := make([]model.Company, 0)
	err := cache.Get(ctx, p.cache, getCompaniesKey(), &list, func() ([]model.Company, error) {
		return p.storage.GetCompanies(ctx)
	}, 0)
	if err != nil {
		return nil, fmt.Errorf("get companies: %v", err)
	}

	return list, nil
}

// GetCompaniesMap returns all companies map
func (p *Provider) GetCompaniesMap(ctx context.Context) (map[int32]model.Company, error) {
	companies, err := p.GetCompanies(ctx)
	if err != nil {
		return nil, fmt.Errorf("get companies: %v", err)
	}

	m := make(map[int32]model.Company, len(companies))
	for _, c := range companies {
		m[c.ID] = c
	}

	return m, nil
}

// GetTopCompanies returns top companies by type
func (p *Provider) GetTopCompanies(ctx context.Context, companyType string, limit int64) ([]model.Company, error) {
	list := make([]model.Company, 0)
	err := cache.Get(ctx, p.cache, getTopCompaniesKey(companyType, limit), &list, func() ([]model.Company, error) {
		switch companyType {
		case model.CompanyTypeDeveloper:
			return p.storage.GetTopDevelopers(ctx, limit)
		case model.CompanyTypePublisher:
			return p.storage.GetTopPublishers(ctx, limit)
		}
		return nil, fmt.Errorf("unsopported company type: %s", companyType)
	}, 0)
	if err != nil {
		return nil, fmt.Errorf("get top companies: %v", err)
	}

	return list, nil
}

// GetCompanyByID returns company by id
func (p *Provider) GetCompanyByID(ctx context.Context, id int32) (model.Company, error) {
	company, err := p.storage.GetCompanyByID(ctx, id)
	if err != nil {
		if apperr.IsStatusCode(err, apperr.NotFound) {
			return model.Company{}, err
		}
		return model.Company{}, fmt.Errorf("get company by id %d: %v", id, err)
	}

	return company, nil
}
