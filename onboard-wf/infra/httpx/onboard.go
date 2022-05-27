package httpx

import (
	"github.com/PaackEng/paackit-domain/workflows/onboard/usecase"
	"github.com/PaackEng/paackit/httpx"
	"github.com/PaackEng/paackit/response"
	"net/http"
)

const (
	OnboardPath = "onboard"
)

type OnboardDI struct {
	Middleware []httpx.Middleware
	Usecase    usecase.OnboardUsecase
}

type onboard struct {
	middleware []httpx.Middleware
	usecase    usecase.OnboardUsecase
}

func NewOnboard(di OnboardDI) httpx.Service {
	return &onboard{
		middleware: di.Middleware,
		usecase:    di.Usecase,
	}
}

func (o *onboard) Path() string {
	return OnboardPath
}

func (o *onboard) Method() httpx.Method {
	return httpx.MethodPost
}

func (o *onboard) Handler() http.HandlerFunc {
	return httpx.WithMiddleware(o.handlerFunc, o.middleware...)
}

func (o *onboard) handlerFunc(w http.ResponseWriter, r *http.Request) {

	authorization := r.Header.Get("x-authorization")

	reqByte, _ := httpx.Decode(r)

	dto := usecase.OnboardDTO{}
	err := dto.UnmarshalBinary(reqByte)
	if err != nil {
		_ = httpx.Encode(w, r, http.StatusBadRequest, response.New(nil, err))
		return
	}

	res, err := o.usecase.Onboard(dto, authorization)

	status := http.StatusCreated
	if err != nil {
		status = http.StatusInternalServerError // TODO
	}

	_ = httpx.Encode(w, r, status, response.New(res, err))
}
