package clients

import (
	"encoding/json"
	"github.com/PaackEng/paackit-domain/auth0/client"
	auth0ClientUC "github.com/PaackEng/paackit-domain/auth0/client/usecase"
	"github.com/PaackEng/paackit/httpx"
	"workflow-tests/onboard-wf/util"
)

type Auth0ServiceClient interface {
	auth0ClientUC.ClientCreateUsecase
	auth0ClientUC.ClientDeleteUsecase
}

type auth0ServiceClient struct {
	auth0ServiceUrl string
	httpClient      httpx.HTTPClient
}

func NewAuth0ServiceClient(auth0ServiceUrl string, httpClient httpx.HTTPClient) Auth0ServiceClient {
	return &auth0ServiceClient{auth0ServiceUrl: auth0ServiceUrl, httpClient: httpClient}
}

func (a auth0ServiceClient) Create(dto auth0ClientUC.ClientCreateDTO) (*client.Client, error) {
	payload, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	request := httpx.HTTPClientRequest{
		Method:  httpx.MethodPost,
		URL:     a.auth0ServiceUrl + "/clients",
		Headers: nil,
		Payload: payload,
		Retry:   false,
	}

	res, err := a.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	var client *client.Client
	err = util.ParseResponse(res, client)

	return client, err
}

func (a auth0ServiceClient) Delete(dto auth0ClientUC.ClientDeleteDTO) error {
	payload, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	request := httpx.HTTPClientRequest{
		Method:  httpx.MethodDelete,
		URL:     a.auth0ServiceUrl + "/clients",
		Headers: nil,
		Payload: payload,
		Retry:   false,
	}

	res, err := a.httpClient.Send(request)
	if err != nil {
		return err
	}

	err = util.ParseResponse(res, nil)

	return err
}
