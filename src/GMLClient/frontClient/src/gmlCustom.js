//根据TemplateObj与StrategyContent重新填充Form
function reFillStrategyFormByTemplateObjAndStrategyContent(templateObj, StrategyContent) {

    let arr = templateObj;
    let baseContainer = document.getElementById("StrategyCustomData");
    baseContainer.innerHTML = "";//清空原有内容
    let nodeType = "", tmpValue = "";
    let d1 = null;
    arr.forEach((v, i) => {
        nodeType = v.type;
        tmpValue = StrategyContent[v.name];
        switch (nodeType) {
            case "string":
                d1 = makeStringNode(v, "s_", tmpValue);
                break;
            case "bool":
                d1 = makeBoolNode(v, "s_", tmpValue);
                break;
            case "int":
                d1 = makeIntNode(v, "s_", tmpValue);
                break;
            case "float":
                d1 = makeFloatNode(v, "s_", tmpValue);
                break;
            case "map":
                d1 = makeMapNode(v, "s_", tmpValue);
                break;
            case "array<string>":
                d1 = makeAndRefillArrayNode(v, "string", "s_", tmpValue);
                break;
            case "array<int>":
                d1 = makeAndRefillArrayNode(v, "int", "s_", tmpValue);
                break;
            case "array<float>":
                d1 = makeAndRefillArrayNode(v, "float", "s_", tmpValue);
                break;
            case "array<map>":
                d1 = makeAndRefillArrayNode(v, "map", "s_", tmpValue);
                break;
        }
        if (d1) {
            baseContainer.appendChild(d1);
        }
        d1 = null;
    })

}

