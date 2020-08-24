package router_test

//nolint:lll
//go:generate mockgen -package router_test -destination mock_test.go -source router.go -mock_names sprintTemplateService=MockSprintTemplateService,sprintStore=MockSprintStore,taskStore=MockTaskStore
