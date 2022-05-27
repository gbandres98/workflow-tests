package clients

import (
	"encoding/json"
	"github.com/PaackEng/paackit-domain/rems"
	remsUC "github.com/PaackEng/paackit-domain/rems/usecase"
	"github.com/PaackEng/paackit/httpx"
	"workflow-tests/onboard-wf/util"
)

type RetailerServiceClient interface {
	remsUC.AccountCreateUsecase
	remsUC.AccountAddAuth0ClientIDUsecase
	remsUC.AccountDeleteUsecase
}

type retailerServiceClient struct {
	retailerServiceUrl string
	httpClient         httpx.HTTPClient
}

func NewRetailerServiceClient(retailerServiceUrl string, httpClient httpx.HTTPClient) RetailerServiceClient {
	return &retailerServiceClient{
		retailerServiceUrl: retailerServiceUrl,
		httpClient:         httpClient,
	}
}

func (c *retailerServiceClient) Create(dto remsUC.AccountCreateDTO) (*rems.Account, error) {
	payload, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	request := httpx.HTTPClientRequest{
		Method:  httpx.MethodPost,
		URL:     c.retailerServiceUrl + "/accounts",
		Headers: nil,
		Payload: payload,
		Retry:   false,
	}

	res, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	var account *rems.Account
	err = util.ParseResponse(res, account)

	return account, err
}

func (c *retailerServiceClient) AddAuth0ClientID(dto remsUC.AccountAddAuth0ClientIDDTO) (*rems.Account, error) {
	payload, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	request := httpx.HTTPClientRequest{
		Method:  httpx.MethodPost,
		URL:     c.retailerServiceUrl + "/accounts/" + dto.ID + "/auth0ClientID",
		Headers: nil,
		Payload: payload,
		Retry:   false,
	}

	res, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	var account *rems.Account
	err = util.ParseResponse(res, account)

	return account, err
}

func (c *retailerServiceClient) Delete(dto remsUC.AccountDeleteDTO) (*rems.Account, error) {
	payload, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	request := httpx.HTTPClientRequest{
		Method:  httpx.MethodDelete,
		URL:     c.retailerServiceUrl + "/accounts",
		Headers: nil,
		Payload: payload,
		Retry:   false,
	}

	res, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	err = util.ParseResponse(res, nil)

	return nil, err
}
