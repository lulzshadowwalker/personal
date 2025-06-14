# Run templ generation in watch mode
templ:
	templ generate --watch --proxy="http://localhost:3000" --open-browser=true

# Run air for Go hot reload
server:
	air \
	--build.cmd "go build -o tmp/bin/main ./cmd/http/main.go" \
	--build.bin "tmp/bin/main" \
	--build.delay "50" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# Watch Tailwind CSS changes
tailwind:
	tailwindcss -i ./internal/http/template/css/app.css -o ./cmd/http/public/css/style.css --watch --minify

# Start development server with all watchers
dev:
	make -j3 tailwind templ server

templui:
	@echo "Installing TemplUI"
	go install github.com/axzilla/templui/cmd/templui@latest
	@echo "Checking installation"
	templui -v