function makeAndRefillArrayNode(v, t, parentName, defValue) {
    let d1 = document.createElement("div");
    d1.setAttribute("class", "layui-form-item");
    let l1 = document.createElement("label");
    l1.innerHTML = v.des + '<i class="layui-icon layui-icon-add-circle-fine" style="cursor: pointer;margin-left:5px;font-size: 20px;line-height:20px; color: #009688;" onClick="javascript:void(0);onArrayAdd(this);"  gstruct="' + encodeURI((JSON.stringify(v))) + '" pn="' + parentName + '" gtype="' + t + '"></i>';
    l1.setAttribute("title", v.name);
    l1.setAttribute("class", "gml-form-label");//采用多行模式
    let sd1 = document.createElement("div");
    sd1.setAttribute("class", "gml-input-block");
    //遍历v.def_v
    let sub = null, sub_tb = null;
    let arr = null;
    if (t == "map") {
        if (defValue && v.def_v[0]) {
            //有值，则用值与模板进行填充
            arr = defValue;
            let sv = v.def_v[0];
            arr.forEach((dv, i) => {
                sub = makeAndRefillSuperMapNode(sv, v.name, parentName, dv);
                if (sub) {
                    //添加操作按钮
                    sub_btn = document.createElement("i");
                    sub_btn.setAttribute("class", "layui-icon layui-icon-close-fill");
                    sub_btn.setAttribute("style", "cursor: pointer;font-size: 24px;line-height:38px; color: #ff5722;float:left;");
                    sub_btn.setAttribute("onClick", "javascript:void(0);onArrayDel(this);");
                    sd1.appendChild(sub_btn);

                    sd1.appendChild(sub);
                    sub = null;
                }
            })
        } else {
            //无值，则用模板默认数据填充
            arr = v.def_v || [];
            arr.forEach((sv, i) => {
                sub = makeSuperMapNode(sv, v.name, parentName);
                if (sub) {
                    //添加操作按钮
                    sub_btn = document.createElement("i");
                    sub_btn.setAttribute("class", "layui-icon layui-icon-close-fill");
                    sub_btn.setAttribute("style", "cursor: pointer;font-size: 24px;line-height:38px; color: #ff5722;float:left;");
                    sub_btn.setAttribute("onClick", "javascript:void(0);onArrayDel(this);");
                    sd1.appendChild(sub_btn);

                    sd1.appendChild(sub);
                    sub = null;
                }
            })
        }
    } else {
        arr = defValue || (v.def_v || []);
        arr.forEach((sv, i) => {
            switch (t) {
                case "string":
                    sub = document.createElement("div");
                    sub.setAttribute("class", "gml-input-block");
                    sub.setAttribute("style", "margin-top:5px;");
                    sub_tb = document.createElement("input");
                    sub_tb.setAttribute("type", "text");
                    sub_tb.setAttribute("lay-verify", "required");
                    sub_tb.setAttribute("placeholder", sv);
                    sub_tb.setAttribute("autocomplete", "off");
                    sub_tb.setAttribute("class", "layui-input");
                    sub_tb.setAttribute("name", parentName + v.name);
                    sub_tb.value = sv;
                    sub.appendChild(sub_tb);
                    break;
                case "int":
                    sub = document.createElement("div");
                    sub.setAttribute("class", "gml-input-block");
                    sub.setAttribute("style", "margin-top:5px;");
                    sub_tb = document.createElement("input");
                    sub_tb.setAttribute("type", "text");
                    sub_tb.setAttribute("lay-verify", "required|number");
                    sub_tb.setAttribute("placeholder", sv);
                    sub_tb.setAttribute("autocomplete", "off");
                    sub_tb.setAttribute("class", "layui-input");
                    sub_tb.setAttribute("name", parentName + v.name);
                    sub_tb.value = sv;
                    sub.appendChild(sub_tb);
                    break;
                case "float":
                    sub = document.createElement("div");
                    sub.setAttribute("class", "gml-input-block");
                    sub.setAttribute("style", "margin-top:5px;");
                    sub_tb = document.createElement("input");
                    sub_tb.setAttribute("type", "text");
                    sub_tb.setAttribute("lay-verify", "required|number");
                    sub_tb.setAttribute("placeholder", sv);
                    sub_tb.setAttribute("autocomplete", "off");
                    sub_tb.setAttribute("class", "layui-input");
                    sub_tb.setAttribute("name", parentName + v.name);
                    sub_tb.value = sv;
                    sub.appendChild(sub_tb);
                    break;
            }
            if (sub) {
                //添加操作按钮
                sub_btn = document.createElement("i");
                sub_btn.setAttribute("class", "layui-icon layui-icon-close-fill");
                sub_btn.setAttribute("style", "cursor: pointer;font-size: 24px;line-height:38px; color: #ff5722;float:left;");
                sub_btn.setAttribute("onClick", "javascript:void(0);onArrayDel(this);");
                sd1.appendChild(sub_btn);

                sd1.appendChild(sub);
                sub = null;
            }
        })
    }

    d1.appendChild(l1);
    d1.appendChild(sd1);
    return d1;
}
function makeAndRefillSuperMapNode(v, name, parentName, defaultValue) {
    let sd1 = document.createElement("div");
    sd1.style.borderRight = "3px";
    sd1.style.borderRightColor = "#009688"
    sd1.style.borderRightStyle = "solid";
    let baseKey = parentName + name
    sd1.setAttribute("name", baseKey);
    baseKey += "_"
    //遍历v.def_v
    let arr = v || [];
    let sub = null, nodeType = "", tValue = null;
    arr.forEach((sv, i) => {
        nodeType = sv.type;
        tValue = defaultValue[sv.name]
        switch (nodeType) {
            case "string":
                sub = makeStringNode(sv, baseKey, tValue);
                break;
            case "bool":
                sub = makeBoolNode(sv, baseKey, tValue);
                break;
            case "int":
                sub = makeIntNode(sv, baseKey, tValue);
                break;
            case "float":
                sub = makeFloatNode(sv, baseKey, tValue);
                break;
            case "map":
                sub = makeMapNode(sv, baseKey, tValue);
                break;
        }
        if (sub) {
            sd1.appendChild(sub);
        }
        sub = null;
    })
    return sd1;
}


