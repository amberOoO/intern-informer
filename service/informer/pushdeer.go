package informer

import (
	"internInformer/utils"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

type PushdeerInformer struct {
	logger    *zap.Logger
	pushkey   string
	serverUrl string
}

func NewDefaultPushdeer() *PushdeerInformer {
	vi := utils.GetDefaultViper()
	cfg := vi.GetStringMapString("pushdeer")
	return &PushdeerInformer{
		logger:    utils.GetZapLogger(zap.String("service", "informer"), zap.String("informer", "pushdeer")),
		pushkey:   cfg["pushkey"],
		serverUrl: cfg["url"],
	}
}

func (pi *PushdeerInformer) Send(msg string) error {
	parsedUrl, err := url.Parse(pi.serverUrl)
	if err != nil {
		pi.logger.Error("parse url failed", zap.Error(err))
		return err
	}

	params := url.Values{}
	params.Add("pushkey", pi.pushkey)
	params.Add("text", msg)

	parsedUrl.RawQuery = params.Encode()
	resp, err := http.Get(parsedUrl.String())
	s, _ := io.ReadAll(resp.Body)
	// 记录一下返回的response
	pi.logger.Info(string(s))
	if err != nil {
		pi.logger.Error("failed to send request", zap.Error(err))
	}
	return err
}
