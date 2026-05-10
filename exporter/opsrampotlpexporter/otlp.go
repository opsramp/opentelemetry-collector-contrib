// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opsrampotlpexporter // import "go.opentelemetry.io/collector/exporter/otlpexporter"

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "github.com/opsramp/go-proxy-dialer/connect" // implemetation for http connect proxy
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	"go.uber.org/zap"
	"golang.org/x/net/http/httpproxy"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	tokenRenewInProgress bool
	credentials          Credentials
	tokenMutex           sync.Mutex // Add mutex to protect global variables
)

type opsrampOTLPExporter struct {
	// Input configuration.
	config *Config

	// gRPC clients and connection.
	traceExporter  ptraceotlp.GRPCClient
	metricExporter pmetricotlp.GRPCClient
	logExporter    plogotlp.GRPCClient
	clientConn     *grpc.ClientConn
	metadata       metadata.MD
	callOptions    []grpc.CallOption

	settings component.TelemetrySettings
	mut      sync.Mutex

	// Default user-agent header.
	userAgent   string
	accessToken string

	//logger *zap.Logger
}

// Crete new exporter and start it. The exporter will begin connecting, but
// this function may return before the connection is established.
func newExporter(cfg component.Config, set exporter.Settings) (*opsrampOTLPExporter, error) {
	oCfg := cfg.(*Config)

	accessToken, err := getAuthToken(oCfg.Security)
	if err != nil {
		tokenRenewInProgress = false
		return nil, fmt.Errorf("access token isn't available: %w", err)
	}
	tokenRenewInProgress = false

	if oCfg.Endpoint == "" {
		return nil, errors.New("OTLP exporter config requires an Endpoint")
	}

	userAgent := fmt.Sprintf("%s/%s (%s/%s)",
		set.BuildInfo.Description, set.BuildInfo.Version, runtime.GOOS, runtime.GOARCH)

	// REMOVE after testing...
	//logger := initLogger()

	return &opsrampOTLPExporter{config: oCfg, settings: set.TelemetrySettings,
		userAgent: userAgent, accessToken: accessToken}, nil
}