//根据TemplateObj生成StrategyContent内容
function createStrategyContentByTemplateObj(templateObj) {
    let arr = templateObj;
    let nodeType = "", key = "";
    let d1 = null;
    let resultObject = {};
    arr.forEach((v, i) => {
        nodeType = v.type;
        key = v.name;
        switch (nodeType) {
            case "string":
                d1 = document.getElementById("s_" + key);
                resultObject[key] = d1.value || "";
                break;
            case "bool":
                d1 = document.getElementById("s_" + key);
                resultObject[key] = (!!d1.checked) ? 1 : 0;
                break;
            case "int":
                d1 = document.getElementById("s_" + key);
                resultObject[key] = d1.value ? parseInt(d1.value) : 0;
                break;
            case "float":
                d1 = document.getElementById("s_" + key);
                resultObject[key] = d1.value ? parseFloat(d1.value) : 0;
                break;
            case "map":
                resultObject[key] = createStrategyContentByMapNode(v, "s_");
                break;
            case "array<string>":
                resultObject[key] = [];
                d1 = document.getElementsByName("s_" + key) || [];
                d1.forEach((tmpv, tmpi) => {
                    resultObject[key].push(tmpv.value || "")
                })
                break;
            case "array<int>":
                resultObject[key] = [];
                d1 = document.getElementsByName("s_" + key) || [];
                d1.forEach((tmpv, tmpi) => {
                    resultObject[key].push(tmpv.value ? parseInt(tmpv.value) : 0)
                })
                break;
            case "array<float>":
                resultObject[key] = [];
                d1 = document.getElementsByName("s_" + key) || [];
                d1.forEach((tmpv, tmpi) => {
                    resultObject[key].push(tmpv.value ? parseFloat(tmpv.value) : 0)
                })
                break;
            case "array<map>":
                resultObject[key] = createStrategyContentByArrayMapNode(v, "s_");
                break;
        }
        d1 = null;
    })
    return resultObject;
}

function createStrategyContentByArrayMapNode(v, parentName) {
    let resultArr = []
    let baseName = parentName + v.name;
    let tempMap = {};
    let keys = [];
    let arr = v.def_v || [];
    let nodeKey = "", key = ""
    //获取有map数组的数量
    let baseArr = document.getElementsByName(baseName);
    if (arr.length > 0 && baseArr.length > 0) {
        baseArr.forEach((arg1, arg2) => {
            resultArr.push({});
        })

        //填充数据
        let currentData = arr[0];
        currentData.forEach((sub, i) => {
            // keys.push(subKey);
            // tempMap[subKey] = document.getElementsByName(baseName + "_" + subKey)
            nodeType = sub.type;
            nodeKey = baseName + "_" + sub.name;
            key = sub.name;
            switch (nodeType) {
                case "string":
                    d1 = document.getElementsByName(nodeKey) || [];
                    d1.forEach((d1sv, ti) => {
                        resultArr[ti][key] = d1sv.value || ""
                    })
                    break;
                case "int":
                    d1 = document.getElementsByName(nodeKey) || [];
                    d1.forEach((d1sv, ti) => {
                        resultArr[ti][key] = d1sv.value ? parseInt(d1sv.value) : 0;
                    })
                    break;
                case "float":
                    d1 = document.getElementsByName(nodeKey) || [];
                    d1.forEach((d1sv, ti) => {
                        resultArr[ti][key] = d1sv.value ? parseFloat(d1sv.value) : 0;
                    })
                    break;
                case "bool":
                    d1 = document.getElementsByName(nodeKey) || [];
                    d1.forEach((d1sv, ti) => {
                        resultArr[ti][key] = (!!d1sv.checked) ? 1 : 0;
                    })
                    break;
            }
        })
    }

    return resultArr;
}

