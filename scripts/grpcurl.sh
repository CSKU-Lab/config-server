#!/bin/bash

# grpcurl -plaintext -d '{"id" : "test_1", "name": "test", "version" : "1"}' \
#   localhost:8081 \
#   config.v1.ConfigService.UpdateLanguage


# grpcurl -plaintext -d '{"include_name" : true, "include_version" : true}' \
#   localhost:8081 \
#   config.v1.ConfigService.GetLanguages
