package core

import (
	"context"

	"go.uber.org/zap"

	"github.com/5gsec/sentryflow/speculator/pkg/config"
	"github.com/5gsec/sentryflow/speculator/pkg/database"
	"github.com/5gsec/sentryflow/speculator/pkg/util"
)

type Manager struct {
	Ctx       context.Context
	Logger    *zap.SugaredLogger
	DBHandler *database.Handler
}

func Run(ctx context.Context, configFilePath string) {
	mgr := &Manager{
		Ctx:    ctx,
		Logger: util.GetLogger(),
	}

	mgr.Logger.Info("starting speculator")

	_, err := config.New(configFilePath, mgr.Logger)
	if err != nil {
		mgr.Logger.Error(err)
		return
	}

	//dbHandler, err := database.New(mgr.Ctx, cfg.Database)
	//if err != nil {
	//	mgr.Logger.Error(err)
	//	return
	//}
	//mgr.DBHandler = dbHandler
	//defer func() {
	//	if err := mgr.DBHandler.Disconnect(); err != nil {
	//		mgr.Logger.Errorf("failed to disconnect to database: %v", err)
	//	}
	//}()

}
