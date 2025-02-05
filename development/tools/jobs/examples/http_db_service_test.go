package examples_test

import (
	"testing"

	"github.com/kyma-project/test-infra/development/tools/jobs/tester"
	"github.com/kyma-project/test-infra/development/tools/jobs/tester/preset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpDbServiceJobsPresubmit(t *testing.T) {
	// WHEN
	jobConfig, err := tester.ReadJobConfig("../../../../prow/jobs/examples/examples-buildpack.yaml")
	// THEN
	require.NoError(t, err)

	assert.Len(t, jobConfig.PresubmitsStatic, 1)
	kymaPresubmits := jobConfig.AllStaticPresubmits([]string{"kyma-project/examples"})
	expName := "pre-main-examples-http-db-service"

	actualPresubmit := tester.FindPresubmitJobByName(kymaPresubmits, expName)
	assert.Equal(t, expName, actualPresubmit.Name)
	assert.Equal(t, []string{"^master$", "^main$"}, actualPresubmit.Branches)
	assert.Equal(t, 10, actualPresubmit.MaxConcurrency)
	assert.False(t, actualPresubmit.SkipReport)
	assert.True(t, actualPresubmit.Decorate)
	tester.AssertThatHasExtraRefTestInfra(t, actualPresubmit.JobBase.UtilityConfig, "main")
	tester.AssertThatHasPresets(t, actualPresubmit.JobBase, preset.DindEnabled, preset.DockerPushRepoKyma, preset.GcrPush)
	assert.Equal(t, "^http-db-service/", actualPresubmit.RunIfChanged)
	assert.Equal(t, tester.ImageGolangBuildpack1_14, actualPresubmit.Spec.Containers[0].Image)
	assert.Equal(t, []string{"/home/prow/go/src/github.com/kyma-project/test-infra/prow/scripts/build-generic.sh"}, actualPresubmit.Spec.Containers[0].Command)
	assert.Equal(t, []string{"/home/prow/go/src/github.com/kyma-project/examples/http-db-service", "ci-pr"}, actualPresubmit.Spec.Containers[0].Args)
}

func TestHttpDbServiceJobPostsubmit(t *testing.T) {
	// WHEN
	jobConfig, err := tester.ReadJobConfig("../../../../prow/jobs/examples/examples-buildpack.yaml")
	// THEN
	require.NoError(t, err)

	assert.Len(t, jobConfig.PostsubmitsStatic, 1)
	kymaPost := jobConfig.AllStaticPostsubmits([]string{"kyma-project/examples"})
	expName := "post-main-examples-http-db-service"

	actualPost := tester.FindPostsubmitJobByName(kymaPost, expName)
	assert.Equal(t, expName, actualPost.Name)
	assert.Equal(t, []string{"^master$", "^main$"}, actualPost.Branches)

	assert.Equal(t, 10, actualPost.MaxConcurrency)
	assert.True(t, actualPost.Decorate)
	tester.AssertThatHasExtraRefTestInfra(t, actualPost.JobBase.UtilityConfig, "main")
	tester.AssertThatHasPresets(t, actualPost.JobBase, preset.DindEnabled, preset.DockerPushRepoKyma, preset.GcrPush)
	assert.Equal(t, "^http-db-service/", actualPost.RunIfChanged)
	assert.Equal(t, tester.ImageGolangBuildpack1_14, actualPost.Spec.Containers[0].Image)
	assert.Equal(t, []string{"/home/prow/go/src/github.com/kyma-project/test-infra/prow/scripts/build-generic.sh"}, actualPost.Spec.Containers[0].Command)
	assert.Equal(t, []string{"/home/prow/go/src/github.com/kyma-project/examples/http-db-service", "ci-main"}, actualPost.Spec.Containers[0].Args)
}