function createStrategyContentByMapNode(v, parentName) {
    let resultObject = {};
    let baseName = parentName + v.name;
    baseName += "_";
    //遍历v.def_v
    let arr = v.def_v || [];
    let key = "", nodeKey = "";
    let d1 = null;
    arr.forEach((sv, i) => {
        nodeType = sv.type;
        key = sv.name;
        nodeKey = baseName + sv.name
        switch (nodeType) {
            case "string":
                d1 = document.getElementById(nodeKey);
                resultObject[key] = d1.value || "";
                break;
            case "bool":
                d1 = document.getElementById(nodeKey);
                resultObject[key] = (!!d1.checked) ? 1 : 0;
                break;
            case "int":
                d1 = document.getElementById(nodeKey);
                resultObject[key] = d1.value ? parseInt(d1.value) : 0;
                break;
            case "float":
                d1 = document.getElementById(nodeKey);
                resultObject[key] = d1.value ? parseFloat(d1.value) : 0;
                break;
            case "map":
                resultObject[key] = createStrategyContentByMapNode(sv, nodeKey);
                break;
            case "array<string>":
                resultObject[key] = [];
                d1 = document.getElementsByName(nodeKey) || [];
                d1.forEach((tmpv, tmpi) => {
                    resultObject[key].push(tmpv.value || "")
                })
                break;
            case "array<int>":
                resultObject[key] = [];
                d1 = document.getElementsByName(nodeKey) || [];
                d1.forEach((tmpv, tmpi) => {
                    resultObject[key].push(tmpv.value || 0)
                })
                break;
            case "array<float>":
                resultObject[key] = [];
                d1 = document.getElementsByName(nodeKey) || [];
                d1.forEach((tmpv, tmpi) => {
                    resultObject[key].push(tmpv.value || 0)
                })
                break;
            case "array<map>":
                resultObject[key] = createStrategyContentByArrayMapNode(sv, baseName);
                break;
        }
        d1 = null;
    })
    return resultObject;
}

//生成并重新填充panel_addStrategy表单的StrategyCustomData容器 的内容
function makeAndRefillStrategyForm(templateObj) {
    let arr = templateObj;
    let baseContainer = document.getElementById("StrategyCustomData");
    baseContainer.innerHTML = "";//清空原有内容
    let nodeType = ""
    let d1 = null;
    arr.forEach((v, i) => {
        nodeType = v.type;
        switch (nodeType) {
            case "string":
                d1 = makeStringNode(v, "s_");
                break;
            case "bool":
                d1 = makeBoolNode(v, "s_");
                break;
            case "int":
                d1 = makeIntNode(v, "s_");
                break;
            case "float":
                d1 = makeFloatNode(v, "s_");
                break;
            case "map":
                d1 = makeMapNode(v, "s_");
                break;
            case "array<string>":
                d1 = makeArrayNode(v, "string", "s_");
                break;
            case "array<int>":
                d1 = makeArrayNode(v, "int", "s_");
                break;
            case "array<float>":
                d1 = makeArrayNode(v, "float", "s_");
                break;
            case "array<map>":
                d1 = makeArrayNode(v, "map", "s_");
                break;
        }
        if (d1) {
            baseContainer.appendChild(d1);
        }
        d1 = null;
    })
}

