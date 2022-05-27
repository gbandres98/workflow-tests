package clients

import (
	"encoding/json"
	"github.com/PaackEng/paackit-domain/rems"
	remsUC "github.com/PaackEng/paackit-domain/rems/usecase"
	"github.com/PaackEng/paackit/httpx"
	"workflow-tests/onboard-wf/util"
)

type ServiceConfigClient interface {
	remsUC.AccountServiceConfigCreateUsecase
}

type serviceConfigClient struct {
	serviceConfigServiceUrl string
	httpClient              httpx.HTTPClient
}

func NewServiceConfigClient(serviceConfigServiceUrl string, httpClient httpx.HTTPClient) ServiceConfigClient {
	return &serviceConfigClient{serviceConfigServiceUrl: serviceConfigServiceUrl, httpClient: httpClient}
}

func (s *serviceConfigClient) Create(dto remsUC.AccountServiceConfigCreateDTO) (*rems.AccountServiceConfig, error) {
	payload, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	request := httpx.HTTPClientRequest{
		Method:  httpx.MethodPost,
		URL:     s.serviceConfigServiceUrl + "/RemsAccountServiceConfig",
		Headers: nil,
		Payload: payload,
		Retry:   false,
	}

	res, err := s.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	var accountServiceConfig *rems.AccountServiceConfig
	err = util.ParseResponse(res, accountServiceConfig)

	return accountServiceConfig, err
}
