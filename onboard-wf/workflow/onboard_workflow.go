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

	createAccountUsecase       remsUC.AccountCreateUsecase
	deleteAccountUsecase       remsUC.AccountDeleteUsecase
	addAuth0ClientIDUsecase    remsUC.AccountAddAuth0ClientIDUsecase
	createServiceConfigUsecase remsUC.AccountServiceConfigCreateUsecase

	createAuth0ClientUsecase auth0ClientUC.ClientCreateUsecase
	addAuth0GrantsUsecase    auth0ClientUC.ClientGrantCreateUsecase
	deleteAuth0ClientUsecase auth0ClientUC.ClientDeleteUsecase

	createOMSRetailerUsecase OMSRetailerCreateUsecase
}

func NewOnboardWorkflow(
	labelerAudience string,
	omsAudience string,

	createAccountUsecase remsUC.AccountCreateUsecase,
	deleteAccountUsecase remsUC.AccountDeleteUsecase,
	addAuth0ClientIDUsecase remsUC.AccountAddAuth0ClientIDUsecase,
	createServiceConfigUsecase remsUC.AccountServiceConfigCreateUsecase,

	createAuth0ClientUsecase auth0ClientUC.ClientCreateUsecase,
	deleteAuth0ClientUsecase auth0ClientUC.ClientDeleteUsecase,
	addAuth0GrantsUsecase auth0ClientUC.ClientGrantCreateUsecase,

	createOMSRetailerClient OMSRetailerCreateUsecase,
) onboardUC.OnboardUsecase {
	return &onboardWorkflow{
		labelerAudience:            labelerAudience,
		omsAudience:                omsAudience,
		createAccountUsecase:       createAccountUsecase,
		createAuth0ClientUsecase:   createAuth0ClientUsecase,
		addAuth0GrantsUsecase:      addAuth0GrantsUsecase,
		addAuth0ClientIDUsecase:    addAuth0ClientIDUsecase,
		createServiceConfigUsecase: createServiceConfigUsecase,
		createOMSRetailerUsecase:   createOMSRetailerClient,
		deleteAccountUsecase:       deleteAccountUsecase,
		deleteAuth0ClientUsecase:   deleteAuth0ClientUsecase,
	}
}

func (o *onboardWorkflow) Onboard(dto onboardUC.OnboardDTO, authorization string) (*onboardUC.OnboardResultDTO, error) {
	account, err := o.createAccountUsecase.Create(dto.Account)
	if err != nil {
		return nil, err
	}

	setClientName(&dto.Auth0Client, account)

	auth0Client, err := o.createAuth0ClientUsecase.Create(dto.Auth0Client)
	if err != nil {
		o.rollback(account, nil)

		return nil, err
	}

	if o.labelerAudience != "" {
		go o.addAuth0GrantsUsecase.Create(auth0ClientUC.ClientGrantCreateDTO{
			ClientID: auth0Client.ID,
			Audience: o.labelerAudience,
			Scope:    defaultGrantScope,
		})
	}

	if o.omsAudience != "" {
		go o.addAuth0GrantsUsecase.Create(auth0ClientUC.ClientGrantCreateDTO{
			ClientID: auth0Client.ID,
			Audience: o.omsAudience,
			Scope:    defaultGrantScope,
		})
	}

	account, err = o.addAuth0ClientIDUsecase.AddAuth0ClientID(remsUC.AccountAddAuth0ClientIDDTO{
		ID:            account.ID,
		Auth0ClientID: auth0Client.ID,
	})
	if err != nil {
		o.rollback(account, auth0Client)

		return nil, err
	}

	accountServiceConfig, err := o.createServiceConfigUsecase.Create(getAccountSeviceConfigDTO(account.ID))
	if err != nil {
		o.rollback(account, auth0Client)

		return nil, err
	}

	err = o.createOMSRetailerUsecase.Create(account.ID, account.Name)

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
		go o.deleteAccountUsecase.Delete(remsUC.AccountDeleteDTO{ID: account.ID})
	}
	if auth0Client != nil {
		go o.deleteAuth0ClientUsecase.Delete(auth0ClientUC.ClientDeleteDTO{ID: auth0Client.ID})
	}
}
