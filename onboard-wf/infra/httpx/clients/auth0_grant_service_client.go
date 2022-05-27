package clients

import (
	"encoding/json"
	"github.com/PaackEng/paackit-domain/auth0/client"
	auth0ClientUC "github.com/PaackEng/paackit-domain/auth0/client/usecase"
	"github.com/PaackEng/paackit/httpx"
	"workflow-tests/onboard-wf/util"
)

type Auth0GrantServiceClient interface {
	auth0ClientUC.ClientGrantCreateUsecase
}

type auth0GrantServiceClient struct {
	auth0ServiceUrl string
	httpClient      httpx.HTTPClient
}

func NewAuth0GrantServiceClient(auth0ServiceUrl string, httpClient httpx.HTTPClient) Auth0GrantServiceClient {
	return auth0GrantServiceClient{auth0ServiceUrl: auth0ServiceUrl, httpClient: httpClient}
}

func (a auth0GrantServiceClient) Create(dto auth0ClientUC.ClientGrantCreateDTO) (*client.ClientGrant, error) {
	payload, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	request := httpx.HTTPClientRequest{
		Method:  httpx.MethodPost,
		URL:     a.auth0ServiceUrl + "/client-grants",
		Headers: nil,
		Payload: payload,
		Retry:   false,
	}

	res, err := a.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	var clientGrant *client.ClientGrant
	err = util.ParseResponse(res, clientGrant)

	return clientGrant, err
}
