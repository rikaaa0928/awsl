package tracing

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/profiler"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/global"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func Start(ctx context.Context, conf config.Configs) {
	gcp, err := conf.GetBool("gcp")
	if err == nil {
		global.GCP = gcp
	}
	bypass, err := conf.GetSlice("trace_bypass_tags")
	if err == nil {
		global.TraceBypassTags = bypass
	}
	tracing, err := conf.GetBool("tracing")
	if err != nil {
		tracing = true
	}
	if !tracing {
		global.Tracing = false
	} else if global.GCP {
		cfg := profiler.Config{
			Service:        "awsl",
			ServiceVersion: "1.0.0",
			// ProjectID must be set if not running on GCP.
			// ProjectID: "my-project",

			// For OpenCensus users:
			// To see Profiler agent spans in APM backend,
			// set EnableOCTelemetry to true
			// EnableOCTelemetry: true,
			MutexProfiling: true,
		}

		// Profiler initialization, best done as early as possible.
		if err := profiler.Start(cfg); err != nil {
			// TODO: Handle error.
			log.Println(err)
		}

		projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
		exporter, err := texporter.NewExporter(texporter.WithProjectID(projectID))
		if err != nil {
			log.Fatalf("texporter.NewExporter: %v", err)
		}
		defer exporter.Shutdown(ctx) // flushes any pending spans

		bsp := sdktrace.NewBatchSpanProcessor(exporter)
		defer bsp.Shutdown(ctx)
		tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp))
		otel.SetTracerProvider(tp)
		global.Tracing = true
	} else if len(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")) != 0 {
		projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
		exporter, err := texporter.NewExporter(texporter.WithProjectID(projectID))
		if err != nil {
			log.Fatalf("texporter.NewExporter: %v", err)
		}
		defer exporter.Shutdown(ctx) // flushes any pending spans

		bsp := sdktrace.NewBatchSpanProcessor(exporter)
		defer bsp.Shutdown(ctx)
		tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp))
		otel.SetTracerProvider(tp)
		global.Tracing = true
	} else {
		global.Tracing = false
	}
}
