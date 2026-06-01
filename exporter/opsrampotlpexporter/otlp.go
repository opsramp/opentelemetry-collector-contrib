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
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/http/httpproxy"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gopkg.in/natefinch/lumberjack.v2"
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

	logger *zap.Logger
}

// Crete new exporter and start it. The exporter will begin connecting, but
// this function may return before the connection is established.
func newExporter(cfg component.Config, set exporter.Settings) (*opsrampOTLPExporter, error) {
	oCfg := cfg.(*Config)

	// Extract TLS options from the gRPC client config
	tlsOpts := TLSOptions{
		Insecure:           oCfg.ClientConfig.TLS.Insecure,
		InsecureSkipVerify: oCfg.ClientConfig.TLS.InsecureSkipVerify,
	}

	accessToken, err := getAuthToken(oCfg.Security, tlsOpts)
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

	logger := initLogger()
	return &opsrampOTLPExporter{config: oCfg, settings: set.TelemetrySettings,
		userAgent: userAgent, accessToken: accessToken, logger: logger}, nil
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

// TLSOptions holds TLS configuration for HTTP/gRPC connections.
type TLSOptions struct {
	// Insecure disables TLS entirely (plaintext connection).
	Insecure bool
	// InsecureSkipVerify skips server certificate verification.
	InsecureSkipVerify bool
}

func getAuthToken(cfg SecuritySettings, tlsOpts TLSOptions) (string, error) {
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

	// Configure TLS based on options
	var tlsConfig *tls.Config
	if !tlsOpts.Insecure {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: tlsOpts.InsecureSkipVerify,
			MinVersion:         tls.VersionTLS12,
		}
	}

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
			TLSClientConfig:       tlsConfig,
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
		// Only fallback to TLS disabled mode if InsecureSkipVerify wasn't already true
		// and the error is TLS-related
		if !tlsOpts.InsecureSkipVerify &&
			(strings.Contains(err.Error(), "x509: certificate signed by unknown authority") ||
				strings.Contains(err.Error(), "TLS handshake timeout") ||
				strings.Contains(err.Error(), "certificate") ||
				strings.Contains(err.Error(), "tls:")) {
			// Retry with TLS verification disabled as a fallback
			return getAuthTokenWithTlsDisabled(cfg)
		}
		tokenRenewInProgress = false
		return "", err
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil {
		tokenRenewInProgress = false
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		tokenRenewInProgress = false
		return "", fmt.Errorf("AUTH_DIAG: OAuth token request failed with status=%d body=%s", resp.StatusCode, string(jsonResp))
	}

	if err := json.Unmarshal(jsonResp, &credentials); err != nil {
		tokenRenewInProgress = false
		return "", fmt.Errorf("AUTH_DIAG: failed to unmarshal OAuth response: %w body=%s", err, string(jsonResp))
	}

	if credentials.AccessToken == "" {
		tokenRenewInProgress = false
		return "", fmt.Errorf("AUTH_DIAG: OAuth returned empty access_token, response body=%s", string(jsonResp))
	}
	tokenRenewInProgress = false
	return credentials.AccessToken, nil
}

