package service

import (
	"context"
	"fmt"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/validate"
)

func (a *App) Analyze(ctx context.Context, springID model.ID, from, to time.Time) (model.AnalysisSummary, error) {
	if err := validate.AnalysisWindow(from, to); err != nil {
		return model.AnalysisSummary{}, err
	}
	exists, err := a.SpringExists(ctx, springID)
	if err != nil {
		return model.AnalysisSummary{}, err
	}
	if !exists {
		return model.AnalysisSummary{}, fmt.Errorf("spring %d does not exist", springID)
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	result, err := a.store.Analysis(dbctx, springID, from, to)
	if err != nil {
		return model.AnalysisSummary{}, a.dbError("analyze spring", err)
	}
	return result, nil
}

func (a *App) DefaultAnalysis(ctx context.Context, springID model.ID) (model.AnalysisSummary, error) {
	to := a.Now()
	return a.Analyze(ctx, springID, to.Add(-24*time.Hour), to)
}