type Creds struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func getAuthToken(cfg SecuritySettings) (string, error) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	if tokenRenewInProgress {
		for tokenRenewInProgress {
			tokenMutex.Unlock()
			time.Sleep(time.Second * 1) // Reduced from 10s to 1s for better performance
			tokenMutex.Lock()
		}
		// Return the refreshed token, don't set tokenRenewInProgress again
		return credentials.AccessToken, nil
	}
	tokenRenewInProgress = true

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return httpproxy.FromEnvironment().ProxyFunc()(req.URL)
			},
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	data := url.Values{}
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("grant_type", grantType)
	request, err := http.NewRequest("POST", cfg.OAuthServiceURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		tokenRenewInProgress = false
		return "", err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(request)
	if err != nil {
		if strings.Contains(err.Error(), "x509: certificate signed by unknown authority") || strings.Contains(err.Error(), "TLS handshake timeout") {
			// If the error is due to an untrusted certificate, we can try to get the token with TLS verification disabled.
			return getAuthTokenWithTlsDisabled(cfg)
		}
		return "", err
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil {
		tokenRenewInProgress = false
		return "", err
	}

	if err := json.Unmarshal(jsonResp, &credentials); err != nil {
		tokenRenewInProgress = false
		return "", err
	}
	tokenRenewInProgress = false
	return credentials.AccessToken, nil
}

// start actually creates the gRPC connection. The client construction is deferred till this point as this
// is the only place we get hold of Extensions which are required to construct auth round tripper.
func (e *opsrampOTLPExporter) start(ctx context.Context, host component.Host) (err error) {
	if serverName := endpointServerName(e.config.Endpoint); serverName != "" {
		configuredServerName := strings.TrimSpace(e.config.ClientConfig.TLSSetting.ServerName)
		if configuredServerName == "" || strings.Contains(configuredServerName, ":") {
			e.config.ClientConfig.TLSSetting.ServerName = serverName
		}
	}

	// Save original TLS settings for logging
	originalTLSInsecure := e.config.ClientConfig.TLSSetting.Insecure
	originalTLSInsecureSkipVerify := e.config.ClientConfig.TLSSetting.InsecureSkipVerify

	// Log TLS configuration
	e.settings.Logger.Info("opsrampotlp TLS configuration",
		zap.Bool("insecure", originalTLSInsecure),
		zap.Bool("insecure_skip_verify", originalTLSInsecureSkipVerify),
		zap.String("endpoint", e.config.Endpoint),
	)

	// Proxy tunnel approach:
	// - When insecure=true (plaintext mode): proxy environment vars are set, gRPC uses proxy for plaintext
	// - When insecure=false (TLS enabled): gRPC handles TLS natively, no custom proxy dialer needed
	// The proxy env vars (HTTP_PROXY, HTTPS_PROXY) are only set when insecure=true (see otelLogConfig.go)

	e.clientConn, err = e.config.ClientConfig.ToClientConn(
		ctx,
		host,
		e.settings,
		configgrpc.WithGrpcDialOption(grpc.WithUserAgent(e.userAgent)),
		configgrpc.WithGrpcDialOption(
			grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
				dialAddr := sanitizeDialAddress(addr)
				e.settings.Logger.Debug("opsrampotlp dial attempt",
					zap.String("addr", addr),
					zap.String("dial_addr", dialAddr),
					zap.String("endpoint", e.config.Endpoint),
				)

				if shouldBypassProxy(dialAddr) {
					e.settings.Logger.Info("opsrampotlp dialing direct (bypass)",
						zap.String("dial_addr", dialAddr),
					)
					return (&net.Dialer{}).DialContext(ctx, "tcp", dialAddr)
				}

				targetURL := &url.URL{Scheme: proxyLookupScheme(e.config.Endpoint, e.config.ClientConfig.TLSSetting.Insecure), Host: dialAddr}
				proxyURL, err := httpproxy.FromEnvironment().ProxyFunc()(targetURL)
				if err != nil {
					e.settings.Logger.Debug("opsrampotlp proxy resolution failed",
						zap.String("target", targetURL.String()),
						zap.Error(err),
					)
					return nil, fmt.Errorf("failed to resolve proxy from environment: %w", err)
				}

				// No proxy set — connect directly
				if proxyURL == nil {
					e.settings.Logger.Info("opsrampotlp dialing direct (no proxy from env)",
						zap.String("dial_addr", dialAddr),
						zap.String("lookup_target", targetURL.String()),
					)
					return (&net.Dialer{}).DialContext(ctx, "tcp", dialAddr)
				}

				e.settings.Logger.Debug("opsrampotlp dialing via proxy",
					zap.String("dial_addr", dialAddr),
					zap.String("proxy_host", proxyURL.Host),
					zap.String("lookup_target", targetURL.String()),
				)

				// Step 1: TCP connect to Squid proxy
				proxyAddr := proxyURL.Host // 172.16.7.37:3128
				conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", proxyAddr)
				if err != nil {
					e.settings.Logger.Debug("opsrampotlp proxy tcp connect failed",
						zap.String("proxy_addr", proxyAddr),
						zap.Error(err),
					)
					return nil, fmt.Errorf("failed to connect to proxy %s: %w", proxyAddr, err)
				}

				// Step 2: Send HTTP/1.1 CONNECT to Squid
				connectReq := &http.Request{
					Method:     "CONNECT",
					URL:        &url.URL{Host: dialAddr},
					Host:       dialAddr,
					Header:     make(http.Header),
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
				}

				// Add proxy auth if present in URL
				if proxyURL.User != nil {
					username := proxyURL.User.Username()
					password, _ := proxyURL.User.Password()
					connectReq.Header.Set(
						"Proxy-Authorization",
						"Basic "+basicAuth(username, password),
					)
				}

				// Write CONNECT request to proxy
				if err := connectReq.Write(conn); err != nil {
					e.settings.Logger.Debug("opsrampotlp proxy CONNECT write failed",
						zap.String("dial_addr", dialAddr),
						zap.String("proxy_addr", proxyAddr),
						zap.Error(err),
					)
					conn.Close()
					return nil, fmt.Errorf("failed to write CONNECT request: %w", err)
				}

				// Step 3: Read Squid's response
				// Use a shared reader so any bytes read ahead during response parsing are
				// preserved for gRPC's HTTP/2/TLS handshakes.
				br := bufio.NewReader(conn)
				resp, err := http.ReadResponse(br, connectReq)
				if err != nil {
					e.settings.Logger.Debug("opsrampotlp proxy CONNECT read failed",
						zap.String("dial_addr", dialAddr),
						zap.String("proxy_addr", proxyAddr),
						zap.Error(err),
					)
					conn.Close()
					return nil, fmt.Errorf("failed to read CONNECT response: %w", err)
				}

				e.settings.Logger.Debug("opsrampotlp proxy CONNECT response",
					zap.String("dial_addr", dialAddr),
					zap.String("proxy_addr", proxyAddr),
					zap.Int("status_code", resp.StatusCode),
					zap.String("status", resp.Status),
				)

				// Step 4: Check tunnel established
				if resp.StatusCode != http.StatusOK {
					resp.Body.Close()
					conn.Close()
					return nil, fmt.Errorf("proxy CONNECT failed: %s", resp.Status)
				}

				// Do not close resp.Body for successful CONNECT. It represents the live
				// tunnel; closing it here will close the underlying socket.

				// Create buffered conn to preserve any bytes read ahead by bufio.Reader
				tunnelConn := &bufferedConn{Conn: conn, reader: br}

				// Step 5: Perform manual TLS handshake through the proxy tunnel
				// Extract server name for TLS
				serverName := dialAddr
				if colonIdx := strings.LastIndex(dialAddr, ":"); colonIdx != -1 {
					serverName = dialAddr[:colonIdx]
				}

				// insecure_skip_verify controls certificate verification independently
				// insecure flag only controls connection mode (proxy) and TLS version fallback
				skipVerify := e.config.ClientConfig.TLSSetting.InsecureSkipVerify

				// Try TLS 1.2 first (mandatory minimum)
				tlsConfig := &tls.Config{
					ServerName:         serverName,
					InsecureSkipVerify: skipVerify,
					MinVersion:         tls.VersionTLS12,
				}

				tlsConn := tls.Client(tunnelConn, tlsConfig)
				tlsErr := tlsConn.HandshakeContext(ctx)

				// Fallback to TLS 1.0+ only when insecure=true and TLS 1.2 handshake fails
				if tlsErr != nil && e.config.ClientConfig.TLSSetting.Insecure {
					e.settings.Logger.Warn("opsrampotlp TLS 1.2 handshake failed, attempting fallback to TLS 1.0+",
						zap.String("dial_addr", dialAddr),
						zap.Error(tlsErr),
					)

					// Close the failed connection and re-establish tunnel for retry
					tunnelConn.Close()

					// Re-connect to proxy
					conn, err = (&net.Dialer{}).DialContext(ctx, "tcp", proxyAddr)
					if err != nil {
						return nil, fmt.Errorf("failed to reconnect to proxy for TLS fallback: %w", err)
					}

					// Re-send CONNECT request
					connectReq = &http.Request{
						Method:     "CONNECT",
						URL:        &url.URL{Host: dialAddr},
						Host:       dialAddr,
						Header:     make(http.Header),
						Proto:      "HTTP/1.1",
						ProtoMajor: 1,
						ProtoMinor: 1,
					}
					if proxyURL.User != nil {
						username := proxyURL.User.Username()
						password, _ := proxyURL.User.Password()
						connectReq.Header.Set("Proxy-Authorization", "Basic "+basicAuth(username, password))
					}
					if err = connectReq.Write(conn); err != nil {
						conn.Close()
						return nil, fmt.Errorf("failed to write CONNECT request for TLS fallback: %w", err)
					}
					br = bufio.NewReader(conn)
					resp, err = http.ReadResponse(br, connectReq)
					if err != nil {
						conn.Close()
						return nil, fmt.Errorf("failed to read CONNECT response for TLS fallback: %w", err)
					}
					if resp.StatusCode != http.StatusOK {
						resp.Body.Close()
						conn.Close()
						return nil, fmt.Errorf("proxy CONNECT failed for TLS fallback: %s", resp.Status)
					}
					tunnelConn = &bufferedConn{Conn: conn, reader: br}

					// Retry with TLS 1.0+ (fallback)
					tlsConfig = &tls.Config{
						ServerName:         serverName,
						InsecureSkipVerify: true,
						MinVersion:         tls.VersionTLS10,
					}
					tlsConn = tls.Client(tunnelConn, tlsConfig)
					tlsErr = tlsConn.HandshakeContext(ctx)
				}

				if tlsErr != nil {
					e.settings.Logger.Error("opsrampotlp TLS handshake failed",
						zap.String("dial_addr", dialAddr),
						zap.Error(tlsErr),
					)
					tunnelConn.Close()
					return nil, fmt.Errorf("TLS handshake failed: %w", tlsErr)
				}

				e.settings.Logger.Debug("opsrampotlp proxy tunnel + TLS established",
					zap.String("dial_addr", dialAddr),
					zap.String("proxy_addr", proxyAddr),
					zap.String("tls_version", tlsVersionString(tlsConn.ConnectionState().Version)),
					zap.Bool("insecure_skip_verify", skipVerify),
				)

				return tlsConn, nil
			})),
	)

	if err != nil {
		e.settings.Logger.Error("opsrampotlp proxy CONNECT read failed",
			zap.Error(err),
		)
		return err
	}

	e.traceExporter = ptraceotlp.NewGRPCClient(e.clientConn)
	e.metricExporter = pmetricotlp.NewGRPCClient(e.clientConn)
	e.logExporter = plogotlp.NewGRPCClient(e.clientConn)
	headers := map[string]string{}
	for k, v := range e.config.ClientConfig.Headers {
		headers[k] = string(v)
	}
	e.metadata = metadata.New(headers)
	e.metadata.Set("Authorization", fmt.Sprintf("Bearer %s", e.accessToken))
	e.callOptions = []grpc.CallOption{
		grpc.MaxCallSendMsgSize(e.config.Security.OtelExporterSetting.GrpcMaxSendSize),
		grpc.MaxCallRecvMsgSize(e.config.Security.OtelExporterSetting.GrpcMaxRecvSize),
		grpc.WaitForReady(e.config.ClientConfig.WaitForReady),
	}
	return
}

