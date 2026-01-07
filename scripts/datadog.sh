#!/bin/bash

# DataDog Agent Installation Script
# Installs the latest DataDog agent for the detected platform

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if DD_API_KEY is set
check_api_key() {
    if [[ -z "${DD_API_KEY}" ]]; then
        log_error "DD_API_KEY environment variable is required but not set"
        log_info "Get your API key from: https://app.datadoghq.com/organization-settings/api-keys"
        log_info "Export it with: export DD_API_KEY=your_api_key_here"
        exit 1
    fi
    log_info "DD_API_KEY is set"
}

# Detect platform
detect_platform() {
    case "$(uname -s)" in
        Linux*)
            PLATFORM="linux"
            ;;
        Darwin*)
            PLATFORM="macos"
            ;;
        *)
            log_error "Unsupported platform: $(uname -s)"
            exit 1
            ;;
    esac
    log_info "Detected platform: $PLATFORM"
}

# Install DataDog agent
install_datadog() {
    local dd_site="${DD_SITE:-datadoghq.com}"
    log_info "Installing DataDog agent for $PLATFORM"
    log_info "Using DD_SITE: $dd_site"

    case "$PLATFORM" in
        linux)
            log_info "Downloading and running DataDog installation script for Linux..."
            bash -c "$(curl -L https://s3.amazonaws.com/dd-agent/scripts/install_script.sh)"
            ;;
        macos)
            log_info "Downloading and running DataDog installation script for macOS..."
            bash -c "$(curl -L https://s3.amazonaws.com/dd-agent/scripts/install_mac_os.sh)"
            ;;
    esac
}

# Configure StatsD
configure_statsd() {
    log_info "Configuring StatsD for Temporal metrics..."
    
    case "$PLATFORM" in
        linux)
            local config_file="/etc/datadog-agent/datadog.yaml"
            ;;
        macos)
            local config_file="/opt/datadog-agent/etc/datadog.yaml"
            ;;
    esac

    if [[ -f "$config_file" ]]; then
        log_info "Enabling DogStatsD in $config_file"
        # Backup original config
        sudo cp "$config_file" "$config_file.backup"
        
        # Enable DogStatsD if not already enabled
        if ! grep -q "^dogstatsd_port:" "$config_file"; then
            echo "dogstatsd_port: 8125" | sudo tee -a "$config_file" > /dev/null
        fi
        if ! grep -q "^use_dogstatsd:" "$config_file"; then
            echo "use_dogstatsd: true" | sudo tee -a "$config_file" > /dev/null
        fi
        
        log_info "DogStatsD configuration added"
    else
        log_warn "DataDog config file not found at $config_file"
    fi
}

# Start DataDog agent
start_agent() {
    log_info "Starting DataDog agent..."
    
    case "$PLATFORM" in
        linux)
            sudo systemctl start datadog-agent
            sudo systemctl enable datadog-agent
            log_info "DataDog agent started and enabled"
            ;;
        macos)
            sudo launchctl load -w /Library/LaunchDaemons/com.datadoghq.agent.plist
            log_info "DataDog agent started"
            ;;
    esac
}

# Check agent status
check_status() {
    log_info "Checking DataDog agent status..."
    
    case "$PLATFORM" in
        linux)
            sudo datadog-agent status
            ;;
        macos)
            sudo /opt/datadog-agent/bin/agent/agent status
            ;;
    esac
}

# Main installation flow
main() {
    log_info "DataDog Agent Installation Script"
    log_info "================================="
    
    check_api_key
    detect_platform
    install_datadog
    configure_statsd
    start_agent
    
    log_info ""
    log_info "Installation complete! 🎉"
    log_info ""
    log_info "Next steps:"
    log_info "1. Run 'make metrics-dogstatsd' to start Temporal with DogStatsD metrics"
    log_info "2. Visit your DataDog dashboard to see metrics"
    log_info "3. Check agent status with the command below:"
    log_info ""
    
    case "$PLATFORM" in
        linux)
            log_info "   sudo datadog-agent status"
            ;;
        macos)
            log_info "   sudo /opt/datadog-agent/bin/agent/agent status"
            ;;
    esac
}

# Show usage if help requested
if [[ "$1" == "--help" || "$1" == "-h" ]]; then
    echo "DataDog Agent Installation Script"
    echo ""
    echo "Usage: $0"
    echo ""
    echo "Environment variables:"
    echo "  DD_API_KEY  - DataDog API key (required)"
    echo "  DD_SITE     - DataDog site (optional, defaults to datadoghq.com)"
    echo ""
    echo "Example:"
    echo "  export DD_API_KEY=your_api_key_here"
    echo "  export DD_SITE=datadoghq.eu  # for EU site"
    echo "  $0"
    exit 0
fi

# Run main function
main "$@"
