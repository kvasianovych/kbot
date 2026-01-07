/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	telebot "gopkg.in/telebot.v3"
)

var (
	TeleToken = os.Getenv("TELE_TOKEN")

	// Metrics
	messagesProcessed metric.Int64Counter
	messagesDuration  metric.Float64Histogram
)

// kbotCmd represents the kbot command
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Initialize telemetry
		telemetryShutdown, err := InitTelemetry(ctx)
		if err != nil {
			slog.Error("failed to initialize telemetry", "error", err)
			// Continue without telemetry
		} else {
			defer telemetryShutdown.Shutdown(ctx)
		}

		// Initialize metrics
		if err := initMetrics(); err != nil {
			slog.Error("failed to initialize metrics", "error", err)
		}

		slog.Info("kbot started", "version", appVersion)

		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			slog.Error("failed to create bot", "error", err, "hint", "check TELE_TOKEN env variable")
			os.Exit(1)
		}

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {
			return handleMessage(ctx, m)
		})

		// Graceful shutdown
		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			<-sigChan
			slog.Info("shutting down...")
			kbot.Stop()
			cancel()
		}()

		kbot.Start()
	},
}

func initMetrics() error {
	meter := otel.Meter("kbot")

	var err error
	messagesProcessed, err = meter.Int64Counter(
		"kbot.messages.processed",
		metric.WithDescription("Total number of messages processed"),
		metric.WithUnit("{message}"),
	)
	if err != nil {
		return err
	}

	messagesDuration, err = meter.Float64Histogram(
		"kbot.messages.duration",
		metric.WithDescription("Duration of message processing"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return err
	}

	return nil
}

func handleMessage(ctx context.Context, m telebot.Context) error {
	start := time.Now()

	// Create a span for this message handling
	tracer := otel.Tracer("kbot")
	ctx, span := tracer.Start(ctx, "handleMessage",
		trace.WithAttributes(
			attribute.String("telegram.chat.type", string(m.Chat().Type)),
			attribute.Int64("telegram.chat.id", m.Chat().ID),
			attribute.String("telegram.message.payload", m.Message().Payload),
		),
	)
	defer span.End()

	// Log with trace context
	logWithTrace(ctx, "received message",
		"payload", m.Message().Payload,
		"text", m.Text(),
		"chat_id", m.Chat().ID,
	)

	payload := m.Message().Payload
	var err error

	switch payload {
	case "hello":
		err = m.Send(fmt.Sprintf("Hello I'm kbot %s!", appVersion))
		span.SetAttributes(attribute.String("response.type", "hello"))
	default:
		span.SetAttributes(attribute.String("response.type", "none"))
	}

	if err != nil {
		span.RecordError(err)
		logWithTrace(ctx, "failed to send message", "error", err)
	}

	// Record metrics
	duration := float64(time.Since(start).Milliseconds())
	attrs := metric.WithAttributes(
		attribute.String("payload", payload),
		attribute.Bool("success", err == nil),
	)

	if messagesProcessed != nil {
		messagesProcessed.Add(ctx, 1, attrs)
	}
	if messagesDuration != nil {
		messagesDuration.Record(ctx, duration, attrs)
	}

	return err
}

// logWithTrace logs a message with trace context for correlation
func logWithTrace(ctx context.Context, msg string, args ...any) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		args = append(args,
			"traceId", span.SpanContext().TraceID().String(),
			"spanId", span.SpanContext().SpanID().String(),
		)
	}
	slog.Info(msg, args...)
}

func init() {
	rootCmd.AddCommand(kbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
