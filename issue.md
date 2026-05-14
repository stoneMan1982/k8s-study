
# #179748 缺陷修改

venus-install-tools/scripts/base.sh

```bash
function utils::common_install() {
    # 升级过程让服务尽快ready，框架推送完镜像后都只有带有arch的镜像
    PushImages

    # 更新标准yaml
    local componentImage=${VENUS_CLUSTER_REGISTRY_ADDR}/venus/${COMP_NAME}:${VENUS_TARGET_VERSION}
    if [[ ${STANDARD_TEMPLATE} == "yes" ]];then
        sed -e "s#{{imageTemplate}}#${componentImage}#g" -e "s#{{venusVersion}}#${VENUS_TARGET_VERSION}#g" ${COMP_YAML_TEMPLATE_PATH} >${COMP_YAML_PATH} || return 1
    else
        if [ -f ${COMP_YAML_TEMPLATE_PATH} ];then
            cp -f ${COMP_YAML_TEMPLATE_PATH} ${COMP_YAML_PATH} || return 1
        fi
    fi

    # 部署标准yaml前调用组件回调函数做一些处理
    common_callback $CALLBACK_FUNC_INSTALL_PRE

    if [[ "$APP_NAMESPACE" == ""  || "$STANDARD_TEMPLATE" != "yes" ]]; then
        kubectl apply -f ${COMP_YAML_PATH}
    else
        kubectl apply -f ${COMP_YAML_PATH} -n $APP_NAMESPACE
    fi

    # 部署标准yaml后调用组件回调函数做一些处理
    common_callback $CALLBACK_FUNC_INSTALL_POST
    return 0
}
```

venus-plugin-license/pkg/license/license.go

```go
type License struct {
	User    User    `json:"user"`
	Product Product `json:"product"`
	Grants  Grants  `json:"grants"`
	Version string  `json:"version"`
}
```

venus-plugin-license/pkg/license/license.go

```go
func (raw Raw) Decode() (License, error) {
	if len(raw) <= 24 {
		return License{}, errors.New("授权码长度不符合要求")
	}
	lr := string(raw)
	sCiphertext := lr[:len(lr)-24]
	sNonce := lr[len(lr)-24:]
	bCiphertext, err := hex.DecodeString(sCiphertext)
	if err != nil {
		return License{}, errors.New("授权码无效")
	}
	bNonce, err := hex.DecodeString(sNonce)
	if err != nil {
		return License{}, errors.New("授权码无效")
	}
	plaintext, err := aesgcm.AesGcmDecrypt(aesgcm.GetDefaultKeyByte(), bCiphertext, bNonce)
	if err != nil {
		return License{}, errors.New("授权码无效")
	}
	var fp string
	var exp int64
	var master int32
	var node int32
	if err := gob.GobUnmarshal([]byte(plaintext), &fp, &exp, &master, &node); err != nil {
		return License{}, errors.New("授权码无效")
	}
	license := NewDefaultLicense()
	license.Product.FingerPrint = fp
	license.Grants.Expire = time.Unix(exp, 0).In(time.Local)
	license.Grants.NodeLimits.Master = int(master)
	license.Grants.NodeLimits.Node = int(node)

	version := os.Getenv("VENUS_VERSION")
	if version != "" {
		license.Version = version
	}
	return license, nil

}
```