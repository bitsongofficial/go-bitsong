"use strict";
var __read = (this && this.__read) || function (o, n) {
    var m = typeof Symbol === "function" && o[Symbol.iterator];
    if (!m) return o;
    var i = m.call(o), r, ar = [], e;
    try {
        while ((n === void 0 || n-- > 0) && !(r = i.next()).done) ar.push(r.value);
    }
    catch (error) { e = { error: error }; }
    finally {
        try {
            if (r && !r.done && (m = i["return"])) m.call(i);
        }
        finally { if (e) throw e.error; }
    }
    return ar;
};
Object.defineProperty(exports, "__esModule", { value: true });
var react_1 = require("react");
var EMPTY_ARRAY = [];
function useUnmount(fn) {
    react_1.useEffect(function () { return fn; }, EMPTY_ARRAY);
}
exports.useUnmount = useUnmount;
function useForceUpdate() {
    var _a = __read(react_1.useState(0), 2), setTick = _a[1];
    var update = react_1.useCallback(function () {
        setTick(function (tick) { return tick + 1; });
    }, []);
    return update;
}
exports.useForceUpdate = useForceUpdate;
function isPlainObject(value) {
    if (!value || typeof value !== "object") {
        return false;
    }
    var proto = Object.getPrototypeOf(value);
    return !proto || proto === Object.prototype;
}
exports.isPlainObject = isPlainObject;
function getSymbol(name) {
    if (typeof Symbol === "function") {
        return Symbol.for(name);
    }
    return "__$mobx-react " + name + "__";
}
exports.getSymbol = getSymbol;
var mockGlobal = {};
function getGlobal() {
    if (typeof window !== "undefined") {
        return window;
    }
    if (typeof global !== "undefined") {
        return global;
    }
    if (typeof self !== "undefined") {
        return self;
    }
    return mockGlobal;
}
exports.getGlobal = getGlobal;