func (e *opsrampOTLPExporter) shutdown(context.Context) error {
	if e.clientConn == nil {
		return nil
	}
	return e.clientConn.Close()
}

func (e *opsrampOTLPExporter) pushTraces(ctx context.Context, td ptrace.Traces) error {
	req := ptraceotlp.NewExportRequestFromTraces(td)
	_, err := e.traceExporter.Export(e.enhanceContext(ctx), req, e.callOptions...)
	// trying to get new access token in case of expiration
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.Unauthenticated {
			if err = e.updateExpiredToken(); err != nil {
				return fmt.Errorf("couldn't retrieve new token instead of expired: %w", err)
			}
			_, err = e.traceExporter.Export(e.enhanceContext(ctx), req, e.callOptions...)
			if err != nil {
				return err
			}
		}
		return processError(err)
	}
	return nil
}

func (e *opsrampOTLPExporter) pushMetrics(ctx context.Context, md pmetric.Metrics) error {
	req := pmetricotlp.NewExportRequestFromMetrics(md)
	_, err := e.metricExporter.Export(e.enhanceContext(ctx), req, e.callOptions...)
	// trying to get new access token in case of expiration
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.Unauthenticated {
			if err = e.updateExpiredToken(); err != nil {
				return fmt.Errorf("couldn't retrieve new token instead of expired: %w", err)
			}

			_, err = e.metricExporter.Export(e.enhanceContext(ctx), req, e.callOptions...)
			if err != nil {
				return err
			}
		}
		return processError(err)
	}

	return processError(err)
}

