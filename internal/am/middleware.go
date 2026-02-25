package am

func MessageStreamWithMiddleware(stream MessageStream, mws ...MessageStreamMiddleware) MessageStream {
	s := stream

	// middleware are applied in reverse
	for i := len(mws) - 1; i >= 0; i-- {
		s = mws[i](s)
	}

	return s
}

func MessagePublisherWithMiddleware(publisher MessagePublisher, mws ...MessagePublisherMiddleware) MessagePublisher {
	p := publisher

	// middleware are applied in reverse
	for i := len(mws) - 1; i >= 0; i-- {
		p = mws[i](p)
	}

	return p
}

func MessageHandlerWithMiddleware(handler MessageHandler, mws ...MessageHandlerMiddleware) MessageHandler {
	h := handler

	// middleware are applied in reverse
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}
