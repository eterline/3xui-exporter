package server

import "net/http"

type MetricsServer struct {
	srv *http.Server
}

func NewMetricsServer(h http.Handler, addr string) *MetricsServer {
	metrics := &MetricsServer{
		srv: &http.Server{
			Addr:    addr,
			Handler: h,
		},
	}

	return metrics
}

func (ms *MetricsServer) Listen(cert, key string) error {
	if cert == "" || key == "" {
		return ms.srv.ListenAndServe()
	}
	return ms.srv.ListenAndServeTLS(cert, key)
}

func (ms *MetricsServer) Stop() {
	ms.srv.Close()
}
