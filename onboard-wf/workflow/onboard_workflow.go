package workflow

import (
	"github.com/PaackEng/paackit-domain/auth0/client"
	auth0ClientUC "github.com/PaackEng/paackit-domain/auth0/client/usecase"
	"github.com/PaackEng/paackit-domain/rems"
	remsUC "github.com/PaackEng/paackit-domain/rems/usecase"
	onboardUC "github.com/PaackEng/paackit-domain/workflows/onboard/usecase"
)

var defaultGrantScope = []string{"default"}

type OMSRetailerCreateUsecase interface { // Should be in paackit-domain
	Create(id string, name string) error
}

type onboardWorkflow struct {
	labelerAudience string
	omsAudience     string

	createAccountClient       remsUC.AccountCreateUsecase
	deleteAccountClient       remsUC.AccountDeleteUsecase
	addAuth0ClientIDClient    remsUC.AccountAddAuth0ClientIDUsecase
	createServiceConfigClient remsUC.AccountServiceConfigCreateUsecase

	createAuth0ClientClient auth0ClientUC.ClientCreateUsecase
	addAuth0GrantsClient    auth0ClientUC.ClientGrantCreateUsecase
	deleteAuth0ClientClient auth0ClientUC.ClientDeleteUsecase

	createOMSRetailerClient OMSRetailerCreateUsecase
}

func NewOnboardWorkflow(
	labelerAudience string,
	omsAudience string,

	createAccountClient remsUC.AccountCreateUsecase,
	deleteAccountClient remsUC.AccountDeleteUsecase,
	addAuth0ClientIDClient remsUC.AccountAddAuth0ClientIDUsecase,
	createServiceConfigClient remsUC.AccountServiceConfigCreateUsecase,

	createAuth0ClientClient auth0ClientUC.ClientCreateUsecase,
	deleteAuth0ClientClient auth0ClientUC.ClientDeleteUsecase,
	addAuth0GrantsClient auth0ClientUC.ClientGrantCreateUsecase,

	createOMSRetailerClient OMSRetailerCreateUsecase,
) onboardUC.OnboardUsecase {
	return &onboardWorkflow{
		labelerAudience:           labelerAudience,
		omsAudience:               omsAudience,
		createAccountClient:       createAccountClient,
		createAuth0ClientClient:   createAuth0ClientClient,
		addAuth0GrantsClient:      addAuth0GrantsClient,
		addAuth0ClientIDClient:    addAuth0ClientIDClient,
		createServiceConfigClient: createServiceConfigClient,
		createOMSRetailerClient:   createOMSRetailerClient,
		deleteAccountClient:       deleteAccountClient,
		deleteAuth0ClientClient:   deleteAuth0ClientClient,
	}
}

func (o *onboardWorkflow) Onboard(dto onboardUC.OnboardDTO, authorization string) (*onboardUC.OnboardResultDTO, error) {
	account, err := o.createAccountClient.Create(dto.Account)
	if err != nil {
		return nil, err
	}

	setClientName(&dto.Auth0Client, account)

	auth0Client, err := o.createAuth0ClientClient.Create(dto.Auth0Client)
	if err != nil {
		o.rollback(account, nil)

		return nil, err
	}

	if o.labelerAudience != "" {
		go o.addAuth0GrantsClient.Create(auth0ClientUC.ClientGrantCreateDTO{
			ClientID: auth0Client.ID,
			Audience: o.labelerAudience,
			Scope:    defaultGrantScope,
		})
	}

	if o.omsAudience != "" {
		go o.addAuth0GrantsClient.Create(auth0ClientUC.ClientGrantCreateDTO{
			ClientID: auth0Client.ID,
			Audience: o.omsAudience,
			Scope:    defaultGrantScope,
		})
	}

	account, err = o.addAuth0ClientIDClient.AddAuth0ClientID(remsUC.AccountAddAuth0ClientIDDTO{
		ID:            account.ID,
		Auth0ClientID: auth0Client.ID,
	})
	if err != nil {
		o.rollback(account, auth0Client)

		return nil, err
	}

	accountServiceConfig, err := o.createServiceConfigClient.Create(getAccountSeviceConfigDTO(account.ID))
	if err != nil {
		o.rollback(account, auth0Client)

		return nil, err
	}

	err = o.createOMSRetailerClient.Create(account.ID, account.Name)

	omsResponse := "OMS Retailer Created"
	if err != nil {
		omsResponse = "OMS Retailer NOT Created"
	}

	return &onboardUC.OnboardResultDTO{
		Account:              *account,
		AccountServiceConfig: *accountServiceConfig,
		OMSResponse:          omsResponse,
	}, err
}

func setClientName(dto *auth0ClientUC.ClientCreateDTO, account *rems.Account) {
	dto.Name = "retailer-" + dto.Name + "-" + account.Country.Code + "-m2m"
}

func getAccountSeviceConfigDTO(accountID string) remsUC.AccountServiceConfigCreateDTO {
	return remsUC.AccountServiceConfigCreateDTO{
		AccountID: accountID,
	}
}

func (o *onboardWorkflow) rollback(account *rems.Account, auth0Client *client.Client) {
	if account != nil {
		go o.deleteAccountClient.Delete(remsUC.AccountDeleteDTO{ID: account.ID})
	}
	if auth0Client != nil {
		go o.deleteAuth0ClientClient.Delete(auth0ClientUC.ClientDeleteDTO{ID: auth0Client.ID})
	}
}