// start actually creates the gRPC connection. The client construction is deferred till this point as this
// is the only place we get hold of Extensions which are required to construct auth round tripper.
func (e *opsrampOTLPExporter) start(ctx context.Context, host component.Host) (err error) {
	if serverName := endpointServerName(e.config.Endpoint); serverName != "" {
		configuredServerName := strings.TrimSpace(e.config.ClientConfig.TLS.ServerName)
		if configuredServerName == "" || strings.Contains(configuredServerName, ":") {
			e.config.ClientConfig.TLS.ServerName = serverName
		}
	}

	// Save original TLS settings — we handle TLS manually in the dialer so that
	// gRPC's own TLS+proxy interception is bypassed completely. By forcing
	// Insecure=true here gRPC treats the connection as plain TCP and never
	// tries to establish its own HTTP-CONNECT tunnel, which would conflict with
	// our manual CONNECT implementation below.
	originalInsecure := e.config.ClientConfig.TLS.Insecure
	originalSkipVerify := e.config.ClientConfig.TLS.InsecureSkipVerify
	e.config.ClientConfig.TLS.Insecure = true

	e.clientConn, err = e.config.ClientConfig.ToClientConn(
		ctx,
		host.GetExtensions(),
		e.settings,
		configgrpc.WithGrpcDialOption(grpc.WithUserAgent(e.userAgent)),
		configgrpc.WithGrpcDialOption(
			grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
				dialAddr := sanitizeDialAddress(addr)

				// establishTLSConnection wraps a raw TCP dialer with a TLS handshake,
				// falling back to InsecureSkipVerify on x509 errors when allowed.
				establishTLSConnection := func(dialAddr string, redial func() (net.Conn, error)) (net.Conn, error) {
					conn, err := redial()
					if err != nil {
						return nil, err
					}
					if originalInsecure {
						return conn, nil
					}
					serverName := dialAddr
					if colonIdx := strings.LastIndex(dialAddr, ":"); colonIdx != -1 {
						serverName = dialAddr[:colonIdx]
					}
					tlsConn := tls.Client(conn, &tls.Config{
						ServerName:         serverName,
						InsecureSkipVerify: originalSkipVerify, // #nosec G402
					})
					if err := tlsConn.HandshakeContext(ctx); err != nil {
						conn.Close()
						if originalSkipVerify {
							e.settings.Logger.Warn("TLS handshake failed (InsecureSkipVerify=true), bypassing TLS",
								zap.String("addr", dialAddr), zap.Error(err))
							conn, err = redial()
							if err != nil {
								return nil, err
							}
							return conn, nil
						}
						if strings.Contains(err.Error(), "x509:") {
							e.settings.Logger.Warn("TLS x509 error, retrying with InsecureSkipVerify",
								zap.String("addr", dialAddr), zap.Error(err))
							conn, err = redial()
							if err != nil {
								return nil, err
							}
							tlsConn = tls.Client(conn, &tls.Config{
								ServerName:         serverName,
								InsecureSkipVerify: true, // #nosec G402
							})
							if err := tlsConn.HandshakeContext(ctx); err != nil {
								conn.Close()
								return nil, fmt.Errorf("TLS handshake failed (InsecureSkipVerify fallback): %w", err)
							}
							return tlsConn, nil
						}
						return nil, fmt.Errorf("TLS handshake failed: %w", err)
					}
					e.settings.Logger.Debug("TLS established (direct)",
						zap.Uint16("version", tlsConn.ConnectionState().Version))
					return tlsConn, nil
				}

				// Bypass proxy for loopback addresses.
				if shouldBypassProxy(dialAddr) {
					return establishTLSConnection(dialAddr, func() (net.Conn, error) {
						return (&net.Dialer{}).DialContext(ctx, "tcp", dialAddr)
					})
				}

				// Resolve proxy URL from environment (HTTPS_PROXY for TLS, HTTP_PROXY for plain).
				targetURL := &url.URL{Scheme: proxyLookupScheme(e.config.Endpoint, originalInsecure), Host: dialAddr}
				proxyURL, err := httpproxy.FromEnvironment().ProxyFunc()(targetURL)
				if err != nil {
					return nil, fmt.Errorf("proxy resolution failed: %w", err)
				}
				if proxyURL == nil {
					return establishTLSConnection(dialAddr, func() (net.Conn, error) {
						return (&net.Dialer{}).DialContext(ctx, "tcp", dialAddr)
					})
				}

				e.settings.Logger.Debug("connecting via proxy",
					zap.String("proxy", proxyURL.Host), zap.String("target", dialAddr),
					zap.Bool("insecure", originalInsecure), zap.Bool("skipVerify", originalSkipVerify))

				// establishProxyTunnel manually issues HTTP CONNECT with Proxy-Authorization.
				// Using the standard library directly (rather than golang.org/x/net/proxy)
				// ensures the Proxy-Authorization header is always sent correctly.
				establishProxyTunnel := func() (net.Conn, error) {
					conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", proxyURL.Host)
					if err != nil {
						return nil, fmt.Errorf("proxy dial failed: %w", err)
					}
					connectReq := &http.Request{
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
					if err := connectReq.Write(conn); err != nil {
						conn.Close()
						return nil, fmt.Errorf("proxy CONNECT write failed: %w", err)
					}
					br := bufio.NewReader(conn)
					resp, err := http.ReadResponse(br, connectReq)
					if err != nil {
						conn.Close()
						return nil, fmt.Errorf("proxy CONNECT read failed: %w", err)
					}
					resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						conn.Close()
						return nil, fmt.Errorf("proxy CONNECT rejected: %s", resp.Status)
					}
					return &bufferedConn{Conn: conn, reader: br}, nil
				}

				tunnelConn, err := establishProxyTunnel()
				if err != nil {
					return nil, err
				}

				if originalInsecure {
					return tunnelConn, nil
				}

				// Manual TLS handshake over the proxy tunnel.
				serverName := dialAddr
				if colonIdx := strings.LastIndex(dialAddr, ":"); colonIdx != -1 {
					serverName = dialAddr[:colonIdx]
				}
				tlsConn := tls.Client(tunnelConn, &tls.Config{
					ServerName:         serverName,
					InsecureSkipVerify: originalSkipVerify, // #nosec G402
				})
				if err := tlsConn.HandshakeContext(ctx); err != nil {
					tunnelConn.Close()
					if originalSkipVerify {
						e.settings.Logger.Warn("TLS via proxy failed (InsecureSkipVerify=true), bypassing TLS",
							zap.String("target", dialAddr), zap.Error(err))
						tunnelConn, err = establishProxyTunnel()
						if err != nil {
							return nil, err
						}
						return tunnelConn, nil
					}
					return nil, fmt.Errorf("TLS handshake via proxy failed: %w", err)
				}
				e.settings.Logger.Debug("TLS established (via proxy)",
					zap.Uint16("version", tlsConn.ConnectionState().Version),
					zap.Bool("insecure", originalInsecure),
					zap.Bool("skipVerify", originalSkipVerify))
				return tlsConn, nil
			})),
	)
	if err != nil {
		return err
	}

	e.traceExporter = ptraceotlp.NewGRPCClient(e.clientConn)
	e.metricExporter = pmetricotlp.NewGRPCClient(e.clientConn)
	e.logExporter = plogotlp.NewGRPCClient(e.clientConn)
	headers := map[string]string{}
	for _, pair := range e.config.ClientConfig.Headers {
		headers[pair.Name] = string(pair.Value)
	}
	e.metadata = metadata.New(headers)
	e.metadata.Set("Authorization", fmt.Sprintf("Bearer %s", e.accessToken))

	// AUTH_DIAG: Log metadata keys being sent with gRPC calls
	var mdKeys []string
	for k, v := range e.metadata {
		if k == "authorization" {
			if len(v) > 0 && len(v[0]) > 15 {
				mdKeys = append(mdKeys, fmt.Sprintf("%s=len(%d)", k, len(v[0])))
			} else {
				mdKeys = append(mdKeys, fmt.Sprintf("%s=EMPTY_OR_SHORT", k))
			}
		} else {
			mdKeys = append(mdKeys, fmt.Sprintf("%s=%v", k, v))
		}
	}
	if e.logger != nil {
		e.logger.Debug("AUTH_DIAG: exporter start() metadata", zap.Strings("metadata_entries", mdKeys))
	}
	e.settings.Logger.Debug("AUTH_DIAG: exporter start() metadata", zap.Strings("metadata_entries", mdKeys))
	e.settings.Logger.Debug("AUTH_DIAG: exporter endpoints",
		zap.String("grpc_endpoint", e.config.ClientConfig.Endpoint),
		zap.String("oauth_url", e.config.Security.OAuthServiceURL),
		zap.Int("access_token_len", len(e.accessToken)),
	)

	e.callOptions = []grpc.CallOption{
		grpc.MaxCallSendMsgSize(e.config.Security.OtelExporterSetting.GrpcMaxSendSize),
		grpc.MaxCallRecvMsgSize(e.config.Security.OtelExporterSetting.GrpcMaxRecvSize),
		grpc.WaitForReady(e.config.ClientConfig.WaitForReady),
	}

	//e.logger.Debug(
	//	"OTLP Exporter started",
	//	zap.String("Calloptions", fmt.Sprintf("%v", e.callOptions)),
	//	zap.Int("Calloptions -> grpc.MaxCallSendMsgSize", e.config.Security.OtelExporterSetting.GrpcMaxSendSize),
	//	zap.Int("Calloptions -> grpc.MaxCallRecvMsgSize", e.config.Security.OtelExporterSetting.GrpcMaxRecvSize),
	//)

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

	//start := time.Now()
	//e.logger.Debug("Exporter: pushLogs Started processing logs", zap.String("start at", start.Format(time.RFC3339)))

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
	//
	//data, _ := req.MarshalJSON()
	//e.logger.Debug("request details",
	//	zap.String("started_at", start.Format(time.RFC3339Nano)),
	//	zap.Int("ResourceLogsCount", ld.ResourceLogs().Len()),
	//	zap.Int("TotalLogRecordCount", ld.LogRecordCount()),
	//	zap.Int("RequestSizeBytes", len(data)),
	//)
	//
	//beforePushEndTime := time.Now()

	_, err := e.logExporter.Export(e.enhanceContext(context.Background()), req, e.callOptions...)

	//end := time.Now()
	//e.logger.Debug("Exporter: pushLogs: completed processing logs",
	//	zap.String("Stage 1", "before push"),
	//	zap.String("ended_at", beforePushEndTime.Format(time.RFC3339Nano)),
	//	zap.Float64("duration_seconds", beforePushEndTime.Sub(start).Seconds()),
	//	zap.Int64("duration_ms", beforePushEndTime.Sub(start).Milliseconds()),
	//	zap.Float64("duration_seconds", beforePushEndTime.Sub(start).Seconds()),
	//	zap.String("Stage 1", "before push - end"),
	//
	//	zap.String("Stage 2", "after push"),
	//	zap.String("ended_at", end.Format(time.RFC3339Nano)),
	//	zap.Float64("duration_seconds", end.Sub(start).Seconds()),
	//	zap.Int64("duration_ms", end.Sub(start).Milliseconds()),
	//	zap.Float64("duration_seconds", end.Sub(start).Seconds()),
	//	zap.String("Stage 2", "after push - end"),
	//)

	// trying to get a new access token in case of expiration
	if err != nil {
		st := status.Convert(err)
		e.settings.Logger.Warn("AUTH_DIAG: pushLogs Export failed",
			zap.String("grpc_code", st.Code().String()),
			zap.String("error", err.Error()),
			zap.Int("metadata_len", e.metadata.Len()),
			zap.Bool("has_tenantId", len(e.metadata.Get("tenantid")) > 0),
			zap.String("tenantId_value", strings.Join(e.metadata.Get("tenantid"), ",")),
			zap.Bool("has_authorization", len(e.metadata.Get("authorization")) > 0),
			zap.Int("access_token_len", len(e.accessToken)),
		)
		if e.logger != nil {
			e.logger.Warn("AUTH_DIAG: pushLogs Export failed",
				zap.String("grpc_code", st.Code().String()),
				zap.String("error", err.Error()),
				zap.Int("metadata_len", e.metadata.Len()),
				zap.Bool("has_tenantId", len(e.metadata.Get("tenantid")) > 0),
				zap.String("tenantId_value", strings.Join(e.metadata.Get("tenantid"), ",")),
				zap.Bool("has_authorization", len(e.metadata.Get("authorization")) > 0),
				zap.Int("access_token_len", len(e.accessToken)),
			)
		}
		if st.Code() == codes.Unauthenticated {
			if err = e.updateExpiredToken(); err != nil {
				e.settings.Logger.Error("AUTH_DIAG: token refresh failed", zap.Error(err))
				return fmt.Errorf("couldn't retrieve new token instead of expired: %w", err)
			}
			e.settings.Logger.Debug("AUTH_DIAG: token refreshed, retrying Export",
				zap.Int("new_token_len", len(e.accessToken)),
			)

			_, err = e.logExporter.Export(e.enhanceContext(context.Background()), req, e.callOptions...)
			if err != nil {
				e.settings.Logger.Error("AUTH_DIAG: retry Export also failed after token refresh",
					zap.Error(err),
					zap.String("tenantId_value", strings.Join(e.metadata.Get("tenantid"), ",")),
					zap.Int("access_token_len", len(e.accessToken)),
				)
				return err
			}
		}
		return processError(err)
	}
	return nil
}

