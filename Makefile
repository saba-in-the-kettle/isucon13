# ↓↓↓当日いじる↓↓↓
# アプリケーション
BUILD_DIR:=./go# ローカルマシンのMakefileからの相対パス
BIN_NAME:=isupipe
SERVER_BINARY_DIR:=~/webapp/go
SERVER_SQL_DIR:=/home/isucon/webapp
SERVICE_NAME:=isupipe-go.service
APP_SERVER_NAMES="s1\ns2\ns3\n"
SERVER_NAMES="s1\ns2\ns3"
# ↑↑↑ここまで↑↑↑

# colors
ESC=$(shell printf '\033')
RESET="${ESC}[0m"
BOLD="${ESC}[1m"
RED="${ESC}[31m"
GREEN="${ESC}[32m"
BLUE="${ESC}[33m"

# commands
START_ECHO=echo "$(GREEN)$(BOLD)[INFO] start $@ $$s $(RESET)"

.PHONY: build
build:
	@ $(START_ECHO);\
	cd $(BUILD_DIR); \
	GOOS=linux GOARCH=amd64 go build -o $(BIN_NAME) *.go

.PHONY: deploy-app
deploy-app: build
	echo ${APP_SERVER_NAMES} | xargs -P 0 -I {} make deploy-app-one s={}

deploy-app-one:
	@$(START_ECHO)
	ssh $s "sudo systemctl daemon-reload & sudo systemctl stop $(SERVICE_NAME)"
	\rsync -avz $(BUILD_DIR)/$(BIN_NAME) $s:$(SERVER_BINARY_DIR)/
	ssh $s "chmod +x $(SERVER_BINARY_DIR)/$(BIN_NAME)"
	ssh $s "sudo systemctl start $(SERVICE_NAME)"

.PHONY: deploy-sql
deploy-sql:
	echo ${APP_SERVER_NAMES} | xargs -P 0 -I {} make deploy-sql-one s={}

deploy-sql-one:
	@$(START_ECHO)
	\rsync -avz sql $s:${SERVER_SQL_DIR}

.PHONY: deploy-config
deploy-config:
	echo ${SERVER_NAMES} | xargs -P 0 -I {} bash $(CURDIR)/{}/deploy.sh

.PHONY: deploy
deploy: deploy-config deploy-sql deploy-app


# hargo用
HARGO_ORIGINAL_NAME:=original.har
HARGO_FILTERED_NAME:=filtered.har

.PHONY: hargo-filter
hargo-filter:
	@$(START_ECHO)
	hargo filter $(HARGO_ORIGINAL_NAME) $(HARGO_FILTERED_NAME)

.PHONY: hargo-fetch
hargo-fetch:
	@$(START_ECHO)
	hargo fetch $(HARGO_FILTERED_NAME) -dir hargo-results