/**
 * 响应用户点击
 * 向数组中添加元素
*/
function onArrayAdd(sender) {

    let gstruct = sender.getAttribute("gstruct");
    if (!gstruct)
        return;
    gstruct = decodeURI(gstruct);
    let parentName = sender.getAttribute("pn");//获取父级名称
    let t = sender.getAttribute("gtype");
    let v = JSON.parse(gstruct)//获取模板数据结构
    if (v && parentName) {
        console.log("arrayAdd===>v=", v, "  parentName=", parentName);
        let sd1 = sender.parentNode.parentNode.children[1];
        if (!sd1)
            return;
        //遍历v.def_v
        let arr = (v.def_v && v.def_v.length > 0) ? [v.def_v[0]] : [];
        let sub = null, sub_tb = null, sub_btn = null;
        arr.forEach((sv, i) => {
            switch (t) {
                case "string":
                    sub = document.createElement("div");
                    sub.setAttribute("class", "gml-input-block");
                    sub.setAttribute("style", "margin-top:5px;");
                    sub_tb = document.createElement("input");
                    sub_tb.setAttribute("type", "text");
                    sub_tb.setAttribute("lay-verify", "required");
                    sub_tb.setAttribute("placeholder", sv);
                    sub_tb.setAttribute("autocomplete", "off");
                    sub_tb.setAttribute("class", "layui-input");
                    sub_tb.setAttribute("name", parentName + v.name);
                    sub.appendChild(sub_tb);
                    break;
                case "int":
                    sub = document.createElement("div");
                    sub.setAttribute("class", "gml-input-block");
                    sub.setAttribute("style", "margin-top:5px;");
                    sub_tb = document.createElement("input");
                    sub_tb.setAttribute("type", "text");
                    sub_tb.setAttribute("lay-verify", "required|number");
                    sub_tb.setAttribute("placeholder", sv);
                    sub_tb.setAttribute("autocomplete", "off");
                    sub_tb.setAttribute("class", "layui-input");
                    sub_tb.setAttribute("name", parentName + v.name);
                    sub.appendChild(sub_tb);
                    break;
                case "float":
                    sub = document.createElement("div");
                    sub.setAttribute("class", "gml-input-block");
                    sub.setAttribute("style", "margin-top:5px;");
                    sub_tb = document.createElement("input");
                    sub_tb.setAttribute("type", "text");
                    sub_tb.setAttribute("lay-verify", "required|number");
                    sub_tb.setAttribute("placeholder", sv);
                    sub_tb.setAttribute("autocomplete", "off");
                    sub_tb.setAttribute("class", "layui-input");
                    sub_tb.setAttribute("name", parentName + v.name);
                    sub.appendChild(sub_tb);
                    break;
                case "map":
                    sub = makeSuperMapNode(sv, v.name, parentName);
                    break;
            }
            if (sub) {
                //添加操作按钮
                sub_btn = document.createElement("i");
                sub_btn.setAttribute("class", "layui-icon layui-icon-close-fill");
                sub_btn.setAttribute("style", "cursor: pointer;font-size: 24px;line-height:38px; color: #ff5722;float:left;");
                sub_btn.setAttribute("onClick", "javascript:void(0);onArrayDel(this);");
                sd1.appendChild(sub_btn);

                sd1.appendChild(sub);//添加元素
                sub = null;
            }
        })
        //刷新表单UI的渲染
        form.render();
    }
}

/**
 * 响应用户点击
 * 从数组中删除元素
*/
function onArrayDel(sender) {
    //删除用户信息
    layer.open({
        title: "删除选项",
        content: "确定要删除吗？",
        area: ["300px", "200px"],
        resize: false,
        btn: ["确定", "取消"],
        yes: (index, obj) => {
            //确定操作
            let arr = sender.parentNode.children;
            if (arr) {
                let j = arr.length;
                let waitDelNode = null;
                for (i = 0; i < j; i++) {
                    if (arr[i] == sender) {
                        waitDelNode = arr[i + 1];
                        if (waitDelNode)
                            sender.parentNode.removeChild(waitDelNode);//移除 item项
                        sender.parentNode.removeChild(sender);//移除 “删除按钮”
                        layer.alert("删除成功")
                        break;
                    }
                }
            }
        },
        btn2: (index, obj) => {
            //取消操作,不用谢任何代码，即可默认关闭layer
        },
        btnAlign: "c"/*居中对齐*/
    })

}