func (e *opsrampOTLPExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	if ld.LogRecordCount() <= 0 {
		return nil
	}

	if e.config.Masking != nil {
		e.applyMasking(ld)
	}
	if e.config.ExpirationSkip != 0 {
		e.skipExpired(ld)
	}
	if ld.LogRecordCount() <= 0 {
		return nil
	}

	if e.config.Masking != nil {
		e.applyMasking(ld)
	}

	if e.config.ExpirationSkip != 0 {
		e.skipExpired(ld)
	}

	req := plogotlp.NewExportRequestFromLogs(ld)
	_, err := e.logExporter.Export(e.enhanceContext(context.Background()), req, e.callOptions...)

	// trying to get a new access token in case of expiration
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.Unauthenticated {
			if err = e.updateExpiredToken(); err != nil {
				return fmt.Errorf("couldn't retrieve new token instead of expired: %w", err)
			}

			_, err = e.logExporter.Export(e.enhanceContext(context.Background()), req, e.callOptions...)
			if err != nil {
				return err
			}
		}
		return processError(err)
	}
	return nil
}

func (e *opsrampOTLPExporter) updateExpiredToken() error {
	accessToken, err := getAuthToken(e.config.Security)
	if err != nil {
		return err
	}
	e.mut.Lock()
	defer e.mut.Unlock()
	e.accessToken = accessToken
	// Always update the metadata when token is refreshed
	e.metadata.Set("Authorization", fmt.Sprintf("Bearer %s", e.accessToken))
	return nil
}

