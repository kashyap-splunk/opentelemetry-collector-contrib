// Copyright 2020 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !windows
// +build !windows

package podmanreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/podmanreceiver"

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/pdata"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
	"go.uber.org/zap"
)

type receiver struct {
	config        *Config
	set           component.ReceiverCreateSettings
	clientFactory clientFactory
	client        client

	metricsComponent component.MetricsReceiver
	logsConsumer     consumer.Logs
	metricsConsumer  consumer.Metrics
}

func newReceiver(
	_ context.Context,
	set component.ReceiverCreateSettings,
	config *Config,
	clientFactory clientFactory,
) (*receiver, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	if clientFactory == nil {
		clientFactory = newPodmanClient
	}

	recv := &receiver{
		config:        config,
		clientFactory: clientFactory,
		set:           set,
	}

	return recv, err
}

func (r *receiver) registerMetricsConsumer(mc consumer.Metrics, set component.ReceiverCreateSettings) error {
	r.metricsConsumer = mc
	scrp, err := scraperhelper.NewScraper(typeStr, r.scrape, scraperhelper.WithStart(r.start))
	if err != nil {
		return err
	}
	r.metricsComponent, err = scraperhelper.NewScraperControllerReceiver(&r.config.ScraperControllerSettings, set, mc, scraperhelper.AddScraper(scrp))
	return err
}

func (r *receiver) registerLogsConsumer(lc consumer.Logs) {
	r.logsConsumer = lc
}

func (r *receiver) Start(ctx context.Context, host component.Host) error {
	if r.logsConsumer == nil {
		r.set.Logger.Warn("Logs Receiver is not set")
	}
	if r.metricsConsumer == nil {
		r.set.Logger.Warn("Metrics Receiver is not set")
	}
	if r.metricsConsumer != nil {
		go func() {
			err := r.metricsComponent.Start(ctx, host)
			if err != nil {
				r.set.Logger.Warn("", zap.Error(err))
			}
		}()
	}
	if r.logsConsumer != nil {
		eventBackoff := backoff.NewExponentialBackOff()
		eventBackoff.InitialInterval = 2 * time.Second
		eventBackoff.MaxInterval = 3 * time.Minute
		eventBackoff.Multiplier = 2
		eventBackoff.MaxElapsedTime = 0
		go func() {
			errorWhileRetry := backoff.Retry(func() error {
				err := r.handleEvents(ctx, eventBackoff)
				return err
			}, eventBackoff)
			if errorWhileRetry != nil {
				r.set.Logger.Warn("", zap.Error(errorWhileRetry))
			}
		}()
	}
	return nil
}

func (r *receiver) Shutdown(ctx context.Context) error {
	if r.metricsConsumer != nil {
		err := r.metricsComponent.Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *receiver) start(context.Context, component.Host) error {
	c, err := r.clientFactory(r.set.Logger, r.config)
	if err == nil {
		r.client = c
	}
	return err
}

func (r *receiver) scrape(context.Context) (pdata.Metrics, error) {
	var err error

	stats, err := r.client.stats()
	if err != nil {
		r.set.Logger.Error("error fetching stats", zap.Error(err))
		return pdata.Metrics{}, err
	}

	md := pdata.NewMetrics()
	for i := range stats {
		translateStatsToMetrics(&stats[i], time.Now(), md.ResourceMetrics().AppendEmpty())
	}
	return md, nil
}

func (r *receiver) handleEvents(ctx context.Context, eventBackoff *backoff.ExponentialBackOff) error {
	c, err := r.clientFactory(r.set.Logger, r.config)
	if err == nil {
		r.client = c
	} else {
		r.set.Logger.Error("error fetching/processing events", zap.Error(err))
		return err
	}
	events := make(chan event)
	errorChan := make(chan error)
	err = r.client.events(r.set.Logger, events, errorChan)
	if err != nil {
		r.set.Logger.Error("error fetching stats", zap.Error(err))
		return err
	}
	for {
		select {
		case err := <-errorChan:
			r.set.Logger.Error("Error while fetching/decoding events", zap.Error(err))
			return err
		case eventToTranslate := <-events:
			ld, er := traslateEventsToLogs(r.set.Logger, eventToTranslate)
			if er != nil {
				r.set.Logger.Error("Failed to translate into logs", zap.Error(er))
				return er
			}
			transferErr := r.logsConsumer.ConsumeLogs(ctx, ld)
			if transferErr != nil {
				r.set.Logger.Error("Something went wrong while transferring it to the next component", zap.Error(transferErr))
				return transferErr
			}
			eventBackoff.Reset()
		}
	}
}
