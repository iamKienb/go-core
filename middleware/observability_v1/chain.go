package observability_v1

// func ServerInterceptors(logger *slog.Logger, tail ...connect.Interceptor) []connect.Interceptor {
// 	interceptors := make([]connect.Interceptor, 0, 5+len(tail))
// 	interceptors = appendTracingInterceptor(interceptors, logger)
// 	interceptors = append(interceptors,
// 		authx.RequestContextInterceptor(),
// 		ErrorResponseInterceptor(logger),
// 		RecoveryInterceptor(logger),
// 		LoggingInterceptor(logger),
// 	)
// 	interceptors = append(interceptors, tail...)
// 	return interceptors
// }

// func InternalServerInterceptors(logger *slog.Logger) []connect.Interceptor {
// 	return ServerInterceptors(
// 		logger,
// 		authx.AuthInternalInterceptor(),
// 		ValidationRequestInterceptor(),
// 	)
// }

// func ClientInterceptors(logger *slog.Logger, interceptors ...connect.Interceptor) []connect.Interceptor {
// 	result := make([]connect.Interceptor, 0, 2+len(interceptors))
// 	result = appendTracingInterceptor(result, logger)
// 	result = append(result, interceptors...)
// 	result = append(result, LoggingInterceptor(logger))
// 	return result
// }

// func ServerOption(logger *slog.Logger, tail ...connect.Interceptor) connect.Option {
// 	return connect.WithInterceptors(ServerInterceptors(logger, tail...)...)
// }

// func InternalServerOption(logger *slog.Logger) connect.Option {
// 	return connect.WithInterceptors(InternalServerInterceptors(logger)...)
// }

// func ClientOption(logger *slog.Logger, interceptors ...connect.Interceptor) connect.Option {
// 	return connect.WithInterceptors(ClientInterceptors(logger, interceptors...)...)
// }

// func appendTracingInterceptor(interceptors []connect.Interceptor, logger *slog.Logger) []connect.Interceptor {
// 	tracingInterceptor, err := TracingInterceptor()
// 	if err != nil {
// 		if logger != nil {
// 			logger.Error("failed to initialize tracing interceptor", slog.Any("error", err))
// 		}
// 		return interceptors
// 	}

// 	return append(interceptors, tracingInterceptor)
// }