func (e *opsrampOTLPExporter) updateExpiredToken() error {
	// Extract TLS options from the gRPC client config
	tlsOpts := TLSOptions{
		Insecure:           e.config.ClientConfig.TLS.Insecure,
		InsecureSkipVerify: e.config.ClientConfig.TLS.InsecureSkipVerify,
	}
	accessToken, err := getAuthToken(e.config.Security, tlsOpts)
	if err != nil {
		return err
	}
	e.mut.Lock()
	defer e.mut.Unlock()
	e.settings.Logger.Debug("AUTH_DIAG: updateExpiredToken refreshed",
		zap.Int("old_token_len", len(e.accessToken)),
		zap.Int("new_token_len", len(accessToken)),
		zap.Bool("token_changed", e.accessToken != accessToken),
	)
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

	req, err := http.NewRequest("POST", cfg.OAuthServiceURL, bytes.NewBufferString(form.Encode()))
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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AUTH_DIAG: OAuth (TLS disabled) token request failed with status=%d body=%s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &credentials); err != nil {
		return "", err
	}

	if credentials.AccessToken == "" {
		return "", fmt.Errorf("AUTH_DIAG: OAuth (TLS disabled) returned empty access_token, response body=%s", string(body))
	}

	return credentials.AccessToken, nil
}

func basicAuth(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
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

func initLogger() *zap.Logger {
	writer := &lumberjack.Logger{
		Filename:   "/var/log/opsramp/exporter-info.log", // or any path you prefer
		MaxSize:    10,                                   // megabytes
		MaxBackups: 5,                                    // number of old files to keep
		MaxAge:     30,                                   // days to keep
		Compress:   true,                                 // gzip
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(writer),
		zap.DebugLevel,
	)

	return zap.New(core)
}
