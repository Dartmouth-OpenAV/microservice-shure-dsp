#!/bin/bash

# Shure DSP Microservice API Test Script
# This script tests all endpoints from the Shure DSP Microservice Postman collection

# Configuration variables - Update these values as needed
MICROSERVICE_URL="localhost:8080"
DEVICE_FQDN="shure-dsp.local"
CHANNEL_ID="11"
MATRIX_INPUT="01"
MATRIX_OUTPUT="17"

echo "Starting Shure DSP Microservice API Tests..."
echo "Microservice URL: $MICROSERVICE_URL"
echo "Device FQDN: $DEVICE_FQDN"
echo "Channel ID: $CHANNEL_ID"
echo "Matrix Input: $MATRIX_INPUT"
echo "Matrix Output: $MATRIX_OUTPUT"
echo "=============================================="

# GET Volume
echo "Testing GET Volume for Channel $CHANNEL_ID..."
curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/volume/$CHANNEL_ID"
sleep 1

# GET Audiomute
echo "Testing GET Audiomute for Channel $CHANNEL_ID..."
curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/audiomute/$CHANNEL_ID"
sleep 1

# GET Matrixmute
echo "Testing GET Matrixmute (Input: $MATRIX_INPUT, Output: $MATRIX_OUTPUT)..."
curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/matrixmute/$MATRIX_INPUT/$MATRIX_OUTPUT"
sleep 1

# GET Matrixvolume
echo "Testing GET Matrixvolume (Input: $MATRIX_INPUT, Output: $MATRIX_OUTPUT)..."
curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/matrixvolume/$MATRIX_INPUT/$MATRIX_OUTPUT"
sleep 1

echo "=============================================="
echo "Starting SET/PUT operations..."
echo "=============================================="

# SET Volume
echo "Testing SET Volume for Channel $CHANNEL_ID (35)..."
curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/volume/$CHANNEL_ID" \
     -H "Content-Type: application/json" \
     -d "\"35\""
sleep 1

# SET Audiomute to true
echo "Testing SET Audiomute for Channel $CHANNEL_ID (true)..."
curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/audiomute/$CHANNEL_ID" \
     -H "Content-Type: application/json" \
     -d "\"true\""
sleep 1

# SET Audiomute to false
echo "Testing SET Audiomute for Channel $CHANNEL_ID (false)..."
curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/audiomute/$CHANNEL_ID" \
     -H "Content-Type: application/json" \
     -d "\"false\""
sleep 1

# SET Matrixmute to true
echo "Testing SET Matrixmute (Input: $MATRIX_INPUT, Output: $MATRIX_OUTPUT) to true..."
curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/matrixmute/$MATRIX_INPUT/$MATRIX_OUTPUT" \
     -H "Content-Type: application/json" \
     -d "\"true\""
sleep 1

# SET Matrixmute to false
echo "Testing SET Matrixmute (Input: $MATRIX_INPUT, Output: $MATRIX_OUTPUT) to false..."
curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/matrixmute/$MATRIX_INPUT/$MATRIX_OUTPUT" \
     -H "Content-Type: application/json" \
     -d "\"false\""
sleep 1

# SET Matrixvolume
echo "Testing SET Matrixvolume (Input: $MATRIX_INPUT, Output: $MATRIX_OUTPUT) to 35..."
curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/matrixvolume/$MATRIX_INPUT/$MATRIX_OUTPUT" \
     -H "Content-Type: application/json" \
     -d "\"35\""
sleep 1

echo "=============================================="
echo "Testing additional volume levels..."
echo "=============================================="

# Test different volume levels
echo "Testing SET Volume for Channel $CHANNEL_ID (50)..."
curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/volume/$CHANNEL_ID" \
     -H "Content-Type: application/json" \
     -d "\"50\""
sleep 1

echo "Testing SET Volume for Channel $CHANNEL_ID (25)..."
curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/volume/$CHANNEL_ID" \
     -H "Content-Type: application/json" \
     -d "\"25\""
sleep 1

echo "Testing SET Matrixvolume (Input: $MATRIX_INPUT, Output: $MATRIX_OUTPUT) to 50..."
curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/matrixvolume/$MATRIX_INPUT/$MATRIX_OUTPUT" \
     -H "Content-Type: application/json" \
     -d "\"50\""
sleep 1

echo "Testing SET Matrixvolume (Input: $MATRIX_INPUT, Output: $MATRIX_OUTPUT) to 25..."
curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/matrixvolume/$MATRIX_INPUT/$MATRIX_OUTPUT" \
     -H "Content-Type: application/json" \
     -d "\"25\""
sleep 1

echo "=============================================="
echo "Final state check - Getting current values..."
echo "=============================================="

# Final state check
echo "Final GET Volume for Channel $CHANNEL_ID..."
curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/volume/$CHANNEL_ID"
sleep 1

echo "Final GET Audiomute for Channel $CHANNEL_ID..."
curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/audiomute/$CHANNEL_ID"
sleep 1

echo "Final GET Matrixmute (Input: $MATRIX_INPUT, Output: $MATRIX_OUTPUT)..."
curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/matrixmute/$MATRIX_INPUT/$MATRIX_OUTPUT"
sleep 1

echo "Final GET Matrixvolume (Input: $MATRIX_INPUT, Output: $MATRIX_OUTPUT)..."
curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/matrixvolume/$MATRIX_INPUT/$MATRIX_OUTPUT"
sleep 1

echo "=============================================="
echo "All Shure DSP Microservice API tests completed!"
echo "=============================================="
