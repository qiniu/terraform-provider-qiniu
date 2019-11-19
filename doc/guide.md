# Qiniu Cloud Provider 用户指南

## 构建 Qiniu Cloud Provider

下载 terraform-provider-qiniu 源码，构建后放入当前用户的 Terraform 插件目录

```bash
git clone git@github.com:qiniu/terraform-provider-qiniu.git --recurse-submodules
cd terraform-provider-qiniu
make
mkdir -p ~/.terraform.d/plugins
mv bin/terraform-provider-qiniu ~/.terraform.d/plugins
```

在 Windows 中，当前用户的 Terraform 插件目录路径为 `%APPDATA%\terraform.d\plugins`。

## 准备 Terraform 执行计划

1. 根据 [官方指南](https://www.terraform.io/downloads.html) 安装 Terraform
2. 创建一个 Terraform 执行计划目录
3. 创建 `main.tf` 文件，写入资源创建信息（见[使用案例](#使用案例)）
4. 在执行计划目录执行 `terraform init` 命令初始化执行计划
5. 在执行计划目录执行 `terraform plan` 命令准备执行计划
6. 在执行计划目录执行 `terraform apply` 命令正式执行计划

## 使用案例

```hcl
# Configure Qiniu Account
provider "qiniu" {
  access_key = "<Qiniu Access Key>"
  secret_key = "<Qiniu Secret Key>"
}

# Create Qiniu Bucket
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"  # Bucket Name
  region_id = "z0"                      # Bucket Region, "z0" means East China
  private   = false                     # Public bucket
}

# Create Qiniu Object
resource "qiniu_bucket_object" "basic_object" {
  bucket    = "basic-test-terraform-1"  # Bucket Name
  key       = "qiniu-key"               # File Key
  source    = "/path/to/file"           # File Path to upload
}

# Qiniu Regions Data Source
data "qiniu_regions" "all" {
}

# Qiniu Buckets Data Source
data "qiniu_buckets" "z1" {
  name_regex = "^bucket-"
  region_id = "z1"
}

# Qiniu Buckets Objects Data Source
data "qiniu_buckets_objects" "all" {
  bucket = "basic-test-terraform-1"
}
```

## 认证

支持静态认证和环境变量认证两种方式

### 静态认证

```hcl
provider "qiniu" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}
```

### 环境变量认证

```hcl
provider "qiniu" {
}
```

```bash
export QINIU_ACCESS_KEY="anaccesskey"
export QINIU_SECRET_KEY="asecretkey"
terraform plan
```

## 私有云配置

在私有云情况下，需要手工配置各个私有云服务器的地址

```hcl
provider "qiniu" {
  use_https      = true
  central_rs_url = "https://<中心化 RS 域名>"
  rs_url         = "https://<RS 域名>"
  rsf_url        = "https://<RSF 域名>"
  up_url         = "https://<UP 域名>"
  api_url        = "https://<API 域名>"
  io_url         = "https://<IO 域名>"
}
```

## 资源创建

### 创建 Bucket

创建 Bucket 时，`name` 和 `region_id` 是必须指定的

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"  # Bucket 名称
  region_id = "z0"                      # Bucket 区域, "z0" 表示华东地区
}
```

可以额外指定 `private` 参数创建私有 Bucket

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"
  region_id = "z0"
  private   = true
}
```

可以额外指定 `image_url` 和 `image_host` 参数指定镜像回源地址

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"
  region_id = "z0"
  image_url = "http://target_url"
}
```

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name       = "basic-test-terraform-1"
  region_id  = "z0"
  image_url  = "http://target_url"
  image_host = "www.qiniu.com"          # 该参数可选，仅在给出了 image_url 参数之后才生效
}
```

可以额外指定 `index_page_on` 参数设置默认首页

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"
  region_id = "z0"
  index_page_on = true
}
```

可以额外指定 `lifecycle_rules` 参数设置生命周期规则

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"
  region_id = "z0"
  lifecycle_rules {
    name = "rule_1"                      # 规则名称，该参数必填
    prefix = "users/"                    # 规则匹配的对象名称前缀
    to_line_after_days = 20              # 用户新创建的文件将在指定天数后自动转为低频存储
    delete_after_days = 30               # 用户新创建的文件将在指定天数后自动删除
  }

  lifecycle_rules {                      # 可以设置多条生命周期规则
    name = "rule_2"                      # 规则名称，该参数必填
    prefix = "admin/"                    # 规则匹配的对象名称前缀
    to_line_after_days = 40              # 用户新创建的文件将在指定天数后自动转为低频存储
    delete_after_days = 60               # 用户新创建的文件将在指定天数后自动删除
  }
}
```

