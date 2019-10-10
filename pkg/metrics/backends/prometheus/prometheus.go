// Copyright (c) 2019 Kien Nguyen-Tuan <kiennt2609@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"context"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	prometheusclient "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// Backend implements a metric backend for Prometheus.
type Backend struct {
	prometheus prometheus.API
	logger     log.Logger
}

// New returns a new client for talking to a Prometheus Backend, or an error
func New(address string, logger log.Logger) (*Backend, error) {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	client, err := prometheusclient.NewClient(prometheusclient.Config{
		Address: address,
	})
	if err != nil {
		return nil, errors.Wrap(err, "instantiating prometheus client")
	}

	api := prometheus.NewAPI(client)
	return &Backend{
		prometheus: api,
		logger:     logger,
	}, nil
}

// QueryInstant performs instant query and returns results in model.Vector type.
func (b *Backend) QueryInstant(ctx context.Context, query string, ts time.Time) (model.Vector, error) {
	level.Debug(b.logger).Log("msg", "querying instant", "query", query)
	val, warns, err := b.prometheus.Query(ctx, query, ts)
	if err != nil {
		return nil, errors.Wrap(err, "querying instant")
	}
	if len(warns) > 0 {
		level.Warn(b.logger).Log("msg", "querying instant warning", strings.Join(warns, ", "), "query", query)
	}

	switch v := val.(type) {
	case model.Vector:
		return v, nil
	default:
		return nil, errors.Errorf("unknown supported type: '%q'", v)
	}
}