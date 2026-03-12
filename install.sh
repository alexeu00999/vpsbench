#!/usr/bin/env bash

# --- VPSBENCH (vpsbench) Installer ---
# Установка в одну команду: wget -qO- https://bench.sh | bash

set -euo pipefail

# Константы
REPO="user/vpsbench"
BINARY_NAME="vpsbench"
INSTALL_DIR="/usr/local/bin"
GITHUB_API="https://api.github.com/repos/${REPO}/releases/latest"

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 1. Определение ОС и Архитектуры
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "${ARCH}" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)
        log_error "Неподдерживаемая архитектура: ${ARCH}"
        exit 1
        ;;
esac

case "${OS}" in
    linux|darwin) ;;
    *)
        log_error "Неподдерживаемая ОС: ${OS}"
        exit 1
        ;;
esac

log_info "Обнаружена система: ${OS}/${ARCH}"

# Для тестов: если установлена переменная DRY_RUN, выходим после определения системы
if [ "${DRY_RUN:-}" == "true" ]; then
    log_success "Тестовый прогон: система определена успешно."
    exit 0
fi

# 2. Получение версии и URL для загрузки
log_info "Поиск последней версии..."
RELEASE_INFO=$(curl -s "${GITHUB_API}")
# Переносимый способ извлечения tag_name без grep -P
VERSION=$(echo "${RELEASE_INFO}" | grep '"tag_name":' | head -n 1 | cut -d '"' -f 4)

if [ -z "${VERSION}" ]; then
    log_warn "Не удалось определить версию через API, использую 'latest'"
    VERSION="latest"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}-${OS}-${ARCH}"
else
    log_info "Найдена версия: ${VERSION}"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}-${OS}-${ARCH}"
fi

# 3. Загрузка бинарного файла
TMP_DIR=$(mktemp -d)
trap 'rm -rf "${TMP_DIR}"' EXIT

log_info "Загрузка ${BINARY_NAME}..."
if ! curl -L -o "${TMP_DIR}/${BINARY_NAME}" "${DOWNLOAD_URL}"; then
    log_error "Ошибка при загрузке. Проверьте соединение или наличие релиза."
    exit 1
fi

chmod +x "${TMP_DIR}/${BINARY_NAME}"

# 4. Установка
log_info "Установка в ${INSTALL_DIR}..."
if [ -w "${INSTALL_DIR}" ]; then
    mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    log_warn "Требуются права sudo для установки в ${INSTALL_DIR}"
    sudo mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

log_success "Установка завершена! Запустите '${BINARY_NAME}' для начала бенчмарка."
