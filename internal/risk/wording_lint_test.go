package risk

import (
	"testing"

	"github.com/iokuru/sentry/internal/model"
)

func TestValidateResultWording_AllCategoriesAndStates(t *testing.T) {
	categories := []model.Category{
		model.CategoryBank,
		model.CategoryDelivery,
		model.CategoryUtility,
		model.CategoryPrize,
		model.CategoryTechSupport,
		model.CategoryRelativeInNeed,
		model.CategoryRomance,
		model.CategoryTax,
		model.CategoryOther,
	}

	levels := []model.RiskLevel{
		model.RiskLow,
		model.RiskMedium,
		model.RiskHigh,
	}

	for _, cat := range categories {
		for _, lvl := range levels {
			t.Run(string(cat)+"_"+string(lvl), func(t *testing.T) {
				out := BuildSeniorOutput(cat, lvl, "US", false)
				errs := ValidateResultWording(out)
				if len(errs) > 0 {
					for _, e := range errs {
						t.Errorf("wording lint error for cat=%s, lvl=%s: %s", cat, lvl, e)
					}
				}
			})
		}
	}

	t.Run("damage_control", func(t *testing.T) {
		out := BuildSeniorOutput(model.CategoryBank, model.RiskHigh, "US", true)
		errs := ValidateResultWording(out)
		if len(errs) > 0 {
			for _, e := range errs {
				t.Errorf("damage control wording lint error: %s", e)
			}
		}
	})
}
