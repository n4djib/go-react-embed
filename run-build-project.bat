@echo off
cd frontend && npm run build && cd.. && go build ./cmd/go-react-embed
