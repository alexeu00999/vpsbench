#!/usr/bin/env bash

# --- Тест скрипта установки (Симуляция ОС и Архитектуры) ---

set -euo pipefail

# Создаем временную директорию для тестов
TEST_DIR=$(mktemp -d)
trap 'rm -rf "${TEST_DIR}"' EXIT

# Копируем скрипт установки
cp install.sh "${TEST_DIR}/install.sh"
chmod +x "${TEST_DIR}/install.sh"

# Функция для запуска теста с подменой uname
test_detection() {
    local mock_os=$1
    local mock_arch=$2
    local expected_msg=$3

    echo "----------------------------------------------------"
    echo "Тест: ОС=${mock_os}, Архитектура=${mock_arch}"
    
    # Создаем временный файл-обертку для uname
    echo -e "#!/bin/bash\nif [ \"\$1\" == \"-s\" ]; then echo \"${mock_os}\"; elif [ \"\$1\" == \"-m\" ]; then echo \"${mock_arch}\"; fi" > "${TEST_DIR}/uname"
    chmod +x "${TEST_DIR}/uname"

    # Запускаем скрипт, подменяя PATH, чтобы наш uname был первым
    local output
    output=$(DRY_RUN=true PATH="${TEST_DIR}:${PATH}" bash "${TEST_DIR}/install.sh" 2>&1)
    
    if echo "${output}" | grep -q "${expected_msg}"; then
        echo -e "\033[0;32m[PASS]\033[0m Найдено ожидаемое сообщение: ${expected_msg}"
    else
        echo -e "\033[0;31m[FAIL]\033[0m Сообщение '${expected_msg}' не найдено!"
        echo "ВЫВОД СКРИПТА:"
        echo "${output}"
        exit 1
    fi
}

# Сценарии тестов
test_detection "Linux" "x86_64" "Обнаружена система: linux/amd64"
test_detection "Darwin" "arm64" "Обнаружена система: darwin/arm64"
test_detection "Linux" "aarch64" "Обнаружена система: linux/arm64"

echo "----------------------------------------------------"
echo -e "\033[0;32mВсе тесты обнаружения системы пройдены!\033[0m"

# Проверка ShellCheck (если установлен)
if command -v shellcheck &> /dev/null; then
    echo "Запуск ShellCheck..."
    shellcheck install.sh && echo -e "\033[0;32mShellCheck пройден!\033[0m" || exit 1
else
    echo "ShellCheck не установлен, пропуск статического анализа."
fi
