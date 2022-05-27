package clients

import (
	"encoding/json"
	"github.com/PaackEng/paackit/httpx"
	"workflow-tests/onboard-wf/workflow"
)

type omsRetailerCreationDTO struct { // Should be in paackit-domain
	ID            string   `json="id"`
	Name          string   `json="name"`
	MaxAttempts   int      `json="max_attempts"`
	ServiceTypes  []string `json="service_types"`
	DeliveryTypes []string `json="delivery_types"`
}

type OMSClient interface {
	workflow.OMSRetailerCreateUsecase
}

type omsClient struct {
	omsUrl     string
	httpClient httpx.HTTPClient
}

func NewOMSClient(omsUrl string, httpClient httpx.HTTPClient) OMSClient {
	return omsClient{omsUrl: omsUrl, httpClient: httpClient}
}

func (o omsClient) Create(id string, name string) error {
	payload, err := json.Marshal(getRetailerCreationDTO(id, name))
	if err != nil {
		return err
	}

	request := httpx.HTTPClientRequest{
		Method:  httpx.MethodPost,
		URL:     o.omsUrl + "/retailers",
		Headers: nil,
		Payload: payload,
		Retry:   false,
	}

	_, err = o.httpClient.Send(request)

	return err
}

func getRetailerCreationDTO(id string, name string) omsRetailerCreationDTO {
	return omsRetailerCreationDTO{
		ID:            id,
		Name:          name,
		MaxAttempts:   3,
		ServiceTypes:  []string{"ST2", "ST4", "STZ", "SF2", "SF4", "SFZ", "NT2", "NT4", "NTH", "NTZ", "NF2", "NF4", "NFA", "XF4", "NFZ", "PT2", "PT4", "PTH", "PTA", "PTZ", "PF2", "PF4", "PFA", "PFZ", "CT2", "CT4", "CTH", "CTZ", "XFZ", "CF2", "CF4", "CFA", "CFZ", "WT2", "WT4", "WTH", "WTZ", "WF2", "WF4", "WFA", "WFZ", "XTZ", "XF2", "ST3", "DDD", "ND", "SD"},
		DeliveryTypes: []string{"direct", "b2b", "reverse"},
	}
}
