#!/bin/bash

# grpcurl -plaintext -d '{"name" : "C", "version": "98", "build_script" "run_script": "python3 main.py"}' \
#   localhost:8081 \
#   config.v1.ConfigService.AddLanguage

# grpcurl -plaintext -d '{"id" : "python_3.11.2", "run_script": "#!/bin/bash\n\npython3 main.py"}' \
#   localhost:8081 \
#   config.v1.ConfigService.UpdateLanguage


grpcurl -plaintext -d '{"include_name" : true, "include_version" : true}' \
  localhost:8081 \
  config.v1.ConfigService.GetLanguages


# grpcurl -plaintext -d '{"id" : "test_1"}' \
#   localhost:8081 \
#   config.v1.ConfigService.DeleteLanguage