function makeArrayNode(v, t, parentName) {
    let d1 = document.createElement("div");
    d1.setAttribute("class", "layui-form-item");
    let l1 = document.createElement("label");

    l1.innerHTML = v.des + '<i class="layui-icon layui-icon-add-circle-fine" style="cursor: pointer;margin-left:5px;font-size: 20px;line-height:20px; color: #009688;" onClick="javascript:void(0);onArrayAdd(this);" gstruct="' + encodeURI((JSON.stringify(v))) + '" pn="' + parentName + '" gtype="' + t + '"></i>';
    l1.setAttribute("title", v.name);
    l1.setAttribute("class", "gml-form-label");//采用多行模式
    let sd1 = document.createElement("div");
    sd1.setAttribute("class", "gml-input-block");
    //遍历v.def_v
    let arr = v.def_v || [];
    let sub = null, sub_tb = null, sub_btn = null;
    arr.forEach((sv, i) => {
        switch (t) {
            case "string":
                sub = document.createElement("div");
                sub.setAttribute("class", "gml-input-block");
                sub.setAttribute("style", "margin-top:5px;");
                sub_tb = document.createElement("input");
                sub_tb.setAttribute("type", "text");
                sub_tb.setAttribute("lay-verify", "required");
                sub_tb.setAttribute("placeholder", sv);
                sub_tb.setAttribute("autocomplete", "off");
                sub_tb.setAttribute("class", "layui-input");
                sub_tb.setAttribute("name", parentName + v.name);
                sub.appendChild(sub_tb);
                break;
            case "int":
                sub = document.createElement("div");
                sub.setAttribute("class", "gml-input-block");
                sub.setAttribute("style", "margin-top:5px;");
                sub_tb = document.createElement("input");
                sub_tb.setAttribute("type", "text");
                sub_tb.setAttribute("lay-verify", "required|number");
                sub_tb.setAttribute("placeholder", sv);
                sub_tb.setAttribute("autocomplete", "off");
                sub_tb.setAttribute("class", "layui-input");
                sub_tb.setAttribute("name", parentName + v.name);
                sub.appendChild(sub_tb);
                break;
            case "float":
                sub = document.createElement("div");
                sub.setAttribute("class", "gml-input-block");
                sub.setAttribute("style", "margin-top:5px;");
                sub_tb = document.createElement("input");
                sub_tb.setAttribute("type", "text");
                sub_tb.setAttribute("lay-verify", "required|number");
                sub_tb.setAttribute("placeholder", sv);
                sub_tb.setAttribute("autocomplete", "off");
                sub_tb.setAttribute("class", "layui-input");
                sub_tb.setAttribute("name", parentName + v.name);
                sub.appendChild(sub_tb);
                break;
            case "map":
                sub = makeSuperMapNode(sv, v.name, parentName);
                break;
        }
        if (sub) {
            //添加操作按钮
            sub_btn = document.createElement("i");
            sub_btn.setAttribute("class", "layui-icon layui-icon-close-fill");
            sub_btn.setAttribute("style", "cursor: pointer;font-size: 24px;line-height:38px; color: #ff5722;float:left;");
            sub_btn.setAttribute("onClick", "javascript:void(0);onArrayDel(this);");
            sd1.appendChild(sub_btn);

            sd1.appendChild(sub);//添加元素
            sub = null;
        }
    })
    d1.appendChild(l1);
    d1.appendChild(sd1);
    return d1;
}
function makeSuperMapNode(v, name, parentName) {
    let sd1 = document.createElement("div");
    sd1.style.borderRight = "3px";
    sd1.style.borderRightColor = "#009688"
    sd1.style.borderRightStyle = "solid";
    let baseKey = parentName + name
    sd1.setAttribute("name", baseKey);
    baseKey += "_"
    //遍历v.def_v
    let arr = v || [];
    let sub = null, nodeType = "";
    arr.forEach((sv, i) => {
        nodeType = sv.type;
        switch (nodeType) {
            case "string":
                sub = makeStringNode(sv, baseKey);
                break;
            case "bool":
                sub = makeBoolNode(sv, baseKey);
                break;
            case "int":
                sub = makeIntNode(sv, baseKey);
                break;
            case "float":
                sub = makeFloatNode(sv, baseKey);
                break;
            case "map":
                sub = makeMapNode(sv, baseKey);
                break;
        }
        if (sub) {
            sd1.appendChild(sub);
        }
        sub = null;
    })
    return sd1;
}

