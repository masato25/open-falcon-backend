find ./  -type f -name "*.go" |grep -v "./package_cache_tmp"|grep -v "./vendor"|grep -v "./scripts" | \
  xargs sed -i 's/open-falcon\/falcon-plus/Cepave\/open-falcon-backend/g'

find ./  -type f -name "*.go" | \
  xargs sed -i 's/modules\/api\//modules\/f2e-api\//g'
find ./  -type f -name "*.go" | \
    xargs sed -i 's/github.com\/gin-gonic\/gin/gopkg.in\/gin-gonic\/gin.v1/g'
