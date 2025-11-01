#!/bin/sh
set -e

echo "Installing dependencies..."
apk add --no-cache curl jq

echo "Waiting for Hydra to be ready..."
until curl -s "$HYDRA_ADMIN_URL/health/ready" | jq -e '.status=="ok"' > /dev/null; do
  sleep 2
done

echo "Checking if OAuth2 client already exists..."
HTTP_RESPONSE=$(curl --silent --location --write-out "%{http_code}" --output /tmp/client.json \
  --header "Accept: application/json" \
  "$HYDRA_ADMIN_URL/admin/clients/$HYDRA_CLIENT_ID")

if [ "$HTTP_RESPONSE" = "200" ]; then
  echo "Client $HYDRA_CLIENT_ID already exists. Skipping creation."
  exit 0
elif [ "$HTTP_RESPONSE" = "404" ]; then
  echo "Client not found. Creating new client..."
  
  CREATE_RESPONSE=$(curl --silent --location --write-out "%{http_code}" --output /tmp/create_client.json \
    --header "Content-Type: application/json" \
    --header "Accept: application/json" \
    --data '{
      "client_id": "'"$HYDRA_CLIENT_ID"'",
      "client_name": "'"$HYDRA_CLIENT_NAME"'",
      "client_secret": "'"$HYDRA_CLIENT_SECRET"'",
      "grant_types": ["client_credentials"],
      "scope": "offline",
      "token_endpoint_auth_method": "client_secret_post"
    }' \
    "$HYDRA_ADMIN_URL/admin/clients")

  if [ "$CREATE_RESPONSE" = "201" ]; then
    echo "Client created successfully!"
  else
    echo "Failed to create client. HTTP Response: $CREATE_RESPONSE"
    echo "Error details:"
    cat /tmp/create_client.json
    exit 1
  fi
else
  echo "Error: Unexpected response from Hydra ($HTTP_RESPONSE). Check Hydra logs for details."
  echo "Error details:"
  cat /tmp/client.json
  exit 1
fi
