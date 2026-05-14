package model

import (
	"context"
	"strings"
)

type middlewareContextKey string

const (
	middlewareStateKey middlewareContextKey = "github.com/stellarlinkco/agentsdk-go/middleware-state"
	// MiddlewareStateKey exposes the context key so other packages can attach middleware state.
	MiddlewareStateKey = middlewareStateKey
)

// TextDeltaHandler 是流式文本增量的回调签名。当模型流式输出文本时,
// runtime 通过 context 传递此回调,供 CompleteStream 内部调用。
type TextDeltaHandler func(delta string) error

type textDeltaHandlerKey string

const (
	textDeltaHandlerCtxKey textDeltaHandlerKey = "github.com/stellarlinkco/agentsdk-go/text-delta-handler"
	// TextDeltaHandlerKey 供外部调用方将 TextDeltaHandler 存入 context。
	TextDeltaHandlerKey = textDeltaHandlerCtxKey
)

// MiddlewareState is the minimal contract required for model providers to
// surface request/response data to middleware consumers without depending on
// the middleware package (which would cause an import cycle).
type MiddlewareState interface {
	SetModelInput(any)
	SetModelOutput(any)
	SetValue(string, any)
}

func middlewareState(ctx context.Context) MiddlewareState {
	if ctx == nil {
		return nil
	}
	state, ok := ctx.Value(middlewareStateKey).(MiddlewareState)
	if !ok {
		return nil
	}
	return state
}

func recordModelRequest(ctx context.Context, req Request) {
	if state := middlewareState(ctx); state != nil {
		state.SetModelInput(req)
	}
}

func recordModelResponse(ctx context.Context, resp *Response) {
	if resp == nil {
		return
	}
	if state := middlewareState(ctx); state != nil {
		state.SetModelOutput(resp)
		state.SetValue("model.response", resp)
		state.SetValue("model.usage", resp.Usage)
		if trimmed := strings.TrimSpace(resp.StopReason); trimmed != "" {
			state.SetValue("model.stop_reason", trimmed)
		}
	}
}
