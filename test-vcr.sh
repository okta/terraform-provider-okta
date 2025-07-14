#!/bin/bash

# Test script for VCR functionality with resource set datasources
# This script helps test the VCR recording and playback functionality

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== VCR Test Script for Resource Set Datasources ===${NC}"

# Function to run a test
run_test() {
    local test_name=$1
    local mode=$2
    local cassette=$3

    echo -e "\n${YELLOW}Running test: ${test_name} in ${mode} mode${NC}"

    # Set environment variables
    export TF_LOG=DEBUG
    export TF_LOG_PATH="../tflog.log"

    if [ "$mode" = "live" ]; then
        unset OKTA_VCR_CASSETTE
        unset OKTA_VCR_TF_ACC
        echo "Running live test (no VCR)"
    elif [ "$mode" = "record" ]; then
        export OKTA_VCR_CASSETTE="$cassette"
        export OKTA_VCR_TF_ACC="record"
        echo "Recording VCR cassette: $cassette"
    elif [ "$mode" = "play" ]; then
        export OKTA_VCR_CASSETTE="$cassette"
        export OKTA_VCR_TF_ACC="play"
        echo "Playing VCR cassette: $cassette"
    fi

    # Clean up previous logs
    rm -f test-output.log
    rm -f tflog.log

    # Run the test
    echo "Executing test..."
    op run --account="NDM32WGBUND3XEEL5OEN22MMOE" --env-file="./env_trial-5309542.okta.com.env" -- \
        make testacc TEST=./okta TESTARGS="-run=${test_name}" 2>&1 | tee test-output.log

    # Check if test passed
    if grep -q "PASS" test-output.log; then
        echo -e "${GREEN}✓ Test ${test_name} passed in ${mode} mode${NC}"
        return 0
    else
        echo -e "${RED}✗ Test ${test_name} failed in ${mode} mode${NC}"
        return 1
    fi
}

# Function to clean up VCR cassettes
cleanup_cassettes() {
    local test_name=$1
    echo -e "\n${YELLOW}Cleaning up VCR cassettes for ${test_name}${NC}"
    rm -rf "test/fixtures/vcr/idaas/${test_name}"
}

# Test cases
test_cases=(
    "TestAccDataSourceOktaResourceSets_read"
    "TestAccDataSourceOktaResourceSets_multiple"
    "TestAccDataSourceOktaResourceSet_"
    "TestAccDataSourceOktaResourceSetResources_"
)

# Main test loop
for test_name in "${test_cases[@]}"; do
    echo -e "\n${YELLOW}=== Testing ${test_name} ===${NC}"

    # Clean up any existing cassettes
    cleanup_cassettes "$test_name"

    # Test 1: Run live (no VCR)
    if run_test "$test_name" "live" ""; then
        echo -e "${GREEN}Live test passed${NC}"
    else
        echo -e "${RED}Live test failed, skipping VCR tests${NC}"
        continue
    fi

    # Test 2: Record VCR
    if run_test "$test_name" "record" "oie-00"; then
        echo -e "${GREEN}VCR recording passed${NC}"
    else
        echo -e "${RED}VCR recording failed${NC}"
        continue
    fi

    # Test 3: Play VCR
    if run_test "$test_name" "play" "oie-00"; then
        echo -e "${GREEN}VCR playback passed${NC}"
    else
        echo -e "${RED}VCR playback failed${NC}"
    fi

    echo -e "${GREEN}✓ All tests for ${test_name} completed${NC}"
done

echo -e "\n${GREEN}=== VCR Test Script Completed ===${NC}"
