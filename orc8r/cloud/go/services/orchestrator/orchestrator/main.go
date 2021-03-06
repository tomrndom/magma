/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/analytics"
	analytics_protos "magma/orc8r/cloud/go/services/analytics/protos"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	exporter_protos "magma/orc8r/cloud/go/services/metricsd/protos"
	"magma/orc8r/cloud/go/services/orchestrator"
	orc8r_analytics "magma/orc8r/cloud/go/services/orchestrator/analytics"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	"magma/orc8r/cloud/go/services/orchestrator/servicers"
	indexer_protos "magma/orc8r/cloud/go/services/state/protos"
	streamer_protos "magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, orchestrator.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating orchestrator service %s", err)
	}

	obsidian.AttachHandlers(srv.EchoServer, handlers.GetObsidianHandlers())

	var exporterServicer exporter_protos.MetricsExporterServer

	var serviceConfig orchestrator.Config
	_, _, err = config.GetStructuredServiceConfig(orc8r.ModuleName, orchestrator.ServiceName, &serviceConfig)
	if err != nil {
		glog.Infof("err %v failed parsing the config file ", err)
		return
	}

	if serviceConfig.UseGRPCExporter == true {
		grpcAddress := serviceConfig.PrometheusGRPCPushAddress
		exporterServicer = servicers.NewGRPCPushExporterServicer(grpcAddress)
	} else {
		exporterServicer = servicers.NewPushExporterServicer(serviceConfig.PrometheusPushAddresses)
	}

	builder_protos.RegisterMconfigBuilderServer(srv.GrpcServer, servicers.NewBuilderServicer())
	exporter_protos.RegisterMetricsExporterServer(srv.GrpcServer, exporterServicer)
	indexer_protos.RegisterIndexerServer(srv.GrpcServer, servicers.NewIndexerServicer())
	streamer_protos.RegisterStreamProviderServer(srv.GrpcServer, servicers.NewProviderServicer())
	analytics_protos.RegisterAnalyticsCollectorServer(srv.GrpcServer,
		analytics.NewCollectorService(analytics.GetPrometheusClient(), orc8r_analytics.GetAnalyticsCalculations(&serviceConfig)))

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %s", err)
	}
}