可以额外指定 `cors_rules` 参数设置跨域规则（如果不设置，则默认允许任何跨域请求）

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"
  region_id = "z0"
  cors_rules {
    allowed_origins = ["http://www.test1.com"]     # 允许的域名列表，该参数必填，支持通配符 *
    allowed_methods = ["GET", "POST"]              # 允许的 HTTP 方法列表，该参数必填，不支持通配符
    allowed_headers = ["X-Reqid", "Content-Type"]  # 允许的 HTTP 头列表，支持通配符 *，但只能是用单独的 * 表示允许全部 HTTP 头，而不能部分匹配，如果为空或不填，则表示不允许任何 HTTP 头
    exposed_headers = ["X-Test-1", "X-Test-2"]     # 暴露的 HTTP 头列表，不支持通配符，X-Log，X-Reqid 是默认的会暴露的 HTTP 头
    max_age = 20                                   # 结果可以缓存的时间，单位为秒。如果为空或不填，则不缓存
  }

  cors_rules {                                     # 可以设置多条跨域规则，但不能超过 10 条
    allowed_origins = ["http://www.test2.com"]     # 允许的域名列表，该参数必填，支持通配符 *
    allowed_methods = ["GET", "POST", "HEAD"]      # 允许的 HTTP 方法列表，该参数必填，不支持通配符
                                                   # 其他参数均可以省略
  }
}
```

可以额外指定 `anti_leech_mode`，`allow_empty_referer`，`referer_pattern`，`only_enable_anti_leech_for_cdn` 等参数设置 Referer 防盗链

`anti_leech_mode` 参数表示设置的防盗链模式，总共有两种模式可以选择，白名单模式和黑名单模式

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"
  region_id = "z0"
  anti_leech_mode = "blacklist"                             # 设置 Referer 黑名单，表示凡是能匹配 Referer 规则的域名均被禁止访问，该参数必填
  referer_pattern = ["foo.com", "*.bar.com", "sub.foo.com"] # Referer 规则，支持通配符 *
  allow_empty_referer = true                                # 是否允许空 Referer
  only_enable_anti_leech_for_cdn = false                    # 是否开启源站防盗链，默认只会为 CDN 请求配置防盗链
}
```

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"
  region_id = "z0"
  anti_leech_mode = "whitelist"                             # 设置 Referer 白名单，表示只有能匹配 Referer 规则的域名才被允许访问，该参数必填
  referer_pattern = ["*"]                                   # Referer 规则，支持通配符 *
                                                            # 其他参数均可以省略
}
```

可以额外指定 `max_age` 参数设置文件客户端缓存时间

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"
  region_id = "z0"
  max_age = 31536000                       # 允许客户端缓存时长，单位为秒
}
```

可以额外指定 `tagging` 参数设置 Bucket 标签

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"
  region_id = "z0"
  tagging = {
      env = "test"                         # 设置 Bucket 第一个标签
      kind = "basic"                       # 设置 Bucket 第二个标签
                                           # 每个 Bucket 不能设置超过十对标签
  }
}
```

### 上传文件

可以将指定路径的文件上传至指定 Bucket

```hcl
resource "qiniu_bucket" "basic_object" {
  bucket        = "basic-test-terraform-1"  # Bucket 名称
  key           = "keyname"                 # 对象名称
  source        = "/path/to/file"           # 源文件路径
  storage_type  = "infrequent"              # 存储类型，只能填写 infrequent 表示低频存储，如果为空或不填则表示普通存储
}
```

## 数据源

### 列出区域

列出所有七牛区域

```hcl
data "qiniu_regions" "all" {
}
```

### 列出 Bucket

列出当前账户所有 Bucket

```hcl
data "qiniu_buckets" "all" {
}
```

列出匹配指定正则表达式的 Bucket

```hcl
data "qiniu_buckets" "all" {
  name_regex = "^data-"
}
```

列出指定区域的 Bucket

```hcl
data "qiniu_buckets" "all" {
  region_id = "z0"
}
```

可以获取到所有 Bucket 的名称

```hcl
${data.qiniu_buckets.all.names}
```

也可以获取到所有 Bucket 的详细信息

```hcl
${data.qiniu_buckets.all.buckets}
```

### 列出文件

列出指定 Bucket 的文件

```hcl
data "qiniu_buckets_objects" "all" {
  bucket = "basic-test-terraform-1"
}
```

列出指定 Bucket 中所有文件名以指定字符串开头的文件

```hcl
data "qiniu_buckets_objects" "all" {
  bucket = "basic-test-terraform-1"
  prefix = "data-"
}
```

可以获取到所有文件的名称

```hcl
${data.qiniu_buckets_objects.all.keys}
```

也可以获取到所有文件的详细信息

```hcl
${data.qiniu_buckets_objects.all.key_infos}
```
