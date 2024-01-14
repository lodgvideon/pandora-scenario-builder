package scenario

import (
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func TestOperation(t *testing.T) {
	t.Run("Abstract Operation has name", func(t *testing.T) {
		name := "Name 134"
		build := NewAbstractOperationBuilder(name).Build()
		assert.Equal(t, build.Name(), name)
	})

	t.Run("Abstract Operation can be Executed", func(t *testing.T) {
		name := "Name 134"
		build := NewAbstractOperationBuilder(name).SetOperation(func(op *Operation, ctx *Context) *OperationResult {
			//Something
			time.Sleep(100 * time.Microsecond)
			m := map[string]string{"Key1": "val1", "Key2": "val2"}
			ctx.Vars().Put("RESULT", m)
			result := AcquireResult()
			result.Response = m
			return &OperationResult{
				Request:  "Do Something Stupid",
				Response: m,
			}
		}).Build()
		scenarioContext := NewContext()
		build.Exec(scenarioContext)

		result := build.Result()

		s := result.Req().(string)
		assert.Equal(t, s, "Do Something Stupid")
		resultMap := result.Resp().(map[string]string)

		assert.Equal(t, resultMap["Key1"], "val1")
		assert.Equal(t, resultMap["Key2"], "val2")

		build.ReleaseResult()
	})
}