function makeMapNode(v, parentName, defValue) {
    let arrayType = (!!defValue) ? "refill" : "make";
    let def = defValue || {};
    let d1 = document.createElement("div");
    d1.setAttribute("class", "layui-form-item");
    let baseKey = parentName + v.name
    d1.setAttribute("id", baseKey);
    baseKey += "_"
    let l1 = document.createElement("label");
    l1.innerText = v.des;
    l1.setAttribute("title", v.name);
    l1.setAttribute("class", "gml-form-label");//采用多行模式
    let sd1 = document.createElement("div");
    sd1.setAttribute("class", "gml-input-block");
    //遍历v.def_v
    let arr = v.def_v || [];
    let sub = null, tValue = null;
    arr.forEach((sv, i) => {
        nodeType = sv.type;
        tValue = def[sv.name];
        switch (nodeType) {
            case "string":
                sub = makeStringNode(sv, baseKey, tValue);
                break;
            case "bool":
                sub = makeBoolNode(sv, baseKey, tValue);
                break;
            case "int":
                sub = makeIntNode(sv, baseKey, tValue);
                break;
            case "float":
                sub = makeFloatNode(sv, baseKey, tValue);
                break;
            case "map":
                sub = makeMapNode(sv, baseKey, tValue);
                break;
            case "array<string>":
                if (arrayType == "refill")
                    sub = makeAndRefillArrayNode(sv, "string", baseKey, tValue);//创建并用值进行填充
                else
                    sub = makeArrayNode(sv, "string", baseKey);//创建并用 模板默认值填充
                break;
            case "array<bool>":
                if (arrayType == "refill")
                    sub = makeAndRefillArrayNode(sv, "bool", baseKey, tValue);//创建并用值进行填充
                else
                    sub = makeArrayNode(sv, "bool", baseKey);//创建并用 模板默认值填充
                break;
            case "array<int>":
                if (arrayType == "refill")
                    sub = makeAndRefillArrayNode(sv, "int", baseKey, tValue);//创建并用值进行填充
                else
                    sub = makeArrayNode(sv, "int", baseKey);//创建并用 模板默认值填充
                break;
            case "array<float>":
                if (arrayType == "refill")
                    sub = makeAndRefillArrayNode(sv, "float", baseKey, tValue);//创建并用值进行填充
                else
                    sub = makeArrayNode(sv, "float", baseKey);//创建并用 模板默认值填充
                break;
            case "array<map>":
                if (arrayType == "refill")
                    sub = makeAndRefillArrayNode(sv, "map", baseKey, tValue);//创建并用值进行填充
                else
                    sub = makeArrayNode(sv, "map", baseKey);//创建并用 模板默认值填充
                break;
        }
        if (sub) {
            sd1.appendChild(sub);
        }
        sub = null;
    })
    d1.appendChild(l1);
    d1.appendChild(sd1);
    return d1;
}

function makeFloatNode(v, parentName, defValue) {
    let d1 = document.createElement("div");
    d1.setAttribute("class", "layui-form-item");
    let l1 = document.createElement("label");
    l1.innerText = v.des;
    l1.setAttribute("title", v.name);
    let sd1 = document.createElement("div");
    if (v.des && v.des.length > 4) {
        l1.setAttribute("class", "gml-form-label");//采用多行模式
        sd1.setAttribute("class", "gml-input-block");
    }
    else {
        l1.setAttribute("class", "layui-form-label");//采用单行模式
        sd1.setAttribute("class", "layui-input-block");
    }
    let tb1 = document.createElement("input");
    let min = v.min_v || 0;
    let max = v.max_v || 1.0;
    tb1.setAttribute("type", "text");
    tb1.setAttribute("lay-verify", "required|qujianNumber");
    tb1.setAttribute("placeholder", v.def_v);
    tb1.setAttribute("autocomplete", "off");
    tb1.setAttribute("class", "layui-input");
    tb1.setAttribute("id", parentName + v.name);
    tb1.setAttribute("name", parentName + v.name);
    tb1.setAttribute("min", min);
    tb1.setAttribute("max", max);
    if (defValue)
        tb1.value = defValue;
    sd1.appendChild(tb1);
    d1.appendChild(l1);
    d1.appendChild(sd1);
    return d1;
}

function makeIntNode(v, parentName, defValue) {
    let d1 = document.createElement("div");
    d1.setAttribute("class", "layui-form-item");
    let l1 = document.createElement("label");
    l1.innerText = v.des;
    l1.setAttribute("title", v.name);
    let sd1 = document.createElement("div");
    if (v.des && v.des.length > 4) {
        l1.setAttribute("class", "gml-form-label");//采用多行模式
        sd1.setAttribute("class", "gml-input-block");
    }
    else {
        l1.setAttribute("class", "layui-form-label");//采用单行模式
        sd1.setAttribute("class", "layui-input-block");
    }
    let tb1 = document.createElement("input");
    let min = v.min_v || 0;
    let max = v.max_v || 100000000;
    tb1.setAttribute("type", "text");
    tb1.setAttribute("lay-verify", "required|qujianNumber");
    tb1.setAttribute("placeholder", v.def_v);
    tb1.setAttribute("autocomplete", "off");
    tb1.setAttribute("class", "layui-input");
    tb1.setAttribute("id", parentName + v.name);
    tb1.setAttribute("name", parentName + v.name);
    tb1.setAttribute("min", min);
    tb1.setAttribute("max", max);
    if (defValue)
        tb1.value = defValue;
    sd1.appendChild(tb1);
    d1.appendChild(l1);
    d1.appendChild(sd1);
    return d1;
}

