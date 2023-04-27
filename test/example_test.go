package test

import (
	"log"
	"testing"

	bdd "github.com/go-bdd/gobdd"
	"github.com/go-bdd/gobdd/context"
	. "github.com/onsi/gomega"
)

func TestScenarios(runner *testing.T) {
	suite := bdd.NewSuite(runner, bdd.WithFeaturesPath("../features/*.feature"))

	log.SetFlags(0)
	RegisterFailHandler(func(message string, callerSkip ...int) {
		log.Println(message)
	})

	suite.AddStep(
		`I add (\d+) and (\d+)`,
		func(
			test bdd.TestingT,
			ctx context.Context,
			left, right int,
		) {
			ctx.Set("sum", left+right)
		})
	suite.AddStep(
		`the result should equal (\d+)`,
		func(
			test bdd.TestingT,
			ctx context.Context,
			expected int,
		) {
			actual, err := ctx.GetInt("sum")

			if nil != err {
				test.Error(err)
			}

			if !Expect(actual).To(Equal(expected)) {
				test.Fail()
			}
		})
	suite.Run()
}
