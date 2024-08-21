GO_PATH = $$GOPATH

ifeq ($(OS), Windows_NT)
	SHELL := powershell.exe
	.SHELLFLAGS := -NoProfile -Command
	SHELL_VERSION = $(shell (Get-Host | Select-Object Version | Format-Table -HideTableHeaders | Out-String).Trim({}))
	OS = $(shell "{0} {1}" -f "windows", (Get-ComputerInfo -Property OsVersion, OsArchitecture | Format-Table -HideTableHeaders | Out-String).Trim({}))
	PACKAGE = $(shell (Get-Content go.mod -head 1).Split(" ")[1])
	HELP_CMD = Select-String "^[a-zA-Z_-]+:.*?\#\# .*$$" "./Makefile" | Foreach-Object { $$_data = $$_.matches -split ":.*?\#\# "; $$obj = New-Object PSCustomObject; Add-Member -InputObject $$obj -NotePropertyName ('Command') -NotePropertyValue $$_data[0]; Add-Member -InputObject $$obj -NotePropertyName ('Description') -NotePropertyValue $$_data[1]; $$obj } | Format-Table -HideTableHeaders @{Expression={ $$e = [char]27; "$$e[36m$$($$_.Command)$${e}[0m" }}, Description
	RM_F_CMD = Remove-Item -erroraction silentlycontinue -Force
	RM_RF_CMD = ${RM_F_CMD} -Recurse
else
	SHELL := bash
	SHELL_VERSION = $(shell echo $$BASH_VERSION)
	UNAME := $(shell uname -s)
	VERSION_AND_ARCH = $(shell uname -rm)
	ifeq ($(UNAME),Darwin)
		OS = macos ${VERSION_AND_ARCH}
	else ifeq ($(UNAME),Linux)
		OS = linux ${VERSION_AND_ARCH}
	else
    $(error OS not supported by this Makefile)
	endif
	PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')
	HELP_CMD = grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	RM_F_CMD = rm -f
	RM_RF_CMD = ${RM_F_CMD} -r
endif

.DEFAULT_GOAL := help

generate-migration: ## Generate migration file example: generate-migration folder name=migration_name
	@[ "${name}" ] || ( echo "migration name not set"; exit 1 )
	migrate create -ext sql -dir entities/migrations $(name)

migrate-up: ## Migrate Up example: migrate-up database="postgres://username:password@host:port/db_name?sslmode=disable" force=boolean version=xxxx
	@[ "${database}" ] || ( echo "database not set"; exit 1 )
ifeq ($(force), true)
	migrate -path entities/migrations -database "$(database)" force $(version) -verbose up $(N)
else
	migrate -path entities/migrations -database "$(database)" -verbose up $(N)
endif

migrate-down: ## Migrate Down example: migrate-down database="postgres://username:password@host:port/db_name?sslmode=disable" force=boolean version=xxxx
	@[ "${database}" ] || ( echo "database not set"; exit 1 )
ifeq ($(force), true)
	migrate -path entities/migrations -database "$(database)" force $(version) -verbose down $(N)
else
	migrate -path entities/migrations -database "$(database)" -verbose down $(N)
endif

.PHONY: migrate-up migrate-down

run: ## Run app example: make run app=app_name watch=boolean
ifeq ($(watch), true)
	gow run main.go
else
	go run main.go
endif