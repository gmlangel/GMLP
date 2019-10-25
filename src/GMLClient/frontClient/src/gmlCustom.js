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

function makeArrayNode(v, t, parentName) {
    let d1 = document.createElement("div");
    d1.setAttribute("class", "layui-form-item");
    let l1 = document.createElement("label");
    l1.innerText = v.des;
    l1.setAttribute("title", v.name);
    l1.setAttribute("class", "gml-form-label");//采用多行模式
    let sd1 = document.createElement("div");
    sd1.setAttribute("class", "gml-input-block");
    //遍历v.def_v
    let arr = v.def_v || [];
    let sub = null, sub_tb = null;
    arr.forEach((sv, i) => {
        switch (t) {
            case "string": document.getElementsByName
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
            sd1.appendChild(sub);
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
    let sub = null;
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

function makeMapNode(v, parentName) {
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
    let sub = null;
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
            case "array<string>":
                sub = makeArrayNode(sv, "string", baseKey);
                break;
            case "array<bool>":
                sub = makeArrayNode(sv, "bool", baseKey);
                break;
            case "array<int>":
                sub = makeArrayNode(sv, "int", baseKey);
                break;
            case "array<float>":
                sub = makeArrayNode(sv, "float", baseKey);
                break;
            case "array<map>":
                sub = makeArrayNode(sv, "map", baseKey);
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

function makeFloatNode(v, parentName) {
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
    tb1.value = v.def_v;
    sd1.appendChild(tb1);
    d1.appendChild(l1);
    d1.appendChild(sd1);
    return d1;
}

function makeIntNode(v, parentName) {
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
    tb1.value = v.def_v;
    sd1.appendChild(tb1);
    d1.appendChild(l1);
    d1.appendChild(sd1);
    return d1;
}

function makeBoolNode(v, parentName) {
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
    if (v.def_v == 1) {
        tb1.setAttribute("checked", "checked");
    }
    sd1.appendChild(tb1);
    d1.appendChild(l1);
    d1.appendChild(sd1);
    return d1;
}

function makeStringNode(v, parentName) {
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
    sd1.appendChild(tb1);
    d1.appendChild(l1);
    d1.appendChild(sd1);
    return d1;
}
