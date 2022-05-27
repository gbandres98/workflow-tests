package main

import (
	"context"
	"fmt"
	"github.com/PaackEng/paackit"
	"github.com/PaackEng/paackit/config"
	infraHTTPX "github.com/PaackEng/paackit/httpx"
	"github.com/PaackEng/paackit/paack"
	"os"
	"workflow-tests/onboard-wf/infra/httpx"
	"workflow-tests/onboard-wf/infra/httpx/clients"
	"workflow-tests/onboard-wf/util"
	"workflow-tests/onboard-wf/workflow"
)

const (
	serviceName                  = "OnboardWorkflow"
	envRetailerServiceUrl        = "RETAILER_SERVICE_URL"
	envRetailerAccountConfigUrl  = "RETAILER_ACCOUNT_CONFIG_URL"
	envAuth0ManagementServiceUrl = "AUTH0_MGMT_PROXY_SERVICE_URL"
	envOMSServiceUrl             = "OMS_SERVICE_URL"
	envLabelerAudience           = "AUTH0_LABELER_AUDIENCE"
	envOMSAudience               = "AUTH0_OMS_AUDIENCE"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	httpClient := infraHTTPX.NewClient(infraHTTPX.HttpClientConfig{})

	retailerServiceClient := clients.NewRetailerServiceClient(
		util.GetEnvOrPanic(envRetailerServiceUrl),
		httpClient)

	serviceConfigClient := clients.NewServiceConfigClient(
		util.GetEnvOrPanic(envRetailerAccountConfigUrl),
		httpClient)

	auth0ServiceClient := clients.NewAuth0ServiceClient(
		util.GetEnvOrPanic(envAuth0ManagementServiceUrl),
		httpClient)

	auth0GrantServiceClient := clients.NewAuth0GrantServiceClient(
		util.GetEnvOrPanic(envAuth0ManagementServiceUrl),
		httpClient)

	omsClient := clients.NewOMSClient(
		util.GetEnvOrPanic(envOMSServiceUrl),
		httpClient)

	workflow := workflow.NewOnboardWorkflow(
		config.GetEnvOrDefault(envLabelerAudience, ""),
		config.GetEnvOrDefault(envOMSAudience, ""),

		retailerServiceClient,
		retailerServiceClient,
		retailerServiceClient,

		serviceConfigClient,

		auth0ServiceClient,
		auth0ServiceClient,
		auth0GrantServiceClient,

		omsClient,
	)

	http := infraHTTPX.NewHTTPX(infraHTTPX.HttpConfig{
		BasePath: "/api/v3/",
		Port:     "8080",
	})

	onboardHTTP := httpx.NewOnboard(httpx.OnboardDI{
		Middleware: []infraHTTPX.Middleware{},
		Usecase:    workflow,
	})

	http.Register([]infraHTTPX.Service{
		onboardHTTP,
	})

	p := paackit.New(paackit.Config{
		Ctx:  ctx,
		Name: serviceName,
	})

	err := p.RegisterTransporter([]paack.Transporter{
		http,
	})
	if err != nil {
		return err
	}

	err = p.Start()
	if err != nil {
		return err
	}

	return nil
}