function makeBoolNode(v, parentName, defValue) {
    let d1 = document.createElement("div");
    d1.setAttribute("class", "layui-form-item");
    let l1 = document.createElement("label");
    l1.innerText = v.des;
    l1.setAttribute("title", v.name);
    let sd1 = document.createElement("div");
    if (v.des && v.des.length > 40) {
        l1.setAttribute("class", "gml-form-label");//采用多行模式
        sd1.setAttribute("class", "gml-input-block");
    }
    else {
        l1.setAttribute("class", "gml-form-label");//采用单行模式
        l1.setAttribute("style", "float:left;width:auto;padding-right:15px;");//采用单行模式
        sd1.setAttribute("class", "layui-input-block");
    }
    let tb1 = document.createElement("input");
    tb1.setAttribute("type", "checkbox");
    tb1.setAttribute("lay-skin", "switch");
    tb1.setAttribute("id", parentName + v.name);
    tb1.setAttribute("name", parentName + v.name);
    if (defValue != undefined) {
        tb1.checked = defValue == 1;
    } else if (v.def_v == 1) {
        tb1.setAttribute("checked", "checked");
    }
    sd1.appendChild(tb1);
    d1.appendChild(l1);
    d1.appendChild(sd1);
    return d1;
}

function makeStringNode(v, parentName, defValue) {
    let d1 = document.createElement("div");
    d1.setAttribute("class", "layui-form-item");
    let l1 = document.createElement("label");
    l1.innerText = v.des;
    l1.setAttribute("title", v.name);
    let sd1 = document.createElement("div");
    if (v.des && v.des.length > 4) {
        l1.setAttribute("class", "gml-form-label");//采用多行模式
        sd1.setAttribute("class", "gml-input-block");
    }
    else {
        l1.setAttribute("class", "layui-form-label");//采用单行模式
        sd1.setAttribute("class", "layui-input-block");
    }
    let tb1 = document.createElement("input");
    tb1.setAttribute("type", "text");
    tb1.setAttribute("lay-verify", "required");
    tb1.setAttribute("placeholder", v.def_v);
    tb1.setAttribute("autocomplete", "off");
    tb1.setAttribute("class", "layui-input");
    tb1.setAttribute("id", parentName + v.name);
    tb1.setAttribute("name", parentName + v.name);
    if (defValue)
        tb1.value = defValue;
    sd1.appendChild(tb1);
    d1.appendChild(l1);
    d1.appendChild(sd1);
    return d1;
}

/**
 * 
 * 生成 2019-10-11 18:30:00的日期字符传
*/
function makeDateTimeStr(date, fmt) {
    if (!date) {
        return "1900-01-01 00:00:00"
    }
    var o = {
        "M+": date.getMonth() + 1,
        "d+": date.getDate(),
        "H+": date.getHours(),
        "m+": date.getMinutes(),
        "s+": date.getSeconds(),
        "S+": date.getMilliseconds()
    };
    //因为date.getFullYear()出来的结果是number类型的,所以为了让结果变成字符串型，下面有两种方法：
    if (/(y+)/.test(fmt)) {
        //第一种：利用字符串连接符“+”给date.getFullYear()+""，加一个空字符串便可以将number类型转换成字符串。
        fmt = fmt.replace(RegExp.$1, (date.getFullYear() + "").substr(4 - RegExp.$1.length));
    }
    for (var k in o) {
        if (new RegExp("(" + k + ")").test(fmt)) {
            //第二种：使用String()类型进行强制数据类型转换String(date.getFullYear())，这种更容易理解。
            fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(String(o[k]).length)));
        }
    }
    return fmt;
};