func (e *opsrampOTLPExporter) enhanceContext(ctx context.Context) context.Context {
	e.mut.Lock()
	defer e.mut.Unlock()
	if e.metadata.Len() > 0 {
		// Create a copy of the metadata to avoid race conditions during gRPC validation
		mdCopy := e.metadata.Copy()
		return metadata.NewOutgoingContext(ctx, mdCopy)
	}
	return ctx
}

// Send a trace or metrics request to the server. "perform" function is expected to make
// the actual gRPC unary call that sends the request. This function implements the
// common OTLP logic around request handling such as retries and throttling.
func processError(err error) error {
	if err == nil {
		// Request is successful, we are done.
		return nil
	}

	// We have an error, check gRPC status code.

	st := status.Convert(err)
	if st.Code() == codes.OK {
		// Not really an error, still success.
		return nil
	}

	// Now, this is this a real error.

	retryInfo := getRetryInfo(st)
	if !shouldRetry(st.Code(), retryInfo) {
		// It is not a retryable error, we should not retry.
		return consumererror.NewPermanent(err)
	}

	// Check if server returned throttling information.
	throttleDuration := getThrottleDuration(retryInfo)
	if throttleDuration != 0 {
		// We are throttled. Wait before retrying as requested by the server.
		return exporterhelper.NewThrottleRetry(err, throttleDuration)
	}

	// Need to retry.

	return err
}

func shouldRetry(code codes.Code, retryInfo *errdetails.RetryInfo) bool {
	switch code {
	case codes.Canceled,
		codes.DeadlineExceeded,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unavailable,
		codes.Unknown,
		codes.PermissionDenied,
		codes.Internal,
		codes.DataLoss:
		// These are retryable errors.
		return true
	case codes.ResourceExhausted:
		// Retry only if RetryInfo was supplied by the server.
		// This indicates that the server can still recover from resource exhaustion.
		return retryInfo != nil
	}
	// Don't retry on any other code.
	return false
}

