package main

import (
	"go.uber.org/zap"
	"time"
)
import "pandoraScript/scenario"

func main() {

	production, _ := zap.NewProduction()

	scenarioBuilder := scenario.NewScenarioBuilder("My First scenario")
	scenarioBuilder.SetLogger(production)

	scenarioBuilder.
		AddOperation(scenario.NewOperationBuilder("MyFirstOperation").
			BeforeOperation(func(op *scenario.Operation, ctx *scenario.Context) error {
				//Тут можно разместить логику по обработки Переод каждой операцией;
				// Операций Before может быть несколько
				ctx.L().Info("Before Operation", zap.String("Operation", op.Name()))
				return nil
			}).
			SetOperation(func(op *scenario.Operation, ctx *scenario.Context) *scenario.OperationResult {
				//Тут любая Кастомная операция
				ctx.L().Info("Running Operation", zap.String("name", op.Name()))
				ctx.L().Info("Pause For 10 sec", zap.String("name", op.Name()))
				ctx.L().Info("Pause For 10 sec", zap.String("name", op.Name()))

				time.Sleep(10 * time.Second)
				ctx.L().Info("Finish Operation", zap.String("name", op.Name()))
				return &scenario.OperationResult{
					Request:  "Request",
					Response: "Response",
				}
			}).
			After(func(op *scenario.Operation, ctx *scenario.Context) error {
				ctx.Log().Info("After Operation", zap.String("operation", op.Name()), zap.String("prev", ctx.Prev().Res().(string)))

				nextOperation := op.Next()
				ctx.Log().Info("Next Op have to be", zap.String("operation", nextOperation.Name()), zap.String("prev", ctx.Prev().Res().(string)))

				//Любая активность после операции.
				// Тут можно провести валидацию операции
				return nil
			}).Build())

	s := scenarioBuilder.Build()
	l, _ := zap.NewProduction()
	//Сценарный контекст - это контекст исполнения текущей итерации сценария
	newContext := scenario.NewContext(l)

	s.Run(newContext)

}