func getRetryInfo(status *status.Status) *errdetails.RetryInfo {
	for _, detail := range status.Details() {
		if t, ok := detail.(*errdetails.RetryInfo); ok {
			return t
		}
	}
	return nil
}

func getThrottleDuration(t *errdetails.RetryInfo) time.Duration {
	if t == nil || t.RetryDelay == nil {
		return 0
	}
	if t.RetryDelay.Seconds > 0 || t.RetryDelay.Nanos > 0 {
		return time.Duration(t.RetryDelay.Seconds)*time.Second + time.Duration(t.RetryDelay.Nanos)*time.Nanosecond
	}
	return 0
}

func (e *opsrampOTLPExporter) applyMasking(ld plog.Logs) {
	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		resLogs := ld.ResourceLogs().At(i)
		for k := 0; k < resLogs.ScopeLogs().Len(); k++ {
			scopedLog := resLogs.ScopeLogs().At(k)
			for z := 0; z < scopedLog.LogRecords().Len(); z++ {
				log := scopedLog.LogRecords().At(z)
				for _, setting := range e.config.Masking {
					rExp := regexp.MustCompile(setting.Regexp)
					log.Body().SetStr(rExp.ReplaceAllString(log.Body().AsString(), setting.Placeholder))
				}
			}
		}
	}
}

func (e *opsrampOTLPExporter) skipExpired(ld plog.Logs) {
	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		resLogs := ld.ResourceLogs().At(i)

		for k := 0; k < resLogs.ScopeLogs().Len(); k++ {
			resLogs.ScopeLogs().At(k).LogRecords().RemoveIf(func(el plog.LogRecord) bool {
				fmt.Println(el.Timestamp().AsTime().String(), time.Now().Add(-e.config.ExpirationSkip).String())
				return el.Timestamp().AsTime().Before(time.Now().Add(-e.config.ExpirationSkip))
			})
		}
	}
}

func getAuthTokenWithTlsDisabled(cfg SecuritySettings) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return httpproxy.FromEnvironment().ProxyFunc()(req.URL)
			},
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}

	form := url.Values{}
	form.Set("client_id", cfg.ClientID)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", cfg.OAuthServiceURL, bytes.NewBufferString(form.Encode())) //nolint:usestdlibvars
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &credentials); err != nil {
		return "", err
	}

	return credentials.AccessToken, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

type bufferedConn struct {
	net.Conn
	reader *bufio.Reader
}

func (c *bufferedConn) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

func proxyLookupScheme(endpoint string, insecure bool) string {
	if strings.HasPrefix(endpoint, "http://") {
		return "http"
	}
	if strings.HasPrefix(endpoint, "https://") {
		return "https"
	}
	if insecure {
		return "http"
	}

	return "https"
}

func sanitizeDialAddress(addr string) string {
	for _, prefix := range []string{"dns:///", "passthrough:///"} {
		if strings.HasPrefix(addr, prefix) {
			return strings.TrimPrefix(addr, prefix)
		}
	}

	return addr
}

func shouldBypassProxy(addr string) bool {
	host := addr
	if h, _, err := net.SplitHostPort(addr); err == nil {
		host = h
	}
	host = strings.Trim(host, "[]")

	if strings.EqualFold(host, "localhost") {
		return true
	}

	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func endpointServerName(endpoint string) string {
	if endpoint == "" {
		return ""
	}

	if strings.Contains(endpoint, "://") {
		u, err := url.Parse(endpoint)
		if err == nil {
			return u.Hostname()
		}
	}

	host := endpoint
	if idx := strings.Index(host, "/"); idx >= 0 {
		host = host[:idx]
	}

	if h, _, err := net.SplitHostPort(host); err == nil {
		return strings.Trim(h, "[]")
	}

	return strings.Trim(host, "[]")
}

func tlsVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("unknown (0x%04x)", version)
	}
}
